package pyooq

type Stringer interface {
	GoAppend(out []byte) []byte
	GoString() string
	Append(out []byte) []byte
	String() string
}

type Expr interface {
	IsValid() bool
	Kind() Kind
	Visit(Visitor) error
	Stringer
}

type ValueExpr interface {
	Expr
	Type() Type
}
