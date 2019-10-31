package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
)

func applyLoadTemplates(buf []byte, fn string, cfg *SK8config, prefix string) ([]byte, error) {
	if bytes.ContainsAny(buf, "{{") {
		if cfg != nil {
			log.Debugf("%sTemplating defaults in %s", prefix, fn)
			tmpl := template.New(fn).Funcs(FuncMap())
			tmpl, err := tmpl.Parse(string(buf))
			if err != nil {
				return nil, err
			}
			var bufStream bytes.Buffer
			err = tmpl.Execute(&bufStream, cfg)
			if err != nil {
				return nil, err
			}
			buf = bufStream.Bytes()
		} else {
			log.Debugf("%sTemplating generics in %s", prefix, fn)
			tmpl := template.New(fn).Funcs(FuncMap())
			tmpl, err := tmpl.Parse(string(buf))
			if err != nil {
				return nil, err
			}
			var bufStream bytes.Buffer
			err = tmpl.Execute(&bufStream, nil)
			if err != nil {
				return nil, err
			}
			buf = bufStream.Bytes()
		}
	}
	return buf, nil
}

func splitDocuments(input []byte) [][]byte {
	var bufs [][]byte

	currbuf := []byte{}
	rd := bytes.NewReader(input)
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "---") {
			if len(currbuf) > 0 {
				bufs = append(bufs, currbuf)
				currbuf = []byte{}
			}
		} else {
			currbuf = append(currbuf, line...)
			currbuf = append(currbuf, '\n')
		}
	}
	if len(currbuf) > 0 {
		bufs = append(bufs, currbuf)
	}

	return bufs
}

func loadFile(fn string, depth int, cfg *SK8config) ([]*SK8config, error) {

	var result []*SK8config

	if strings.HasPrefix(fn, "?") {
		fn = fn[1:]
		if _, err := os.Stat(fn); err != nil {
			log.Debugf("Skipping optional parent: %s", fn)
			return nil, nil
		}
	}

	prefix := strings.Repeat("| ", depth)
	log.Debugf("%sLoading from %s", prefix, fn)
	prefix = prefix + "+-"

	filedata, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, fmt.Errorf("readfile error %q: %q", fn, err.Error())
	}

	filedata, err = applyLoadTemplates(filedata, fn, cfg, prefix)
	if err != nil {
		return nil, fmt.Errorf("apply template-error %q: %q", fn, err.Error())
	}

	bufs := splitDocuments(filedata)
	log.Infof("%s splitted into %d documents", prefix, len(bufs))

	for _, buf := range bufs {

		var o = SK8config{}
		var raw interface{}

		//err = yaml_mapstr.Unmarshal(buf, &raw)
		err = yamlUnmarshal(buf, &raw)
		if err != nil {
			fmt.Println(string(buf))
			return nil, fmt.Errorf("yaml parse error %q: %q", fn, err.Error())
		}

		buf, err = json.Marshal(raw)
		json.Unmarshal(buf, &o)
		if err != nil {
			return nil, err
		}

		if depth < 5 && len(o.Parents) > 0 {
			parent := &SK8config{}
			parents := o.Parents
			o.Parents = nil
			log.Debugf("About to load parents: %v", parents)

			for _, pn := range parents {
				log.Debugf("%sInherit from %s", prefix, pn)
				gps, err := loadFile(pn, depth+1, cfg)
				if err != nil {
					return nil, err
				}
				if len(gps) != 1 {
					return nil, fmt.Errorf("multiple documents in parent %s", pn)
				}
				gp := gps[0]
				if args.Debug && args.Verbose {
					log.Debugf("Dump of intermediate parent: %s", pn)
					dump(gp)
				}
				parent.mergeWith(gp)
			}
			parent.mergeWith(&o)
			o = *parent
		}

		result = append(result, &o)
	}
	return result, nil
}

func (cfg *SK8config) mergeWith(copyfrom *SK8config) *SK8config {

	if copyfrom == nil {
		return cfg
	}

	features := make(map[string]bool)
	for _, f := range cfg.Features {
		features[f] = true
	}
	for _, f := range copyfrom.Features {
		features[f] = true
	}

	v, _ := json.Marshal(copyfrom)
	err := json.Unmarshal(v, cfg)
	if err != nil {
		panic(err)
	}
	var list []string
	for f := range features {
		list = append(list, f)
	}
	cfg.Features = list
	return cfg
}

func (cfg *SK8config) fixFile() error {

	if cfg.Image == "" {
		cfg.Image = cfg.Name
	}

	cwd, _ := os.Getwd()
	for t, f := range cfg.Templates {
		if !path.IsAbs(f) {
			f2 := path.Join(cwd, f)
			cfg.Templates[t] = f2
		}
	}

	return nil
}
