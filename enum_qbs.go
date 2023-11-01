package pyooq

import (
	"fmt"
)

type queryBuilderState byte

const (
	qbsIdle queryBuilderState = iota
	qbsJoin
	qbsAfterJoin
	qbsAfterSelect
	qbsAfterWhere
)

const qbsSize = 5

var qbsName = [qbsSize]string{
	"qbsIdle",
	"qbsJoin",
	"qbsAfterJoin",
	"qbsAfterSelect",
	"qbsAfterWhere",
}

func (qbs queryBuilderState) IsValid() bool {
	return qbs < qbsSize
}
func (qbs queryBuilderState) GoString() string {
	if qbs.IsValid() {
		return qbsName[qbs]
	}
	return fmt.Sprintf("queryBuilderState(%d)", uint(qbs))
}
func (qbs queryBuilderState) String() string {
	return qbs.GoString()
}

var (
	_ fmt.GoStringer = queryBuilderState(0)
	_ fmt.Stringer   = queryBuilderState(0)
)
