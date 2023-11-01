package pyooq

import (
	"encoding"
	"fmt"
	"strings"

	"github.com/chronos-tachyon/pyooq/internal"
)

type Operator byte

const (
	OpIdentity Operator = iota
	OpNot
	OpAnd
	OpNand
	OpOr
	OpNor
	OpXor
	OpXnor
	OpAndNot
	OpEQ
	OpNE
	OpLT
	OpLE
	OpGE
	OpGT
	OpNeg
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpSplat
)

const NumOperators Operator = 22
const InvalidOperator = ^Operator(0)

const uintMax = ^uint(0)

var operatorGoNames = [NumOperators]string{
	"pyooq.OpIdentity",
	"pyooq.OpNot",
	"pyooq.OpAnd",
	"pyooq.OpNand",
	"pyooq.OpOr",
	"pyooq.OpNor",
	"pyooq.OpXor",
	"pyooq.OpXnor",
	"pyooq.OpAndNot",
	"pyooq.OpEQ",
	"pyooq.OpNE",
	"pyooq.OpLT",
	"pyooq.OpLE",
	"pyooq.OpGE",
	"pyooq.OpGT",
	"pyooq.OpNeg",
	"pyooq.OpAdd",
	"pyooq.OpSub",
	"pyooq.OpMul",
	"pyooq.OpDiv",
	"pyooq.OpMod",
	"pyooq.OpSplat",
}

var operatorNames = [NumOperators]string{
	"Identity",
	"Not",
	"And",
	"Nand",
	"Or",
	"Nor",
	"Xor",
	"Xnor",
	"AndNot",
	"EQ",
	"NE",
	"LT",
	"LE",
	"GE",
	"GT",
	"Neg",
	"Add",
	"Sub",
	"Mul",
	"Div",
	"Mod",
	"Splat",
}

var operatorMap map[string]Operator

func (enum Operator) IsValid() bool {
	return enum < NumOperators
}
func (enum Operator) GoString() string {
	if enum.IsValid() {
		return operatorGoNames[enum]
	}
	return fmt.Sprintf("pyooq.Operator(%d)", uint(enum))
}
func (enum Operator) String() string {
	if enum.IsValid() {
		return operatorNames[enum]
	}
	return fmt.Sprintf("%%!ERR[invalid pyooq.Operator %d]", uint(enum))
}
func (enum Operator) GoAppend(out []byte) []byte {
	return appendString(out, enum.GoString())
}
func (enum Operator) Append(out []byte) []byte {
	return appendString(out, enum.String())
}
func (enum Operator) MarshalText() ([]byte, error) {
	return []byte(enum.String()), nil
}
func (enum *Operator) Parse(input string) error {
	if op, found := operatorMap[input]; found {
		*enum = op
		return nil
	}
	if op, found := operatorMap[strings.ToLower(input)]; found {
		*enum = op
		return nil
	}
	*enum = InvalidOperator
	return fmt.Errorf("failed to parse %q as pyooq.Operator", input)
}
func (enum *Operator) UnmarshalText(input []byte) error {
	return enum.Parse(string(input))
}

var (
	_ Stringer                 = Operator(0)
	_ encoding.TextMarshaler   = Operator(0)
	_ encoding.TextUnmarshaler = (*Operator)(nil)
)

func init() {
	m := internal.BuildEnumMap[Operator](
		NumOperators,
		operatorGoNames[:],
		operatorNames[:],
	)

	internal.FinishEnumMap(m)
	operatorMap = m
}
