package schema

import (
	"fmt"
)

type Table struct {
	Name    string    `json:"name"`
	Columns []*Column `json:"columns"`
}

func (t Table) Append(out []byte) []byte {
	out = append(out, t.Name...)
	out = append(out, '(')
	for index, c := range t.Columns {
		if index > 0 {
			out = append(out, ',', ' ')
		}
		out = c.Append(out)
	}
	out = append(out, ')')
	return out
}

func (t Table) String() string {
	var scratch [256]byte
	return string(t.Append(scratch[:0]))
}

var _ fmt.Stringer = Table{}
