package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"yaml_mapstr"
)

func loadFile(fn string, depth int, cfg *SK8config) (*SK8config, error) {

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

	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, fmt.Errorf("ReadFile-error %q: %q", fn, err.Error())
	}

	if cfg != nil && bytes.ContainsAny(buf, "{{") {
		log.Debugf("%sTemplating defaults in %s", prefix, fn)
		tmpl := template.New(fn)
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
	}

	var o SK8config
	var raw interface{}

	err = yaml_mapstr.Unmarshal(buf, &raw)
	if err != nil {
		return nil, fmt.Errorf("yaml-parse error %q: %q", fn, err.Error())
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
			gp, err := loadFile(pn, depth+1, cfg)
			if err != nil {
				return nil, err
			}
			if args.Debug && args.Verbose {
				log.Debugf("Dump of intermediate parent: %s", pn)
				dump(gp)
			}
			parent.mergeWith(gp)
		}
		parent.mergeWith(&o)
		o = *parent
	}

	return &o, nil
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
