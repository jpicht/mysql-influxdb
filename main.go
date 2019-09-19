package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ftloc/exception"
	_ "github.com/go-sql-driver/mysql"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/jmoiron/sqlx"
	"github.com/jpicht/logger"
	"github.com/jpicht/mysql-influxdb/influxsender"
)

type (
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
		User     string `db:"USER"`
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
	app struct {
		now         time.Time
		wg          *sync.WaitGroup
		sender      influxsender.InfluxSender
		currentHost *RunningHostInfo
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

func (a *app) send(log logger.Logger, measurement string, tags map[string]string, values map[string]interface{}) {
	p, err := influx.NewPoint(
		measurement,
		tags,
		values,
		a.now,
	)

	if err == nil {
		a.sender.AddPoint(p)
	} else {
		log.WithData("error", err).Warning("Error creating point")
	}
}

func (a *app) globalStatus(log logger.Logger, filter *regexp.Regexp, failed *int32) {
	defer a.wg.Done()
	exception.Try(func() {
		rows := make([]Line, 0)
		err := a.currentHost.Connection.Select(&rows, "SHOW GLOBAL STATUS")
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

		a.send(log, "mysql", a.currentHost.Tags, values)
	}).CatchAll(func(interface{}) {
		atomic.AddInt32(failed, 1)
	}).Go()
}

func (a *app) procList(log logger.Logger, failed *int32) {
	defer a.wg.Done()
	exception.Try(func() {
		rows := make([]ProcesslistAbbrevLine, 0)
		err := a.currentHost.Connection.Select(
			&rows,
			"SELECT "+
				"IFNULL(USER, '') USER, "+
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

			key := row.User + row.Command + row.Database + row.Host + row.State
			if _, ok := tmp[key]; !ok {
				tmp[key] = newTmpPointProcList(a.currentHost.Tags, row.User, row.Command, row.Database, row.Host, state)
			}

			tmp[key].values["threads"] = row.Count
		}

		for _, tmpP := range tmp {
			a.send(log, "threads", tmpP.tags, tmpP.values)
		}
	}).CatchAll(func(interface{}) {
		atomic.AddInt32(failed, 1)
	}).Go()
}

func tags(base map[string]string, additional map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(additional))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range additional {
		out[k] = v
	}
	return out
}

func (a *app) main() {
	if len(os.Args) != 2 {
		fail("Please specify config file")
	}

	c := loadConfigFile(os.Args[1])

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

	s, err := influxsender.NewSender(&c.Influx)
	failOnError(err, "Error: %s", err)
	a.sender = s

	tick := time.NewTicker(interval)
	sendTick := time.NewTicker(c.Influx.Interval)
	log := logger.NewStdoutLogger()

	a.wg = &sync.WaitGroup{}
	failCount := 0
	for {
		select {
		case <-tick.C:
			failed := int32(0)
			a.wg.Add(4 * len(cons))

			for i := range cons {
				a.now = time.Now()
				a.currentHost = &(cons[i])
				locallog := log.WithData("server", a.currentHost.Name)

				a.globalStatus(locallog, filter, &failed)
				a.procList(locallog, &failed)
				a.innoStatus(locallog, &failed)
				a.masterStatus(locallog, &failed)
			}
			// synchronous at the moment, but whatever
			a.wg.Wait()

			if failed > 0 {
				failCount++
				log.Warningf("Some items failed (%d/%d items, %d consecutive rounds)", failed, len(cons)*4, failCount)
				if failCount > 10 {
					time.Sleep(60 * time.Second)
				} else if failCount > 5 {
					time.Sleep(10 * time.Second)
				} else {
					time.Sleep(2 * time.Second)
				}
			} else {
				failed = 0
			}

		case <-sendTick.C:
			go a.sender.Send()
		case <-signals:
			tick.Stop()
			sendTick.Stop()
			a.sender.Send()
			return
		}
	}
}

func main() {
	a := &app{}
	a.main()
}
