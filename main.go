package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type (
	Config struct {
		Interval time.Duration
		Filter   string
		Influx   Influx
		Hosts    []Host
	}
	Host struct {
		Name string
		DSN  string
		Tags map[string]string
	}
	Influx struct {
		Interval time.Duration
		Hostname string
		Database string
	}
	Line struct {
		Name  string `db:"Variable_name"`
		Value string `db:"Value"`
	}
	RunningHostInfo struct {
		Name       string
		Connection *sqlx.DB
		Tags       map[string]string
	}

	InfluxSender struct {
		lock   sync.Mutex
		data   influx.BatchPoints
		client influx.Client
	}
)

func fail(f string, data ...interface{}) {
	fmt.Fprintf(os.Stderr, f+"\n", data...)
	os.Exit(1)
}

func failOnError(err error, f string, data ...interface{}) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, f+"\n", data...)
	os.Exit(1)
}

func newSender(c Influx) *InfluxSender {
	batchPoints, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: c.Database,
	})
	failOnError(err, "Cannot create influx batch: %s", err)

	sender, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:    "http://" + c.Hostname + ":8086",
		Timeout: 5 * time.Second,
	})
	failOnError(err, "Cannot create influxdb output: %s", err)

	sender.Query(influx.NewQuery(
		"CREATE DATABASE "+c.Database,
		"",
		"",
	))

	return &InfluxSender{
		lock:   sync.Mutex{},
		data:   batchPoints,
		client: sender,
	}
}

func (s *InfluxSender) AddPoint(point *influx.Point) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data.AddPoint(point)
}

func (s *InfluxSender) Send() {
	defer s.lock.Unlock()
	s.lock.Lock()
	batchPoints, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: s.data.Database(),
	})
	failOnError(err, "Cannot create influx batch: %s", err)

	data := s.data
	s.data = batchPoints

	go func () {
		err := s.client.Write(data)
		fmt.Printf("sent: %d %s\n", len(data.Points()), err)
	}()
}

func (c *Config) UnmarshalJSON(data []byte) error {
	var temp struct {
		Interval string
		Filter   string
		Influx   Influx
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

func (i* Influx) UnmarshalJSON(data []byte) error {
	var temp struct {
		Interval string
		Hostname string
		Database string
	}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	i.Hostname = temp.Hostname
	i.Database = temp.Database
	tmp_duration, err := time.ParseDuration(temp.Interval)
	if err != nil {
		return err
	}
	i.Interval = tmp_duration
	return nil
}

func (i *Influx) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}

func main() {
	if len(os.Args) != 2 {
		fail("Please specify config file")
	}

	c := &Config{}
	d, err := ioutil.ReadFile(os.Args[1])
	failOnError(err, "Cannot read config file: %s", err)
	err = json.Unmarshal(d, c)
	failOnError(err, "Cannot parse config file: %s", err)

	var interval = c.Interval

	if interval < 1 * time.Second{
		interval = 1 * time.Second
	}

	filter, err := regexp.Compile(c.Filter)
	failOnError(err, "Invalid filter: %s", err)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	cons := make([]RunningHostInfo, len(c.Hosts))
	for i := range c.Hosts {
		con, err := sqlx.Connect("mysql", c.Hosts[i].DSN)
		failOnError(err, "Cannot connect to '%s': %s", c.Hosts[i].Name, err)
		tags := c.Hosts[i].Tags
		if tags == nil {
			tags = make(map[string]string)
		}
		if _, ok := tags["host"]; !ok {
			tags["host"] = c.Hosts[i].Name
		}
		cons[i] = RunningHostInfo{
			Name:       c.Hosts[i].Name,
			Connection: con,
			Tags:       tags,
		}
	}

	sender := newSender(c.Influx)

	tick := time.NewTicker(interval)
	sendTick := time.NewTicker(c.Influx.Interval)

	for {
		select {
		case <-tick.C:
			wg := sync.WaitGroup{}
			wg.Add(len(cons))

			for i := range cons {
				go func() {
					defer wg.Done()
					rows := make([]Line, 0)
					err = cons[i].Connection.Select(&rows, "SHOW GLOBAL STATUS")
					failOnError(err, "Error: %s", err)

					values := map[string]interface{}{}

					for ii := range rows {
						v := rows[ii]

						if !filter.MatchString(v.Name) {
							continue
						}

						values[v.Name], _ = strconv.Atoi(v.Value)
					}

					p, err := influx.NewPoint(
						"mysql",
						cons[i].Tags,
						values,
						time.Now(),
					)

					failOnError(err, "Error: %s", err)

					sender.AddPoint(p)
				}()
			}
			wg.Wait()
		case <-sendTick.C:
			go sender.Send()
		case <-signals:
			tick.Stop()
			return
		}
	}
}
