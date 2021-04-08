package main

import (
	"bytes"
	"fmt"
	// "encoding/json"
	// "yaml_mapstr"
	"github.com/mickep76/encoding"
	_ "github.com/mickep76/encoding/json"
	_ "github.com/mickep76/encoding/toml"
	_ "github.com/mickep76/encoding/yaml"
)

var codecs struct {
	yaml encoding.Codec
	json encoding.Codec
	toml encoding.Codec
}

func yamlCodec() encoding.Codec {
	if codecs.yaml == nil {
		c, err := encoding.NewCodec("yaml", encoding.WithMapString())
		if err != nil {
			panic(err)
		}
		codecs.yaml = c
	}
	return codecs.yaml
}

func jsonCodec() encoding.Codec {
	if codecs.json == nil {
		c, err := encoding.NewCodec("json", encoding.WithIndent("  "))
		if err != nil {
			panic(err)
		}
		codecs.json = c
	}
	return codecs.json
}

func tomlCodec() encoding.Codec {
	if codecs.toml == nil {
		c, err := encoding.NewCodec("toml")
		if err != nil {
			panic(err)
		}
		codecs.toml = c
	}
	return codecs.toml
}

func yaml_Unmarshal(buf []byte, dest interface{}) error {
	//return yamlCodec().Decode(buf, dest)

	var res interface{}
	if err := yamlCodec().Decode(buf, &res); err != nil {
		return err
	}

	*dest.(*interface{}) = cleanupMapValue(res)
	return nil

	//return yaml_mapstr.Unmarshal(buf, dest)
}

func json_Marshal(o interface{}) ([]byte, error) {
	return jsonCodec().Encode(o)
}

func json_Unmarshal(buf []byte, o interface{}) error {
	return jsonCodec().Decode(buf, o)
}

// ToYaml takes an interface, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func ToYaml(v interface{}) string {
	data, err := yamlCodec().Encode(v)
	//data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}

// FromYaml converts a YAML document into a map[string]interface{}.
//
// This is not a general-purpose YAML parser, and will not parse all valid
// YAML documents. Additionally, because its intended use is within templates
// it tolerates errors. It will insert the returned error message string into
// m["Error"] in the returned map.
func FromYaml(str string) map[string]interface{} {
	m := map[string]interface{}{}

	if err := yamlCodec().Decode([]byte(str), &m); err != nil {
		//if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

// ToToml takes an interface, marshals it to toml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func ToToml(v interface{}) string {
	b := bytes.NewBuffer(nil)
	//e := toml.NewEncoder(b)
	e, err := tomlCodec().NewEncoder(b)
	if err != nil {
		return err.Error()
	}
	err = e.Encode(v)
	if err != nil {
		return err.Error()
	}
	return b.String()
}

// ToJson takes an interface, marshals it to json, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
// TODO: change the function signature in Helm 3
func ToJson(v interface{}) string { // nolint
	data, err := jsonCodec().Encode(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}

// FromJson converts a JSON document into a map[string]interface{}.
//
// This is not a general-purpose JSON parser, and will not parse all valid
// JSON documents. Additionally, because its intended use is within templates
// it tolerates errors. It will insert the returned error message string into
// m["Error"] in the returned map.
// TODO: change the function signature in Helm 3
func FromJson(str string) map[string]interface{} { // nolint
	m := map[string]interface{}{}

	if err := jsonCodec().Decode([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

// copied from encoding/yaml/map_string.go
func cleanupInterfaceArray(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = cleanupMapValue(v)
	}
	return res
}

func cleanupInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[fmt.Sprintf("%v", k)] = cleanupMapValue(v)
	}
	return res
}

func cleanupMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanupInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanupInterfaceMap(v)
	default:
		return v
	}
}
