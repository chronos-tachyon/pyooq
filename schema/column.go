package schema

import (
	"fmt"
)

type Column struct {
	Name string `json:"name"`
	Type Type   `json:"type"`
}

func (c Column) Append(out []byte) []byte {
	out = append(out, c.Name...)
	out = append(out, ':', ' ')
	out = append(out, c.Type.String()...)
	return out
}

func (c Column) String() string {
	var scratch [64]byte
	return string(c.Append(scratch[:0]))
}

var _ fmt.Stringer = Column{}
