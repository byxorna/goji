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

}
