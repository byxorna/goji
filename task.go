package main

import (
	"encoding/json"
	"fmt"
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

type TaskClient struct {
	host   string
	port   int
	client *http.Client
}

type Task struct {
	Id        string `json:"id"`
	Ports     []int  `json:"ports"`
	Host      string `json:"host"`
	stagedAt  string `json:"stagedAt"`
	startedAt string `json:"startedAt"`
	version   string `json:"version"`
}

func NewTaskClient(cfg *Config) *TaskClient {
	return &TaskClient{
		host:   cfg.MarathonHost,
		port:   cfg.MarathonPort,
		client: &http.Client{},
	}
}

//TODO this may be more efficient to hit /v2/tasks?status=running
// and filter for the apps we care about
func (c *TaskClient) GetTasks(appId string) ([]Task, error) {
	url := fmt.Sprintf("http://%s:%d/v2/apps%s/tasks", c.host, c.port, appId)
	log.Printf("Getting %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "byxorna/marathon_http_proxy_generator")
	resp, err := c.client.Do(req)

	//TODO this feels awfully wordy. I miss ruby's brevity....

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Got response %d from %s: %s", resp.StatusCode, url, body)
	}

	var js map[string][]Task
	err = json.Unmarshal(body, &js)
	if err != nil {
		return nil, err
	}
	t := js["tasks"]
	return t, nil
}

// populates the list of tasks in each service
func (tc TaskClient) LoadAllAppTasks(services *[]Service) error {
	for _, service := range *services {
		ts, err := tc.GetTasks(service.AppId)
		if err != nil {
			return err
		}
		service.Tasks = ts
	}
	return nil
}
