package main

import (
	"flag"
	"fmt"
	"github.com/byxorna/goji/marathon"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configPath string
	config     Config
	// current state of services we know about
	services ServiceList
	// listen to this channel for update triggers
	eventChan chan string
	client    marathon.Client
	sigChan   chan os.Signal
)

func init() {
	flag.StringVar(&configPath, "conf", "", "config json file")
	flag.Parse()
}

func main() {
	if configPath == "" {
		log.Fatal("You need to pass a -conf")
	}
	var err error
	config, err = LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded config: %s\n", config)
	client = marathon.NewClient(config.MarathonHost, config.MarathonPort)
	services = config.Services

	// lets let people know when we get a signal, so we can clean up
	sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	log.Printf("Loading tasks from marathon\n")
	err = services.LoadAllAppTasks(client)
	if err != nil {
		log.Fatal(err.Error())
	}

	// just print out what apps and tasks we found
	for _, service := range services {
		log.Printf("Found app %s with %d tasks\n", service.AppId, len(*service.Tasks))
		for _, task := range *service.Tasks {
			for _, port := range task.Ports {
				log.Printf("  %s:%d\n", task.Host, port)
			}
		}
	}

	// spit out the first config before we start listening for events
	log.Printf("Templating %s with %d services\n", config.TemplateFile, len(services))
	output, err := Template(services, config.TemplateFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Writing config to %s\n", config.TargetFile)
	err = WriteConfig(output, config.TargetFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Wrote %s!\n", config.TargetFile)

	// every event hits eventChan (buffer 10 events)
	eventChan = make(chan string, 10)

	go func() {
		coalesceEvents(eventChan, time.Duration(config.TemplateDelay)*time.Second, func() {
			err := LoadTasksAndEmitConfig()
			if err != nil {
				log.Printf(err.Error())
			}
		})
	}()

	var listenAddr = fmt.Sprintf(":%d", config.HttpPort)
	log.Printf("Listening for marathon events on %s/event\n", listenAddr)
	log.Fatal(ListenForEvents(listenAddr))

}
