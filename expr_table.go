package pyooq

import (
	"fmt"

	"github.com/chronos-tachyon/pyooq/internal"
)

type TableExpr struct {
	t *Table
	n string
}

func (tx TableExpr) IsValid() bool    { return tx.t.IsValid() }
func (tx TableExpr) Kind() Kind       { return Kind_Table }
func (tx TableExpr) GoString() string { return string(tx.GoAppend(nil)) }
func (tx TableExpr) String() string   { return string(tx.Append(nil)) }
func (tx TableExpr) Table() *Table {
	internal.AssertValid(tx)
	return tx.t
}
func (tx TableExpr) As(alias string) TableExpr {
	return tx.Table().As(alias)
}
func (tx TableExpr) Name() string {
	internal.AssertValid(tx)
	return tx.n
}
func (tx TableExpr) Splat() SplatExpr {
	internal.AssertValid(tx)
	return SplatExpr{tx: tx}
}
func (tx TableExpr) Visit(v Visitor) error {
	if tx.IsValid() {
		return v.VisitTable(tx)
	}
	return fmt.Errorf("%v is not valid", tx)
}
func (tx TableExpr) GoAppend(out []byte) []byte {
	if tx.IsValid() {
		t := tx.Table()
		name := tx.Name()
		return fmt.Appendf(out, "%#v.As(%q)", t, name)
	}
	str := "pyooq.TableExpr{}"
	return append(out, str...)
}
func (tx TableExpr) Append(out []byte) []byte {
	str := "%!INVALID"
	if tx.IsValid() {
		str = tx.Name()
	}
	return append(out, str...)
}

var _ Expr = TableExpr{}
