package scene

type Play[T any] interface {
}

type Handler[T any] func(play Play[T], data T) error

type Step[T any] struct {
	h Handler[T]
}

func NewStep[T any](h Handler[T]) *Step[T] {
	return &Step[T]{
		h: h,
	}
}

func (s *Step[T]) Run(play Play[T], data T) error {
	return s.h(play, data)
}
