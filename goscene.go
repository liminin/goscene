package goscene

type SceneManager[T any] struct {
	scenes map[string]*Scene[T]
	store  Store
}

func NewSceneManager[T any](store Store) *SceneManager[T] {
	return &SceneManager[T]{
		store:  store,
		scenes: map[string]*Scene[T]{},
	}
}

func (s *SceneManager[T]) AddScene(scene *Scene[T]) {
	scene.store = s.store
	s.scenes[scene.Key] = scene
}

func (s *SceneManager[T]) Play(key string, userID int, data T) (err error) {
	scene, ok := s.scenes[key]

	if !ok {
		err = errSceneNotFount
		return
	}

	previousPlay, ok := s.GetUserActivePlay(userID)
	if ok {
		s.store.Info().Set(previousPlay.ID, firstTimeKey, true)
	}

	err = scene.Play(userID)

	if err != nil {
		return
	}

	play, _ := s.GetUserActivePlay(userID)

	play.SetData(data)
	err = play.Execute()

	return
}

func (s *SceneManager[T]) GetUserActivePlay(userID int) (play *ScenePlay[T], ok bool) {
	pi, err := s.store.Info().GetByUserID(userID)
	if err != nil {
		return
	}

	play = &ScenePlay[T]{
		ID:           pi.ID,
		UserID:       userID,
		store:        s.store,
		scene:        *s.scenes[pi.SceneKey],
		currentIndex: pi.CurrentIndex,
		FirstTime:    pi.FirstTime,
		sm:           s,
	}
	ok = true

	return
}
