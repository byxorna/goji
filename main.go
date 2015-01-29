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
	serviceTasksMap, err := tc.loadTasks(config.Services)
	if err != nil {
		log.Fatal(err.Error())
	}

	// just print out what apps and tasks we found
	for service, tasks := range serviceTasksMap {
		appId := config.Services[service]
		log.Printf("App %s at %s.%s:\n", appId, service, config.Domain)
		for _, task := range *tasks {
			for _, port := range task.Ports {
				log.Printf("  %s:%d\n", task.Host, port)
			}
		}

		Template(*tasks, config.TemplateFile)
	}

}

func (tc TaskClient) loadTasks(services ServiceAppIdMap) (ServiceTasksMap, error) {
	res := ServiceTasksMap{}
	for service, appId := range services {
		ts, err := tc.GetTasks(appId)
		if err != nil {
			return nil, err
		}
		res[service] = &ts
	}
	return res, nil

}
