package pyooq

import (
	"github.com/chronos-tachyon/pyooq/internal"
)

type QueryColumn struct {
	Name  string
	Value ValueExpr
}

type QuerySelect interface {
	IsEmpty() bool
	NumColumns() uint
	Column(index uint) QueryColumn
}

type emptySelect struct{}
type singleSelect struct{ only QueryColumn }
type listSelect struct{ list []QueryColumn }

func (emptySelect) IsEmpty() bool {
	return true
}
func (emptySelect) NumColumns() uint {
	return 0
}
func (emptySelect) Column(index uint) QueryColumn {
	panic(internal.RangeErr(index, 0))
}

func (s *singleSelect) IsEmpty() bool {
	return false
}
func (s *singleSelect) NumColumns() uint {
	return 1
}
func (s *singleSelect) Column(index uint) QueryColumn {
	switch index {
	case 0:
		return s.only
	default:
		panic(internal.RangeErr(index, 1))
	}
}

func (s *listSelect) IsEmpty() bool {
	return false
}
func (s *listSelect) NumColumns() uint {
	return uint(len(s.list))
}
func (s *listSelect) Column(index uint) QueryColumn {
	if size := s.NumColumns(); index >= size {
		panic(internal.RangeErr(index, size))
	}
	return s.list[index]
}

func makeQuerySelect(list []QueryColumn) QuerySelect {
	if len(list) <= 0 {
		return emptySelect{}
	}
	if len(list) <= 1 {
		return &singleSelect{only: list[0]}
	}
	list = internal.CopyList(list)
	return &listSelect{list: list}
}

var (
	_ QuerySelect = emptySelect{}
	_ QuerySelect = (*singleSelect)(nil)
	_ QuerySelect = (*listSelect)(nil)
)
