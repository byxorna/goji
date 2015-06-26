package marathon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	host   string
	port   int
	client *http.Client
}

func NewClient(host string, port int) Client {
	return Client{
		host:   host,
		port:   port,
		client: &http.Client{},
	}
}

func (c *Client) doGetRequest(route string) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:%d%s", c.host, c.port, route)
	log.Printf("Getting %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "byxorna/goji")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetAllTasks() (map[AppId]TaskList, error) {
	resp, err := c.doGetRequest("/v2/tasks?status=running")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Got response %d from %s: %s", resp.StatusCode, resp.Request.URL, body)
	}

	var js map[string]TaskList
	err = json.Unmarshal(body, &js)
	if err != nil {
		return nil, err
	}

	//structure map from appid to tasks
	m := make(map[AppId]TaskList)
	for _, v := range js["tasks"] {
		if _, ok := m[v.AppId]; !ok {
			m[v.AppId] = TaskList{v}
		} else {
			m[v.AppId] = append(m[v.AppId], v)
		}
	}

	return m, nil
}

func (c *Client) GetTasks(appId string, appMustExist bool) (TaskList, error) {
	resp, err := c.doGetRequest(fmt.Sprintf("/v2/apps%s/tasks", appId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if !appMustExist && resp.StatusCode == http.StatusNotFound {
		log.Printf("App %s does not exist in marathon; assuming no tasks\n", appId)
		return TaskList{}, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Got response %d from %s: %s", resp.StatusCode, resp.Request.URL, body)
	}

	var js map[string]TaskList
	err = json.Unmarshal(body, &js)
	if err != nil {
		return nil, err
	}
	t := js["tasks"]
	log.Printf("Found %d tasks for appId %s\n", len(t), appId)
	return t, nil
}

func (c *Client) HasCallback(callback string) (bool, error) {
	url := fmt.Sprintf("http://%s:%d/v2/eventSubscriptions", c.host, c.port)
	log.Printf("Fetching event subscriptions from %s\n", url)
	resp, err := c.client.Get(url)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// 400: {"message":"http event callback system is not running on this Marathon instance. Please re-start this instance with \"--event_subscriber http_callback\"."}
	// 200: {"callbackUrls":[]}
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("Error: http event callback system is not running on %s:%d. Please re-start this instance with \"--event_subscriber http_callback\"", c.host, c.port)
	}

	var js map[string][]string
	err = json.Unmarshal(body, &js)
	if err != nil {
		return false, err
	}
	for _, cb := range js["callbackUrls"] {
		if cb == callback {
			return true, nil
		}
	}
	return false, nil
}

func (c *Client) RegisterCallback(callback string) error {
	url := fmt.Sprintf("http://%s:%d/v2/eventSubscriptions?callbackUrl=%s", c.host, c.port, callback)
	log.Printf("Adding event subscription %s to %s\n", callback, url)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "byxorna/goji")
	resp, err := c.client.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error: unable to register callback: %s", body)
	}
	return nil

}

// TODO this could use some DRYing up with RegisterCallback
func (c *Client) RemoveCallback(callback string) error {
	url := fmt.Sprintf("http://%s:%d/v2/eventSubscriptions?callbackUrl=%s", c.host, c.port, callback)
	log.Printf("Removing event subscription %s from %s\n", callback, url)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "byxorna/goji")
	resp, err := c.client.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error: unable to register callback: %s", body)
	}
	return nil
}
