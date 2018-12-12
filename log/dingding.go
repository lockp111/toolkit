package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// DingDingHook ..
type DingDingHook struct {
	Blacklist map[string]bool
	URL       string
	ON        bool
}

var dingding = &DingDingHook{}

// NewDingDingHook ...
func NewDingDingHook(url string, black map[string]bool) *DingDingHook {
	dingding = &DingDingHook{
		URL:       url,
		Blacklist: black,
		ON:        false,
	}

	return dingding
}

// Levels ...
func (hook *DingDingHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// SetDingDingSwitch ...
func SetDingDingSwitch(flag bool) {
	dingding.ON = flag
}

// GetDingDingSwitch ...
func GetDingDingSwitch() bool {
	return dingding.ON
}

// Fire ...
func (hook *DingDingHook) Fire(entry *logrus.Entry) error {
	if !hook.ON {
		return nil
	}

	if entry.Level != logrus.ErrorLevel {
		if _, ok := entry.Data["dingding"]; !ok {
			return nil
		}
	}

	go hook.Alert(entry)
	return nil
}

// Alert ...
func (hook *DingDingHook) Alert(entry *logrus.Entry) {
	type Data struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}

	//entry.Data["log_message"] = entry.Message
	filterlog := hook.Filter(entry.Data)
	filterlog["log_message"] = entry.Message
	jsonstr, err := json.Marshal(filterlog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[logrus] Can't Marshal dingding (%v): %v\n", string(jsonstr), err)
		return
	}

	var data Data
	data.MsgType = "text"
	data.Text.Content = string(jsonstr)

	bs, _ := json.Marshal(&data)
	_, err = http.Post(hook.URL, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[logrus] Can't Post dingding (%v): %v\n", string(jsonstr), err)
	}
}

// Filter ...
func (hook *DingDingHook) Filter(files logrus.Fields) logrus.Fields {
	data := make(logrus.Fields)

	for k, v := range files {
		if !hook.Blacklist[k] {
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
	return data
}
