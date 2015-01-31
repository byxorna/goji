Marathon HTTP Proxy Generator
===================

Reverse proxy generator for apps running in Marathon

## TODO

* fix deadlock in receiving from channel (http handler is blocking waiting on updateChan when getting event)
* have http listener push messages into channel to trigger templating
* Start an http server and register with marathon for callbacks
* Run a command after generation (check and reload)
* Write documentation

https://mesosphere.github.io/marathon/docs/event-bus.html
