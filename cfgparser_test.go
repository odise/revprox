package main

import (
	"reflect"
	"testing"
)

func TestConfigParsing(t *testing.T) {
	expected := &Config{

		Proxies: []ProxyConfig{
			ProxyConfig{
				Name:       "api.example.com",
				Secret:     "mysecret",
				AuthPath:   []string{"/"},
				PublicPath: []string{"/"},
				Targets:    []string{"forward1.host.example.com", "forward2.host.example.com"},
				Overwrite:  true,
				Hostname:   "service.example.com",
			},
		},
	}

	config, err := ParseConfig(testConfig)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(config, expected) {
		t.Error("Config structure differed from expectation")
	}
}

const testConfig = `

proxy "api.example.com" {
  secret = "mysecret"
  authpath = [ "/" ]
  publicpath = [ "/" ]
  target = ["forward1.host.example.com", "forward2.host.example.com"]
  overwritehost = true
  hostname = "service.example.com"
  }

`
