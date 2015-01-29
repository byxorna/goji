package main

import (
	"fmt"
	"github.com/byxorna/marathon_http_proxy_generator/config"
	"log"
	"net/http"
)

var (
	host   string
	port   int
	client *http.Client
)

func Configure(cfg *config.Config) {
	host = cfg.MarathonHost
	port = cfg.MarathonPort
	client = &http.Client{}
}

/*
GET /v2/apps/{appId}/tasks

List all running tasks for application appId.

Example (as JSON)

Request:

GET /v2/apps/my-app/tasks HTTP/1.1
Accept: application/json
Accept-Encoding: gzip, deflate, compress
Content-Type: application/json; charset=utf-8
Host: localhost:8080
User-Agent: HTTPie/0.7.2

HTTP/1.1 200 OK
Content-Type: application/json
Server: Jetty(8.y.z-SNAPSHOT)
Transfer-Encoding: chunked

{
    "tasks": [
        {
            "host": "agouti.local",
            "id": "my-app_1-1396592790353",
            "ports": [
                31336,
                31337
            ],
            "stagedAt": "2014-04-04T06:26:30.355Z",
            "startedAt": "2014-04-04T06:26:30.860Z",
            "version": "2014-04-04T06:26:23.051Z"
        },
        {
            "host": "agouti.local",
            "id": "my-app_0-1396592784349",
            "ports": [
                31382,
                31383
            ],
            "stagedAt": "2014-04-04T06:26:24.351Z",
            "startedAt": "2014-04-04T06:26:24.919Z",
            "version": "2014-04-04T06:26:23.051Z"
        }
    ]
}
*/

//func GetTasks(string appId) ([]string, error) {
func GetTasks(appId string) error {
	url := fmt.Sprintf("http://%s:%d/v2/apps%s/tasks", host, port, appId)
	log.Printf("Getting %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "byxorna/marathon_http_proxy_generator")
	resp, err := client.Do(req)
	res, err := http.Get(url)
}
