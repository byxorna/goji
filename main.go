package main

import (
	"flag"
	"log"
)

var (
	configPath string
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
	log.Printf("Querying marathon for tasks\n")
	appIdTasksMap, err := tc.LoadAllAppTasks(config.Services)
	if err != nil {
		log.Fatal(err.Error())
	}

	// just print out what apps and tasks we found
	for _, service := range config.Services {
		tasks, ok := appIdTasksMap[service.AppId]
		if !ok {
			log.Fatal("Unable to find tasks for app id %s", service.AppId)
		}
		log.Printf("App %s at %s.%s:\n", service.AppId, service.Vhost)
		for _, task := range *tasks {
			for _, port := range task.Ports {
				log.Printf("  %s:%d\n", task.Host, port)
			}
		}

		Template(*tasks, config.TemplateFile)
	}

	log.Printf("Listening for events on :%d\n", config.HttpPort)
	log.Fatal(ListenForEvents(config))

}
