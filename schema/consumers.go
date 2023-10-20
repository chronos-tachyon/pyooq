package schema

type Consumer[T any] interface {
	Consume(T)
}

type FuncConsumer[T any] func(T)

func (fn FuncConsumer[T]) Consume(value T) {
	fn(value)
}

var _ Consumer[Table] = FuncConsumer[Table](nil)

func AppendConsumer[T any](ptr *[]T, preds ...Predicate[T]) Consumer[T] {
	return appendConsumerImpl[T]{ptr: ptr, preds: preds}
}

type appendConsumerImpl[T any] struct {
	ptr   *[]T
	preds []Predicate[T]
}

func (c appendConsumerImpl[T]) Consume(value T) {
	for _, pred := range c.preds {
		if !pred.Matches(value) {
			return
		}
	}
	*c.ptr = append(*c.ptr, value)
}

var _ Consumer[Table] = appendConsumerImpl[Table]{}
