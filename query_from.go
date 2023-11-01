package pyooq

import (
	"github.com/chronos-tachyon/pyooq/internal"
)

type QueryJoin struct {
	Type      Join
	Table     TableExpr
	Condition ValueExpr
}

type QueryFrom interface {
	IsEmpty() bool
	Base() TableExpr
	NumJoins() uint
	Join(index uint) QueryJoin
}

type emptyFrom struct{}
type flatFrom struct{ only TableExpr }
type joinFrom struct {
	first TableExpr
	rest  []QueryJoin
}

func (from emptyFrom) IsEmpty() bool {
	return true
}
func (from emptyFrom) Base() TableExpr {
	return TableExpr{}
}
func (from emptyFrom) NumJoins() uint {
	return 0
}
func (from emptyFrom) Join(index uint) QueryJoin {
	panic(internal.RangeErr(index, 0))
}

func (from flatFrom) IsEmpty() bool {
	return false
}
func (from flatFrom) Base() TableExpr {
	return from.only
}
func (from flatFrom) NumJoins() uint {
	return 0
}
func (from flatFrom) Join(index uint) QueryJoin {
	panic(internal.RangeErr(index, 0))
}

func (from *joinFrom) IsEmpty() bool {
	return false
}
func (from *joinFrom) Base() TableExpr {
	return from.first
}
func (from *joinFrom) NumJoins() uint {
	return uint(len(from.rest))
}
func (from *joinFrom) Join(index uint) QueryJoin {
	if size := from.NumJoins(); index >= size {
		panic(internal.RangeErr(index, size))
	}
	return from.rest[index]
}

func makeQueryFrom(first TableExpr, rest []QueryJoin) QueryFrom {
	if !first.IsValid() {
		return emptyFrom{}
	}
	if len(rest) <= 0 {
		return flatFrom{only: first}
	}
	rest = internal.CopyList(rest)
	return &joinFrom{first: first, rest: rest}
}

var (
	_ QueryFrom = emptyFrom{}
	_ QueryFrom = flatFrom{}
	_ QueryFrom = (*joinFrom)(nil)
)
