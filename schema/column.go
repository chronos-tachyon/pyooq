package schema

import (
	"github.com/chronos-tachyon/pyooq/types"
)

type Type = types.Type
type ColumnPredicate = Predicate[Column]
type ColumnConsumer = Consumer[Column]

type Column interface {
	Namer
	Typer
}

func MakeColumn(colName string, colType Type) Column {
	return &genericColumn{colName: colName, colType: colType}
}

type genericColumn struct {
	colName string
	colType Type
}

func (c *genericColumn) Name() string { return c.colName }
func (c *genericColumn) Type() Type   { return c.colType }

var _ Column = (*genericColumn)(nil)
