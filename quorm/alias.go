package quorm

import (
	"github.com/doug-martin/goqu/v8"
	"github.com/doug-martin/goqu/v8/exp"
)

// A ...
func A(table, alias string, schema ...string) *Alias {
	t := goqu.T(table)
	if len(schema) != 0 {
		t = t.Schema(schema[0])
	}

	c := &Alias{
		t.As(alias),
		goqu.T(alias),
	}

	return c
}

// Alias ...
type Alias struct {
	table exp.AliasedExpression
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

	return a.GetTable() + "." + col
}

// Table ...
func (a *Alias) Table() exp.AliasedExpression {
	return a.table
}
