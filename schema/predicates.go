package schema

type Predicate[T any] interface {
	Matches(T) bool
}

type Namer interface{ Name() string }
type namerPredicate[T Namer] struct{ n string }

func (pred *namerPredicate[T]) Matches(value T) bool {
	n := value.Name()
	return n == pred.n
}

func Named[T Namer](n string) Predicate[T] { return &namerPredicate[T]{n: n} }

type Typer interface{ Type() Type }
type typerPredicate[T Typer] struct{ t Type }

func (pred *typerPredicate[T]) Matches(value T) bool {
	t := value.Type()
	return t == pred.t
}

func Typed[T Typer](t Type) Predicate[T] { return &typerPredicate[T]{t: t} }
