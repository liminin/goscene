package goscene

type Handler[T any] func(play *ScenePlay[T], data T) error

type Scene[T any] struct {
	Key      string
	handlers []Handler[T]
	store    Store
}

func NewScene[T any](key string) *Scene[T] {
	return &Scene[T]{
		Key:      key,
		handlers: []Handler[T]{},
	}
}

func (s *Scene[T]) AddHandler(h Handler[T]) {
	s.handlers = append(s.handlers, h)
}

func (s *Scene[T]) Play(userID int) (err error) {
	err = s.store.Play().New(s.Key, userID)

	return
}
