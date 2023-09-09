package iteratorhelper

type Builder[K, V any] struct {
	iter func(yield func(k K, v V) bool)
}

func NewBuilder[K, V any](iter func(yield func(k K, v V) bool)) *Builder[K, V] {
	return &Builder[K, V]{
		iter: iter,
	}
}
