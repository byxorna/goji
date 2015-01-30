package main

import (
	"bytes"
	"text/template"
)

// given services and a template, do the templating
func Template(services []Service, templateFile string) (string, error) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, services)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
