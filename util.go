package pyooq

import (
	"github.com/chronos-tachyon/pyooq/internal"
)

const primaryKeyName = "*PRIMARY KEY"

func appendString(out []byte, str string) []byte {
	return append(out, str...)
}

func assertKind(x Expr, expect Kind) {
	internal.AssertValid(x)
	internal.AssertIs("Expr has wrong kind", x.Kind(), expect)
}

func assertKindIn(x Expr, list ...Kind) {
	internal.AssertValid(x)
	internal.AssertOneOf("Expr has wrong kind", x.Kind(), list...)
}

func assertValue(x Expr) ValueExpr {
	assertKindIn(x, valueKinds...)
	return x.(ValueExpr)
}

func assertType(x Expr, expect Type) {
	v := assertValue(x)
	internal.AssertIs("ValueExpr has wrong type", v.Type(), expect)
}

func assertTypeIn(x Expr, list ...Type) {
	v := assertValue(x)
	internal.AssertOneOf("ValueExpr has wrong type", v.Type(), list...)
}

var valueKinds = []Kind{
	Kind_Column,
	Kind_Literal,
	Kind_Op,
}
