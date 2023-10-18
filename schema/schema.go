package schema

import (
	"fmt"
)

type Schema struct {
	Tables []*Table `json:"tables"`
}

func (s *Schema) Append(out []byte) []byte {
	out = append(out, '[')
	for index, t := range s.Tables {
		if index > 0 {
			out = append(out, ';', ' ')
		}
		out = t.Append(out)
	}
	out = append(out, ']')
	return out
}

func (s *Schema) String() string {
	var scratch [1024]byte
	return string(s.Append(scratch[:0]))
}

var _ fmt.Stringer = (*Schema)(nil)
