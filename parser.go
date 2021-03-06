package main

import (
	"bufio"
	"bytes"
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

func isKubectl(buf []byte) bool {
	if len(buf) < 10 {
		return false
	}
	first := string(buf[0:9])
	return first == "#!kubectl"
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
	if len(bufs) > 1 {
		log.Infof("%s %s splitted into %d documents", prefix, fn, len(bufs))
	}

	for _, buf := range bufs {
		var o = SK8config{}

		if isKubectl(buf) {
			o.cfgType = typeKubectl
			o.RawYAML = buf
		} else {
			o.cfgType = typeSK8
			var raw interface{}

			//err = yaml_mapstr.Unmarshal(buf, &raw)
			err = yaml_Unmarshal(buf, &raw)
			if err != nil {
				fmt.Println(string(buf))
				return nil, fmt.Errorf("yaml parse error %q: %q", fn, err.Error())
			}

			buf2, err := json_Marshal(raw)
			if err != nil {
				fmt.Println(string(buf))
				return nil, fmt.Errorf("json_Marshal(raw) error: %v", err)
			}
			err = json_Unmarshal(buf2, &o)
			if err != nil {
				return nil, fmt.Errorf("json_Unmarshal(buf2,o) error: %v", err)
			}
			if o.Kind != "" {
				if o.RawMetadata == nil {
					log.Fatalf("Metadata missing in raw object: %s in %s", o.Kind, fn)
				}
				if o.RawMetadata.Name == "" {
					log.Fatalf("Name missing in raw object: %s in %s", o.Kind, fn)
				}
				if o.RawMetadata.Namespace == "" {
					log.Fatalf("Namespace missing in raw object: %s %q in %s ", o.Kind, o.RawMetadata.Name, fn)
				}
				o = SK8config{
					cfgType: typeKubectl,
					RawYAML: buf,
					Kind:    o.Kind,
				}
			} else {

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
						if len(gps) > 0 {
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
					}
					parent.mergeWith(&o)
					o = *parent
					o.cfgType = typeSK8
				}
			}
		}
		result = append(result, &o)
	}
	return result, nil
}

func (cfg *SK8config) mergeWith(copyfrom *SK8config) *SK8config {

	if copyfrom == nil {
		return cfg
	}

	override := cfg.Override
	cfg.Override = nil

	features := make(map[string]bool)
	for _, f := range cfg.Features {
		features[f] = true
	}
	for _, f := range copyfrom.Features {
		features[f] = true
	}
	containers := getContainers(cfg.Containers, copyfrom.Containers)
	volumes := getVolumes(cfg.Volume, copyfrom.Volume)

	overrideEnv(cfg, copyfrom)

	v, _ := json_Marshal(copyfrom)
	err := json_Unmarshal(v, cfg)
	if err != nil {
		panic(err)
	}
	var list []string
	for f := range features {
		list = append(list, f)
	}
	cfg.Features = list
	cfg.Containers = containers
	cfg.Volume = volumes

	if override != nil {
		return cfg.mergeWith(override)
	}

	return cfg
}

func getContainers(cfg, copyfrom []Container) []Container {
	result := make([]Container, 0, len(cfg)+len(copyfrom))
	names := make(map[string]Container)

	for _, c := range cfg {
		names[c.Name] = c
	}
	for _, c := range copyfrom {
		if org, found := names[c.Name]; found {
			v, _ := json_Marshal(c) // merge
			err := json_Unmarshal(v, &org)
			if err != nil {
				log.Fatalf("Container '%s' could not be merged with override: %s", c.Name, err)
			}
			log.Infof("Container '%s' already exists, merge attempted", c.Name)
			names[c.Name] = org
		} else {
			names[c.Name] = c
		}
	}
	for _, c := range names {
		result = append(result, c)
	}
	return result
}

func getVolumes(cfg, copyfrom []VolumeType) []VolumeType {
	result := make([]VolumeType, 0, len(cfg)+len(copyfrom))
	names := make(map[string]VolumeType)

	for _, c := range cfg {
		names[c.Name] = c
	}
	for _, c := range copyfrom {
		if _, found := names[c.Name]; found {
			log.Fatalf("Volume '%s' already exists, cannot overwrite", c.Name)
		}
		names[c.Name] = c
	}
	for _, c := range names {
		result = append(result, c)
	}
	return result
}

func overrideEnv(cfg, copyfrom *SK8config) {
	for k := range copyfrom.Env.Values {
		delete(cfg.Env.Config, k)
		delete(cfg.Env.Secret, k)
	}
	for k := range copyfrom.Env.Config {
		delete(cfg.Env.Values, k)
		delete(cfg.Env.Secret, k)
	}
	for k := range copyfrom.Env.Secret {
		delete(cfg.Env.Values, k)
		delete(cfg.Env.Config, k)
	}
}

func (cfg *SK8config) fixFile() {

	if cfg.Image == "" {
		cfg.Image = cfg.Name
	}

	log.Debugf("%s: fixFile()", cfg.Name)
	cwd, _ := os.Getwd()
	for t, f := range cfg.Templates {
		if !path.IsAbs(f) {
			f2 := path.Join(cwd, f)
			log.Debugf("change template %q -> %q", f, f2)
			cfg.Templates[t] = f2
		} else {
			log.Debugf("template %q is already absolute", f)
		}
	}
}
