package datasource

import (
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

type (
	DataPoint struct {
		name   string
		tags   map[string]string
		values map[string]interface{}
	}
)

func NewDataPoint(n string, t map[string]string, v map[string]interface{}) *DataPoint {
	if t == nil {
		t = make(map[string]string)
	}
	if v == nil {
		v = make(map[string]interface{})
	}
	return &DataPoint{n, t, v}
}

func (dp *DataPoint) AddTagsIfNotExist(t map[string]string) {
	for k, v := range t {
		if _, ok := dp.tags[k]; !ok {
			dp.tags[k] = v
		}
	}
}

func (dp *DataPoint) InfluxPoint(t ...time.Time) (*influx.Point, error) {
	return influx.NewPoint(
		dp.name,
		dp.tags,
		dp.values,
		t...,
	)
}
