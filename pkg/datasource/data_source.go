package datasource

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type (
	BaseDataSource struct {
		sink       chan *DataPoint
		connection *sqlx.DB
	}
	DataSource interface {
		C() <-chan *DataPoint
		Close()
		Run() error
	}
)

func New(db *sqlx.DB) *BaseDataSource {
	return &BaseDataSource{
		connection: db,
		sink:       make(chan *DataPoint),
	}
}

func (ds *BaseDataSource) C() <-chan *DataPoint {
	return ds.sink
}

func (ds *BaseDataSource) Close() {
	close(ds.sink)
}
