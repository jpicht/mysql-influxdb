package datasource

import (
	"regexp"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type (
	globalStatusLine struct {
		Name  string `db:"Variable_name"`
		Value string `db:"Value"`
	}
	GlobalStatus struct {
		BaseDataSource

		Filter *regexp.Regexp
	}
)

func NewGlobalStatus(db *sqlx.DB, filter *regexp.Regexp) *GlobalStatus {
	return &GlobalStatus{
		BaseDataSource: *New(db),
		Filter:         filter,
	}
}

func (ds *GlobalStatus) Run() error {
	rows := make([]globalStatusLine, 0)
	err := ds.connection.Select(&rows, "SHOW GLOBAL STATUS")
	if err != nil {
		return errors.Wrap(err, "global status query failed")
	}

	values := map[string]interface{}{}

	var e error
	for _, v := range rows {
		if !ds.Filter.MatchString(v.Name) {
			i, err := strconv.Atoi(v.Value)
			if err != nil {
				e = multierror.Append(e, err)
				continue
			}
			values[v.Name] = i
		}
	}

	ds.sink <- NewDataPoint("mysql", nil, values)
	return e
}
