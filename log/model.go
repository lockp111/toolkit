package log

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
)

// Logs ...
type Logs struct {
	Level   int
	Message string
	Type    int
	SubType int
	Raw     map[string]interface{}
	Time    time.Time
}

// FilterStruct ..
func FilterStruct(raw map[string]interface{}, log interface{}) error {
	// get fiter str
	last := make(map[string]interface{})
	for key, value := range raw {
		last[key] = value
	}

	for key := range structs.Map(log) {
		delete(last, key)
	}

	jsonstr, err := json.Marshal(last)
	if err != nil {
		return err
	}

	// set log_data
	raw["Data"] = string(jsonstr)
	jsonstr, err = json.Marshal(raw)
	if err != nil {
		return err
	}

	// jsonstr to struct
	err = json.Unmarshal(jsonstr, log)
	if err != nil {
		return err
	}

	return nil
}

// ReplaceKey ..
func ReplaceKey(raw interface{}, key, replace string) (map[string]interface{}, error) {
	jsonstr, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	// jsonstr to struct
	err = json.Unmarshal(jsonstr, &data)
	if err != nil {
		return nil, err
	}

	if _, ok := data[key]; ok {
		data[key] = replace
	}
	return data, nil
}

// NormalLog ..
type NormalLog struct {
	ID        int `json:"-" gorm:"primary_key"`
	Level     int
	Type      int `structs:"LogType"`
	LogIndex  string
	SubType   int `structs:"LogSubType"`
	Message   string
	CreatedAt time.Time
	Data      string
}

// save ..
func save(db *gorm.DB, tablename string, subtable bool, entry *Logs, log interface{}) error {

	err := FilterStruct(entry.Raw, log)
	if err != nil {
		return err
	}

	tablestr := tablename
	if subtable {
		tablestr = tablestr + "_" + time.Now().Format("200601")
	}

	dbmonth, ok := monthmap[tablename]
	if !ok || tablestr != dbmonth {

		if db.HasTable(tablestr) == false {
			if err := db.Table(tablestr).CreateTable(log).Error; err != nil {
				return errors.New("[logrus] Can't create table: " + tablestr + " error :" + err.Error())
			}
		}
		monthmap[tablename] = tablestr
	}

	return db.Table(tablestr).Create(log).Error
}

// InsertNormalLogsFunc ..
func InsertNormalLogsFunc(db *gorm.DB, tablename string, subtable bool, entry *Logs) error {
	log := &NormalLog{
		Level:     entry.Level,
		Type:      entry.Type,
		SubType:   entry.SubType,
		Message:   entry.Message,
		CreatedAt: entry.Time,
	}

	return save(db, tablename, subtable, entry, log)
}

// APILog ..
type APILog struct {
	ID              int `json:"-" gorm:"primary_key"`
	Level           int
	Type            int `structs:"LogType"`
	SubType         int `structs:"LogSubType"`
	Message         string
	CreatedAt       time.Time
	Data            string
	API             string
	Method          string
	ElapsedTime     float64
	Status          int
	UserID          uint64
	UserName        string
	SourceIP        string
	DestIP          string
	Country         string
	RequestContent  string
	ResponseContent string
}

// InsertAPILogsFunc ..
func InsertAPILogsFunc(db *gorm.DB, tablename string, subtable bool, entry *Logs) error {
	log := &APILog{
		Level:     entry.Level,
		Type:      entry.Type,
		Message:   entry.Message,
		CreatedAt: entry.Time,
	}

	return save(db, tablename, subtable, entry, log)
}

// OperLog ..
type OperLog struct {
	ID        int
	Level     int
	Type      int `structs:"LogType"`
	SubType   int `structs:"LogSubType"`
	Message   string
	CreatedAt time.Time
	Data      string
	UserID    string
	UserName  string
	Status    int
	SourceIP  string
}

// InsertOperationLogsFunc ..
func InsertOperationLogsFunc(db *gorm.DB, tablename string, subtable bool, entry *Logs) error {
	log := &OperLog{
		Level:     entry.Level,
		Type:      entry.Type,
		Message:   entry.Message,
		CreatedAt: entry.Time,
	}

	return save(db, tablename, subtable, entry, log)
}
