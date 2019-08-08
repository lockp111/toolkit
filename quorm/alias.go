package quorm

import (
	"sync"

	"github.com/doug-martin/goqu/v8"
	"github.com/doug-martin/goqu/v8/exp"
)

// A ...
func A(table string, schema ...string) *Alias {
	c := &Alias{
		exp:  goqu.T(table),
		cols: make(map[string]string),
	}

	if len(schema) != 0 {
		c.exp = c.exp.Schema(schema[0])
	}

	return c
}

// Alias ...
type Alias struct {
	exp  exp.IdentifierExpression
	cols map[string]string
	sync.Mutex
}

// Expression ...
func (c *Alias) Expression() exp.IdentifierExpression {
	return c.exp
}

// String ...
func (c *Alias) String(col string) string {
	if len(c.exp.GetTable()) == 0 {
		return col
	}

	prefix := c.exp.GetTable()
	if len(c.exp.GetSchema()) != 0 {
		prefix = c.exp.GetSchema() + "." + prefix
	}

	return prefix + "." + col
}

// All ...
func (c *Alias) All() exp.IdentifierExpression {
	return c.All()
}

// Col ...
func (c *Alias) Col(col string, alias ...string) *Alias {
	c.Lock()

	c.cols[col] = col
	if len(alias) != 0 {
		c.cols[col] = alias[0]
	}

	c.Unlock()
	return c
}

func (c *Alias) getCols() []interface{} {
	cols := make([]interface{}, 0, len(c.cols))

	c.Lock()
	for col, alias := range c.cols {
		cols = append(cols, c.exp.Col(col).As(alias))
	}
	c.Unlock()

	return cols
}

// Cols ...
func Cols(alias ...*Alias) []interface{} {
	cols := make([]interface{}, 0, 4)

	for _, a := range alias {
		cs := a.getCols()
		if len(cs) == 0 {
			cols = append(cols, a.All())
			continue
		}

		cols = append(cols, cs...)
	}

	return cols
}
