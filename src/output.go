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

	"OSS/sk8/chartutil"

	"github.com/Masterminds/sprig"
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
		err = cmd.Wait()
		if err != nil {
			log.Errorf("kubectl failed: %s", err.Error())
			return false
		}
	}
	return true
}

func (cfg *SK8config) makeYaml(path string, div *string, output *io.WriteCloser) {
	if path != "" {
		w := *output
		log.Infof("Create YAML from the %q template for %s/%s", path, cfg.Namespace, cfg.Name)
		tmpl := template.New("tmpl").Funcs(FuncMap())
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

// FuncMap returns a mapping of all of the functions that Engine has.
//
// Because some functions are late-bound (e.g. contain context-sensitive
// data), the functions may not all perform identically outside of an
// Engine as they will inside of an Engine.
//
// Known late-bound functions:
//
//	- "include": This is late-bound in Engine.Render(). The version
//	   included in the FuncMap is a placeholder.
//      - "required": This is late-bound in Engine.Render(). The version
//	   included in the FuncMap is a placeholder.
//      - "tpl": This is late-bound in Engine.Render(). The version
//	   included in the FuncMap is a placeholder.
func FuncMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	// Add some extra functionality
	extra := template.FuncMap{
		"indent2":   indent2,
		"getFile":   getFile,
		"multiline": multilineYaml,

		"toToml":   chartutil.ToToml,
		"toYaml":   chartutil.ToYaml,
		"fromYaml": chartutil.FromYaml,
		"toJson":   chartutil.ToJson,
		"fromJson": chartutil.FromJson,

		// This is a placeholder for the "include" function, which is
		// late-bound to a template. By declaring it here, we preserve the
		// integrity of the linter.
		// "include":  func(string, interface{}) string { return "not implemented" },
		// "required": func(string, interface{}) interface{} { return "not implemented" },
		// "tpl":      func(string, interface{}) interface{} { return "not implemented" },
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

func indent2(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return strings.Replace(v, "\n", "\n"+pad, -1)
}

func getFile(name string) string {
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		log.Errorf("getFile(%s) error: %s", name, err.Error())
		return ""
	}
	return string(buf)
}

func multilineYaml(v string) string {
	return "|\n" + v
}
