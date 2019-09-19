package app

import (
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

func (a *App) globalStatus(filter *regexp.Regexp) error {
	defer a.wg.Done()
	rows := make([]Line, 0)
	err := a.currentHost.Connection.Select(&rows, "SHOW GLOBAL STATUS")
	if err != nil {
		return errors.Wrap(err, "global status query failed")
	}

	values := map[string]interface{}{}

	for ii := range rows {
		v := rows[ii]

		if !filter.MatchString(v.Name) {
			continue
		}

		values[v.Name], _ = strconv.Atoi(v.Value)
	}

	return a.send("mysql", a.currentHost.Tags, values)
}
