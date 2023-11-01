package pyooq

import (
	"encoding"
	"fmt"
	"strings"

	"github.com/chronos-tachyon/pyooq/internal"
)

type Join byte

const (
	Join_Inner Join = iota
	Join_Outer
)

const NumJoins Join = 2
const InvalidJoin = ^Join(0)

var joinGoNames = [NumJoins]string{
	"pyooq.Join_Inner",
	"pyooq.Join_Outer",
}

var joinNames = [NumJoins]string{
	"innerJoin",
	"outerJoin",
}

var joinMap map[string]Join

func (enum Join) IsValid() bool {
	return enum < NumJoins
}
func (enum Join) GoString() string {
	if enum.IsValid() {
		return joinGoNames[enum]
	}
	return fmt.Sprintf("pyooq.Join(%d)", uint(enum))
}
func (enum Join) String() string {
	if enum.IsValid() {
		return joinNames[enum]
	}
	return fmt.Sprintf("%%!ERR[invalid pyooq.Join %d]", uint(enum))
}
func (enum Join) GoAppend(out []byte) []byte {
	return appendString(out, enum.GoString())
}
func (enum Join) Append(out []byte) []byte {
	return appendString(out, enum.String())
}
func (enum Join) MarshalText() ([]byte, error) {
	return []byte(enum.String()), nil
}
func (enum *Join) Parse(input string) error {
	if jtype, found := joinMap[input]; found {
		*enum = jtype
		return nil
	}
	if jtype, found := joinMap[strings.ToLower(input)]; found {
		*enum = jtype
		return nil
	}
	*enum = InvalidJoin
	return fmt.Errorf("failed to parse %q as pyooq.Join", input)
}
func (enum *Join) UnmarshalText(input []byte) error {
	return enum.Parse(string(input))
}

var (
	_ Stringer                 = Join(0)
	_ encoding.TextMarshaler   = Join(0)
	_ encoding.TextUnmarshaler = (*Join)(nil)
)

func init() {
	m := internal.BuildEnumMap[Join](
		NumJoins,
		joinGoNames[:],
		joinNames[:],
	)

	m["Inner"] = Join_Inner
	m["InnerJoin"] = Join_Inner

	m["Outer"] = Join_Outer
	m["OuterJoin"] = Join_Outer

	internal.FinishEnumMap(m)
	joinMap = m
}
