package main

import (
	"flag"
	"github.com/byxorna/marathon_http_proxy_generator/config"
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
	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded config: %s\n", config)
	Configure(config)
	for id, service := range config.Services {
		ts, err := GetTasks(id)
		if err != nil {
			log.Printf("Error fetching tasks for %s: %s\n", id, err.Error())
		}
		log.Printf("Application %s at %s.%s:\n", id, service, config.Domain)
		for _, t := range ts {
			for _, p := range t.Ports {
				log.Printf("  %s:%d\n", t.Host, p)
			}
		}
	}
}
