package main

import (
	"os"

	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("bizzdo")
var logFormat = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{level} %{message}%{color:reset}`)
var logger *logging.LogBackend

func initLog(verbose bool) {
	logger = logging.NewLogBackend(os.Stderr, "", 0)
	loggerFmt := logging.NewBackendFormatter(logger, logFormat)
	logLevel := logging.AddModuleLevel(loggerFmt)
	if !verbose {
		logLevel.SetLevel(logging.INFO, "")
	}
	logging.SetBackend(logLevel)
}
