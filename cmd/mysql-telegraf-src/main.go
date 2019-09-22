package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jpicht/mysql-influxdb/pkg/datasource"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		if c.NArg() != 1 {
			return errors.New("No config file specified")
		}
		config, err := loadConfigFile(c.Args().First())
		if err != nil {
			return errors.Wrap(err, "Cannot read config file")
		}

		wg := &sync.WaitGroup{}
		now := time.Now()

		for _, host := range config.Hosts {
			logger := logrus.WithField("server", host.Name)

			conn, err := sqlx.Connect("mysql", host.DSN)
			if err != nil {
				logger.WithError(err).Warn("connect failed")
				continue
			}

			baseTags := map[string]string{
				"host": host.Name,
			}

			for n, ds := range map[string]datasource.DataSource{
				"global":        datasource.NewGlobalStatus(conn, config.filter),
				"process list":  datasource.NewProcessList(conn),
				"innodb":        datasource.NewInnoDB(conn),
				"master status": datasource.NewMasterStatus(conn),
			} {
				fnLogger := logger.WithField("ds", n)
				wg.Add(2)
				go func() {
					defer wg.Done()
					defer ds.Close()
					ds.Run()
				}()
				go func() {
					defer wg.Done()
					for pt := range ds.C() {
						pt.AddTagsIfNotExist(host.Tags)
						pt.AddTagsIfNotExist(baseTags)

						p, err := pt.InfluxPoint(now)
						if err != nil {
							fnLogger.WithError(err).Warn()
						}

						fmt.Println(p.String())
					}
				}()
				wg.Wait()
			}
		}

		wg.Wait()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
