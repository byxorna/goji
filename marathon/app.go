package marathon

/*
{
  "id": "/myapp",
  "cmd": "env && sleep 60",
  "args": null,
  "user": null,
  "env": {
    "LD_LIBRARY_PATH": "/usr/local/lib/myLib"
  },
  "instances": 3,
  "cpus": 0.1,
  "mem": 5,
  "disk": 0,
  "executor": "",
  "constraints": [
    [
      "hostname",
      "UNIQUE",
      ""
    ]
  ],
  "uris": [
    "https://raw.github.com/mesosphere/marathon/master/README.md"
  ],
  "storeUrls": [],
  "ports": [
    10013,
    10015
  ],
  "requirePorts": false,
  "backoffSeconds": 1,
  "backoffFactor": 1.15,
  "maxLaunchDelaySeconds": 3600,
  "container": null,
  "healthChecks": [],
  "dependencies": [],
  "upgradeStrategy": {
    "minimumHealthCapacity": 1,
    "maximumOverCapacity": 1
  },
  "labels": {},
  "acceptedResourceRoles": null,
  "version": "2015-09-25T15:13:48.343Z",
  "versionInfo": {
    "lastScalingAt": "2015-09-25T15:13:48.343Z",
    "lastConfigChangeAt": "2015-09-25T15:13:48.343Z"
  },
  "tasksStaged": 0,
  "tasksRunning": 0,
  "tasksHealthy": 0,
  "tasksUnhealthy": 0,
  "deployments": [
    {
      "id": "9538079c-3898-4e32-aa31-799bf9097f74"
    }
  ]
}
*/

// for now, just stub out the bits we could possibly care about for a app
// which boils down to env and labels, so we can detect an app as being eligible
// for goji
type App struct {
	Id     string            `json:"id"`
	Env    map[string]string `json:"env"`
	Labels map[string]string `json:"labels"`
}
type AppList []App

func (a AppList) Len() int           { return len(a) }
func (a AppList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AppList) Less(i, j int) bool { return a[i].Id < a[j].Id }
func (a *App) String() string {
	return a.Id
}
