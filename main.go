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
  // override the target in config
	target string
	config     Config
	// current state of services we know about
	services ServiceList
	// listen to this channel for update triggers
	eventChan chan string
	client    marathon.Client
	sigChan   chan os.Signal
	server    bool
)

func init() {
	flag.StringVar(&configPath, "conf", "", "config json file")
	flag.BoolVar(&server, "server", false, "start a HTTP server listening for Marathon events")
	flag.StringVar(&target, "target", "", "target file to write to")
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
  // allow -target to override the target specified in config for CLI testing
  if target != "" {
    config.TargetFile = target
  }
	log.Printf("Loaded config: %s\n", config)
	client = marathon.NewClient(config.MarathonHost, config.MarathonPort)
	services, err = NewServiceList(config.Services)
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, s := range services {
		log.Printf("Loaded %s service %s at port %d\n", s.Protocol, s.AppId, s.Port)
	}

	// lets let people know when we get a signal, so we can clean up
	sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

  err = LoadTasksAndEmitConfig()
  if err != nil {
    log.Fatal(err.Error())
  }

	if server {
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

}
