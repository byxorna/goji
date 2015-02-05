Goji: Marathon Task Proxy Config Generator
===================

```goji``` is a server that registers with a Marathon instance, consumes events, and emits templated configs containing information about running tasks for a set of apps that you care about.

## Run

### One-shot mode

Just hit up marathon for the task list, generate a config, and write it out.

```./goji -conf myconfig.json```

### Server

Generates a config like one-shot mode, but ```goji``` will then listen on ```http-port``` for events from Marathon.

```./goji -conf myconfig.json -server```

## Configuration

```goji``` takes a config file, formatted in json, as the ```-conf``` option. It tells ```goji``` information about your marathon instance, what services (app IDs) to query for tasks, and where and with what template to write configs out.

```
{
  "marathon-host":"marathon1.tumblr.net",
  "marathon-port":8080,
  "services": [
    { "app-id": "/sre/byxorna/site", "vhost": "pipefail.service.iata.tumblr.net" },
    { "app-id": "/sre/byxorna/app1", "vhost": "app1.service.iata.tumblr.net" }
  ],
  "template": "templates/haproxy.tmpl",
  "target": "/tmp/haproxy.cfg",
  "http-port": 8000,
  "delay": 5
}
```

* ```marathon-host```: Hostname of the marathon instance to connect to (required)
* ```marathon-port```: Port the marathon service is running on (optional, default: 8080)
* ```template```: Template config file to feed services and tasks into, ```text/template``` format (required)
* ```target```: Write templated configuration to this location (required)
* ```http-port```: What port to start an HTTP event listener on to register and receive event messages from marathon (optional, default: 8000)
* ```delay```: Coalesce events within this window before triggering a task get and config emit (optional, default: 0)
* ```listen```: Register eventhandler with Marathon, and listen on ```http-port``` for events (optional, default: true)
* ```services```: List of Services. A service is an object with a ```app-id``` key of the marathon app ID you want tasks from, and a ```vhost``` that will be passed into your template for each service.


## Templates

You can use ```goji``` to emit whatever configs you care about. Common usecases would be HAproxy or Nginx configurations. The ```vhost``` attribute could be useful for doing nginx/apache vhosting to identify a request by ```Host:``` header.

Here is an example:

```
This is a test template
{{ range $index, $service := . }}AppId {{ $service.AppId }} at {{ $service.Vhost }}{{ range $i, $task := $service.Tasks }}
{{ range $j, $port := $task.Ports }}  task: {{ $task.Id }} {{ $task.Host }}:{{ $port }}{{ end }}{{ end }}{{ end }}
Sweet!
```

## Build

```
$ go build
$ ./goji -help
Usage of ./goji:
  -conf="": config json file
```

## More information

Similar in function to https://github.com/QubitProducts/bamboo, but with fewer features, and not HAproxy specific.
https://mesosphere.github.io/marathon/docs/event-bus.html

## TODO

* Fix Service to specify TCP or HTTP
* Fix up specifying a health-check in a service (does this work?)
* event coalescing is a bit wonky and fires on very first event
* Run a command after generation (check and reload)

