package main

import (
	"fmt"
	"github.com/byxorna/goji/marathon"
	"sort"
	"strings"
)

type Service struct {
	Name            string
	AppId           marathon.AppId
	tasks           *marathon.TaskList
	healthCheckPath string
	Protocol        ProtocolType
	Port            int
	Options         map[string]string
}

type ProtocolType string

const (
	HTTP ProtocolType = "HTTP"
	TCP               = "TCP"
	UDP               = "UDP"
)

type ServiceList []Service

// populates the list of tasks in each service
// and clobber each service in the ServiceList's Tasks with a new set
func (services *ServiceList) LoadAllAppTasks(c marathon.Client) error {
	marathonApps, err := c.GetAllTasks(marathon.TaskAny)
	if err != nil {
		return err
	}

	for i, service := range *services {
		ts, ok := marathonApps[service.AppId]
		if !ok {
			if appPresenceRequired {
				return fmt.Errorf("%s does not exist in marathon, -app-required is true", service.AppId)
			} else {
				// we should assume an empty task set
				marathonApps[service.AppId] = marathon.TaskList{}
			}
		}
		// I still really dont grok how go's pointers work for mutability
		// but this works...
		// Make sure we sort tasks, so configs have a predictable ordering and dont change every run
		sort.Sort(ts)
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
	}
	// just accept the protocol given if we dont recognize it? Could be someone using a strange
	// proto for SRV records like ldap or so on.
	if cfg.Port == 0 {
		// just assume port 80 if not specified
		cfg.Port = 80
	}
	return Service{
		Name:            cfg.Name,
		AppId:           cfg.AppId,
		Protocol:        cfg.Protocol,
		healthCheckPath: cfg.HealthCheckPath,
		Port:            cfg.Port,
		Options:         cfg.Options,
	}, nil
}

// replaces / with ::, useful for creating haproxy identifiers
func (s *Service) EscapeAppIdColon() string {
	return strings.Replace(string(s.AppId), "/", "::", -1)
}

// replaces / with _, useful for creating nginx identifiers
func (s *Service) EscapeAppIdUnderscore() string {
	return strings.Replace(string(s.AppId), "/", "_", -1)
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

func (s *Service) HTTPProtocol() bool {
	return s.Protocol == HTTP
}

func (s *Service) TCPProtocol() bool {
	return s.Protocol == TCP
}

func (s *Service) UDPProtocol() bool {
	return s.Protocol == UDP
}
