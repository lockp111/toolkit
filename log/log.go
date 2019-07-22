package log

import (
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/lockp111/toolkit/utils"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Config ...
type Config struct {
	Path           string
	Level          string
	DBUrl          string
	Logstash       string
	DingDingURL    string
	TelegramURL    string
	TelegramChatID int32
	FileContext    bool
}

// Fields ...
type Fields map[string]interface{}

// LogTable
const (
	LogTableNone = iota
	LogTableOper
	LogTableAPI
	LogTableNormal
)

var (
	// std is the name of the standard logger in stdlib `log`
	std = logrus.StandardLogger()
)

// SettingLog ...
func SettingLog(config *Config) error {
	if config.Level != "" {
		lvl, err := logrus.ParseLevel(config.Level)
		if err != nil {
			return err
		}
		logrus.SetLevel(lvl)
	}

	if config.Path != "" {
		dir := getDir(config.Path)
		if isPathNotExist(dir) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				utils.ErrExit(err)
			}
		}

		w, err := rotatelogs.New(config.Path+"-%Y%m%d.log",
			rotatelogs.WithRotationCount(30),
			rotatelogs.WithRotationTime(24*time.Hour),
		)

		if err != nil {
			return err
		}

		logrus.SetOutput(os.Stdout)
		logrus.AddHook(lfshook.NewHook(lfshook.WriterMap{
			logrus.DebugLevel: w, // 为不同级别设置不同的输出目的
			logrus.InfoLevel:  w,
			logrus.WarnLevel:  w,
			logrus.ErrorLevel: w,
			logrus.FatalLevel: w,
			logrus.PanicLevel: w,
		},
			&logrus.TextFormatter{},
		))
	}

	blacklist := map[string]bool{}
	blacklist["dingding"] = true
	blacklist["telegram"] = true

	if config.FileContext {
		logrus.AddHook(ContextHook{})

		blacklist["file"] = true
		blacklist["line"] = true
	}

	if config.DingDingURL != "" {
		dingblack := make(map[string]bool)
		for key, value := range blacklist {
			dingblack[key] = value
		}

		dingblack["LogType"] = true
		dingblack["LogSubtype"] = true
		dingblack["LogIndex"] = true
		logrus.AddHook(NewDingDingHook(config.DingDingURL, dingblack))
	}

	if config.TelegramURL != "" {
		telegramblack := make(map[string]bool)
		for key, value := range blacklist {
			telegramblack[key] = value
		}

		telegramblack["LogType"] = true
		telegramblack["LogSubtype"] = true
		telegramblack["LogIndex"] = true
		logrus.AddHook(NewTelegramHook(config.TelegramURL, config.TelegramChatID, telegramblack))
	}

	if config.DBUrl != "" {
		gormDB, err := gorm.Open("mysql", config.DBUrl)
		if err != nil {
			return err
		}

		operlog := &DBInfo{
			Tablename: "oper_log",
			//Subtable:  true,
			Function:  InsertOperationLogsFunc,
			Blacklist: blacklist,
		}

		apilog := &DBInfo{
			Tablename: "api_log",
			//Subtable:  true,
			Function:  InsertAPILogsFunc,
			Blacklist: blacklist,
		}

		normallog := &DBInfo{
			Tablename: "normal_log",
			//Subtable:  true,
			Function:  InsertNormalLogsFunc,
			Blacklist: blacklist,
		}

		datamap := map[int]*DBInfo{
			LogTableOper:   operlog,
			LogTableAPI:    apilog,
			LogTableNormal: normallog,
		}

		logrus.AddHook(NewDBHook(gormDB, datamap))
	}

	std = logrus.StandardLogger()
	return nil
}

func isPathNotExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return true
		}
	}
	return false
}

func getDir(path string) string {
	paths := strings.Split(path, "/")
	return strings.Join(
		paths[:len(paths)-1],
		"/",
	)
}
