package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// SendOutput starts 'kubectl' and pipes the templated result into an 'apply'
func SendOutput(cfgs []*SK8config) bool {
	div := ""

	var cmd *exec.Cmd
	stopIn := make(chan bool)
	stopErr := make(chan bool)

	var output io.WriteCloser
	output = os.Stdout
	if args.Apply {
		cmd = exec.Command("kubectl", "apply", "-f", "-")
		overrideIn, err := cmd.StdinPipe()
		if err != nil {
			panic(err)
		}

		overrideErr, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}

		overrideOut, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		scannerIn := bufio.NewScanner(overrideOut)
		go readStuff(scannerIn, stopIn, false)

		scannerErr := bufio.NewScanner(overrideErr)
		go readStuff(scannerErr, stopErr, true)

		output = overrideIn
	}

	for _, o := range cfgs {
		if o.cfgType == typeKubectl {
			log.Infof("Send raw k8-configs: %s", o.Kind)
			output.Write([]byte("---\n"))
			_, _ = output.Write(o.RawYAML)
		} else {
			if args.AllTemplates {
				for _, path := range o.Templates {
					o.makeYaml(path, &div, &output)
				}
			} else {
				for _, tmplArg := range args.Templates {
					for tmplID, path := range o.Templates {
						if strings.HasPrefix(tmplID, tmplArg) {
							o.makeYaml(path, &div, &output)
						}
					}
					// path := o.Templates[tmplkey]
				}
			}
		}
	}

	if args.Apply {
		output.Close()
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		<-stopIn
		<-stopErr
		err = cmd.Wait()
		if err != nil {
			log.Errorf("kubectl failed: %s", err.Error())
			return false
		}
	}
	return true
}

func (cfg *SK8config) makeYaml(path string, div *string, output *io.WriteCloser) {
	fi, err := os.Stat(path)
	if err != nil {
		cwd, _ := os.Getwd()
		log.Debugf("cwd = %s", cwd)
		log.Warningf("Warning: template %s error: %v", path, err)
		return
	}
	if fi.IsDir() {
		log.Debugf("Ignore template %s -> directory", path)
		return
	}
	if path != "" {
		w := *output
		log.Infof("Create YAML from the %q template for %s/%s", path, cfg.Namespace, cfg.Name)
		tmpl := template.New("tmpl").Funcs(FuncMap())
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		tmpl, err = tmpl.Parse(string(buf))
		//tmpl, err := tmpl.ParseFiles(os.Args[2])
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, *div)
		*div = "---\n"
		//err = tmpl.Execute(os.Stdout, cfg)
		err = tmpl.Execute(w, cfg)
		if err != nil {
			panic(err)
		}
	}
}

func dump(o interface{}) {
	out, err := json.MarshalIndent(&o, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	fmt.Println("##############################################")
}
