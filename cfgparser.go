package main

import (
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
)

// Config type
type Config struct {
	Proxies []ProxyConfig `hcl:"proxy"`
}

// ProxyConfig type
type ProxyConfig struct {
	Name       string   `hcl:",key"`
	Secret     string   `hcl:"secret"`
	AuthPath   []string `hcl:"authpath"`
	PublicPath []string `hcl:"publicpath"`
	Targets    []string `hcl:"target"`
	Overwrite  bool     `hcl:"overwritehost"`
	Hostname   string   `hcl:"hostname"`
}

// ParseConfig parse the given HCL string into a Config struct.
func ParseConfig(hclText string) (*Config, error) {
	result := &Config{}
	var errors *multierror.Error

	hclParseTree, err := hcl.Parse(hclText)
	if err != nil {
		return nil, err
	}

	if err := hcl.DecodeObject(&result, hclParseTree); err != nil {
		return nil, err
	}

	log.Printf("%+v\n", result)

	return result, errors.ErrorOrNil()
}
