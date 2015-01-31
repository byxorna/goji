package main

import (
	"encoding/json"
	"fmt"
	"github.com/byxorna/marathon_http_proxy_generator/marathon"
	"io/ioutil"
	"log"
	"net/http"
)

func ListenForEvents(listenAddr string) error {
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
	switch e.EventType {
	case "api_post_event":
		eventChan <- e.EventType
	}

}
