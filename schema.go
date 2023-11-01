package pyooq

import (
	"fmt"

	"github.com/chronos-tachyon/pyooq/internal"
)

type Schema struct {
	list  []*Table
	names map[string]*Table
}

func (s *Schema) IsValid() bool    { return s != nil }
func (s *Schema) GoString() string { return string(s.GoAppend(nil)) }
func (s *Schema) String() string   { return string(s.Append(nil)) }
func (s *Schema) NumTables() uint {
	if s.IsValid() {
		return uint(len(s.list))
	}
	return 0
}
func (s *Schema) Table(index uint) *Table {
	if size := s.NumTables(); index >= size {
		panic(fmt.Errorf("%v: %w", s, internal.RangeErr(index, size)))
	}
	return s.list[index]
}
func (s *Schema) ByName(name string) *Table {
	if s.IsValid() {
		if t, found := s.names[name]; found {
			return t
		}
	}
	panic(fmt.Errorf("%v: table %q does not exist", s, name))
}
func (s *Schema) GoAppend(out []byte) []byte {
	if s.IsValid() {
		return fmt.Appendf(out, "&pyooq.Schema[%p]", s)
	}
	str := "(*pyooq.Schema)(nil)"
	return append(out, str...)
}
func (s *Schema) Append(out []byte) []byte {
	return s.GoAppend(out)
}

var _ Stringer = (*Schema)(nil)
