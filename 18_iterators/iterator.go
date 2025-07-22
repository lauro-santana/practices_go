package iterators

type Iterator[T any] struct {
	Slice []T
	index int
}

func (i *Iterator[T]) FilterIterator(filter func(T) bool) *T {
	if filter(i.Slice[i.index]) {
		return &i.Slice[i.index]
	}
	return nil
}

func (i *Iterator[E]) Next() bool {
	i.index++
	return i.index < len(i.Slice)
}

func FilterYield[Slice ~[]E, E any](s Slice,
	filter func(E) bool) func(yield func(E) bool) {
	return func(yield func(E) bool) {
		for _, v := range s {
			if filter(v) && !yield(v) {
				return
			}
		}
	}
}

func FilterSlice[E any](s []E, filter func(E) bool) []E {
	sl := make([]E, 0, len(s))
	for _, v := range s {
		if filter(v) {
			sl = append(sl, v)
		}
	}
	return sl
}
