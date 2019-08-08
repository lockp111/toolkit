package quorm

import (
	"github.com/doug-martin/goqu/v8"
	"github.com/doug-martin/goqu/v8/exp"
)

// A ...
func A(table string, schema ...string) *Alias {
	c := &Alias{
		goqu.T(table),
	}

	if len(schema) != 0 {
		c.IdentifierExpression = c.Schema(schema[0])
	}

	return c
}

// Alias ...
type Alias struct {
	exp.IdentifierExpression
}

// I ...
func (a *Alias) I(col ...string) exp.IdentifierExpression {
	if len(col) != 0 {
		return a.Col(col[0])
	}

	return a.IdentifierExpression
}

// C ...
func (a *Alias) C(col string) string {
	if len(a.GetTable()) == 0 {
		return col
	}

	prefix := a.GetTable()
	if len(a.GetSchema()) != 0 {
		prefix = a.GetSchema() + "." + prefix
	}

	return prefix + "." + col
}
