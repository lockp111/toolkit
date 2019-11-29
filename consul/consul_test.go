package consul

import (
	"testing"
)

func TestBytes(t *testing.T) {
	var bs = make([]byte, 0)

	if len(bs) != 0 {
		t.Error("bs len not zero")
		// return
	}

	t.Log(bs)

	s := string(bs)
	if len(s) != 0 {
		t.Error("s len not zero")
	}

	t.Log(s)
}

func TestMap(t *testing.T) {
	var m = make(map[string]int)
	m["1"] = 1
	m["2"] = 2
	m["3"] = 3
	m["4"] = 4

	for k := range m {
		delete(m, k)
	}

	t.Log(m)
	t.Fail()
}
