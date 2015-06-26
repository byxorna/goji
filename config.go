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
	MarathonPort int             `json:"marathon-port"`
	Services     []ConfigService `json:"services,omitempty"`
	TemplateFile string          `json:"template,omitempty"`
	TargetFile   string          `json:"target,omitempty"`
	// port upon which to listen for events from marathon
	HttpPort      int    `json:"http-port"`
	TemplateDelay int    `json:"delay"`
	Command       string `json:"command"`
}

// the user defined representation of a service
// passed into NewService to create an actual Service struct which does validation
type ConfigService struct {
	Name            string            `json:"name"`
	AppId           string            `json:"app-id"`
	HealthCheckPath string            `json:"health-check"`
	Protocol        ProtocolType      `json:"protocol"`
	Port            int               `json:"port"`
	Options         map[string]string `json:"options"`
}

func LoadConfig(configPath string) (Config, error) {
	c := Config{}
	f, err := os.Open(configPath)
	if err != nil {
		return c, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return c, err
	}
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
