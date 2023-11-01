package pyooq

import (
	"sync"

	"github.com/chronos-tachyon/assert"

	"github.com/chronos-tachyon/pyooq/internal"
)

var gSchemaPool = sync.Pool{New: func() any { return &sbImpl{} }}

type SchemaBuilder struct{ ptr *sbImpl }
type TableBuilder struct{ ptr *sbImpl }
type IndexBuilder struct{ ptr *sbImpl }

func BuildSchema(opts ...SchemaOption) SchemaBuilder {
	ptr := gSchemaPool.Get().(*sbImpl)
	ptr.resetBuilder(sbsSchema, opts)
	return SchemaBuilder{ptr: ptr}
}

func (sb SchemaBuilder) Table(name string, opts ...TableOption) TableBuilder {
	ptr := sb.ptr.assertState(sbsSchema)
	ptr.assertNewTable(name)

	ptr.resetTable(sbsTable, name, opts)
	return TableBuilder{ptr: ptr}
}

func (tb TableBuilder) Column(name string, t Type, opts ...ColumnOption) TableBuilder {
	ptr := tb.ptr.assertState(sbsTable)
	ptr.assertNewColumn(name)
	internal.AssertValid(t)

	_ = applyOptions(opts)
	c := &Column{
		n:  name,
		ct: t,
	}
	ptr.colList = append(ptr.colList, c)
	ptr.colNames[name] = c
	return TableBuilder{ptr: ptr}
}

func (tb TableBuilder) Index(name string, opts ...IndexOption) IndexBuilder {
	ptr := tb.ptr.assertState(sbsTable)
	ptr.assertNewIndex(name)

	ptr.resetIndex(sbsIndex, name, Index_NonUnique, opts)
	return IndexBuilder{ptr: ptr}
}

func (tb TableBuilder) Unique(name string, opts ...IndexOption) IndexBuilder {
	ptr := tb.ptr.assertState(sbsTable)
	ptr.assertNewIndex(name)

	ptr.resetIndex(sbsIndex, name, Index_Unique, opts)
	return IndexBuilder{ptr: ptr}
}

func (tb TableBuilder) PrimaryKey(opts ...IndexOption) IndexBuilder {
	ptr := tb.ptr.assertState(sbsTable)
	ptr.assertNewPrimaryKey()

	ptr.resetIndex(sbsIndex, primaryKeyName, Index_PrimaryKey, opts)
	return IndexBuilder{ptr: ptr}
}

func (ib IndexBuilder) BuildIndex() TableBuilder {
	ptr := ib.ptr.assertState(sbsIndex)

	i := &Index{
		n:    ptr.idxName,
		it:   ptr.idxType,
		list: internal.CopyList(ptr.idxColList),
	}
	for _, ic := range ptr.idxColList {
		ic.i = i
	}
	ptr.idxList = append(ptr.idxList, i)
	ptr.idxNames[ptr.idxName] = i

	ptr.resetIndex(sbsTable, "", InvalidIndexType, nil)
	return TableBuilder{ptr: ptr}
}

func (ib IndexBuilder) Type(it IndexType) IndexBuilder {
	ptr := ib.ptr.assertState(sbsIndex)
	internal.AssertValid(it)

	ptr.idxType = it
	return ib
}

func (ib IndexBuilder) Add(name string, opts ...IndexedColumnOption) IndexBuilder {
	ptr := ib.ptr.assertState(sbsIndex)

	c := ptr.colNames[name]
	assert.Assertf(c != nil, "column %q does not exist", name)
	ptr.assertNewIndexColumn(c)

	o := applyOptions(opts)
	ic := &IndexedColumn{
		c:      c,
		isDesc: o.isDesc,
	}
	ptr.idxColList = append(ptr.idxColList, ic)
	ptr.idxColSet[c] = struct{}{}
	return ib
}

