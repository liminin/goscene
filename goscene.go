package goscene

type SceneManager struct {
	scenes map[string]*Scene
	store  Store
}

func NewSceneManager(store Store) *SceneManager {
	return &SceneManager{
		store:  store,
		scenes: map[string]*Scene{},
	}
}

func (s *SceneManager) AddScene(scene *Scene) {
	scene.store = s.store
	s.scenes[scene.Key] = scene
}

func (s *SceneManager) Play(key string, userID int, data interface{}) (result interface{}, err error) {
	scene, ok := s.scenes[key]

	if !ok {
		err = errSceneNotFount
		return
	}

	err = scene.Play(userID)

	if err != nil {
		return
	}

	play, _ := s.GetUserActivePlay(userID)

	result = play.Execute(data)

	return
}

func (s *SceneManager) GetUserActivePlay(userID int) (play *ScenePlay, ok bool) {
	pi, err := s.store.Info().GetByUserID(userID)
	if err != nil {
		return
	}

	play = &ScenePlay{
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
