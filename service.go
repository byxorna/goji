package main

import (
	"fmt"
	"github.com/byxorna/goji/marathon"
	"strings"
)

type Service struct {
	Vhost           string
	AppId           string
	tasks           *[]marathon.Task
	healthCheckPath string
	Protocol        ProtocolType
	Port            int
}

type ProtocolType string

const (
	HTTP ProtocolType = "HTTP"
	TCP               = "TCP"
)

type ServiceList []Service

// populates the list of tasks in each service
// and clobber each service in the ServiceList's Tasks with a new set
func (services *ServiceList) LoadAllAppTasks(c marathon.Client) error {
	for i, service := range *services {
		ts, err := c.GetTasks(service.AppId)
		if err != nil {
			return err
		}
		// I still really dont grok how go's pointers work for mutability
		// but this works...
		(*services)[i].tasks = &ts
	}
	return nil
}

func NewServiceList(configservices []ConfigService) (ServiceList, error) {
	svcs := make(ServiceList, len(configservices))
	for i, s := range configservices {
		ns, err := NewService(s)
		if err != nil {
			return svcs, err
		}
		svcs[i] = ns
	}
	return svcs, nil
}
func NewService(cfg ConfigService) (Service, error) {
	//TODO do validation of healthcheck here as well
	if cfg.Protocol == "" {
		cfg.Protocol = HTTP
	} else if cfg.Protocol != HTTP && cfg.Protocol != TCP {
		return Service{}, fmt.Errorf("Unknown protocol %s", cfg.Protocol)
	}
	if cfg.Port == 0 {
		// just assume port 80 if not specified
		cfg.Port = 80
	}
	return Service{
		Vhost:           cfg.Vhost,
		AppId:           cfg.AppId,
		Protocol:        cfg.Protocol,
		healthCheckPath: cfg.HealthCheckPath,
		Port:            cfg.Port,
	}, nil
}

// replaces / with ::, useful for creating haproxy identifiers
func (s *Service) EscapeAppIdColon() string {
	return strings.Replace(s.AppId, "/", "::", -1)
}

// replaces / with _, useful for creating nginx identifiers
func (s *Service) EscapeAppIdUnderscore() string {
	return strings.Replace(s.AppId, "/", "_", -1)
}

// returns a copy of the list of tasks, or [] if tasks is a nil pointer
func (s *Service) Tasks() []marathon.Task {
	if s.tasks == nil {
		return []marathon.Task{}
	} else {
		return *s.tasks
	}
}

// if the service is HTTP protocol, and defined a health check, return it, else nil
func (s *Service) HealthCheckPath() string {
	if s.Protocol == HTTP && s.healthCheckPath != "" && strings.HasPrefix(s.healthCheckPath, "/") {
		return s.healthCheckPath
	} else {
		return ""
	}
}

// is this service using the HTTP protocol?
func (s *Service) HTTPProtocol() bool {
	return s.Protocol == HTTP
}

// is this service using the TCP protocol?
func (s *Service) TCPProtocol() bool {
	return s.Protocol == TCP
}
