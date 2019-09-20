package app

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"time"

	"github.com/jpicht/mysql-influxdb/influxsender"
)

type (
	Config struct {
		Interval time.Duration
		Filter   string
		Influx   influxsender.Config
		Hosts    []Host

		filter *regexp.Regexp
	}
)

func loadConfigFile(fn string) *Config {
	c := &Config{}

	d, err := ioutil.ReadFile(fn)
	failOnError(err, "Cannot read config file: %s", err)

	err = json.Unmarshal(d, c)
	failOnError(err, "Cannot parse config file: %s", err)

	c.filter, err = regexp.Compile(c.Filter)
	failOnError(err, "invalid filter regex: %s", err)

	if c.Interval < 1*time.Second {
		c.Interval = 1 * time.Second
	}

	return c
}

func (c *Config) UnmarshalJSON(data []byte) error {
	var temp struct {
		Interval string
		Filter   string
		Influx   influxsender.Config
		Hosts    []Host
	}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	c.Filter = temp.Filter
	c.Influx = temp.Influx
	c.Hosts = temp.Hosts
	tmp_duration, err := time.ParseDuration(temp.Interval)
	if err != nil {
		return err
	}
	c.Interval = tmp_duration
	return nil
}
