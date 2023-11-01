package pyooq

import (
	"fmt"

	"github.com/chronos-tachyon/assert"
	"github.com/chronos-tachyon/pyooq/internal"
)

type Column struct {
	t  *Table
	ct Type
	n  string
}

func (c *Column) IsValid() bool    { return c != nil && c.t != nil }
func (c *Column) GoString() string { return string(c.GoAppend(nil)) }
func (c *Column) String() string   { return string(c.Append(nil)) }
func (c *Column) Schema() *Schema  { return c.Table().Schema() }
func (c *Column) Table() *Table {
	internal.AssertValid(c)
	return c.t
}
func (c *Column) Name() string {
	internal.AssertValid(c)
	return c.n
}
func (c *Column) Type() Type {
	internal.AssertValid(c)
	return c.ct
}
func (c *Column) Expr() ColumnExpr {
	internal.AssertValid(c)
	tx := c.Table().Expr()
	return ColumnExpr{c: c, tx: tx}
}
func (c *Column) ExprUsing(tx TableExpr) ColumnExpr {
	internal.AssertValid(c)
	internal.AssertValid(tx)
	t1 := c.Table()
	t2 := tx.Table()
	assert.Assertf(t2 == t1, "%v belongs to %v, not to %v", c, t1, t2)
	return ColumnExpr{c: c, tx: tx}
}
func (c *Column) GoAppend(out []byte) []byte {
	if c.IsValid() {
		tname := c.t.Name()
		cname := c.Name()
		return fmt.Appendf(out, "&pyooq.Column[%p, %q, %q]", c, tname, cname)
	}
	str := "(*pyooq.Column)(nil)"
	return append(out, str...)
}
func (c *Column) Append(out []byte) []byte {
	if c.IsValid() {
		tname := c.t.Name()
		cname := c.Name()
		return fmt.Appendf(out, "%s.%s", tname, cname)
	}
	str := "%!INVALID"
	return append(out, str...)
}

var _ Stringer = (*Column)(nil)
