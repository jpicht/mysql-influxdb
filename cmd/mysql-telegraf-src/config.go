package main

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
)

type (
	Config struct {
		Filter string
		Hosts  []Host

		filter *regexp.Regexp
	}
	Host struct {
		Name string
		DSN  string
		Tags map[string]string
	}
)

func loadConfigFile(fn string) (*Config, error) {
	c := &Config{}

	d, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(d, c)
	if err != nil {
		return nil, err
	}

	c.filter, err = regexp.Compile(c.Filter)
	if err != nil {
		return nil, err
	}

	return c, nil
}
