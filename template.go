package main

import (
	"os"
	"text/template"
)

// given tasks and a template, do the templating
func Template(tasks []Task, templateFile string) error {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}
	err = tmpl.Execute(os.Stdout, tasks)
	if err != nil {
		return err
	}
	return nil
}
