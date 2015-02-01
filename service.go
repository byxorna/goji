package main

import "github.com/byxorna/goji/marathon"

type Service struct {
	Vhost string           `json:"vhost"`
	AppId string           `json:"app-id"`
	Tasks *[]marathon.Task `json:"-"`
	//TODO add config for healthchecking, type of connection (HTTP/TCP), etc
	//TODO possibly add configurable domains
}

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
		(*services)[i].Tasks = &ts
	}
	return nil
}
