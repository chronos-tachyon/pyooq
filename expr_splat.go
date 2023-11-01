package pyooq

import (
	"fmt"

	"github.com/chronos-tachyon/pyooq/internal"
)

type SplatExpr struct{ tx TableExpr }

func (sx SplatExpr) IsValid() bool    { return sx.tx.IsValid() }
func (sx SplatExpr) Kind() Kind       { return Kind_Splat }
func (sx SplatExpr) GoString() string { return string(sx.GoAppend(nil)) }
func (sx SplatExpr) String() string   { return string(sx.Append(nil)) }
func (sx SplatExpr) Table() *Table    { return sx.TableExpr().Table() }
func (sx SplatExpr) Name() string     { return sx.TableExpr().Name() }
func (sx SplatExpr) TableExpr() TableExpr {
	internal.AssertValid(sx)
	return sx.tx
}
func (sx SplatExpr) Visit(v Visitor) error {
	if sx.IsValid() {
		return v.VisitSplat(sx)
	}
	return fmt.Errorf("%v is not valid", sx)
}
func (sx SplatExpr) GoAppend(out []byte) []byte {
	if sx.IsValid() {
		return fmt.Appendf(out, "%#v.Splat()", sx.TableExpr())
	}
	str := "pyooq.SplatExpr{}"
	return append(out, str...)
}
func (sx SplatExpr) Append(out []byte) []byte {
	if sx.IsValid() {
		return fmt.Appendf(out, "%v.*", sx.TableExpr())
	}
	str := "%!INVALID"
	return append(out, str...)
}

var _ Expr = SplatExpr{}
