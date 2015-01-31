package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	// localhost
	MarathonHost string `json:"marathon-host,omitempty"`
	// 8080
	MarathonPort int         `json:"marathon-port"`
	Services     ServiceList `json:"services,omitempty"`
	TemplateFile string      `json:"template,omitempty"`
	TargetFile   string      `json:"target,omitempty"`
	// port upon which to listen for events from marathon
	HttpPort      int `json:"http-port"`
	TemplateDelay int `json:"delay"`
}

func LoadConfig(configPath string) (Config, error) {
	c := Config{}
	f, err := os.Open(configPath)
	if err != nil {
		return c, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&c)
	if c.MarathonPort == 0 {
		c.MarathonPort = 8080
	}
	if c.HttpPort == 0 {
		c.HttpPort = 8000
	}
	if c.TemplateFile == "" {
		return c, fmt.Errorf("template is required")
	}
	if c.TargetFile == "" {
		return c, fmt.Errorf("target is required")
	}
	if c.MarathonHost == "" {
		return c, fmt.Errorf("marathon-host is required")
	}
	if len(c.Services) == 0 {
		return c, fmt.Errorf("At least one service is required in `services`")
	}
	return c, nil
}
