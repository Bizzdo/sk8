package main

import (
	"bytes"
	"io/ioutil"

	"github.com/dimchansky/utfbom"
)

func getFile(name string) string {
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		log.Errorf("getFile(%s) error: %s", name, err.Error())
		return ""
	}
	return string(buf)
}

func getTextfile(name string) string {
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		log.Errorf("getTextfile(%s) error: %s", name, err.Error())
		return ""
	}

	sr, _ := utfbom.Skip(bytes.NewReader(buf))

	buf, err = ioutil.ReadAll(sr)
	if err != nil {
		log.Errorf("getTextfile(%s) error: %s", name, err.Error())
		return ""
	}

	output := make([]byte, 0, len(buf))
	for _, ch := range buf {
		switch ch {
		case '\r':

		case '\t':
			output = append(output, ' ')

		default:
			output = append(output, ch)
		}
	}

	return string(output)
}
