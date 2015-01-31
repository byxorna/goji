package main

type Service struct {
	Vhost string  `json:"vhost"`
	AppId string  `json:"app-id"`
	Tasks *[]Task `json:"-"`
	//TODO add config for healthchecking, type of connection (HTTP/TCP), etc
	//TODO possibly add configurable domains
}

type ServiceList []Service
