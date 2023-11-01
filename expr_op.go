package pyooq

import (
	"fmt"

	"github.com/chronos-tachyon/pyooq/internal"
)

type OpExpr struct {
	op   Operator
	t    Type
	impl opImpl
}

func (ox OpExpr) IsValid() bool    { return ox.impl != nil }
func (ox OpExpr) Kind() Kind       { return Kind_Op }
func (ox OpExpr) GoString() string { return string(ox.GoAppend(nil)) }
func (ox OpExpr) String() string   { return string(ox.Append(nil)) }
func (ox OpExpr) Type() Type {
	internal.AssertValid(ox)
	return ox.t
}
func (ox OpExpr) Operator() Operator {
	internal.AssertValid(ox)
	return ox.op
}
func (ox OpExpr) NumOperands() uint {
	internal.AssertValid(ox)
	return ox.impl.Size()
}
func (ox OpExpr) Operand(index uint) Expr {
	internal.AssertValid(ox)
	return ox.impl.At(index)
}
func (ox OpExpr) Visit(v Visitor) error {
	if ox.IsValid() {
		return v.VisitOp(ox)
	}
	return fmt.Errorf("%#v is not valid", ox)
}
func (ox OpExpr) GoAppend(out []byte) []byte {
	if ox.IsValid() {
		out = fmt.Appendf(out, "pyooq.OpExpr{%#v, %#v", ox.op, ox.t)
		size := ox.NumOperands()
		for i := uint(0); i < size; i++ {
			x := ox.Operand(i)
			out = fmt.Appendf(out, ", %#v", x)
		}
		out = fmt.Appendf(out, "}")
		return out
	}
	str := "pyooq.OpExpr{}"
	return append(out, str...)
}
func (ox OpExpr) Append(out []byte) []byte {
	if ox.IsValid() {
		out = fmt.Appendf(out, "[%#v, %#v", ox.op, ox.t)
		size := ox.NumOperands()
		for i := uint(0); i < size; i++ {
			x := ox.Operand(i)
			out = fmt.Appendf(out, ", %#v", x)
		}
		out = fmt.Appendf(out, "]")
		return out
	}
	str := "%!INVALID"
	return append(out, str...)
}

type opImpl interface {
	Size() uint
	At(index uint) Expr
}

type nilaryOpImpl struct{}

func (nilaryOpImpl) Size() uint {
	return 0
}
func (nilaryOpImpl) At(index uint) Expr {
	panic(internal.RangeErr(index, 0))
}
func makeNilaryOp(op Operator, t Type) OpExpr {
	return OpExpr{op, t, nilaryOpImpl{}}
}

type unaryOpImpl[T ValueExpr] struct{ a T }

func (impl *unaryOpImpl[T]) Size() uint {
	return 1
}
func (impl *unaryOpImpl[T]) At(index uint) Expr {
	switch index {
	case 0:
		return impl.a
	default:
		panic(internal.RangeErr(index, 1))
	}
}
func makeUnaryOp[T ValueExpr](op Operator, t Type, a T) OpExpr {
	return OpExpr{op, t, &unaryOpImpl[T]{a}}
}

type binaryOpImpl[TA Expr, TB Expr] struct {
	a TA
	b TB
}

func (impl *binaryOpImpl[TA, TB]) Size() uint {
	return 2
}
func (impl *binaryOpImpl[TA, TB]) At(index uint) Expr {
	switch index {
	case 0:
		return impl.a
	case 1:
		return impl.b
	default:
		panic(internal.RangeErr(index, 2))
	}
}
func makeBinaryOp[TA Expr, TB Expr](op Operator, t Type, a TA, b TB) OpExpr {
	return OpExpr{op, t, &binaryOpImpl[TA, TB]{a, b}}
}

type ternaryOpImpl[TA Expr, TB Expr, TC Expr] struct {
	a TA
	b TB
	c TC
}

func (impl *ternaryOpImpl[TA, TB, TC]) Size() uint {
	return 3
}
func (impl *ternaryOpImpl[TA, TB, TC]) At(index uint) Expr {
	switch index {
	case 0:
		return impl.a
	case 1:
		return impl.b
	case 2:
		return impl.c
	default:
		panic(internal.RangeErr(index, 3))
	}
}
func makeTernaryOp[TA Expr, TB Expr, TC Expr](op Operator, t Type, a TA, b TB, c TC) OpExpr {
	return OpExpr{op, t, &ternaryOpImpl[TA, TB, TC]{a, b, c}}
}

type listOpImpl[T Expr] struct{ list []T }

func (impl *listOpImpl[T]) Size() uint {
	return uint(len(impl.list))
}
func (impl *listOpImpl[T]) At(index uint) Expr {
	if size := impl.Size(); index >= size {
		panic(internal.RangeErr(index, size))
	}
	return impl.list[index]
}
func makeMultiOp[T Expr](op Operator, t Type, list ...T) OpExpr {
	return OpExpr{op, t, &listOpImpl[T]{list}}
}

var _ ValueExpr = OpExpr{}
