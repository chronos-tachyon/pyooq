package pyooq

import (
	"fmt"
)

type schemaBuilderState byte

const (
	sbsIdle schemaBuilderState = iota
	sbsSchema
	sbsTable
	sbsColumn
	sbsIndex
)

const sbsSize = 5

var sbsNames = [sbsSize]string{
	"sbsIdle",
	"sbsSchema",
	"sbsTable",
	"sbsColumn",
	"sbsIndex",
}

func (sbs schemaBuilderState) IsValid() bool {
	return sbs < sbsSize
}
func (sbs schemaBuilderState) GoString() string {
	if sbs.IsValid() {
		return sbsNames[sbs]
	}
	return fmt.Sprintf("schemaBuilderState(%d)", uint(sbs))
}
func (sbs schemaBuilderState) String() string {
	return sbs.String()
}

var (
	_ fmt.GoStringer = schemaBuilderState(0)
	_ fmt.Stringer   = schemaBuilderState(0)
)
