package goscene

type ScenePlayInfo struct {
	ID           int
	UserID       int    `mapstructure:"user_id"`
	SceneKey     string `mapstructure:"scene_key"`
	CurrentIndex int    `mapstructure:"current"`
	FirstTime    bool   `mapstructure:"first_time"`
}

type ScenePlay[T any] struct {
	ID           int
	UserID       int
	FirstTime    bool
	store        Store
	sm           *SceneManager[T]
	scene        Scene[T]
	currentIndex int
	data         T
}

func (s *ScenePlay[T]) Get(key string) *StateCmd {
	return NewStateCmd(s.store.State().Get(s.ID, key))
}

func (s *ScenePlay[T]) Set(key string, value interface{}) (ok bool) {
	err := s.store.State().Set(s.ID, key, value)

	if err == nil {
		ok = true
	}

	return
}

func (s *ScenePlay[T]) SetData(data T) {
	s.data = data
}

func (s *ScenePlay[T]) Execute() error {
	if s.FirstTime {
		s.store.Info().Set(s.ID, firstTimeKey, false)
	}

	return s.scene.handlers[s.currentIndex](s, s.data)
}

func (s *ScenePlay[T]) Next() error {
	if s.currentIndex >= len(s.scene.handlers)-1 {
		s.Exit()

		return nil
	}

	return s.Go(s.currentIndex + 1)
}

func (s *ScenePlay[T]) Back() error {
	if s.currentIndex == 0 {
		s.Exit()

		return nil
	}

	return s.Go(s.currentIndex - 1)
}

func (s *ScenePlay[T]) Go(i int) error {
	if i < 0 || i >= len(s.scene.handlers) {
		return errInvalidCurrentIndex
	}

	s.currentIndex = i
	s.store.Info().Set(s.ID, currentIndexKey, i)
	s.FirstTime = true
	s.store.Info().Set(s.ID, firstTimeKey, true)

	return s.Execute()
}

func (s *ScenePlay[T]) Exit() error {
	s.store.Play().End(s.ID)

	p, isFound := s.sm.GetUserActivePlay(s.UserID)
	if isFound {
		p.FirstTime = true
		return p.Execute()
	}

	return nil
}
