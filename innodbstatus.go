package main

import (
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/ftloc/exception"
	"github.com/jpicht/logger"
)

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

var (
	reInnoDbTransactionsTableStatus = regexp.MustCompile("^mysql tables in use ([0-9]+), locked ([0-9]+)")
	reInnoDbTransactionsMySQL       = regexp.MustCompile("^MySQL thread id ([0-9]+), OS thread handle (0x[0-9a-fA-F]+|[0-9]+), query id ([0-9]+) ([^ \t]+) ([a-zA-Z0-9_]+) ?(.*)$")
	reInnoDbTransactionsStatus      = regexp.MustCompile("([0-9]+) lock struct.s., heap size ([0-9]+), ([0-9]+) row lock.s.(, undo log entries ([0-9]+))?")

	//reInnoDbGlobalStatusItem = regexp.MustCompile("(^|, )([A-Za-z_ ]+[a-z]):? +([0-9]+)(,|$)")
	reInnoDbGlobalStatusItem = regexp.MustCompile("^([A-Za-z ]+[a-z]) +([0-9]+)$")
)

func (t *transaction) send(a *app, log logger.Logger) {
	a.send(log, "transactions", tags(
		a.currentHost.Tags,
		map[string]string{
			"client_host": t.Host,
			"user":        t.User,
			"status":      t.Status,
		},
	), map[string]interface{}{
		"active":        t.Time,
		"tables_used":   t.Using,
		"tables_locked": t.Locked,
	})
}

func MustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	exception.ThrowOnError(err, err)
	return i
}

type outputvalues map[string]interface{}

func (ov outputvalues) merge(other outputvalues) {
	if other == nil {
		return
	}
	for k, v := range other {
		ov[k] = v
	}
}

func (ov outputvalues) get() map[string]interface{} {
	return ov
}

func (a *app) innoStatusTransactions(lines []string, log logger.Logger, failed *int32) outputvalues {
	transactionsTotal := 0
	transactionsActive := 0
	var t *transaction
	t = nil

	for _, line := range lines {
		exception.Try(func() {
			if len(line) > 14 && line[0:14] == "---TRANSACTION" {
				if t != nil {
					t.send(a, log)
				}
				transactionsTotal++

				t = &transaction{
					Active: strings.Contains(line, " ACTIVE "),
				}
				if t.Active {
					transactionsActive++
					secsStrs := strings.Split(strings.Split(line, " ACTIVE ")[1], " ")
					for i, s := range secsStrs[1:] {
						if len(s) >= 3 && strings.ToLower(s[0:3]) == "sec" {
							t.Time = MustAtoi(secsStrs[i])
							break
						}
					}
					/*if secsStrs[0] == "(PREPARED)" {
						t.Time = MustAtoi(secsStrs[0])
					} else {
						t.Time = MustAtoi(secsStrs[1])
					}*/
				}
			}
			if matches := reInnoDbTransactionsTableStatus.FindStringSubmatch(line); matches != nil {
				t.Using = MustAtoi(matches[1])
				t.Locked = MustAtoi(matches[2])
			} else if matches := reInnoDbTransactionsMySQL.FindStringSubmatch(line); matches != nil {
				t.Host = matches[4]
				t.User = matches[5]
				if t.User != "root" {
					log.Alertf("%#v", matches)
				}
				t.Status = matches[6]
			} else if strings.Contains(line, " lock struct") {
				matches := reInnoDbTransactionsStatus.FindStringSubmatch(line)
				if matches == nil {
					log.Alertf("Could not parse: '%s'", line)
					return
				}
				// 1 lock struct(s), heap size 1136, 0 row lock(s), undo log entries 2
				t.LockStructs = MustAtoi(matches[1])
				t.HeapSize = MustAtoi(matches[2])
				t.RowLocks = MustAtoi(matches[3])
				if len(matches) > 4 && len(matches[4]) > 0 {
					t.UndoEntries = MustAtoi(matches[5])
				}
			}
		}).CatchAll(func(err interface{}) {
			log.Warningf("Error parsing line '%s': %#v", line, err)
		}).Go()
	}
	if t != nil {
		t.send(a, log)
	}

	return outputvalues{
		"transactions_active": transactionsActive,
		"transactions_total":  transactionsTotal,
	}
}

