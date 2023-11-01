package pyooq

import (
	"github.com/chronos-tachyon/pyooq/internal"
)

type OrderCriterion struct {
	Index  uint
	IsDesc bool
}

type QueryOrderBy interface {
	IsEmpty() bool
	NumCriteria() uint
	Criterion(index uint) OrderCriterion
}

type emptyOrderBy struct{}
type singleOrderBy struct{ only OrderCriterion }
type listOrderBy struct{ list []OrderCriterion }

func (emptyOrderBy) IsEmpty() bool {
	return true
}
func (emptyOrderBy) NumCriteria() uint {
	return 0
}
func (emptyOrderBy) Criterion(index uint) OrderCriterion {
	panic(internal.RangeErr(index, 0))
}

func (o *singleOrderBy) IsEmpty() bool {
	return false
}
func (o *singleOrderBy) NumCriteria() uint {
	return 1
}
func (o *singleOrderBy) Criterion(index uint) OrderCriterion {
	switch index {
	case 0:
		return o.only
	default:
		panic(internal.RangeErr(index, 1))
	}
}

func (o *listOrderBy) IsEmpty() bool {
	return false
}
func (o *listOrderBy) NumCriteria() uint {
	return uint(len(o.list))
}
func (o *listOrderBy) Criterion(index uint) OrderCriterion {
	if size := o.NumCriteria(); index >= size {
		panic(internal.RangeErr(index, size))
	}
	return o.list[index]
}

func makeOrderBy(list []OrderCriterion) QueryOrderBy {
	if len(list) <= 0 {
		return emptyOrderBy{}
	}
	if len(list) <= 1 {
		return &singleOrderBy{only: list[0]}
	}
	list = internal.CopyList(list)
	return &listOrderBy{list: list}
}

var (
	_ QueryOrderBy = emptyOrderBy{}
	_ QueryOrderBy = (*singleOrderBy)(nil)
	_ QueryOrderBy = (*listOrderBy)(nil)
)
