package main

import (
	"io"
	"net/http"
)

func ListenForEvents(listenAddr string) error {
	http.HandleFunc("/event", handleEvent)
	return http.ListenAndServe(listenAddr, nil)
}

func handleEvent(res http.ResponseWriter, req *http.Request) {
	//TODO check for POST, what appId it is, what type of event, etc
	io.WriteString(res, "hello, world!\n")
	res.Header().Set("Content-Type", "text/plain")
	updateChan <- "fixme"
}
