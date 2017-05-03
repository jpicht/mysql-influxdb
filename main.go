package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ftloc/exception"
	_ "github.com/go-sql-driver/mysql"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/jmoiron/sqlx"
	"github.com/jpicht/logger"
	"github.com/jpicht/mysql-influxdb/influxsender"
)

type (
	Config struct {
		Interval time.Duration
		Filter   string
		Influx   influxsender.Config
		Hosts    []Host
	}
	Host struct {
		Name string
		DSN  string
		Tags map[string]string
	}
	Line struct {
		Name  string `db:"Variable_name"`
		Value string `db:"Value"`
	}
	ProcesslistAbbrevLine struct {
		Command  string `db:"COMMAND"`
		Database string `db:"DB"`
		Host     string `db:"REMOTE"`
		State    string `db:"STATE"`
		Count    int    `db:"COUNT"`
	}
	EventWaitsHostsLine struct {
		Host    string  `db:"HOST"`
		Event   string  `db:"EVENT_NAME"`
		Count   int64   `db:"COUNT_STAR"`
		SumWait float64 `db:"SUM_TIMER_WAIT"`
		MinWait float64 `db:"MIN_TIMER_WAIT"`
		AvgWait float64 `db:"AVG_TIMER_WAIT"`
		MaxWait float64 `db:"MAX_TIMER_WAIT"`
	}
	RunningHostInfo struct {
		Name       string
		Connection *sqlx.DB
		Tags       map[string]string
	}
	TmpPoint struct {
		tags   map[string]string
		values map[string]interface{}
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

func newTmpPointProcList(tags map[string]string, command string, db string, host string, state string) *TmpPoint {
	copy := make(map[string]string)
	for k, v := range tags {
		copy[k] = v
	}
	copy["command"] = command
	copy["db"] = db
	copy["client"] = host
	copy["state"] = state
	return &TmpPoint{
		tags:   copy,
		values: make(map[string]interface{}),
	}
}

func newTmpPoint(tags map[string]string, additional map[string]string, values map[string]interface{}) *TmpPoint {
	copy := make(map[string]string)
	for k, v := range tags {
		copy[k] = v
	}
	for k, v := range additional {
		copy[k] = v
	}
	return &TmpPoint{
		tags:   copy,
		values: values,
	}
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

func send(log logger.Logger, now time.Time, sender influxsender.InfluxSender, measurement string, tags map[string]string, values map[string]interface{}) {
	p, err := influx.NewPoint(
		measurement,
		tags,
		values,
		now,
	)

	if err == nil {
		sender.AddPoint(p)
	} else {
		log.WithData("error", err).Warning("Error creating point")
	}
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

	if interval < 1*time.Second {
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

	sender, err := influxsender.NewSender(c.Influx)
	failOnError(err, "Error: %s", err)

	tick := time.NewTicker(interval)
	sendTick := time.NewTicker(c.Influx.Interval)
	log := logger.NewStdoutLogger()

	failed := 0
	for {
		select {
		case <-tick.C:
			wg := sync.WaitGroup{}
			wg.Add(2 * len(cons))

			for i := range cons {
				now := time.Now()
				now_ts := now.Unix()

				go func(now time.Time) {
					defer wg.Done()
					rows := make([]Line, 0)
					err = cons[i].Connection.Select(&rows, "SHOW GLOBAL STATUS")
					if err != nil {
						log.WithData("error", err).Warning("MySQL-Query-Error")
						return
					}

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
						now,
					)

					if err == nil {
						sender.AddPoint(p)
					} else {
						log.WithData("error", err).Warning("Error creating point")
					}
				}(now)

				go func(now time.Time) {
					defer wg.Done()
					rows := make([]ProcesslistAbbrevLine, 0)
					err = cons[i].Connection.Select(
						&rows,
						"SELECT "+
							"IFNULL(COMMAND, '') COMMAND, "+
							"IFNULL(DB, '') DB, "+
							"SUBSTRING_INDEX(HOST, ':', 1) AS REMOTE, "+
							"IFNULL(STATE, '') STATE, "+
							"COUNT(*) AS `COUNT` "+
							"FROM PROCESSLIST "+
							"WHERE USER NOT IN ('system user', 'repl') "+
							"GROUP BY COMMAND, DB, REMOTE, STATE;",
					)
					if err != nil {
						log.WithData("error", err).Warning("MySQL-Query-Error")
						return
					}

					tmp := map[string]*TmpPoint{}

					for _, row := range rows {
						state := strings.Replace(row.State, " ", "_", -1)
						if state == "" {
							state = row.Command
						}

						key := row.Command + row.Database + row.Host + row.State
						if _, ok := tmp[key]; !ok {
							tmp[key] = newTmpPointProcList(cons[i].Tags, row.Command, row.Database, row.Host, state)
						}

						tmp[key].values["threads"] = row.Count
					}

					for _, tmpP := range tmp {
						p, err := influx.NewPoint(
							"threads",
							tmpP.tags,
							tmpP.values,
							now,
						)

						if err == nil {
							sender.AddPoint(p)
						} else {
							log.WithData("error", err).Warning("Error creating point")
						}
					}
				}(now)

				/*
					go func(now time.Time) {
						defer wg.Done()
						rows := make([]EventWaitsHostsLine, 0)
						err = cons[i].Connection.Select(
							&rows,
							"SELECT * FROM events_waits_summary_by_host_by_event_name WHERE COUNT_STAR > 0;",
						)
						if err != nil {
							log.WithData("error", err).Warning("MySQL-Query-Error")
							return
						}

						for _, row := range rows {
							tmpP := newTmpPoint(
								cons[i].Tags,
								map[string]string{
									"client": row.Host,
									"event":  row.Event,
								},
								map[string]interface{}{
									"count":    row.Count,
									"sum_wait": row.SumWait,
									"min_wait": row.MinWait,
									"max_wait": row.MaxWait,
									"avg_wait": row.AvgWait,
								},
							)
							p, err := influx.NewPoint(
								"waits",
								tmpP.tags,
								tmpP.values,
								now,
							)

							if err == nil {
								sender.AddPoint(p)
							} else {
								log.WithData("error", err).Warning("Error creating point")
							}
						}
					}(now)
				*/
			}
			wg.Wait()
		case <-sendTick.C:
			go sender.Send()
		case <-signals:
			tick.Stop()
			sendTick.Stop()
			sender.Send()
			return
		}
	}
}
