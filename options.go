package pyooq

type Option[T any] interface {
	ApplyTo(*T)
}

type baseOptions struct {
	dummy struct{}
}

type SchemaOptions struct {
	baseOptions
}
type SchemaOption = Option[SchemaOptions]

type TableOptions struct {
	baseOptions
}
type TableOption = Option[TableOptions]

type ColumnOptions struct {
	baseOptions
}
type ColumnOption = Option[ColumnOptions]

type IndexOptions struct {
	baseOptions
}
type IndexOption = Option[IndexOptions]

type IndexedColumnOptions struct {
	baseOptions
	isDesc bool
}
type IndexedColumnOption = Option[IndexedColumnOptions]

type Desc bool

func (desc Desc) ApplyTo(o *IndexedColumnOptions) {
	o.isDesc = bool(desc)
}

var _ IndexedColumnOption = Desc(true)

func applyOptions[T any](opts []Option[T]) T {
	var o T
	for _, opt := range opts {
		opt.ApplyTo(&o)
	}
	return o
}
