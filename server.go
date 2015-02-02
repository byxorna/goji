package main

import (
	"encoding/json"
	"fmt"
	"github.com/byxorna/goji/marathon"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func ListenForEvents(listenAddr string) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	cb := fmt.Sprintf("http://%s:%d/event", hostname, config.HttpPort)
	cbRegistered, err := client.HasCallback(cb)
	if err != nil {
		return err
	}
	if !cbRegistered {
		client.RegisterCallback(cb)
	} else {
		log.Printf("Event callback already registered; skipping registration\n")
	}
	// cleanup registered callback if we catch a signal
	go func() {
		s := <-sigChan
		log.Printf("Cleaning up registered callback after signal %s\n", s)
		err := client.RemoveCallback(cb)
		if err != nil {
			log.Fatal(err.Error())
		}
		os.Exit(0)
	}()
	http.HandleFunc("/event", handleEvent)
	return http.ListenAndServe(listenAddr, nil)
}

func handleEvent(res http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s %s %s\n", req.Method, req.RemoteAddr, req.RequestURI, req.Proto)
	body, err := ioutil.ReadAll(req.Body)
	res.Header().Set("Content-Type", "text/plain")
	if err != nil {
		log.Printf(err.Error())
		fmt.Fprintf(res, "Error reading event body")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Fprintf(res, "Thanks for the event!")
	// lets deal with parsing and identifying the event out of the request handler
	go determineEventRelevancy(body)
}

func determineEventRelevancy(body []byte) {
	e := marathon.Event{}
	err := json.Unmarshal(body, &e)
	if err != nil {
		log.Printf("Unable to decode event body: %s\n", err.Error())
		return
	}
	log.Printf("Received event %s generated at %s\n", e.EventType, e.Time())
	var handleEvent = func() {
		eventChan <- e.EventType
	}
	switch e.EventType {
	case "api_post_event":
		handleEvent()
	case "status_update_event":
		handleEvent()
	default:
		log.Printf("Ignoring event type %s\n", e.EventType)
	}

}
