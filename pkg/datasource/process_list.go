package datasource

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type (
	ProcessList struct {
		DataSource
	}
	processlistAbbrevLine struct {
		User     string `db:"USER"`
		Command  string `db:"COMMAND"`
		Database string `db:"DB"`
		Host     string `db:"REMOTE"`
		State    string `db:"STATE"`
		Count    int    `db:"COUNT"`
	}
)

func NewProcessList(db *sqlx.DB) *ProcessList {
	return &ProcessList{
		*New(db),
	}
}

func (pl *ProcessList) Run() error {
	rows := make([]processlistAbbrevLine, 0)
	err := pl.connection.Select(
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
		pl.sink <- NewDataPoint(
			"threads",
			map[string]string{
				"user":    row.User,
				"command": row.Command,
				"db":      row.Database,
				"client":  row.Host,
				"state":   row.State,
			},
			map[string]interface{}{
				"threads": row.Count,
			},
		)
	}

	return nil
}
