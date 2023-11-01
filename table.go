package pyooq

import (
	"fmt"

	"github.com/chronos-tachyon/pyooq/internal"
)

type Table struct {
	s        *Schema
	n        string
	colList  []*Column
	colNames map[string]*Column
	idxList  []*Index
	idxNames map[string]*Index
}

func (t *Table) IsValid() bool    { return t != nil && t.s != nil }
func (t *Table) GoString() string { return string(t.GoAppend(nil)) }
func (t *Table) String() string   { return string(t.Append(nil)) }
func (t *Table) Schema() *Schema {
	internal.AssertValid(t)
	return t.s
}
func (t *Table) Name() string {
	internal.AssertValid(t)
	return t.n
}
func (t *Table) NumColumns() uint {
	if t.IsValid() {
		return uint(len(t.colList))
	}
	return 0
}
func (t *Table) Column(index uint) *Column {
	if size := t.NumColumns(); index >= size {
		panic(fmt.Errorf("%v: %w", t, internal.RangeErr(index, size)))
	}
	return t.colList[index]
}
func (t *Table) ByName(name string) *Column {
	if t.IsValid() {
		if c, found := t.colNames[name]; found {
			return c
		}
	}
	panic(fmt.Errorf("%v: column %q does not exist", t, name))
}
func (t *Table) NumIndices() uint {
	if t.IsValid() {
		return uint(len(t.idxList))
	}
	return 0
}
func (t *Table) Index(index uint) *Index {
	if size := t.NumIndices(); index >= size {
		panic(fmt.Errorf("%v: %w", t, internal.RangeErr(index, size)))
	}
	return t.idxList[index]
}
func (t *Table) PrimaryKey() *Index {
	if t.IsValid() {
		if i, found := t.idxNames[primaryKeyName]; found {
			return i
		}
	}
	panic(fmt.Errorf("%v: primary key does not exist", t))
}
func (t *Table) IndexByName(name string) *Index {
	if t.IsValid() {
		if i, found := t.idxNames[name]; found {
			return i
		}
	}
	panic(fmt.Errorf("%v: index %q does not exist", t, name))
}
func (t *Table) Expr() TableExpr {
	internal.AssertValid(t)
	name := t.Name()
	return TableExpr{t: t, n: name}
}
func (t *Table) As(alias string) TableExpr {
	internal.AssertValid(t)
	internal.AssertValidName("table alias", alias)
	return TableExpr{t: t, n: alias}
}
func (t *Table) GoAppend(out []byte) []byte {
	if t.IsValid() {
		return fmt.Appendf(out, "&pyooq.Table[%p, %q]", t, t.n)
	}
	str := "(*pyooq.Table)(nil)"
	return append(out, str...)
}
func (t *Table) Append(out []byte) []byte {
	str := "%!INVALID"
	if t.IsValid() {
		str = t.n
	}
	return append(out, str...)
}

var _ Stringer = (*Table)(nil)
