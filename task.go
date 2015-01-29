package main

import (
	"encoding/json"
	"fmt"
	"github.com/byxorna/marathon_http_proxy_generator/config"
	"io/ioutil"
	"log"
	"net/http"
)

/*
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

type Task struct {
	Id        string `json:"id"`
	Ports     []int  `json:"ports"`
	Host      string `json:"host"`
	stagedAt  string `json:"stagedAt"`
	startedAt string `json:"startedAt"`
	version   string `json:"version"`
}

var (
	host   string
	port   int
	client *http.Client
)

// this feels wrong. whats the idiomatic way to configure the module before use?
func Configure(cfg *config.Config) {
	host = cfg.MarathonHost
	port = cfg.MarathonPort
	client = &http.Client{}
}

func GetTasks(appId string) ([]Task, error) {
	url := fmt.Sprintf("http://%s:%d/v2/apps%s/tasks", host, port, appId)
	log.Printf("Getting %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "byxorna/marathon_http_proxy_generator")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	//TODO this feels awfully words. I miss ruby's brevity....

	if err != nil {
		log.Printf("Error fetching tasks for %s: %s\n", appId, err.Error())
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Got response %d: %s", resp.StatusCode, body)
	}

	var js map[string][]Task
	err = json.Unmarshal(body, &js)
	if err != nil {
		return nil, err
	}
	t := js["tasks"]
	return t, nil
}
