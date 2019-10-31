package main

import (
	"yaml_mapstr"
)

func yamlUnmarshal(buf []byte, dest interface{}) error {
	return yaml_mapstr.Unmarshal(buf, dest)
}
