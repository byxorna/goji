package main

import (
	"fmt"
	"io"
	"net/http"
)

func ListenForEvents(c *Config) error {
	http.HandleFunc("/event", handleEvent)
	return http.ListenAndServe(fmt.Sprintf(":%d", c.HttpPort), nil)
}

func handleEvent(res http.ResponseWriter, req *http.Request) {
	//TODO check for POST, what appId it is, what type of event, etc
	io.WriteString(res, "hello, world!\n")
	res.Header().Set("Content-Type", "text/plain")
}
