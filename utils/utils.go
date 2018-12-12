package utils

import (
	"log"
	"math"
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
