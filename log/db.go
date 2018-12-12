package log

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// constants
const (
	BufSize = 8192
)

type insertFunc func(db *gorm.DB, tablename string, subtable bool, entry *Logs) error

// DBHook ..
type DBHook struct {
	db      *gorm.DB
	buf     chan *Logs
	datamap map[int]*DBInfo
}

// DBInfo ..
type DBInfo struct {
	Tablename string
	Subtable  bool
	Function  insertFunc
	Blacklist map[string]bool
}

var monthmap = make(map[string]string)

// NewDBHook ...
func NewDBHook(db *gorm.DB, data map[int]*DBInfo) *DBHook {
	hook := &DBHook{
		db:      db,
		datamap: data,
		buf:     make(chan *Logs, BufSize),
	}

	go hook.fire()
	return hook
}

func (hook *DBHook) fire() {
	for {
		for {
			select {
			case entry := <-hook.buf:

				info, _ := hook.datamap[entry.Type]
				if err := info.Function(hook.db, info.Tablename, info.Subtable, entry); err != nil {
					fmt.Fprintf(os.Stderr, "[logrus] Can't insert entry (%v): %v\n", entry, err)
				}
			}
		}
	}
}

func (hook *DBHook) newEntry(entry *logrus.Entry, info *DBInfo, logtype int) *Logs {
	// Don't modify entry.Data directly, as the entry will used after this hook was fired
	data := map[string]interface{}{}

	for k, v := range entry.Data {
		if !info.Blacklist[k] {
			data[k] = v
			if k == logrus.ErrorKey {
				asError, isError := v.(error)
				_, isMarshaler := v.(json.Marshaler)
				if isError && !isMarshaler {
					data[k] = asError.Error()
				}
			}
		}
	}

	return &Logs{
		Type:    logtype,
		Raw:     data,
		Time:    entry.Time,
		Level:   int(entry.Level),
		Message: entry.Message,
	}
}

// Levels ...
func (hook DBHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire ...
func (hook DBHook) Fire(entry *logrus.Entry) error {
	// table filter
	var (
		logtype int
	)

	if value, ok := entry.Data["LogType"]; ok {
		logtype, _ = value.(int)
	}

	// get db info
	info, ok := hook.datamap[logtype]
	if !ok {
		return nil
	}

	select {
	case hook.buf <- hook.newEntry(entry, info, logtype):
	default:
		return errors.New("[logrus] chan full")
	}

	return nil
}
