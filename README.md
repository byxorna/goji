Marathon HTTP Proxy Generator
===================

Reverse proxy generator for apps running in Marathon

## TODO

* event coalescing is a bit wonky and fires on very first event
* register with marathon for callbacks at startup
* cleanup callback registration when shutting down?
* Run a command after generation (check and reload)
* write flags for all config options
* Write documentation

https://mesosphere.github.io/marathon/docs/event-bus.html
