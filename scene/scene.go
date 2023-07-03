package scene

import "errors"

type Scene[T any] struct {
	Key   string
	steps []Step[T]
}

func NewScene[T any](key string) *Scene[T] {
	return &Scene[T]{
		Key:   key,
		steps: []Step[T]{},
	}
}

func (s *Scene[T]) AddSteps(steps ...Step[T]) {
	s.steps = append(s.steps, steps...)
}

func (s *Scene[T]) GetStep(index int) (step Step[T], err error) {
	if len(s.steps) < index+1 {
		return Step[T]{}, errors.New("out of range")
	}

	return
}

// Count returns the number of steps
func (s *Scene[T]) Count() int {
	return len(s.steps)
}

// Count returns  the index of the last step
func (s *Scene[T]) LastStepIndex() int {
	return len(s.steps) - 1
}
