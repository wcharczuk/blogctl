package engine

import (
	"bytes"
	"html/template"

	sdkTemplate "github.com/blend/go-sdk/template"
)

// ParseTemplate creates a new template from a string
func ParseTemplate(literal string) (*template.Template, error) {
	tmp := template.New("")
	tmp.Funcs(sdkTemplate.Funcs.FuncMap())
	return tmp.Parse(literal)
}

// RenderString renders a template to a string for a given viewmodel.
func RenderString(tmp *template.Template, vm interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	if err := tmp.Execute(buffer, vm); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
