package utils

import (
	"testing"
)

func TestFixedRound(t *testing.T) {
	a := 1.32123
	b := FixedRound(a, 2)
	t.Log(b)
	t.Fail()
}

func TestRandomString(t *testing.T) {
	s := RandomString(32)
	t.Log(s)
	t.Fail()
}
