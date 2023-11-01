package pyooq

import (
	"fmt"

	"github.com/chronos-tachyon/pyooq/internal"
)

type Index struct {
	t    *Table
	n    string
	it   IndexType
	list []*IndexedColumn
}

func (i *Index) IsValid() bool    { return i != nil && i.t != nil }
func (i *Index) Schema() *Schema  { return i.Table().Schema() }
func (i *Index) GoString() string { return string(i.GoAppend(nil)) }
func (i *Index) String() string   { return string(i.Append(nil)) }
func (i *Index) Table() *Table {
	internal.AssertValid(i)
	return i.t
}
func (i *Index) Name() string {
	internal.AssertValid(i)
	return i.n
}
func (i *Index) Type() IndexType {
	internal.AssertValid(i)
	return i.it
}
func (i *Index) NumColumns() uint {
	if i.IsValid() {
		return uint(len(i.list))
	}
	return 0
}
func (i *Index) Column(index uint) *IndexedColumn {
	if size := i.NumColumns(); index >= size {
		panic(fmt.Errorf("%v: %w", i, internal.RangeErr(index, size)))
	}
	return i.list[index]
}
func (i *Index) GoAppend(out []byte) []byte {
	if i.IsValid() {
		tname := i.Table().Name()
		iname := i.Name()
		return fmt.Appendf(out, "&pyooq.Index[%p, %q, %q]", i, tname, iname)
	}
	str := "(*pyooq.Index)(nil)"
	return appendString(out, str)
}
func (i *Index) Append(out []byte) []byte {
	return i.GoAppend(out)
}

var _ Stringer = (*Index)(nil)
