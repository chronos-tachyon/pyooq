package repr

import (
	"fmt"

	"github.com/chronos-tachyon/pyooq/types"
)

type Column struct {
	Name string     `json:"name"`
	Type types.Type `json:"type"`
}

func (c Column) Append(out []byte) []byte {
	return fmt.Appendf(out, "%s: %v", c.Name, c.Type)
}

func (c Column) String() string {
	var scratch [64]byte
	return string(c.Append(scratch[:0]))
}

var _ fmt.Stringer = Column{}
