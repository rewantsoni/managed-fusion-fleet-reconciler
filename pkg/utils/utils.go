package utils

type Pair[L, R any] struct {
	left  L
	right R
}

func NewPair[L, R any](l L, r R) Pair[L, R] {
	return Pair[L, R]{left: l, right: r}
}

func (p *Pair[L, R]) Unpack() (L, R) {
	return p.left, p.right
}

func ToPointer[T any](v T) *T {
	return &v
}
