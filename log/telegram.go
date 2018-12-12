package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// TelegramHook ..
type TelegramHook struct {
	Blacklist map[string]bool
	URL       string
	ChatID    int32
	ON        bool
}

var telegram = &TelegramHook{}

// NewTelegramHook ...
func NewTelegramHook(url string, chatid int32, black map[string]bool) *TelegramHook {
	telegram = &TelegramHook{
		URL:       url,
		Blacklist: black,
		ChatID:    chatid,
		ON:        false,
	}

	return telegram
}

// Levels ...
func (hook *TelegramHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// SetTelegramSwitch ...
func SetTelegramSwitch(flag bool) {
	telegram.ON = flag
}

// GetTelegramSwitch ...
func GetTelegramSwitch() bool {
	return telegram.ON
}

// Fire ...
func (hook *TelegramHook) Fire(entry *logrus.Entry) error {
	if !hook.ON {
		return nil
	}

	if entry.Level != logrus.ErrorLevel {
		if _, ok := entry.Data["telegram"]; !ok {
			return nil
		}
	}

	go hook.Alert(entry)
	return nil
}

// Alert ...
func (hook *TelegramHook) Alert(entry *logrus.Entry) {
	type Data struct {
		ChatID    int32  `json:"chat_id"`
		ParseMode string `json:"parse_mode"`
		Text      string `json:"text"`
	}

	var data Data
	data.ParseMode = "Markdown"
	data.ChatID = hook.ChatID

	filterlog := hook.Filter(entry.Data)
	filterlog["log_message"] = entry.Message
	data.Text = "----------------\n"
	for key, value := range filterlog {
		data.Text += fmt.Sprintf("***%v***: %v\n", key, value)
	}
	data.Text += "----------------"

	bs, _ := json.Marshal(&data)
	rsp, err := http.Post(hook.URL, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[logrus] Can't Post telegram (%v): %v, rsp[%v]\n", data.Text, err, rsp)
	}
}

// Filter ...
func (hook *TelegramHook) Filter(files logrus.Fields) logrus.Fields {
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
