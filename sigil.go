package sigil

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/mgood/go-posix"
)

var (
	TemplateDir string
)

var fnMap = template.FuncMap{}

func Register(fm template.FuncMap) {
	for k, v := range fm {
		fnMap[k] = v
	}
}

func Execute(input string, vars map[string]string) (string, error) {
	var tmplVars string
	for k, v := range vars {
		err := os.Setenv(k, v)
		if err != nil {
			return "", err
		}
		escaped := strings.Replace(v, "\"", "\\\"", -1)
		tmplVars = tmplVars + fmt.Sprintf("{{ $%s := \"%s\" }}", k, escaped)
	}
	tmpl, err := template.New("template").Funcs(fnMap).Parse(tmplVars + input)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		return "", err
	}
	postprocessed, err := posix.ExpandEnv(buf.String())
	if err != nil {
		return "", err
	}
	return postprocessed, nil
}