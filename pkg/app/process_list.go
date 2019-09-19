package app

import (
	"strings"

	"github.com/pkg/errors"
)

func (a *App) procList() error {
	defer a.wg.Done()
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
		return errors.Wrap(err, "cannot select from processlist")
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
		if err := a.send("threads", tmpP.tags, tmpP.values); err != nil {
			return err
		}
	}
	return nil
}
