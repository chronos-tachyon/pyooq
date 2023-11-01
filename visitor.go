package pyooq

import (
	"errors"
	"time"
)

var (
	Skip           = errors.New("skip")
	NotImplemented = errors.New("not implemented")
)

type Visitor interface {
	VisitQuery(*Query) error
	VisitTable(TableExpr) error
	VisitSplat(SplatExpr) error
	VisitColumn(ColumnExpr) error
	VisitOp(OpExpr) error

	VisitNullLiteral() error
	VisitBoolLiteral(bool) error
	VisitIntLiteral(int64) error
	VisitUintLiteral(uint64) error
	VisitFloatLiteral(float64) error
	VisitStringLiteral(string) error
	VisitBytesLiteral([]byte) error
	VisitTimeLiteral(time.Time) error
}

type BaseVisitor struct{}

func (BaseVisitor) VisitQuery(q *Query) error      { return nil }
func (BaseVisitor) VisitTable(x TableExpr) error   { return nil }
func (BaseVisitor) VisitSplat(x SplatExpr) error   { return nil }
func (BaseVisitor) VisitColumn(x ColumnExpr) error { return nil }
func (BaseVisitor) VisitOp(x OpExpr) error         { return nil }

func (BaseVisitor) VisitNullLiteral() error            { return nil }
func (BaseVisitor) VisitBoolLiteral(v bool) error      { return nil }
func (BaseVisitor) VisitIntLiteral(v int64) error      { return nil }
func (BaseVisitor) VisitUintLiteral(v uint64) error    { return nil }
func (BaseVisitor) VisitFloatLiteral(v float64) error  { return nil }
func (BaseVisitor) VisitStringLiteral(v string) error  { return nil }
func (BaseVisitor) VisitBytesLiteral(v []byte) error   { return nil }
func (BaseVisitor) VisitTimeLiteral(v time.Time) error { return nil }

var _ Visitor = BaseVisitor{}
