package consul

import (
	"testing"
)

func init() {
	InitSource(GetAddress(), "test")
}

func TestGet(t *testing.T) {
	m := make(map[string]interface{})
	ConfigGet(&m)

	t.Log(m)
	t.Fail()
}