func (a *app) innoStatusBufferpool(global []string, indiv []string, log logger.Logger, failed *int32) outputvalues {
	data := make(outputvalues)
	if false {
		for _, line := range global {
			if matches := reInnoDbGlobalStatusItem.FindStringSubmatch(line); matches != nil && len(matches) == 3 {
				data["bufferpool_"+strings.Replace(strings.ToLower(matches[1]), " ", "_", -1)] = MustAtoi(matches[2])
			}
			/*
				matches := reInnoDbGlobalStatusItem.FindAllStringSubmatch(line, -1)
				if len(matches) > 0 {
					log.Infof("LINE: %s", line)
					for _, match := range matches {

						//data["bufferpool_"+strings.Replace(strings.ToLower(match[2]), " ", "_", -1)] = MustAtoi(match[3])
						log.Infof("%s: %d", "bufferpool_"+strings.Replace(strings.ToLower(match[2]), " ", "_", -1), MustAtoi(match[3]))
					}
				} else {
					log.Infof("LINE: ERR %s", line)
				}
			*/
		}
	}
	pools := make(map[int][]string)
	cur := -1
	for _, line := range indiv {
		if startsWith(line, "---BUFFER POOL ") {
			cur = MustAtoi(line[15:])
			pools[cur] = make([]string, 0)
		}
		if cur == -1 {
			continue
		}
		pools[cur] = append(pools[cur], line)
	}

	for num, lines := range pools {
		pooldata := make(outputvalues)
		for _, line := range lines {
			if matches := reInnoDbGlobalStatusItem.FindStringSubmatch(line); matches != nil && len(matches) == 3 {
				pooldata[strings.Replace(strings.ToLower(matches[1]), " ", "_", -1)] = MustAtoi(matches[2])
			}
		}
		a.send(log, "innodb_pools", tags(a.currentHost.Tags, map[string]string{
			"pool": strconv.Itoa(num),
		}), pooldata.get())
	}

	return data
}

func startsWith(haystack, needle string) bool {
	if len(needle) > len(haystack) {
		return false
	}
	return haystack[0:len(needle)] == needle
}

func (a *app) innoStatus(log logger.Logger, failed *int32) {
	defer a.wg.Done()
	exception.Try(func() {
		type InnoStatus struct {
			Type   string `db:"Type"`
			Name   string `db:"Name"`
			Status string `db:"Status"`
		}
		data := &InnoStatus{}
		err := a.currentHost.Connection.Get(data, "SHOW ENGINE INNODB STATUS;")
		if err != nil {
			log.WithData("error", err).Warning("MySQL-Query-Error")
			atomic.AddInt32(failed, 1)
			return
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

		values := make(outputvalues)
		values.merge(a.innoStatusTransactions(blocks["TRANSACTIONS"], log, failed))
		values.merge(a.innoStatusBufferpool(blocks["BUFFER POOL AND MEMORY"], blocks["INDIVIDUAL BUFFER POOL INFO"], log, failed))
		a.send(log, "innodb", a.currentHost.Tags, values.get())
	}).CatchAll(func(i interface{}) {
		log.Alertf("Exception caught: %#v", i)
		ok, file, line := exception.GetThrower()
		if ok {
			log.Alertf("Location: %s:%d", file, line)
		} else {
			cs := make([]uintptr, 20)
			amount := runtime.Callers(1, cs)

			for i := 0; i < amount; i++ {
				f := runtime.FuncForPC(cs[i])
				file, line := f.FileLine(cs[i])
				log.Alert("  " + f.Name())
				log.Alertf("      %s:%d", file, line)
			}
		}
		atomic.AddInt32(failed, 1)
	}).Go()
}
