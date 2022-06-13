package main

import (
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

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
		"indent2":     indent2,
		"getFile":     getFile,
		"getTextfile": getTextfile,
		"multiline":   multilineYaml,
		"contains":    contains,

		"toToml":   ToToml,
		"toYaml":   ToYaml,
		"fromYaml": FromYaml,
		"toJson":   ToJson,
		"fromJson": FromJson,

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

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func indent2(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return strings.Replace(v, "\n", "\n"+pad, -1)
}

func multilineYaml(v string) string {
	return "|\n" + v
}
