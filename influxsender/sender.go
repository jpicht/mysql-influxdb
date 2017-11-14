package influxsender

import (
	"sync"
	"time"

	"github.com/ftloc/exception"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/mitchellh/packer/common/json"
	"github.com/pkg/errors"
)

type (
	Config struct {
		Interval    time.Duration
		PackageSize int
		Concurrency int
		Url         string
		Database    string
	}
	InfluxSender interface {
		AddPoint(point *influx.Point)
		Send()
		Stop()
	}
	influxSender struct {
		config     Config
		lock       sync.Mutex
		data       influx.BatchPoints
		sendWorker *sync.WaitGroup
		c          chan influx.BatchPoints
		stop       chan bool
	}
)

func (i *Config) UnmarshalJSON(data []byte) error {
	var temp struct {
		Interval string
		Url      string
		Database string
	}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	i.Url = temp.Url
	i.Database = temp.Database
	if len(temp.Interval) > 0 {
		tmp_duration, err := time.ParseDuration(temp.Interval)
		if err != nil {
			return err
		}
		i.Interval = tmp_duration
	}
	return nil
}

func (i *Config) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}

func NewSender(c *Config) (InfluxSender, error) {
	batchPoints, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: c.Database,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Cannot create influx batch")
	}

	sender, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:    c.Url,
		Timeout: 5 * time.Second,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Cannot create influxdb output")
	}

	sender.Query(influx.NewQuery(
		"CREATE DATABASE "+c.Database,
		"",
		"",
	))

	if c.PackageSize == 0 {
		c.PackageSize = 10000
	}
	if c.Concurrency == 0 {
		c.Concurrency = 4
	}
	if c.Interval == 0 {
		c.Interval = 100 * time.Millisecond
	}

	s := &influxSender{
		config:     *c,
		lock:       sync.Mutex{},
		data:       batchPoints,
		sendWorker: &sync.WaitGroup{},
		c:          make(chan influx.BatchPoints, c.Concurrency),
	}
	s.sendWorker.Add(c.Concurrency)
	for i := 0; i < c.Concurrency; i++ {
		go s.worker()
	}

	// start timed flusher
	s.stop = make(chan bool)
	go func() {
		t := time.NewTicker(100 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				s.Send()
			case <-s.stop:
				return
			}
		}
	}()

	return s, nil
}

func (s *influxSender) AddPoint(point *influx.Point) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data.AddPoint(point)

	// max batch size exceeded
	if len(s.data.Points()) >= s.config.PackageSize {
		s.send(false)
	}
}

func (s *influxSender) Send() {
	s.send(true)
}

func (s *influxSender) send(lock bool) {
	if lock {
		s.lock.Lock()
	}
	if len(s.data.Points()) == 0 {
		if lock {
			s.lock.Unlock()
		}
		return
	}

	batchPoints, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: s.data.Database(),
	})
	exception.ThrowOnError(err, errors.Wrap(err, "Cannot create influx batch"))

	data := s.data
	s.data = batchPoints
	if lock {
		s.lock.Unlock()
	}

	s.c <- data
}

func (s *influxSender) Stop() {
	close(s.c)
	s.stop <- true
	s.Send()
	s.sendWorker.Wait()
}

func (s *influxSender) worker() {
	defer s.sendWorker.Done()
	sender, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:    s.config.Url,
		Timeout: 5 * time.Second,
	})

	if err != nil {
		panic(errors.Wrap(err, "Cannot create influxdb output"))
	}

	for batch := range s.c {
		sender.Write(batch)
	}
}
