package app

import (
	"os"
	"os/signal"
	"regexp"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/jmoiron/sqlx"
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
	RunningHostInfo struct {
		Name       string
		Connection *sqlx.DB
		Tags       map[string]string
	}
	App struct {
		now         time.Time
		wg          *sync.WaitGroup
		sender      influxsender.InfluxSender
		currentHost *RunningHostInfo
		log         logrus.FieldLogger
	}
)

func (a *App) send(measurement string, tags map[string]string, values map[string]interface{}) error {
	p, err := influx.NewPoint(
		measurement,
		tags,
		values,
		a.now,
	)

	if err == nil {
		a.sender.AddPoint(p)
	}

	return err
}

func (a *App) Main() {
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
	log := logrus.StandardLogger()

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
				locallog := log.WithField("server", a.currentHost.Name)

				for _, fn := range []func() error{
					func() error { return a.globalStatus(filter) },
					a.procList,
					a.innoStatus,
					a.masterStatus,
				} {
					if err := fn(); err != nil {
						locallog.WithError(err).Warn()
						failed++
					}
				}
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
