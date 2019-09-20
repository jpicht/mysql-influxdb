package datasource

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type (
	InnoDB struct {
		BaseDataSource
		outputValues map[string]interface{}
	}
)

func NewInnoDB(db *sqlx.DB) *InnoDB {
	return &InnoDB{
		BaseDataSource: *New(db),
		outputValues:   make(map[string]interface{}),
	}
}

var (
	reInnoDbTransactionsTableStatus = regexp.MustCompile("^mysql tables in use ([0-9]+), locked ([0-9]+)")
	reInnoDbTransactionsMySQL       = regexp.MustCompile("^MySQL thread id ([0-9]+), OS thread handle (0x[0-9a-fA-F]+|[0-9]+), query id ([0-9]+) +([^ \t]+) ([a-zA-Z0-9_]+) ?(.*)$")
	reInnoDbTransactionsStatus      = regexp.MustCompile("([0-9]+) lock struct.s., heap size ([0-9]+), ([0-9]+) row lock.s.(, undo log entries ([0-9]+))?")

	reInnoDbGlobalStatusItem = regexp.MustCompile("^([A-Za-z ]+[a-z]) +([0-9]+)$")
)

func (i *InnoDB) Run() error {
	type InnoStatus struct {
		Type   string `db:"Type"`
		Name   string `db:"Name"`
		Status string `db:"Status"`
	}
	data := &InnoStatus{}

	if err := i.connection.Get(data, "SHOW ENGINE INNODB STATUS;"); err != nil {
		return errors.Wrap(err, "cannot get innodb status")
	}

	lines := strings.Split(data.Status, "\n")
	block := ""
	blocks := make(map[string][]string)
	for i, line := range lines {
		if i < 2 {
			continue
		}
		if line == lines[i-2] && line[0] == '-' {
			if block != "" {
				if len(blocks[block]) > 2 {
					blocks[block] = blocks[block][0 : len(blocks[block])-2]
				} else {
					delete(blocks, block)
				}
			}
			block = lines[i-1]
			blocks[block] = make([]string, 0)
			continue
		}
		if block == "" {
			continue
		}
		blocks[block] = append(blocks[block], line)
	}

	i.transactions(blocks["TRANSACTIONS"])
	i.bufferpool(blocks["BUFFER POOL AND MEMORY"], blocks["INDIVIDUAL BUFFER POOL INFO"])

	i.sink <- NewDataPoint("innodb", nil, i.outputValues)

	return nil
}

type transaction struct {
	Active      bool
	Time        int
	Host        string
	User        string
	Using       int
	Locked      int
	Status      string
	LockStructs int
	HeapSize    int
	RowLocks    int
	UndoEntries int
}

func (t *transaction) dataPoint() *DataPoint {
	return NewDataPoint(
		"transactions",
		map[string]string{
			"client_host": t.Host,
			"user":        t.User,
			"status":      t.Status,
			"is_active":   fmt.Sprintf("%#v", t.Active),
		},
		map[string]interface{}{
			"active":        t.Time,
			"tables_used":   t.Using,
			"tables_locked": t.Locked,
		},
	)
}

func (i *InnoDB) transactions(lines []string) {
	transactionsTotal := 0
	transactionsActive := 0
	var t *transaction
	t = nil

	userTransactions := map[string]int{}
	activeUserTransactions := map[string]int{}

	for _, line := range lines {
		if len(line) > 14 && line[0:14] == "---TRANSACTION" {
			if t != nil {
				i.sink <- t.dataPoint()
			}
			transactionsTotal++

			t = &transaction{
				Active: strings.Contains(line, " ACTIVE "),
			}
			if !t.Active {
				continue
			}
			transactionsActive++
			secsStrs := strings.Split(strings.Split(line, " ACTIVE ")[1], " ")
			for i, s := range secsStrs[1:] {
				if len(s) >= 3 && strings.ToLower(s[0:3]) == "sec" {
					t.Time, _ = strconv.Atoi(secsStrs[i])
					break
				}
			}
		}

		log := logrus.WithField("line", line)

		if strings.HasPrefix(line, "mysql tables in use") {
			matches := reInnoDbTransactionsTableStatus.FindStringSubmatch(line)
			if matches == nil {
				log.WithField("re", reInnoDbTransactionsTableStatus).Fatal()
			}
			t.Using, _ = strconv.Atoi(matches[1])
			t.Locked, _ = strconv.Atoi(matches[2])
		} else if strings.HasPrefix(line, "MySQL thread id") {
			matches := reInnoDbTransactionsMySQL.FindStringSubmatch(line)
			if matches == nil {
				if strings.Contains(line, "relay log") {
					continue
				}
				log.WithField("re", reInnoDbTransactionsMySQL).Fatal()
			}
			t.Host = matches[4]
			t.User = matches[5]
			t.Status = matches[6]

			userTransactions[t.User]++
			if t.Active {
				activeUserTransactions[t.User]++
			}
		} else if strings.Contains(line, " lock struct") {
			matches := reInnoDbTransactionsStatus.FindStringSubmatch(line)

			if matches == nil {
				log.WithField("re", reInnoDbTransactionsStatus).Fatal("could not parse line")
				continue
			}

			// 1 lock struct(s), heap size 1136, 0 row lock(s), undo log entries 2
			t.LockStructs, _ = strconv.Atoi(matches[1])
			t.HeapSize, _ = strconv.Atoi(matches[2])
			t.RowLocks, _ = strconv.Atoi(matches[3])
			if len(matches) > 4 && len(matches[4]) > 0 {
				t.UndoEntries, _ = strconv.Atoi(matches[5])
			}
		}

	}

	for user, numTransactions := range userTransactions {
		i.sink <- &DataPoint{
			"user_transactions",
			map[string]string{
				"user": user,
			},
			map[string]interface{}{
				"transactions_count":  numTransactions,
				"transactions_active": activeUserTransactions[user],
			},
		}
	}

	i.outputValues["transactions_active"] = transactionsActive
	i.outputValues["transactions_total"] = transactionsTotal
}

func (i *InnoDB) bufferpool(global []string, indiv []string) {
	pools := make(map[int][]string)
	cur := -1
	for _, line := range indiv {
		if strings.HasPrefix(line, "---BUFFER POOL ") {
			cur, _ = strconv.Atoi(line[15:])
			pools[cur] = make([]string, 0)
		}
		if cur == -1 {
			continue
		}
		pools[cur] = append(pools[cur], line)
	}

	for num, lines := range pools {
		pooldata := make(map[string]interface{})
		for _, line := range lines {
			if matches := reInnoDbGlobalStatusItem.FindStringSubmatch(line); matches != nil && len(matches) == 3 {
				pooldata[strings.Replace(strings.ToLower(matches[1]), " ", "_", -1)], _ = strconv.Atoi(matches[2])
			}
		}
		i.sink <- NewDataPoint(
			"innodb_pools",
			map[string]string{
				"pool": strconv.Itoa(num),
			},
			pooldata,
		)
	}
}
