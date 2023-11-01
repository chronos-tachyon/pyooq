package pyooq

import (
	"fmt"

	"github.com/chronos-tachyon/pyooq/internal"
)

type IndexedColumn struct {
	i      *Index
	c      *Column
	isDesc bool
}

func (ic *IndexedColumn) IsValid() bool    { return ic != nil && ic.i != nil }
func (ic *IndexedColumn) GoString() string { return string(ic.GoAppend(nil)) }
func (ic *IndexedColumn) String() string   { return string(ic.Append(nil)) }
func (ic *IndexedColumn) Schema() *Schema  { return ic.Index().Schema() }
func (ic *IndexedColumn) Table() *Table    { return ic.Index().Table() }
func (ic *IndexedColumn) Index() *Index {
	internal.AssertValid(ic)
	return ic.i
}
func (ic *IndexedColumn) Column() *Column {
	internal.AssertValid(ic)
	return ic.c
}
func (ic *IndexedColumn) IsDesc() bool {
	if ic.IsValid() {
		return ic.isDesc
	}
	return false
}
func (ic *IndexedColumn) GoAppend(out []byte) []byte {
	if ic.IsValid() {
		tname := ic.Table().Name()
		iname := ic.Index().Name()
		cname := ic.Column().Name()
		return fmt.Appendf(out, "&pyooq.IndexedColumn[%p, %q, %q, %q]", ic, tname, iname, cname)
	}
	str := "(*pyooq.IndexedColumn)(nil)"
	return appendString(out, str)
}
func (ic *IndexedColumn) Append(out []byte) []byte {
	return ic.GoAppend(out)
}

var _ Stringer = (*IndexedColumn)(nil)
