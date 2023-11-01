package pyooq

import (
	"encoding"
	"fmt"
	"strings"

	"github.com/chronos-tachyon/pyooq/internal"
)

type IndexType byte

const (
	Index_NonUnique IndexType = iota
	Index_Unique
	Index_PrimaryKey
)

const NumIndexTypes IndexType = 3
const InvalidIndexType = ^IndexType(0)

var indexTypeGoNames = [NumIndexTypes]string{
	"pyooq.Index_NonUnique",
	"pyooq.Index_Unique",
	"pyooq.Index_PrimaryKey",
}

var indexTypeNames = [NumIndexTypes]string{
	"nonUniqueIndex",
	"uniqueIndex",
	"primaryKey",
}

var indexTypeMap map[string]IndexType

func (enum IndexType) IsValid() bool {
	return enum < NumIndexTypes
}
func (enum IndexType) IsUnique() bool {
	return enum == Index_Unique || enum == Index_PrimaryKey
}
func (enum IndexType) GoString() string {
	if enum.IsValid() {
		return indexTypeGoNames[enum]
	}
	return fmt.Sprintf("pyooq.IndexType(%d)", uint(enum))
}
func (enum IndexType) String() string {
	if enum.IsValid() {
		return indexTypeNames[enum]
	}
	return fmt.Sprintf("%%!ERR[invalid pyooq.IndexType %d]", uint(enum))
}
func (enum IndexType) GoAppend(out []byte) []byte {
	return appendString(out, enum.GoString())
}
func (enum IndexType) Append(out []byte) []byte {
	return appendString(out, enum.String())
}
func (enum IndexType) MarshalText() ([]byte, error) {
	return []byte(enum.String()), nil
}
func (enum *IndexType) Parse(input string) error {
	if itype, found := indexTypeMap[input]; found {
		*enum = itype
		return nil
	}
	if itype, found := indexTypeMap[strings.ToLower(input)]; found {
		*enum = itype
		return nil
	}
	*enum = InvalidIndexType
	return fmt.Errorf("failed to parse %q as pyooq.IndexType", input)
}
func (enum *IndexType) UnmarshalText(input []byte) error {
	return enum.Parse(string(input))
}

var (
	_ Stringer                 = IndexType(0)
	_ encoding.TextMarshaler   = IndexType(0)
	_ encoding.TextUnmarshaler = (*IndexType)(nil)
)

func init() {
	m := internal.BuildEnumMap[IndexType](
		NumIndexTypes,
		indexTypeGoNames[:],
		indexTypeNames[:],
	)

	m["Index"] = Index_NonUnique
	m["Key"] = Index_NonUnique
	m["NonUnique"] = Index_NonUnique
	m["NonUniqueIndex"] = Index_NonUnique
	m["NonUniqueKey"] = Index_NonUnique
	m["nonUnique"] = Index_NonUnique
	//m["nonUniqueIndex"] = Index_NonUnique
	m["nonUniqueKey"] = Index_NonUnique
	m["Non Unique"] = Index_NonUnique
	m["Non Unique Index"] = Index_NonUnique
	m["Non Unique Key"] = Index_NonUnique
	m["non Unique"] = Index_NonUnique
	m["non Unique Index"] = Index_NonUnique
	m["non Unique Key"] = Index_NonUnique

	m["Unique"] = Index_Unique
	m["UniqueIndex"] = Index_Unique
	m["UniqueKey"] = Index_Unique
	//m["uniqueIndex"] = Index_Unique
	m["uniqueKey"] = Index_Unique
	m["Unique Index"] = Index_Unique
	m["Unique Key"] = Index_Unique
	m["unique Index"] = Index_Unique
	m["unique Key"] = Index_Unique

	m["Primary"] = Index_PrimaryKey
	m["PrimaryKey"] = Index_PrimaryKey
	m["PrimaryIndex"] = Index_PrimaryKey
	//m["primaryKey"] = Index_PrimaryKey
	m["primaryIndex"] = Index_PrimaryKey
	m["Primary Key"] = Index_PrimaryKey
	m["Primary Index"] = Index_PrimaryKey
	m["primary Key"] = Index_PrimaryKey
	m["primary Index"] = Index_PrimaryKey

	internal.FinishEnumMap(m)
	indexTypeMap = m
}
