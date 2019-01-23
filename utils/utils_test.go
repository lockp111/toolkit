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
