package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"runtime"
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
	RunningHostInfo struct {
		Name       string
		Connection *sqlx.DB
		Tags       map[string]string
	}
	App struct {
		now time.Time
		wg  *sync.WaitGroup
		log logrus.FieldLogger

		c      *Config
		filter *regexp.Regexp

		ctx         context.Context
		cancel      context.CancelFunc
		cons        []RunningHostInfo
		failCount   int
		currentHost *RunningHostInfo

		toSender chan *datasource.DataPoint
	}
)

func (a *App) sender(s influxsender.InfluxSender) {
	for pt := range a.toSender {
		p, err := pt.InfluxPoint(a.now)

		if err != nil {
			a.log.WithError(err).WithField("point", pt).Warn("lost point")
			continue
		}

		s.AddPoint(p)
	}
}

func (a *App) init() {
	a.ctx, a.cancel = context.WithCancel(context.Background())
	if len(os.Args) != 2 {
		fail("Please specify config file")
	}

	logrus.StandardLogger().Level = logrus.DebugLevel
	a.log = logrus.StandardLogger()

	a.c = loadConfigFile(os.Args[1])

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		first := true
		for range signals {
			if first {
				a.cancel()
				first = false
			} else {
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, true)
				fmt.Printf("%s", buf)
				os.Exit(9)
			}
		}
	}()
}

func (a *App) connect() {
	a.cons = make([]RunningHostInfo, len(a.c.Hosts))
	for i, host := range a.c.Hosts {
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
		a.cons[i] = RunningHostInfo{
			Name:       host.Name,
			Connection: con,
			Tags:       tags,
		}
	}
}

func (a *App) Main() {
	a.init()
	a.connect()
	s, err := influxsender.NewSender(&(a.c.Influx))
	if err != nil {
		a.log.WithError(err).Fatal()
	}

	tick := time.NewTicker(a.c.Interval)
	sendTick := time.NewTicker(a.c.Influx.Interval)

	a.wg = &sync.WaitGroup{}
	a.toSender = make(chan *datasource.DataPoint, 1000)
	defer close(a.toSender)
	go a.sender(s)

	for {
		select {
		case <-tick.C:
			a.tick()

		case <-sendTick.C:
			go s.Send()

		case <-a.ctx.Done():
			a.log.Info("Quitting")
			tick.Stop()
			sendTick.Stop()
			a.log.Info("Waiting for workers")
			a.wg.Wait()
			a.log.Info("Pushing data")
			s.Send()
			return
		}
	}
}

func (a *App) tick() {
	a.log.Debug("tick")
	failed := 0

	a.wg.Add(4 * len(a.cons))

	for i := range a.cons {
		if a.ctx.Err() != nil {
			break
		}

		a.now = time.Now()
		a.currentHost = &(a.cons[i])
		serverLog := a.log.WithField("server", a.currentHost.Name)

		type dataSource interface {
			Run() error
			Close()
			C() <-chan *datasource.DataPoint
		}
		dsWrapper := func(ds dataSource) func() error {
			return func() error {
				defer ds.Close()
				defer a.wg.Done()
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
			"global":        dsWrapper(datasource.NewGlobalStatus(a.currentHost.Connection, a.c.filter)),
			"process list":  dsWrapper(datasource.NewProcessList(a.currentHost.Connection)),
			"innodb":        dsWrapper(datasource.NewInnoDB(a.currentHost.Connection)),
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
		a.failCount++
		var sleepTime = 2 * time.Second

		if a.failCount > 10 {
			sleepTime = 60 * time.Second
		} else if a.failCount > 5 {
			sleepTime = 10 * time.Second
		}

		a.log.WithFields(logrus.Fields{
			"failed": failed,
			"round":  a.failCount,
		}).Warn("some items failed")

		if a.ctx.Err() == nil {
			time.Sleep(sleepTime)
		}
	} else {
		a.failCount = 0
	}
}
