package main

import (
	"os"
	"strings"
)

type myArgs struct {
	Apply        bool
	AllTemplates bool
	Debug        bool
	Verbose      bool
	Kubeconfig   bool
	Help         bool
	Templates    []string
	Files        []string
}

var args myArgs

func parseArgs() {
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			if strings.EqualFold(arg, "-all") {
				args.AllTemplates = true
			} else if strings.EqualFold(arg, "-apply") {
				args.Apply = true
			} else if strings.EqualFold(arg, "-kubeconfig") || strings.EqualFold(arg, "-kc") {
				args.Kubeconfig = true
			} else if strings.EqualFold(arg, "-debug") || strings.EqualFold(arg, "-d") {
				args.Debug = true
			} else if strings.EqualFold(arg, "-help") || strings.EqualFold(arg, "-h") {
				args.Help = true
			} else if strings.EqualFold(arg, "-verbose") || strings.EqualFold(arg, "-v") {
				args.Verbose = true
			} else {
				args.Templates = append(args.Templates, arg[1:])
			}
		} else {
			if _, err := os.Lstat(arg); err == nil {
				args.Files = append(args.Files, arg)
			}
		}
	}
}
