package repr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chronos-tachyon/pyooq"
)

type Index struct {
	Name    string           `json:"name"`
	Type    pyooq.IndexType  `json:"type"`
	Columns []*IndexedColumn `json:"columns"`
}

func (i Index) Append(out []byte) []byte {
	return fmt.Appendf(out, "%v/%s()", i.Type, i.Name)
}

func (i Index) String() string {
	var scratch [64]byte
	return string(i.Append(scratch[:0]))
}

var _ fmt.Stringer = Index{}

type IndexedColumn struct {
	Column string
	IsDesc bool
}

func (ic IndexedColumn) MarshalJSON() ([]byte, error) {
	var tmp icJSON
	tmp.From(ic)
	return json.Marshal(&tmp)
}

func (ic *IndexedColumn) Parse(input string) error {
	*ic = IndexedColumn{}
	if strings.HasSuffix(input, ":desc") {
		ic.Column = input[:len(input)-5]
		ic.IsDesc = true
		return nil
	}
	if strings.HasSuffix(input, ":asc") {
		input = input[:len(input)-4]
	}
	ic.Column = input
	return nil
}

func (ic *IndexedColumn) UnmarshalJSON(input []byte) error {
	*ic = IndexedColumn{}

	if len(input) >= 2 && input[0] == '"' {
		var tmp string
		if err := json.Unmarshal(input, &tmp); err != nil {
			return err
		}
		return ic.Parse(tmp)
	}

	var tmp icJSON
	d := json.NewDecoder(bytes.NewReader(input))
	d.DisallowUnknownFields()
	d.UseNumber()
	if err := d.Decode(&tmp); err != nil {
		return err
	}
	tmp.To(ic)
	return nil
}

type icJSON struct {
	Column string `json:"column"`
	IsDesc bool   `json:"is_desc"`
}

func (x *icJSON) From(in IndexedColumn) {
	*x = icJSON{
		Column: in.Column,
		IsDesc: in.IsDesc,
	}
}

func (x icJSON) To(out *IndexedColumn) {
	*out = IndexedColumn{
		Column: x.Column,
		IsDesc: x.IsDesc,
	}
}

var (
	_ json.Marshaler   = IndexedColumn{}
	_ json.Marshaler   = (*IndexedColumn)(nil)
	_ json.Unmarshaler = (*IndexedColumn)(nil)
)
