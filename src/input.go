package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func findDefaults(folder string) (string, bool) {
	// log.Debugf("? findDefaults( %s ) ", folder)

	if len(folder) <= 1 {
		return "", false
	}

	foo := 0
	for folder != "/" {
		foo++
		if foo > 30 {
			// log.Fatal("to deep folders, or possible loop")
			return "", false
		}

		folder = path.Dir(folder)
		// log.Debugf("? checking for .sk8default in %q", folder)
		folder2 := path.Join(folder, ".sk8default")
		fi, err := os.Stat(folder2)
		if err != nil {
			// log.Debugf("? -> %s", err.Error())
		} else {
			// log.Debugf("? => %+v", fi)
			if fi != nil && fi.IsDir() {
				// log.Debugf("Found .sk8default in %q as %q", folder, fi.Name())
				return folder2, true
			}
		}
	}
	return "", false
}

func mergeSettingsFrom(dirname string, base, o *SK8config) {
	files, err := ioutil.ReadDir(dirname)
	if err == nil {
		for _, def := range files {
			if strings.HasSuffix(def.Name(), ".yaml") {
				defName := path.Join(dirname, def.Name())
				log.Debugf("Import defaults: %s", defName)
				overrides, err := loadFile(defName, 0, o)
				if err != nil {
					panic(err)
				}
				if len(overrides) > 1 {
					panic(fmt.Errorf("multiple documents in %s (%s)", def.Name(), dirname))
				}
				base.mergeWith(overrides[0])
			}
		}
	}
}

// LoadInput -
func LoadInput() []*SK8config {
	var cfgs []*SK8config

	currdir, _ := os.Getwd()

	for _, f := range args.Files {
		dir := path.Dir(f)
		item := path.Base(f)
		log.Debugf("File %s: folder %s, name %s", f, dir, item)
		_ = os.Chdir(dir)

		docs, err := loadFile(item, 0, nil)
		if err != nil {
			log.Fatal(err.Error())
			return nil
		}

		fullpath, err := os.Getwd()
		if err == nil {
			fullpath = path.Join(fullpath, item)
		}

		for _, o := range docs {
			log.Debugf("doc(%s) = %v", o.Namespace, o.cfgType)
			if o.cfgType == typeSK8 {
				base := &SK8config{}

				if defaults, found := findDefaults(fullpath); found {
					mergeSettingsFrom(defaults, base, o)
				}

				mergeSettingsFrom(".sk8", base, o)

				o = base.mergeWith(o)

				overrideFile := strings.Replace(item, ".yaml", ".override", -1)
				if _, err := os.Stat(overrideFile); err == nil {
					log.Debugf("Applying override: %s", overrideFile)
					overrides, err := loadFile(overrideFile, 0, o)
					if err != nil {
						panic(err)
					}
					if len(overrides) > 1 {
						panic(fmt.Errorf("multiple documents in override %s", overrideFile))
					}
					o.mergeWith(overrides[0])
				}

				o.fixFile()
			}
			cfgs = append(cfgs, o)
			_ = os.Chdir(currdir)
		}
	}

	return cfgs
}
