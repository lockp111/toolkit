package utils

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	letterBytes   = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

// RandomBytes ...
func RandomBytes(n int, seed ...string) []byte {
	var (
		seedBytes = letterBytes
	)

	if len(seed) != 0 {
		seedBytes = seed[0]
	}

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(seedBytes) {
			b[i] = seedBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}

// RandomString ...
func RandomString(n int, seed ...string) string {
	return string(RandomBytes(n, seed...))
}
