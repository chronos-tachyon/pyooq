package types

import (
	"encoding"
	"fmt"
	"strings"
)

type Type uint

const (
	I64 Type = iota
	I32
	I16
	I8
	U64
	U32
	U16
	U8
	F64
	F32
	Timestamp
	Text
	Binary
	UUID
	IPAddr
)

const typeSize = 15

var typeGoNames = [typeSize]string{
	"types.I64",
	"types.I32",
	"types.I16",
	"types.I8",
	"types.U64",
	"types.U32",
	"types.U16",
	"types.U8",
	"types.F64",
	"types.F32",
	"types.Timestamp",
	"types.Text",
	"types.Binary",
	"types.UUID",
	"types.IPAddr",
}

var typeNames = [typeSize]string{
	"i64",
	"i32",
	"i16",
	"i8",
	"u64",
	"u32",
	"u16",
	"u8",
	"f64",
	"f32",
	"timestamp",
	"text",
	"binary",
	"uuid",
	"ipAddr",
}

var typeMap = map[string]Type{
	"i64":       I64,
	"int64":     I64,
	"i32":       I32,
	"int32":     I32,
	"i16":       I16,
	"int16":     I16,
	"i8":        I8,
	"int8":      I8,
	"u64":       U64,
	"uint64":    U64,
	"u32":       U32,
	"uint32":    U32,
	"u16":       U16,
	"uint16":    U16,
	"u8":        U8,
	"uint8":     U8,
	"f64":       F64,
	"float64":   F64,
	"f32":       F32,
	"float32":   F32,
	"timestamp": Timestamp,
	"text":      Text,
	"string":    Text,
	"str":       Text,
	"binary":    Binary,
	"bin":       Binary,
	"bytes":     Binary,
	"uuid":      UUID,
	"ipAddr":    IPAddr,
	"ip":        IPAddr,
}

func (t Type) IsValid() bool {
	return t < typeSize
}

func (t Type) GoString() string {
	if t.IsValid() {
		return typeGoNames[t]
	}
	return fmt.Sprintf("types.Type(%d)", uint(t))
}

func (t Type) String() string {
	if t.IsValid() {
		return typeNames[t]
	}
	return fmt.Sprintf("%%!ERR[invalid types.Type %d]", uint(t))
}

func (t Type) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Type) Parse(input string) error {
	*t = ^Type(0)
	if x, found := typeMap[input]; found {
		*t = x
		return nil
	}
	lc := strings.ToLower(input)
	if x, found := typeMap[lc]; found {
		*t = x
		return nil
	}
	return fmt.Errorf("failed to parse %q as types.Type", input)
}

func (t *Type) UnmarshalText(input []byte) error {
	return t.Parse(string(input))
}

var (
	_ fmt.GoStringer           = Type(0)
	_ fmt.Stringer             = Type(0)
	_ encoding.TextMarshaler   = Type(0)
	_ encoding.TextUnmarshaler = (*Type)(nil)
)
