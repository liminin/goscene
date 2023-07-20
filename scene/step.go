package scene

import "github.com/liminin/goscene/state"

type Play[T any] interface {
	FirstTime() bool
	State() *state.State
	SetData(data T)
	Execute() error
	Next() error
	Back() error
	Go(i int) error
	Exit() error
}

type Handler[T any] func(play Play[T], data T) error

type Step[T any] struct {
	h Handler[T]
}

func NewStep[T any](h Handler[T]) Step[T] {
	return Step[T]{
		h: h,
	}
}

func (s *Step[T]) Run(play Play[T], data T) error {
	return s.h(play, data)
}
