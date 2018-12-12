package env

import (
	"os"
	"path/filepath"
	"time"

	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/lockp111/toolkit/utils"
	log "github.com/sirupsen/logrus"
)

// Const
var (
	Dir       string
	RunDir    string
	LogDir    string
	LogPath   string
	Pid       int
	Hostname  string
	ServiceID string
	LogLevel  log.Level
)

func init() {
	file, _ := filepath.Abs(os.Args[0])
	dir := filepath.Dir(file)

	Dir = filepath.Dir(dir + "..")

	LogDir = Dir + "/log/"
	if !isPathExist(LogDir) {
		if err := os.MkdirAll(LogDir, os.ModePerm); err != nil {
			utils.ErrExit(err)
		}
	}

	LogPath = LogDir + filepath.Base(os.Args[0]) + ".log"
	RunDir, _ := os.Getwd()

	RunDir, _ = filepath.Abs(RunDir)
	LogDir, _ = filepath.Abs(LogDir)
	LogPath, _ = filepath.Abs(LogPath)

	Pid = os.Getpid()

	hostname, err := os.Hostname()
	utils.ErrExit(err)

	logLevel := os.Getenv("LOGLEVEL")
	log.SetOutput(os.Stdout)
	if logLevel != "" {
		LogLevel, err = log.ParseLevel(logLevel)
		utils.ErrExit(err)

		log.SetLevel(LogLevel)

		if LogLevel != log.DebugLevel {
			w, err := rotatelogs.New(LogPath+".%Y%m%d",
				rotatelogs.WithRotationCount(30),
				rotatelogs.WithRotationTime(24*time.Hour),
			)
			if err != nil {
				utils.ErrExit(err)
			}
			log.SetOutput(w)
		}
	}

	ServiceID = os.Getenv("SERVICEID")
	Hostname = hostname
}

// isPathExist ...
func isPathExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			return true
		}
	}
	return false
}
