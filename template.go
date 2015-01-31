package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"text/template"
)

// given services and a template, do the templating
func Template(services ServiceList, templateFile string) (string, error) {
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

func WriteConfig(cfg string, outputFile string) error {
	return ioutil.WriteFile(outputFile, []byte(cfg), 0644)
}

func LoadTasksAndEmitConfig() error {
	log.Printf("Loading tasks from marathon\n")
	err := services.LoadAllAppTasks(client)
	if err != nil {
		return fmt.Errorf("Unable to load tasks from marathon: %s", err.Error())
	}

	log.Printf("Templating %s with %d services\n", config.TemplateFile, len(services))
	output, err := Template(services, config.TemplateFile)
	if err != nil {
		return fmt.Errorf("Unable to compile template: %s", err.Error())
	}

	log.Printf("Writing config to %s\n", config.TargetFile)
	err = WriteConfig(output, config.TargetFile)
	if err != nil {
		return fmt.Errorf("Unable to write config to %s: %s", config.TargetFile, err.Error())
	}
	log.Printf("Wrote %s\n", config.TargetFile)
	return nil
}
