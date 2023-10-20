package schema

type TablePredicate = Predicate[Table]
type TableConsumer = Consumer[Table]

type Table interface {
	Namer
	NumColumns() uint
	ForEachColumn(consumer ColumnConsumer)
}

func MakeTable(name string, columns ...Column) Table {
	dupe := make([]Column, len(columns))
	copy(dupe, columns)
	return &genericTable{name: name, columns: dupe}
}

type genericTable struct {
	name    string
	columns []Column
}

func (t *genericTable) Name() string {
	return t.name
}

func (t *genericTable) NumColumns() uint {
	return uint(len(t.columns))
}

func (t *genericTable) ForEachColumn(consumer ColumnConsumer) {
	for _, c := range t.columns {
		consumer.Consume(c)
	}
}

var _ Table = (*genericTable)(nil)

func Columns(t Table, preds ...ColumnPredicate) []Column {
	out := make([]Column, 0, t.NumColumns())
	t.ForEachColumn(AppendConsumer(&out, preds...))
	return out
}
