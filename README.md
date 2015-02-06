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
    { "app-id": "/sre/byxorna/app1", "vhost": "app1.service.iata.tumblr.net", "protocol":"TCP" },
    { "app-id": "/sre/byxorna/webapp", "vhost": "web.service.iata.tumblr.net", "protocol":"HTTP", "health-check":"/_health" }
  ],
  "template": "templates/haproxy.tmpl",
  "target": "/tmp/haproxy.cfg",
  "command": "/usr/bin/check_haproxy_config /tmp/haproxy.cfg && cp /tmp/haproxy.cfg /etc/haproxy/haproxy.cfg && service haproxy reload",
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
* ```command```: Run a script after writing out the config (optional, default: empty)
* ```services```: List of Services. A service is an object with a ```app-id``` key of the marathon app ID you want tasks from, and a ```vhost``` that will be passed into your template for each service. See below.

### Service Configuration

```
{
  // marathon app id for your application
  "app-id": "/sre/byxorna/webapp",
  // a vhost that is associated with the service. Useful for doing nginx vhosting for http apps
  "vhost": "web.service.iata.tumblr.net",
  // TCP or HTTP, defaults to HTTP
  "protocol":"HTTP",
  // if the protocol of the service is HTTP, you can specify a health check URI here
  "health-check":"/_health",
  // what service port to use. Defaults to 80, but you can override for TCP services
  "port":80
}
```

## Commands

The ```comand``` field in the config json specifies a command to run upon successful creation of a new config. This can be anything you want, but here are some useful examples:

#### HAProxy Reloading
```
...
"command":"cp /etc/haproxy/haproxy.cfg /etc/haproxy/haproxy.cfg.bak && cp ./haproxy.cfg /etc/haproxy/haproxy.cfg && service haproxy check && service haproxy restart || mv /etc/haproxy/haproxy.cfg.bak /etc/haproxy/haproxy.cfg",
...
```

## Templates

You can use ```goji``` to emit whatever configs you care about. Common usecases would be HAproxy or Nginx configurations. The ```vhost``` attribute could be useful for doing nginx/apache vhosting to identify a request by ```Host:``` header.

Here is an example:

```
This is a test template
{{ range $index, $service := . }}AppId {{ $service.AppId }} at {{ $service.Vhost }}{{ range $i, $task := $service.Tasks }}
{{ range $j, $port := $task.Ports }}  task: {{ $task.Id }} {{ $task.Host }}:{{ $port }}{{ end }}{{ end }}{{ end }}
Sweet!
```

A much more useful template is available in ```example/haproxy.tmpl``` and ```example/nginx.tmpl```

## Build

```
$ go build
$ ./goji -help
Usage of ./goji:
  -conf="": config json file
  -server=false: start a HTTP server listening for Marathon events
  -target="": target file to write to
```

## More information

Similar in function to https://github.com/QubitProducts/bamboo, but with fewer features, and not HAproxy specific.
https://mesosphere.github.io/marathon/docs/event-bus.html

## TODO

* event coalescing is a bit wonky and fires on very first event
* Write test suite
* Setup travis-ci builds

