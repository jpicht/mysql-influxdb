package app

import (
	"context"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"time"

	"github.com/jpicht/mysql-influxdb/influxsender"
	"github.com/jpicht/mysql-influxdb/pkg/datasource"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type (
	Host struct {
		Name string
		DSN  string
		Tags map[string]string
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

		toSender chan *datasource.DataPoint
	}
)

func (a *App) sender(s influxsender.InfluxSender) {
	for pt := range a.toSender {
		a.log.WithField("pt", pt).Info()
		p, err := pt.InfluxPoint(a.now)

		if err != nil {
			a.log.WithError(err).WithField("point", pt).Warn("lost point")
			continue
		}

		s.AddPoint(p)
	}
}

func (a *App) Main() {
	ctx, cancel := context.WithCancel(context.Background())
	if len(os.Args) != 2 {
		fail("Please specify config file")
	}

	logrus.StandardLogger().Level = logrus.DebugLevel
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
	go func() {
		first := true
		for range signals {
			if first {
				cancel()
				first = false
			} else {
				os.Exit(9)
			}
		}
	}()

	cons := make([]RunningHostInfo, len(c.Hosts))
	for i, host := range c.Hosts {
		a.log.Debug(host)
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
	a.toSender = make(chan *datasource.DataPoint, 1000)
	defer close(a.toSender)
	go a.sender(s)

	failCount := 0
	for {
		select {
		case <-tick.C:
			a.log.Debug("tick")
			failed := int32(0)
			a.wg.Add(4 * len(cons))

			for i := range cons {
				a.now = time.Now()
				a.currentHost = &(cons[i])
				serverLog := a.log.WithField("server", a.currentHost.Name)

				type dataSource interface {
					Run() error
					Close()
					C() <-chan *datasource.DataPoint
				}
				dsWrapper := func(ds dataSource) func() error {
					return func() error {
						defer ds.Close()
						go func() {
							for p := range ds.C() {
								p.AddTagsIfNotExist(a.currentHost.Tags)
								a.toSender <- p
							}
						}()
						return ds.Run()
					}
				}

				for n, fn := range map[string]func() error{
					"global":        dsWrapper(datasource.NewGlobalStatus(a.currentHost.Connection, filter)),
					"process list":  a.procList,
					"innodb":        a.innoStatus,
					"master status": dsWrapper(datasource.NewMasterStatus(a.currentHost.Connection)),
				} {
					localLog := serverLog.WithField("fn", n)
					localLog.Debug("tick")
					if err := fn(); err != nil {
						localLog.WithError(err).Warn()
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
		case <-ctx.Done():
			tick.Stop()
			sendTick.Stop()
			a.wg.Wait()
			s.Send()
			return
		}
	}
}
