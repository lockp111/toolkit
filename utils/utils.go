package utils

import (
	"encoding/json"
	"io/ioutil"
	"math"

	log "github.com/sirupsen/logrus"
)

// ErrExit ...
func ErrExit(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// FixedRound ...
func FixedRound(price float64, digits int) float64 {
	exp := 0.5
	if price < 0 {
		exp = -exp
	}
	pip := math.Pow10(digits)

	np := price*pip + exp
	tf := math.Trunc(np)
	return tf / pip
}

// LoadJSONFile ..
func LoadJSONFile(filename string, v interface{}) error {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}
