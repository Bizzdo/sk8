package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
	"yaml_mapstr"
)

func loadFile(fn string, depth int, cfg *sk8config) (*sk8config, error) {
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

	var o sk8config
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
		parent := &sk8config{}
		parents := o.Parents
		o.Parents = nil
		log.Debugf("About to load parents: %v", parents)

		for _, pn := range parents {
			log.Debugf("%sInherit from %s", prefix, pn)
			gp, err := loadFile(pn, depth+1, cfg)
			if err != nil {
				return nil, err
			}
			parent.mergeWith(gp)
		}
		parent.mergeWith(&o)
		o = *parent
	}

	return &o, nil
}

func (cfg *sk8config) mergeWith(copyfrom *sk8config) *sk8config {
	v, _ := json.Marshal(copyfrom)
	err := json.Unmarshal(v, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func (cfg *sk8config) fixFile() error {

	if cfg.Image == "" {
		cfg.Image = cfg.Name
	}

	if cfg.Port == 0 {
		//defaultPort := 80
		//		o.Port = defaultPort
	}

	return nil
}
