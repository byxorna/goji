package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Service struct {
	Vhost string  `json:"vhost"`
	AppId string  `json:"app-id"`
	Tasks *[]Task `json:"-"`
	//TODO add config for healthchecking, type of connection (HTTP/TCP), etc
	//TODO possibly add configurable domains
}

type AppIdTasksMap map[string]*[]Task

type Config struct {
	// localhost
	MarathonHost string `json:"marathon-host,omitempty"`
	// 8080
	MarathonPort int       `json:"marathon-port"`
	Services     []Service `json:"services,omitempty"`
	TemplateFile string    `json:"template,omitempty"`
	// port upon which to listen for events from marathon
	HttpPort int `json:"http-port"`
}

func LoadConfig(configPath string) (*Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	c := Config{}
	err = json.NewDecoder(f).Decode(&c)
	if c.MarathonPort == 0 {
		c.MarathonPort = 8080
	}
	if c.HttpPort == 0 {
		c.HttpPort = 8000
	}
	if c.TemplateFile == "" {
		return nil, fmt.Errorf("template is required")
	}
	if c.MarathonHost == "" {
		return nil, fmt.Errorf("marathon-host is required")
	}
	if len(c.Services) == 0 {
		return nil, fmt.Errorf("At least one service is required in `services`")
	}
	return &c, nil
}
