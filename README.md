Goji: Marathon Task Proxy Config Generator
===================

```goji``` is a server that registers with a Marathon instance, consumes events, and emits templated configs containing information about running tasks for a set of apps that you care about.

[![Build Status](https://travis-ci.org/byxorna/goji.svg)](https://travis-ci.org/byxorna/goji)
[![Build Status](https://drone.io/github.com/byxorna/goji/status.png)](https://drone.io/github.com/byxorna/goji/latest)

## Run

### One-shot mode

Just hit up marathon for the task list, generate a config, and write it out.

```./goji -conf myconfig.json```

### Server

Generates a config like one-shot mode, but ```goji``` will then listen on ```http-port``` for events from Marathon.

```./goji -conf myconfig.json -server```

### Docker

An automated build of master is available at [registry.hub.docker.com/u/byxorna/goji](https://registry.hub.docker.com/u/byxorna/goji/). The container is more useful when run in single shot mode, as opposed to ```-server``` mode. You can use dockers ```-v``` argument to mount directories into the container to provide a config or output directory.

```
# docker run -it byxorna/goji -h
Usage of ./goji:
  -app-required=false: Require marathon applications to exist (assumes no tasks for missing apps if false)
  -conf="": Config JSON file
  -server=false: Start a HTTP server listening for Marathon events
  -target="": Target file to write to
```

## Configuration

```goji``` takes a config file, formatted in json, as the ```-conf``` option. It tells ```goji``` information about your marathon instance, what services (app IDs) to query for tasks, and where and with what template to write configs out.

```
{
  "marathon-host":"marathon1.tumblr.net",
  "marathon-port":8080,
  "services": [
    { "app-id": "/sre/byxorna/site", "name": "pipefail.service.iata.tumblr.net" },
    { "app-id": "/sre/byxorna/app1", "name": "app1.service.iata.tumblr.net", "protocol":"TCP" },
    { "app-id": "/sre/byxorna/webapp", "name": "web.service.iata.tumblr.net", "protocol":"HTTP", "health-check":"/_health" },
    { "app-id": "/sre/byxorna/appwithopts", "name": "myapp", "protocol":"HTTP", "health-check":"/_health",
      "options": { "acl-match": "-m beg", "healthcheck-rate":"100"}}
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
* ```services```: List of Services. A service is an object with a ```app-id``` key of the marathon app ID you want tasks from, and a ```name``` that will be passed into your template for each service. See below.

### Service Configuration

```
{
  // marathon app id for your application
  "app-id": "/sre/byxorna/webapp",
  // a name that is associated with the service. Useful for doing nginx vhosting for http apps, or service name for DNS SRV records. Just a string
  "name": "web.service.iata.tumblr.net",
  // TCP or HTTP, defaults to HTTP
  "protocol":"HTTP",
  // if the protocol of the service is HTTP, you can specify a health check URI here
  "health-check":"/_health",
  // what service port to use. Defaults to 80, but you can override for TCP services
  "port":80
  // options is an optional map[string]string of arbitrary options you can switch on in your templates
  // these are useful to specify behavior logic per service
  options: { "acl-match": "-m beg", "healthcheck-rate":"100" }
}
```

## Commands

The ```comand``` field in the config json specifies a command to run upon successful creation of a new config. This can be anything you want, but here are some useful examples:

#### HAProxy Reloading
```
...
"command":"cp /etc/haproxy/haproxy.cfg /etc/haproxy/haproxy.cfg.bak && cp ./haproxy.cfg /etc/haproxy/haproxy.cfg && service haproxy check && service haproxy reload || (mv /etc/haproxy/haproxy.cfg.bak /etc/haproxy/haproxy.cfg && exit 1)",
...
```

## Templates

You can use ```goji``` to emit whatever configs you care about. Common usecases would be HAproxy or Nginx configurations. The ```name``` attribute could be useful for doing nginx/apache vhosting to identify a request by ```Host:``` header.

Here is an example:

```
This is a test template
{{ range $index, $service := . }}AppId {{ $service.AppId }} at {{ $service.Name }}{{ range $i, $task := $service.Tasks }}
{{ range $j, $port := $task.Ports }}  task: {{ $task.Id }} {{ $task.Host }}:{{ $port }}{{ end }}{{ end }}{{ end }}
Sweet!
```

A much more useful template is available in `example/haproxy.tmpl`, `example/nginx.tmpl`, and `example/named.tmpl`

You may use arbitrary keys and values in `$service.Options` to do clever things per service. For example, this snippet will allow you to modify how the haproxy acl rule works with options stored per service.
```acl {{ $service.EscapeAppIdColon }}-aclrule hdr(host) {{with index $service.Options "acl-opts"}}{{index $service.Options "acl-opts"}} {{end}}{{ $service.Name }}```

### DNS SRV Records

You can generate a named zone with SRV records from marathon backends trivially, given the template:

```
@ IN SOA ns1.example.com. admin.example.com. (
  12345      ; serial
  600        ; refresh
  1800       ; retry
  604800     ; expire
  300        ; minimum
  )

  IN NS ns1.example.com.
  IN NS ns2.example.com.

$ORIGIN goji.example.com.
; service.proto.owner-name     ttl   class   rr    pri   weight    port    target
; _http._tcp.goji.example.com. 60    IN      SRV   0     5         301234  ct-12345.iata.example.com.
{{ range $index, $service := . }}
; SRV for {{$service.Name}}
{{ range $i, $task := $service.Tasks }}{{ range $j, $port := $task.Ports }}_{{$service.Name}}._{{$service.Protocol}}  {{with index $service.Options "ttl"}}{{index $service.Options "ttl"}}{{else}}60{{end}}  IN SRV 0 5 {{$port}} {{$task.Host}}.{{end}}
{{end}}{{end}}
```

And config:

```
{
  "marathon-host":"marathon.example.com",
  "services": [
    { "app-id": "/sre/website", "name": "http", "protocol":"HTTP"},
    { "app-id": "/sre/tracker", "name": "tracker","protocol":"UDP","options":{"ttl":"120"}},
    { "app-id": "/sre/tcpservice", "name": "myservice", "protocol":"TCP"}
  ],
  "template": "example/named.tmpl",
  "target": "./named.zone.cfg",
  "command": "named-checkzone goji.example.com ./named.zone.cfg && cp ./named.zone.cfg /var/named/services.zone && service named reload",
  "delay": 5
}
```


## Build

```
$ go build
$ ./goji -help
Usage of ./goji:
  -app-required=false: Require marathon applications to exist (assumes no tasks for missing apps if false)
  -conf="": Config JSON file
  -server=false: Start a HTTP server listening for Marathon events
  -target="": Target file to write to
```

## More information

Similar in function to https://github.com/QubitProducts/bamboo, but with fewer features, and not HAproxy specific.
https://mesosphere.github.io/marathon/docs/event-bus.html

## TODO

* event coalescing is a bit wonky and fires on very first event
* Write test suite
* Setup travis-ci builds

