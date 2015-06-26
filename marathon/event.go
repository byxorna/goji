package marathon

import (
	"log"
	"time"
)

//TODO: we should have separate event structs for each message type
type Event struct {
	EventType string `json:"eventType"`
	Timestamp string `json:"timestamp"`
	AppId     AppId  `json:"appId,omitempty"`
	TaskId    string `json:"taskId,omitempty"`
}

type StatusUpdateEvent struct {
	EventType  string     `json:"eventType"`
	Timestamp  string     `json:"timestamp"`
	AppId      AppId      `json:"appId"`
	TaskId     string     `json:"taskId"`
	SlaveId    string     `json:"slaveId"`
	TaskStatus TaskStatus `json:"taskStatus"`
	Host       string     `json:"host"`
	Ports      []int      `json:"ports"`
}

type HealthStatusChangedEvent struct {
	EventType string `json:"eventType"`
	Timestamp string `json:"timestamp"`
	AppId     AppId  `json:"appId"`
	TaskId    string `json:"taskId"`
	Alive     bool   `json:"alive"`
}

func (e *Event) Time() time.Time {
	t, err := time.Parse(time.RFC3339, e.Timestamp)
	if err != nil {
		log.Printf(err.Error())
		return time.Time{}
	}
	return t
}
