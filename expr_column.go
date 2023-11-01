package pyooq

import (
	"fmt"

	"github.com/chronos-tachyon/pyooq/internal"
)

type ColumnExpr struct {
	c  *Column
	tx TableExpr
}

func (cx ColumnExpr) IsValid() bool     { return cx.c.IsValid() }
func (cx ColumnExpr) Kind() Kind        { return Kind_Column }
func (cx ColumnExpr) GoString() string  { return string(cx.GoAppend(nil)) }
func (cx ColumnExpr) String() string    { return string(cx.Append(nil)) }
func (cx ColumnExpr) Type() Type        { return cx.Column().Type() }
func (cx ColumnExpr) Name() string      { return cx.Column().Name() }
func (cx ColumnExpr) Table() *Table     { return cx.TableExpr().Table() }
func (cx ColumnExpr) TableName() string { return cx.TableExpr().Name() }
func (cx ColumnExpr) Column() *Column {
	internal.AssertValid(cx)
	return cx.c
}
func (cx ColumnExpr) TableExpr() TableExpr {
	internal.AssertValid(cx)
	return cx.tx
}
func (cx ColumnExpr) Visit(v Visitor) error {
	if cx.IsValid() {
		return v.VisitColumn(cx)
	}
	return fmt.Errorf("%v is not valid", cx)
}
func (cx ColumnExpr) GoAppend(out []byte) []byte {
	if cx.IsValid() {
		c := cx.Column()
		tx := cx.TableExpr()
		return fmt.Appendf(out, "%#v.ExprUsing(%#v)", c, tx)
	}
	str := "pyooq.ColumnExpr{}"
	return append(out, str...)
}
func (cx ColumnExpr) Append(out []byte) []byte {
	if cx.IsValid() {
		tname := cx.TableName()
		cname := cx.Name()
		return fmt.Appendf(out, "%s.%s", tname, cname)
	}
	str := "%!INVALID"
	return append(out, str...)
}

var _ ValueExpr = ColumnExpr{}
