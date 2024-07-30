package monad

type Tupple[T, U any] struct {
	Item1 *T
	Item2 *U
}

type Tupple3[T, U, V any] struct {
	Item1 *T
	Item2 *U
	Item3 *V
}

func NewTupple[T, U any](item1 *T, item2 *U) *Tupple[T, U] {
	return &Tupple[T, U]{
		Item1: item1,
		Item2: item2,
	}
}

func NewTupple3[T, U, V any](item1 *T, item2 *U, item3 *V) *Tupple3[T, U, V] {
	return &Tupple3[T, U, V]{
		Item1: item1,
		Item2: item2,
		Item3: item3,
	}
}
