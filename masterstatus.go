package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ftloc/exception"
	"github.com/jpicht/logger"
	"github.com/jpicht/mysql-influxdb/influxsender"
)

func (a *app) masterStatus(wg *sync.WaitGroup, log logger.Logger, info *RunningHostInfo, sender influxsender.InfluxSender, failed *int32) {
	defer wg.Done()
	exception.Try(func() {
		type helper struct {
			BinlogSize int64 `db:"s"`
		}
		max := &helper{0}
		info.Connection.Get(max, "SELECT @@max_binlog_size AS s")
		type MasterData struct {
			File     string `db:"File"`
			Position int64  `db:"Position"`
			DoDb     string `db:"Binlog_Do_DB"`
			IngoreDb string `db:"Binlog_Ignore_DB"`
			ExGtid   string `db:"Executed_Gtid_Set"`
		}
		row := info.Connection.QueryRowx("SHOW MASTER STATUS")
		data := &MasterData{}
		err := row.StructScan(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		dotPos := strings.LastIndex(data.File, ".")
		if dotPos == -1 {
			fmt.Println("pos -1")
			return
		}
		fileNum, err := strconv.Atoi(data.File[dotPos+1:])
		if err != nil {
			fmt.Println(err)
			return
		}
		position := int64(fileNum)*max.BinlogSize + data.Position
		a.send(log, sender, "replication", info.Tags, map[string]interface{}{
			"master_position": position,
		})
	}).CatchAll(func(i interface{}) {
		log.Alertf("Exception caught: %#v", i)
		ok, file, line := exception.GetThrower()
		if ok {
			log.Alertf("Location: %s:%d", file, line)
		} else {
			cs := make([]uintptr, 20)
			amount := runtime.Callers(1, cs)

			for i := 0; i < amount; i++ {
				f := runtime.FuncForPC(cs[i])
				file, line := f.FileLine(cs[i])
				log.Alert("  " + f.Name())
				log.Alertf("      %s:%d", file, line)
			}
		}
		atomic.AddInt32(failed, 1)
	}).Go()
}