func (tb TableBuilder) BuildTable() SchemaBuilder {
	ptr := tb.ptr.assertState(sbsTable)

	t := &Table{
		n:        ptr.tblName,
		colList:  internal.CopyList(ptr.colList),
		colNames: internal.CopyMap(ptr.colNames),
		idxList:  internal.CopyList(ptr.idxList),
		idxNames: internal.CopyMap(ptr.idxNames),
	}
	for _, c := range ptr.colList {
		c.t = t
	}
	for _, i := range ptr.idxList {
		i.t = t
	}
	ptr.tblList = append(ptr.tblList, t)
	ptr.tblNames[ptr.tblName] = t

	ptr.resetTable(sbsSchema, "", nil)
	return SchemaBuilder{ptr: ptr}
}

func (sb SchemaBuilder) Build() *Schema {
	ptr := sb.ptr.assertState(sbsSchema)

	list := internal.CopyList(ptr.tblList)
	names := internal.CopyMap(ptr.tblNames)
	s := &Schema{list: list, names: names}

	ptr.resetBuilder(sbsIdle, nil)
	gSchemaPool.Put(sb)
	return s
}

type sbImpl struct {
	sbs      schemaBuilderState
	schOpts  SchemaOptions
	tblList  []*Table
	tblNames map[string]*Table

	tblOpts  TableOptions
	tblName  string
	colList  []*Column
	colNames map[string]*Column
	idxList  []*Index
	idxNames map[string]*Index

	idxOpts    IndexOptions
	idxName    string
	idxType    IndexType
	idxColList []*IndexedColumn
	idxColSet  map[*Column]struct{}
}

func (ptr *sbImpl) assertState(expect schemaBuilderState) *sbImpl {
	assert.NotNil(&ptr)
	internal.AssertIs("wrong state", ptr.sbs, expect)
	return ptr
}

func (ptr *sbImpl) assertNewTable(name string) {
	if internal.Contains(ptr.tblNames, name) {
		assert.Raisef("a table already exists with name %q", name)
		return
	}
	internal.AssertValidName("table", name)
}

func (ptr *sbImpl) assertNewColumn(name string) {
	if internal.Contains(ptr.colNames, name) {
		assert.Raisef("a column already exists with name %q", name)
		return
	}
	internal.AssertValidName("column", name)
}

func (ptr *sbImpl) assertNewPrimaryKey() {
	if internal.Contains(ptr.idxNames, primaryKeyName) {
		assert.Raise("a primary key already exists")
		return
	}
}

func (ptr *sbImpl) assertNewIndex(name string) {
	if internal.Contains(ptr.idxNames, name) {
		assert.Raisef("an index already exists with name %q", name)
		return
	}
	internal.AssertValidName("index", name)
}

func (ptr *sbImpl) assertNewIndexColumn(c *Column) {
	if internal.Contains(ptr.idxColSet, c) {
		assert.Raisef("column %q is already part of the current index", c.Name())
		return
	}
}

func (ptr *sbImpl) resetBuilder(sbs schemaBuilderState, opts []SchemaOption) {
	ptr.sbs = sbs
	ptr.schOpts = applyOptions(opts)
	internal.RecycleList(&ptr.tblList, 32)
	internal.RecycleMap(&ptr.tblNames, 32)
}

func (ptr *sbImpl) resetTable(sbs schemaBuilderState, name string, opts []TableOption) {
	ptr.sbs = sbs
	ptr.tblOpts = applyOptions(opts)
	ptr.tblName = name
	internal.RecycleList(&ptr.colList, 16)
	internal.RecycleMap(&ptr.colNames, 16)
	internal.RecycleList(&ptr.idxList, 4)
	internal.RecycleMap(&ptr.idxNames, 4)
}

func (ptr *sbImpl) resetIndex(sbs schemaBuilderState, name string, it IndexType, opts []IndexOption) {
	ptr.sbs = sbs
	ptr.idxOpts = applyOptions(opts)
	ptr.idxName = name
	ptr.idxType = it
	internal.RecycleList(&ptr.idxColList, 4)
	internal.RecycleMap(&ptr.idxColSet, 4)
}
