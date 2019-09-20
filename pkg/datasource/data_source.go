package datasource

import (
	"github.com/jmoiron/sqlx"
)

type (
	DataSource struct {
		sink       chan *DataPoint
		connection *sqlx.DB
	}
)

func New(db *sqlx.DB) *DataSource {
	return &DataSource{
		connection: db,
		sink:       make(chan *DataPoint),
	}
}

func (ds *DataSource) C() <-chan *DataPoint {
	return ds.sink
}

func (ds *DataSource) Close() {
	close(ds.sink)
}
