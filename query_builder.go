package pyooq

import (
	"sync"

	"github.com/chronos-tachyon/assert"

	"github.com/chronos-tachyon/pyooq/internal"
)

var gQueryPool = sync.Pool{New: func() any { return &qbImpl{} }}

type QueryBuilder struct{ ptr *qbImpl }

func NoFrom() QueryBuilder {
	ptr := gQueryPool.Get().(*qbImpl)
	ptr.Reset()
	ptr.qbs = qbsAfterJoin
	return QueryBuilder{ptr: ptr}
}

func From(table TableExpr) QueryBuilder {
	internal.AssertValid(table)
	assertKind(table, Kind_Table)

	ptr := gQueryPool.Get().(*qbImpl)
	ptr.Reset()
	ptr.qbs = qbsJoin
	ptr.fromFirst = table
	ptr.tableNames[table.Name()] = table
	return QueryBuilder{ptr: ptr}
}

func (qb QueryBuilder) Join(join Join, table TableExpr, cond ValueExpr) QueryBuilder {
	ptr := qb.ptr.assertState(qbsJoin)
	internal.AssertValid(join)
	internal.AssertValid(table)
	assertKind(table, Kind_Table)
	if cond != nil {
		assertType(cond, Type_Bool)
	}

	name := table.Name()
	old, found := ptr.tableNames[name]
	assert.Assertf(!found, "%#v with name %q collides with an existing table, %#v", table, name, old)

	ptr.fromRest = append(ptr.fromRest, QueryJoin{Type: join, Table: table, Condition: cond})
	ptr.tableNames[name] = table
	return qb
}

func (qb QueryBuilder) InnerJoin(table TableExpr, cond ValueExpr) QueryBuilder {
	return qb.Join(Join_Inner, table, cond)
}

func (qb QueryBuilder) OuterJoin(table TableExpr, cond ValueExpr) QueryBuilder {
	return qb.Join(Join_Outer, table, cond)
}

func (qb QueryBuilder) Select(selections ...ValueExpr) QueryBuilder {
	ptr := qb.ptr.assertStateIn(qbsJoin, qbsAfterJoin)

	for _, item := range selections {
		switch item.Kind() {
		case Kind_Splat:
			ptr.selections = append(ptr.selections, QueryColumn{})

		default:
			v := assertValue(item)
			ptr.selections = append(ptr.selections, QueryColumn{Value: v})
		}
	}

	ptr.qbs = qbsAfterSelect
	return qb
}

func (qb QueryBuilder) Where(where ValueExpr) QueryBuilder {
	ptr := qb.ptr.assertState(qbsAfterSelect)
	if where != nil {
		assertType(where, Type_Bool)
	}

	ptr.where = where
	ptr.qbs = qbsAfterWhere
	return qb
}

func (qb QueryBuilder) Build() *Query {
	ptr := qb.ptr.assertStateIn(qbsAfterSelect, qbsAfterWhere)

	q := &Query{
		from:  makeQueryFrom(ptr.fromFirst, ptr.fromRest),
		sel:   makeQuerySelect(ptr.selections),
		where: ptr.where,
	}

	ptr.Reset()
	gQueryPool.Put(ptr)
	return q
}

type qbImpl struct {
	qbs queryBuilderState

	fromFirst  TableExpr
	fromRest   []QueryJoin
	tableNames map[string]TableExpr

	selections []QueryColumn

	where ValueExpr
}

func (ptr *qbImpl) Reset() {
	ptr.qbs = qbsIdle
	ptr.fromFirst = TableExpr{}
	ptr.where = nil
	internal.RecycleList(&ptr.fromRest, 4)
	internal.RecycleMap(&ptr.tableNames, 4)
	internal.RecycleList(&ptr.selections, 32)
}

func (ptr *qbImpl) assertState(expect queryBuilderState) *qbImpl {
	assert.NotNil(&ptr)
	internal.AssertIs("wrong state", ptr.qbs, expect)
	return ptr
}

func (ptr *qbImpl) assertStateIn(list ...queryBuilderState) *qbImpl {
	assert.NotNil(&ptr)
	internal.AssertOneOf("wrong state", ptr.qbs, list...)
	return ptr
}
