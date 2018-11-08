package main

import (
	"io/ioutil"
	"os"
	"path"
)

// LoadInput -
func LoadInput() []*SK8config {
	var cfgs []*SK8config

	currdir, _ := os.Getwd()

	for _, f := range args.Files {
		dir := path.Dir(f)
		item := path.Base(f)
		log.Debugf("File %s: folder %s, name %s", f, dir, item)
		os.Chdir(dir)

		o, err := loadFile(item, 0, nil)
		if err != nil {
			log.Fatal(err.Error())
			return nil
		}

		base := &SK8config{}
		dirname := ".sk8"
		files, err := ioutil.ReadDir(dirname)
		if err == nil {
			for _, def := range files {
				defName := path.Join(dirname, def.Name())
				log.Debugf("Import defaults: %s", defName)
				override, err := loadFile(defName, 0, o)
				if err != nil {
					panic(err)
				}
				base.mergeWith(override)
			}
		}

		o = base.mergeWith(o)

		err = o.fixFile()
		if err != nil {
			panic(err)
		}

		cfgs = append(cfgs, o)
	}

	os.Chdir(currdir)

	return cfgs
}
