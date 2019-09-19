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
		currentHost *RunningHostInfo
		log         logrus.FieldLogger

		toSender chan datapoint
	}
	datapoint struct {
		name   string
		tags   map[string]string
		values map[string]interface{}
	}
)

func (a *App) sender(s influxsender.InfluxSender) {
	for pt := range a.toSender {
		p, err := influx.NewPoint(
			pt.name,
			pt.tags,
			pt.values,
			a.now,
		)

		if err != nil {
			a.log.WithError(err).WithField("point", pt).Warn("lost point")
			continue
		}

		s.AddPoint(p)
	}
}

func (a *App) Main() {
	if len(os.Args) != 2 {
		fail("Please specify config file")
	}

	a.log = logrus.StandardLogger()

	c := loadConfigFile(os.Args[1])

	var interval = c.Interval

	if interval < 1*time.Second {
		interval = 1 * time.Second
	}

	filter, err := regexp.Compile(c.Filter)
	if err != nil {
		a.log.WithError(err).Fatal("invalid filter")
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	cons := make([]RunningHostInfo, len(c.Hosts))
	for i, host := range c.Hosts {
		con, err := sqlx.Connect("mysql", host.DSN)
		if err != nil {
			a.log.WithError(err).WithField("server", host.Name).Fatal("cannot connect to host")
		}
		tags := host.Tags
		if tags == nil {
			tags = make(map[string]string)
		}
		if _, ok := tags["host"]; !ok {
			tags["host"] = host.Name
		}
		cons[i] = RunningHostInfo{
			Name:       host.Name,
			Connection: con,
			Tags:       tags,
		}
	}

	s, err := influxsender.NewSender(&c.Influx)
	if err != nil {
		a.log.WithError(err).Fatal()
	}

	tick := time.NewTicker(interval)
	sendTick := time.NewTicker(c.Influx.Interval)

	a.wg = &sync.WaitGroup{}
	a.toSender = make(chan datapoint, 1000)
	defer close(a.toSender)
	go a.sender(s)

	failCount := 0
	for {
		select {
		case <-tick.C:
			failed := int32(0)
			a.wg.Add(4 * len(cons))

			for i := range cons {
				a.now = time.Now()
				a.currentHost = &(cons[i])
				locallog := a.log.WithField("server", a.currentHost.Name)

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
				var sleepTime = 2 * time.Second

				if failCount > 10 {
					sleepTime = 60 * time.Second
				} else if failCount > 5 {
					sleepTime = 10 * time.Second
				}

				a.log.WithFields(logrus.Fields{
					"failed": failed,
					"round":  failCount,
				}).Warn("some items failed")

				time.Sleep(sleepTime)
			}
			failed = 0

		case <-sendTick.C:
			go s.Send()
		case <-signals:
			tick.Stop()
			sendTick.Stop()
			a.wg.Wait()
			s.Send()
			return
		}
	}
}
