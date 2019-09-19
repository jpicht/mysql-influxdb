package app

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func (a *App) masterStatus() error {
	defer a.wg.Done()
	type helper struct {
		BinlogSize int64 `db:"s"`
	}
	max := &helper{0}
	a.currentHost.Connection.Get(max, "SELECT @@max_binlog_size AS s")
	type MasterData struct {
		File     string `db:"File"`
		Position int64  `db:"Position"`
		DoDb     string `db:"Binlog_Do_DB"`
		IngoreDb string `db:"Binlog_Ignore_DB"`
		ExGtid   string `db:"Executed_Gtid_Set"`
	}
	row := a.currentHost.Connection.QueryRowx("SHOW MASTER STATUS")
	data := &MasterData{}
	err := row.StructScan(data)
	if err != nil {
		return errors.Wrap(err, "struct scan failed")
	}
	dotPos := strings.LastIndex(data.File, ".")
	if dotPos == -1 {
		return errors.Wrap(err, "data file does not contain '.'")
	}
	fileNum, err := strconv.Atoi(data.File[dotPos+1:])
	if err != nil {
		return errors.Wrap(err, "cannot parse numeric part of filename")
	}
	position := int64(fileNum)*max.BinlogSize + data.Position
	a.toSender <- datapoint{
		"replication",
		a.currentHost.Tags,
		map[string]interface{}{
			"master_position": position,
		},
	}
	return nil
}
