package pyooq

import (
	"fmt"

	"github.com/chronos-tachyon/pyooq/internal"
)

type Query struct {
	from  QueryFrom
	sel   QuerySelect
	where ValueExpr
	order QueryOrderBy
}

func (q *Query) IsValid() bool    { return q != nil && q.from != nil }
func (q *Query) Kind() Kind       { return Kind_Query }
func (q *Query) GoString() string { return string(q.GoAppend(nil)) }
func (q *Query) String() string   { return string(q.Append(nil)) }
func (q *Query) From() QueryFrom {
	internal.AssertValid(q)
	return q.from
}
func (q *Query) Select() QuerySelect {
	internal.AssertValid(q)
	return q.sel
}
func (q *Query) Where() ValueExpr {
	internal.AssertValid(q)
	return q.where
}
func (q *Query) OrderBy() QueryOrderBy {
	internal.AssertValid(q)
	return q.order
}
func (q *Query) Visit(v Visitor) error {
	if q.IsValid() {
		return v.VisitQuery(q)
	}
	return fmt.Errorf("%#v is not valid", q)
}
func (q *Query) GoAppend(out []byte) []byte {
	if q.IsValid() {
		return fmt.Appendf(out, "&pyooq.Query[%p]", q)
	}
	str := "(*pyooq.Query)(nil)"
	return append(out, str...)
}
func (q *Query) Append(out []byte) []byte {
	if q.IsValid() {
		return fmt.Appendf(out, "Query[%p]", q)
	}
	str := "%!INVALID"
	return append(out, str...)
}

var _ Expr = (*Query)(nil)
