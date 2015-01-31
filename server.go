package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func ListenForEvents(listenAddr string) error {
	http.HandleFunc("/event", handleEvent)
	return http.ListenAndServe(listenAddr, nil)
}

func handleEvent(res http.ResponseWriter, req *http.Request) {
	//TODO check for POST, what appId it is, what type of event, etc
	log.Printf("Got an event!")
	io.WriteString(res, "hello, world!\n")
	res.Header().Set("Content-Type", "text/plain")
	eventChan <- "fixme"
	fmt.Fprintf(res, "Thanks for the event!")
	log.Printf("All done")
}
