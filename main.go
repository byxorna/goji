package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	configPath string
	config     *Config
	// current state of services we know about
	services []Service
)

func init() {
	flag.StringVar(&configPath, "conf", "", "config json file")
	flag.Parse()
}

func main() {
	if configPath == "" {
		log.Fatal("You need to pass a -conf")
	}
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded config: %s\n", config)
	tc := NewTaskClient(config)
	services = config.Services
	log.Printf("Loading tasks from marathon\n")
	err = tc.LoadAllAppTasks(&services)
	if err != nil {
		log.Fatal(err.Error())
	}

	// just print out what apps and tasks we found
	for _, service := range services {
		log.Printf("App %s at %s.%s:\n", service.AppId, service.Vhost)
		for _, task := range service.Tasks {
			for _, port := range task.Ports {
				log.Printf("  %s:%d\n", task.Host, port)
			}
		}

	}
	log.Printf("Templating %s with %d services\n", config.TemplateFile, len(services))
	output, err := Template(services, config.TemplateFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Got output:\n%s\n", output)

	var listenAddr = fmt.Sprintf(":%d", config.HttpPort)
	log.Printf("Listening for events on %s\n", listenAddr)
	log.Fatal(ListenForEvents(listenAddr))

}
