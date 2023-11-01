package pyooq

import (
	"fmt"
	"time"
)

type StringLiteral string

func (lit StringLiteral) Raw() string           { return string(lit) }
func (lit StringLiteral) IsValid() bool         { return true }
func (lit StringLiteral) Kind() Kind            { return Kind_Literal }
func (lit StringLiteral) Type() Type            { return Type_String }
func (lit StringLiteral) Visit(v Visitor) error { return v.VisitStringLiteral(lit.Raw()) }
func (lit StringLiteral) GoString() string      { return string(lit.GoAppend(nil)) }
func (lit StringLiteral) String() string        { return string(lit.Append(nil)) }
func (lit StringLiteral) GoAppend(out []byte) []byte {
	return fmt.Appendf(out, "pyooq.StringLiteral(%q)", lit.Raw())
}
func (lit StringLiteral) Append(out []byte) []byte {
	return fmt.Appendf(out, "%q", lit.Raw())
}

type BytesLiteral []byte

func (lit BytesLiteral) Raw() []byte           { return []byte(lit) }
func (lit BytesLiteral) IsValid() bool         { return true }
func (lit BytesLiteral) Kind() Kind            { return Kind_Literal }
func (lit BytesLiteral) Type() Type            { return Type_Bytes }
func (lit BytesLiteral) Visit(v Visitor) error { return v.VisitBytesLiteral(lit.Raw()) }
func (lit BytesLiteral) GoString() string      { return string(lit.GoAppend(nil)) }
func (lit BytesLiteral) String() string        { return string(lit.Append(nil)) }
func (lit BytesLiteral) GoAppend(out []byte) []byte {
	return fmt.Appendf(out, "pyooq.BytesLiteral(%x)", lit.Raw())
}
func (lit BytesLiteral) Append(out []byte) []byte {
	return fmt.Appendf(out, "%x", lit.Raw())
}

type IntLiteral int64

func (lit IntLiteral) Raw() int64            { return int64(lit) }
func (lit IntLiteral) IsValid() bool         { return true }
func (lit IntLiteral) Kind() Kind            { return Kind_Literal }
func (lit IntLiteral) Type() Type            { return Type_I64 }
func (lit IntLiteral) Visit(v Visitor) error { return v.VisitIntLiteral(lit.Raw()) }
func (lit IntLiteral) GoString() string      { return string(lit.GoAppend(nil)) }
func (lit IntLiteral) String() string        { return string(lit.Append(nil)) }
func (lit IntLiteral) GoAppend(out []byte) []byte {
	return fmt.Appendf(out, "pyooq.IntLiteral(%d)", lit.Raw())
}
func (lit IntLiteral) Append(out []byte) []byte {
	return fmt.Appendf(out, "%d", lit.Raw())
}

type UintLiteral uint64

func (lit UintLiteral) Raw() uint64           { return uint64(lit) }
func (lit UintLiteral) IsValid() bool         { return true }
func (lit UintLiteral) Kind() Kind            { return Kind_Literal }
func (lit UintLiteral) Type() Type            { return Type_U64 }
func (lit UintLiteral) Visit(v Visitor) error { return v.VisitUintLiteral(lit.Raw()) }
func (lit UintLiteral) GoString() string      { return string(lit.GoAppend(nil)) }
func (lit UintLiteral) String() string        { return string(lit.Append(nil)) }
func (lit UintLiteral) GoAppend(out []byte) []byte {
	return fmt.Appendf(out, "pyooq.UintLiteral(%d)", lit.Raw())
}
func (lit UintLiteral) Append(out []byte) []byte {
	return fmt.Appendf(out, "%d", lit.Raw())
}

type FloatLiteral float64

func (lit FloatLiteral) Raw() float64          { return float64(lit) }
func (lit FloatLiteral) IsValid() bool         { return true }
func (lit FloatLiteral) Kind() Kind            { return Kind_Literal }
func (lit FloatLiteral) Type() Type            { return Type_F64 }
func (lit FloatLiteral) Visit(v Visitor) error { return v.VisitFloatLiteral(lit.Raw()) }
func (lit FloatLiteral) GoString() string      { return string(lit.GoAppend(nil)) }
func (lit FloatLiteral) String() string        { return string(lit.Append(nil)) }
func (lit FloatLiteral) GoAppend(out []byte) []byte {
	return fmt.Appendf(out, "pyooq.FloatLiteral(%g)", lit.Raw())
}
func (lit FloatLiteral) Append(out []byte) []byte {
	return fmt.Appendf(out, "%g", lit.Raw())
}

type BoolLiteral bool

var boolMap = map[bool]string{
	false: "FALSE",
	true:  "TRUE",
}

func (lit BoolLiteral) Raw() bool             { return bool(lit) }
func (lit BoolLiteral) IsValid() bool         { return true }
func (lit BoolLiteral) Kind() Kind            { return Kind_Literal }
func (lit BoolLiteral) Type() Type            { return Type_Bool }
func (lit BoolLiteral) Visit(v Visitor) error { return v.VisitBoolLiteral(lit.Raw()) }
func (lit BoolLiteral) GoString() string      { return string(lit.GoAppend(nil)) }
func (lit BoolLiteral) String() string        { return boolMap[lit.Raw()] }
func (lit BoolLiteral) GoAppend(out []byte) []byte {
	return fmt.Appendf(out, "pyooq.BoolLiteral(%t)", lit.Raw())
}
func (lit BoolLiteral) Append(out []byte) []byte {
	str := boolMap[lit.Raw()]
	return appendString(out, str)
}

type NullLiteral struct{}

func (lit NullLiteral) IsValid() bool              { return true }
func (lit NullLiteral) Kind() Kind                 { return Kind_Literal }
func (lit NullLiteral) Type() Type                 { return Type_Null }
func (lit NullLiteral) Visit(v Visitor) error      { return v.VisitNullLiteral() }
func (lit NullLiteral) GoString() string           { return "pyooq.NullLiteral{}" }
func (lit NullLiteral) String() string             { return "NULL" }
func (lit NullLiteral) GoAppend(out []byte) []byte { return appendString(out, "pyooq.NullLiteral{}") }
func (lit NullLiteral) Append(out []byte) []byte   { return appendString(out, "NULL") }

type TimeLiteral struct{ Time time.Time }

func (lit TimeLiteral) Raw() time.Time        { return lit.Time }
func (lit TimeLiteral) IsValid() bool         { return true }
func (lit TimeLiteral) Kind() Kind            { return Kind_Literal }
func (lit TimeLiteral) Type() Type            { return Type_Timestamp }
func (lit TimeLiteral) Visit(v Visitor) error { return v.VisitTimeLiteral(lit.Raw()) }
func (lit TimeLiteral) GoString() string      { return string(lit.GoAppend(nil)) }
func (lit TimeLiteral) String() string        { return string(lit.Append(nil)) }
func (lit TimeLiteral) GoAppend(out []byte) []byte {
	t := lit.Raw()
	return fmt.Appendf(out, "pyooq.TimeLiteral{%#v}", t)
}
func (lit TimeLiteral) Append(out []byte) []byte {
	t := lit.Raw()
	return fmt.Appendf(out, "%q", t.Format(time.RFC3339Nano))
}

var (
	_ ValueExpr = StringLiteral("")
	_ ValueExpr = BytesLiteral(nil)
	_ ValueExpr = IntLiteral(0)
	_ ValueExpr = UintLiteral(0)
	_ ValueExpr = FloatLiteral(0)
	_ ValueExpr = BoolLiteral(false)
	_ ValueExpr = NullLiteral{}
	_ ValueExpr = TimeLiteral{}
)
