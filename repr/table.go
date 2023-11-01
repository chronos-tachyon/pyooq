package repr

import (
	"fmt"
)

type Table struct {
	Name    string    `json:"name"`
	Columns []*Column `json:"columns"`
	Indices []*Index  `json:"indices"`
}

func (t Table) Append(out []byte) []byte {
	out = append(out, t.Name...)
	out = append(out, '(')
	needComma := false
	for _, c := range t.Columns {
		if needComma {
			out = append(out, ',', ' ')
		}
		out = c.Append(out)
		needComma = true
	}
	for _, i := range t.Indices {
		if needComma {
			out = append(out, ',', ' ')
		}
		out = i.Append(out)
		needComma = true
	}
	out = append(out, ')')
	return out
}

func (t Table) String() string {
	var scratch [256]byte
	return string(t.Append(scratch[:0]))
}

var _ fmt.Stringer = Table{}
