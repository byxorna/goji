package main

import (
	"testing"
)

func TestLoadMalformedJson(t *testing.T) {
	_, err := LoadConfig("test/malformed_json.json")
	if err == nil {
		t.Error("Should have barfed on malformed json")
	}
}

func TestLoadSample1Json(t *testing.T) {
	c, err := LoadConfig("test/sample1.json")
	if err != nil {
		t.Error(err)
	}
	if c.MarathonHost != "mymarathonhost" {
		t.Errorf("marathon-host not parsed correctly, found %s", c.MarathonHost)
	}
	if c.MarathonPort != 8080 {
		t.Errorf("marathon-port not set to default 8080, found %d", c.MarathonPort)
	}
	if c.TemplateDelay != 5 {
		t.Errorf("delay parsed incorrectly, found %d", c.TemplateDelay)
	}
	if c.HttpPort != 8000 {
		t.Errorf("http-port not set to default 8000, found %d", c.HttpPort)
	}
}

//TODO test needing at least 1 service in list
//TODO test requiring template file, target file, host
