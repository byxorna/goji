package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/byxorna/goji/marathon"
)

type Config struct {
	// localhost
	MarathonHost string `json:"marathon-host,omitempty"`
	// 8080
	MarathonPort int `json:"marathon-port"`
	//	Services     []ConfigService `json:"services,omitempty"`
	TemplateFile string `json:"template,omitempty"`
	TargetFile   string `json:"target,omitempty"`
	// port upon which to listen for events from marathon
	HttpPort int `json:"http-port"`
	// hostname to request marathon hit with webhook. defaults to os.Hostname()
	CallbackHostname string `json:"callback-hostname"`
	TemplateDelay    int    `json:"delay"`
	Command          string `json:"command"`
}

// the user defined representation of a service
// passed into NewService to create an actual Service struct which does validation
type ConfigService struct {
	Name            string            `json:"name"`
	AppId           marathon.AppId    `json:"app-id"`
	HealthCheckPath string            `json:"health-check"`
	Protocol        ProtocolType      `json:"protocol"`
	Port            int               `json:"port"`
	Options         map[string]string `json:"options"`
}

func MergeConfigWithEnv(some Config) (Config, error) {
	if some.MarathonHost == "" {
		some.MarathonHost = os.Getenv("MARATHON_HOST")
	}
	if some.MarathonPort == 0 {
		v, err := strconv.Atoi(os.Getenv("MARATHON_PORT"))
		if err != nil {
			return some, err
		}
		some.MarathonPort = v
	}
	if some.TemplateFile == "" {
		some.TemplateFile = os.Getenv("TEMPLATE_FILE")
	}
	if some.TargetFile == "" {
		some.TargetFile = os.Getenv("TARGET_FILE")
	}
	if some.HttpPort == 0 {
		v, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
		if err != nil {
			return some, err
		}
		some.HttpPort = v
	}
	if some.CallbackHostname == "" {
		some.CallbackHostname = os.Getenv("CALLBACK_HOSTNAME")
	}
	if some.TemplateDelay == 0 {
		v, err := strconv.Atoi(os.Getenv("TEMPLATE_DELAY"))
		if err != nil {
			return some, err
		}
		some.TemplateDelay = v
	}
	if some.Command == "" {
		some.Command = os.Getenv("COMMAND")
	}

	return some, nil
}

func LoadConfigFromFile(configPath string) (Config, error) {
	c := Config{}
	f, err := os.Open(configPath)
	if err != nil {
		return c, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&c)
	return c, err
}

func (c *Config) ValidateAndSetDefaults() error {
	if c.CallbackHostname == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return err
		}
		c.CallbackHostname = hostname
	}
	if c.MarathonPort == 0 {
		c.MarathonPort = 8080
	}
	if c.HttpPort == 0 {
		c.HttpPort = 8000
	}
	if c.TemplateFile == "" {
		return fmt.Errorf("TemplateFile is required")
	}
	if c.TargetFile == "" {
		return fmt.Errorf("TargetFile is required")
	}
	if c.MarathonHost == "" {
		return fmt.Errorf("MarathonHost is required")
	}
	//if len(c.Services) == 0 {
	//	return fmt.Errorf("At least one service is required in `services`")
	//}

	return nil
}
