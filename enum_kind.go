package pyooq

import (
	"encoding"
	"fmt"
	"strings"

	"github.com/chronos-tachyon/pyooq/internal"
)

type Kind byte

const (
	Kind_Unspec Kind = iota
	Kind_Query
	Kind_Table
	Kind_Splat
	Kind_Column
	Kind_Literal
	Kind_Op
)

const NumKinds Kind = 7
const InvalidKind = ^Kind(0)

var kindGoNames = [NumKinds]string{
	"pyooq.Kind_Unspec",
	"pyooq.Kind_Query",
	"pyooq.Kind_Table",
	"pyooq.Kind_Splat",
	"pyooq.Kind_Column",
	"pyooq.Kind_Literal",
	"pyooq.Kind_Op",
}

var kindNames = [NumKinds]string{
	"Unspecified",
	"Query",
	"Table",
	"Splat",
	"Column",
	"Literal",
	"Op",
}

var kindMap map[string]Kind

func (enum Kind) IsValid() bool {
	return enum < NumKinds
}
func (enum Kind) GoString() string {
	if enum.IsValid() {
		return kindGoNames[enum]
	}
	return fmt.Sprintf("pyooq.Kind(%d)", uint(enum))
}
func (enum Kind) String() string {
	if enum.IsValid() {
		return kindNames[enum]
	}
	return fmt.Sprintf("%%!ERR[invalid pyooq.Kind %d]", uint(enum))
}
func (enum Kind) GoAppend(out []byte) []byte {
	return appendString(out, enum.GoString())
}
func (enum Kind) Append(out []byte) []byte {
	return appendString(out, enum.String())
}
func (enum Kind) MarshalText() ([]byte, error) {
	return []byte(enum.String()), nil
}
func (enum *Kind) Parse(input string) error {
	if xtype, found := kindMap[input]; found {
		*enum = xtype
		return nil
	}
	if xtype, found := kindMap[strings.ToLower(input)]; found {
		*enum = xtype
		return nil
	}
	*enum = InvalidKind
	return fmt.Errorf("failed to parse %q as pyooq.Kind", input)
}
func (enum *Kind) UnmarshalText(input []byte) error {
	return enum.Parse(string(input))
}

var (
	_ Stringer                 = Kind(0)
	_ encoding.TextMarshaler   = Kind(0)
	_ encoding.TextUnmarshaler = (*Kind)(nil)
)

func init() {
	m := internal.BuildEnumMap[Kind](
		NumKinds,
		kindGoNames[:],
		kindNames[:],
	)

	internal.FinishEnumMap(m)
	kindMap = m
}
