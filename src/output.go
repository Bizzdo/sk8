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

// SendOuput -
func SendOutput(cfgs []*sk8config) {
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

	if args.Apply {
		output.Close()
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		<-stopIn
		<-stopErr
		cmd.Wait()
	}
}

func (cfg *sk8config) makeYaml(path string, div *string, output *io.WriteCloser) {
	if path != "" {
		w := *output
		log.Infof("Create YAML from the %q template for %s/%s", path, cfg.Namespace, cfg.Name)
		tmpl := template.New("tmpl")
		buf, err := ioutil.ReadFile(path)
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
