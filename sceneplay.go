package goscene

type ScenePlayInfo struct {
	ID           int
	UserID       int    `mapstructure:"user_id"`
	SceneKey     string `mapstructure:"scene_key"`
	CurrentIndex int    `mapstructure:"current"`
	FirstTime    bool   `mapstructure:"first_time"`
}

type ScenePlay struct {
	ID           int
	UserID       int
	FirstTime    bool
	store        Store
	scene        Scene
	currentIndex int
	lastData     interface{}
}

func (s *ScenePlay) Get(key string) (v interface{}, ok bool) {
	v, err := s.store.State().Get(s.ID, key)

	if err == nil {
		ok = true
	}

	return
}

func (s *ScenePlay) Set(key string, value interface{}) (ok bool) {
	err := s.store.State().Set(s.ID, key, value)

	if err == nil {
		ok = true
	}

	return
}

func (s *ScenePlay) Execute(data interface{}) interface{} {
	if s.FirstTime {
		s.store.Info().Set(s.ID, firstTimeKey, false)
	}

	s.lastData = data

	return s.scene.handlers[s.currentIndex](s, data)
}

func (s *ScenePlay) Next() interface{} {
	data, _ := s.Go(s.currentIndex + 1)
	return data
}

func (s *ScenePlay) Back() interface{} {
	data, _ := s.Go(s.currentIndex - 1)
	return data
}

func (s *ScenePlay) Go(i int) (interface{}, error) {
	if i < 0 || i >= len(s.scene.handlers) {
		s.Exit()
		return nil, errInvalidCurrentIndex
	}

	s.currentIndex = i
	s.store.Info().Set(s.ID, currentIndexKey, i)
	s.FirstTime = true
	s.store.Info().Set(s.ID, firstTimeKey, true)

	return s.Execute(s.lastData), nil
}

func (s *ScenePlay) Exit() {
	s.store.Play().End(s.ID)
}
