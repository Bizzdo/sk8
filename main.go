package main

import (
	"bufio"
	"fmt"
	"os"

	"strings"
)

func main() {

	parseArgs()

	initLog(args.Verbose)

	validKC := validateKubeconfig(args.Kubeconfig || args.Apply)

	if args.Kubeconfig || (!validKC && args.Apply) {
		return
	}

	cfgs := LoadInput()

	if len(cfgs) == 0 || args.Help {
		fmt.Println("Usage: sk8 [option(s)] {file(s)}")
		fmt.Println("")
		fmt.Println("Options:  -apply              Call 'kubectl' to update Kubernetes")
		fmt.Println("          -kubeconfig | -kc   Show the current server according to the KUBECONFIG-file")
		fmt.Println("          -verbose    | -v    Show verbose output")
		fmt.Println("          -{template-tag}     Applies the template indicated by the {template-tag} key (i.e is prefixed with)")
		fmt.Println("          -all                Applies all templates that are specified")
		return
	}

	if (args.AllTemplates || len(args.Templates) > 0) && !args.Debug {
		if !SendOutput(cfgs) {
			os.Exit(1)
		}
	} else {
		for _, o := range cfgs {
			dump(o)
		}
	}
}

func readStuff(scanner *bufio.Scanner, stop chan bool, isErr bool) {
	format := "kubectl> %s"
	for scanner.Scan() {
		txt := scanner.Text()
		if isErr {
			if strings.HasPrefix(txt, "W") {
				log.Warningf(format, txt)
			} else {
				log.Errorf(format, txt)
			}
		} else {
			log.Noticef(format, txt)
		}
	}
	stop <- true
}
