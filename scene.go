package goscene

type Handler func(play *ScenePlay, data interface{}) interface{}

type Scene struct {
	Key      string
	handlers []Handler
	store    Store
}

func NewScene(key string) *Scene {
	return &Scene{
		Key:      key,
		handlers: []Handler{},
	}
}

func (s *Scene) AddHandler(h Handler) {
	s.handlers = append(s.handlers, h)
}

func (s *Scene) Play(userID int) (err error) {
	err = s.store.Play().New(s.Key, userID)

	return
}
