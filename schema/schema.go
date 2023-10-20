package schema

type Schema interface {
	NumTables() uint
	ForEachTable(consumer Consumer[Table])
}

func MakeSchema(tables ...Table) Schema {
	dupe := make([]Table, len(tables))
	copy(dupe, tables)
	return &genericSchema{tables: tables}
}

type genericSchema struct {
	tables []Table
}

func (s *genericSchema) NumTables() uint {
	return uint(len(s.tables))
}

func (s *genericSchema) ForEachTable(consumer TableConsumer) {
	for _, t := range s.tables {
		consumer.Consume(t)
	}
}

var _ Schema = (*genericSchema)(nil)

func Tables(s Schema, preds ...TablePredicate) []Table {
	out := make([]Table, 0, s.NumTables())
	s.ForEachTable(AppendConsumer(&out, preds...))
	return out
}
