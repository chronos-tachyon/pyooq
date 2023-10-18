package schema

import (
	"encoding"
	"fmt"
	"strings"
)

type Type uint

const (
	Type_Void Type = iota
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
	Type_Timestamp
	Type_String
	Type_Bytes
	Type_UUID
	Type_IPAddr
)

const typeSize = 16

var typeGoNames = [typeSize]string{
	"schema.Type_Void",
	"schema.Type_I64",
	"schema.Type_I32",
	"schema.Type_I16",
	"schema.Type_I8",
	"schema.Type_U64",
	"schema.Type_U32",
	"schema.Type_U16",
	"schema.Type_U8",
	"schema.Type_F64",
	"schema.Type_F32",
	"schema.Type_Timestamp",
	"schema.Type_String",
	"schema.Type_Bytes",
	"schema.Type_UUID",
	"schema.Type_IPAddr",
}

var typeNames = [typeSize]string{
	"void",
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
	"str",
	"bytes",
	"uuid",
	"ipAddr",
}

func (t Type) IsValid() bool {
	return t < typeSize
}

func (t Type) GoString() string {
	if t.IsValid() {
		return typeGoNames[t]
	}
	return fmt.Sprintf("schema.Type(%d)", uint(t))
}

func (t Type) String() string {
	if t.IsValid() {
		return typeNames[t]
	}
	return fmt.Sprintf("%%!ERR[invalid schema.Type %d]", uint(t))
}

func (t Type) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Type) Parse(input string) error {
	*t = ^Type(0)
	for index := Type(0); index.IsValid(); index++ {
		if strings.EqualFold(input, typeNames[index]) {
			*t = index
			return nil
		}
	}
	return fmt.Errorf("failed to parse %q as schema.Type", input)
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
