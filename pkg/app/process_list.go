package app

import (
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
			"GROUP BY USER, COMMAND, DB, REMOTE, STATE;",
	)
	if err != nil {
		return errors.Wrap(err, "cannot select from processlist")
	}

	for _, row := range rows {
		a.toSender <- datapoint{
			"threads",
			tags(
				a.currentHost.Tags,
				map[string]string{
					"user":    row.User,
					"command": row.Command,
					"db":      row.Database,
					"client":  row.Host,
					"state":   row.State,
				},
			),
			map[string]interface{}{
				"threads": row.Count,
			},
		}
	}

	return nil
}
