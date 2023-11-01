package pyooq

import (
	"encoding"
	"fmt"
	"strings"

	"github.com/chronos-tachyon/pyooq/internal"
)

type Type byte

const (
	Type_Void Type = iota
	Type_Null
	Type_I64
	Type_I32
	Type_I16
	Type_I8
	Type_U64
	Type_U32
	Type_U16
	Type_U8
	Type_F64
	Type_F32
	Type_Bool
	Type_Timestamp
	Type_String
	Type_Bytes
	Type_UUID
	Type_IPAddr
)

const NumTypes Type = 18
const InvalidType = ^Type(0)

var typeGoNames = [NumTypes]string{
	"pyooq.Type_Void",
	"pyooq.Type_Null",
	"pyooq.Type_I64",
	"pyooq.Type_I32",
	"pyooq.Type_I16",
	"pyooq.Type_I8",
	"pyooq.Type_U64",
	"pyooq.Type_U32",
	"pyooq.Type_U16",
	"pyooq.Type_U8",
	"pyooq.Type_F64",
	"pyooq.Type_F32",
	"pyooq.Type_Bool",
	"pyooq.Type_Timestamp",
	"pyooq.Type_String",
	"pyooq.Type_Bytes",
	"pyooq.Type_UUID",
	"pyooq.Type_IPAddr",
}

var typeNames = [NumTypes]string{
	"Void",
	"Null",
	"I64",
	"I32",
	"I16",
	"I8",
	"U64",
	"U32",
	"U16",
	"U8",
	"F64",
	"F32",
	"Bool",
	"Timestamp",
	"String",
	"Binary",
	"UUID",
	"IPAddr",
}

var typeWidensTo [NumTypes][NumTypes]bool

var typeMap map[string]Type

func (enum Type) IsValid() bool {
	return enum < NumTypes
}
func (enum Type) WidensTo(other Type) bool {
	if enum.IsValid() && other.IsValid() {
		return typeWidensTo[enum][other]
	}
	return false
}
func (enum Type) NarrowsTo(other Type) bool {
	return other.WidensTo(enum)
}
func (enum Type) GoString() string {
	if enum.IsValid() {
		return typeGoNames[enum]
	}
	return fmt.Sprintf("pyooq.Type(%d)", uint(enum))
}
func (enum Type) String() string {
	if enum.IsValid() {
		return typeNames[enum]
	}
	return fmt.Sprintf("%%!ERR[invalid pyooq.Type %d]", uint(enum))
}
func (enum Type) GoAppend(out []byte) []byte {
	return appendString(out, enum.GoString())
}
func (enum Type) Append(out []byte) []byte {
	return appendString(out, enum.String())
}
func (enum Type) MarshalText() ([]byte, error) {
	return []byte(enum.String()), nil
}
func (enum *Type) Parse(input string) error {
	if t, found := typeMap[input]; found {
		*enum = t
		return nil
	}
	if t, found := typeMap[strings.ToLower(input)]; found {
		*enum = t
		return nil
	}
	*enum = InvalidType
	return fmt.Errorf("failed to parse %q as pyooq.Type", input)
}
func (enum *Type) UnmarshalText(input []byte) error {
	return enum.Parse(string(input))
}

var (
	_ Stringer                 = Type(0)
	_ encoding.TextMarshaler   = Type(0)
	_ encoding.TextUnmarshaler = (*Type)(nil)
)

func init() {
	var a [NumTypes][NumTypes]bool
	for t := Type(0); t.IsValid(); t++ {
		a[t][t] = true
	}

	a[Type_I8][Type_I16] = true
	a[Type_I8][Type_I32] = true
	a[Type_I8][Type_I64] = true
	a[Type_I16][Type_I32] = true
	a[Type_I16][Type_I64] = true
	a[Type_I32][Type_I64] = true

	a[Type_U8][Type_U16] = true
	a[Type_U8][Type_U32] = true
	a[Type_U8][Type_U64] = true
	a[Type_U16][Type_U32] = true
	a[Type_U16][Type_U64] = true
	a[Type_U32][Type_U64] = true

	a[Type_F32][Type_F64] = true

	// test for transitivity
	for i := Type(0); i < NumTypes; i++ {
		for j := Type(0); j < NumTypes; j++ {
			for k := Type(0); k < NumTypes; k++ {
				x := a[i][j]
				y := a[j][k]
				z := a[i][k]
				if x && y && !z {
					panic(fmt.Errorf("BUG: i=%d j=%d k=%d: i => j && j => k but !(i => k)", i, j, k))
				}
			}
		}
	}

	typeWidensTo = a
}

func init() {
	m := internal.BuildEnumMap[Type](
		NumTypes,
		typeGoNames[:],
		typeNames[:],
	)

	m["Int64"] = Type_I64
	m["Int32"] = Type_I32
	m["Int16"] = Type_I16
	m["Int8"] = Type_I8
	m["Uint64"] = Type_U64
	m["Uint32"] = Type_U32
	m["Uint16"] = Type_U16
	m["Uint8"] = Type_U8
	m["UInt64"] = Type_U64
	m["UInt32"] = Type_U32
	m["UInt16"] = Type_U16
	m["UInt8"] = Type_U8
	m["Byte"] = Type_U8
	m["Float64"] = Type_F64
	m["Float32"] = Type_F32
	m["Double"] = Type_F64
	m["Single"] = Type_F32
	m["Boolean"] = Type_Bool
	m["Char"] = Type_String
	m["VarChar"] = Type_String
	m["Text"] = Type_String
	m["Str"] = Type_String
	m["VarBinary"] = Type_Bytes
	m["Binary"] = Type_Bytes
	m["Bin"] = Type_Bytes
	m["Blob"] = Type_Bytes
	m["IP"] = Type_IPAddr
	m["ipAddr"] = Type_IPAddr

	internal.FinishEnumMap(m)
	typeMap = m
}
