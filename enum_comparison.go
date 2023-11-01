package pyooq

import (
	"encoding"
	"fmt"
	"strings"

	"github.com/chronos-tachyon/pyooq/internal"
)

type Comparison byte

const (
	Cmp_EQ Comparison = iota
	Cmp_NE
	Cmp_LT
	Cmp_LE
	Cmp_GE
	Cmp_GT

	Cmp_LTE = Cmp_LE
	Cmp_GTE = Cmp_GE
)

const NumComparisons Comparison = 6
const InvalidComparison = ^Comparison(0)

var cmpGoNames = [NumComparisons]string{
	"pyooq.Cmp_EQ",
	"pyooq.Cmp_NE",
	"pyooq.Cmp_LT",
	"pyooq.Cmp_LE",
	"pyooq.Cmp_GE",
	"pyooq.Cmp_GT",
}

var cmpNames = [NumComparisons]string{
	"EQ",
	"NE",
	"LT",
	"LE",
	"GE",
	"GT",
}

var cmpMap map[string]Comparison

func (enum Comparison) IsValid() bool {
	return enum < NumComparisons
}
func (enum Comparison) GoString() string {
	if enum.IsValid() {
		return cmpGoNames[enum]
	}
	return fmt.Sprintf("pyooq.Comparison(%d)", uint(enum))
}
func (enum Comparison) String() string {
	if enum.IsValid() {
		return cmpNames[enum]
	}
	return fmt.Sprintf("%%!ERR[invalid pyooq.Comparison %d]", uint(enum))
}
func (enum Comparison) GoAppend(out []byte) []byte {
	return appendString(out, enum.GoString())
}
func (enum Comparison) Append(out []byte) []byte {
	return appendString(out, enum.String())
}
func (enum Comparison) MarshalText() ([]byte, error) {
	return []byte(enum.String()), nil
}
func (enum *Comparison) Parse(input string) error {
	if compare, found := cmpMap[input]; found {
		*enum = compare
		return nil
	}
	if compare, found := cmpMap[strings.ToLower(input)]; found {
		*enum = compare
		return nil
	}
	*enum = InvalidComparison
	return fmt.Errorf("failed to parse %q as pyooq.Comparison", input)
}
func (enum *Comparison) UnmarshalText(input []byte) error {
	return enum.Parse(string(input))
}
func (enum Comparison) Eval(cmp int) bool {
	switch enum {
	case Cmp_EQ:
		return cmp == 0
	case Cmp_NE:
		return cmp != 0
	case Cmp_LT:
		return cmp < 0
	case Cmp_LE:
		return cmp <= 0
	case Cmp_GE:
		return cmp >= 0
	case Cmp_GT:
		return cmp > 0
	default:
		return false
	}
}

var (
	_ Stringer                 = Comparison(0)
	_ encoding.TextMarshaler   = Comparison(0)
	_ encoding.TextUnmarshaler = (*Comparison)(nil)
)

func init() {
	m := internal.BuildEnumMap[Comparison](
		NumComparisons,
		cmpGoNames[:],
		cmpNames[:],
	)

	m["Equal"] = Cmp_EQ
	m["Same"] = Cmp_EQ
	m["Is"] = Cmp_EQ
	m["=="] = Cmp_EQ
	m["="] = Cmp_EQ

	m["NotEqual"] = Cmp_NE
	m["Unequal"] = Cmp_NE
	m["Different"] = Cmp_NE
	m["IsNot"] = Cmp_NE
	m["!="] = Cmp_NE
	m["≠"] = Cmp_NE

	m["LessThan"] = Cmp_LT
	m["Less"] = Cmp_LT
	m["<"] = Cmp_LT

	m["LTE"] = Cmp_LE
	m["LessThanOrEqual"] = Cmp_LE
	m["LessThanEqual"] = Cmp_LE
	m["LessOrEqual"] = Cmp_LE
	m["LessEqual"] = Cmp_LE
	m["<="] = Cmp_LE
	m["≤"] = Cmp_LE

	m["GTE"] = Cmp_GE
	m["GreaterThanOrEqual"] = Cmp_GE
	m["GreaterThanEqual"] = Cmp_GE
	m["GreaterOrEqual"] = Cmp_GE
	m["GreaterEqual"] = Cmp_GE
	m["MoreThanOrEqual"] = Cmp_GE
	m["MoreThanEqual"] = Cmp_GE
	m["MoreOrEqual"] = Cmp_GE
	m["MoreEqual"] = Cmp_GE
	m["≥"] = Cmp_GE
	m[">="] = Cmp_GE

	m["GreaterThan"] = Cmp_GT
	m["Greater"] = Cmp_GT
	m["MoreThan"] = Cmp_GT
	m["More"] = Cmp_GT
	m[">"] = Cmp_GT

	internal.FinishEnumMap(m)
	cmpMap = m
}
