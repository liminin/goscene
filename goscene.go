package goscene

import (
	"errors"

	"github.com/liminin/goscene/play"
	"github.com/liminin/goscene/scene"
	"github.com/liminin/goscene/store"
)

type SceneManager[T any] struct {
	scenes map[string]*scene.Scene[T]
	store  store.Store
}

func NewSceneManager[T any](store store.Store) *SceneManager[T] {
	return &SceneManager[T]{
		store:  store,
		scenes: map[string]*scene.Scene[T]{},
	}
}

func (s *SceneManager[T]) AddScene(scene *scene.Scene[T]) {
	s.scenes[scene.Key] = scene
}

func (s *SceneManager[T]) Play(key string, userID int, data T) (err error) {
	err = s.startNewUserPlay(userID, key)
	if err != nil {
		return
	}

	return s.execute(userID, data)
}

func (s *SceneManager[T]) execute(userID int, data T) error {
	play, ok := s.GetUserActivePlay(userID)
	if !ok {
		return errors.New("user does not have active play")
	}

	play.SetData(data)

	return play.Execute()
}

func (s *SceneManager[T]) startNewUserPlay(userID int, sceneKey string) error {
	_, ok := s.scenes[sceneKey]

	if !ok {
		return errSceneNotFount
	}

	return s.store.Play().New(sceneKey, userID)
}

func (s *SceneManager[T]) GetUserActivePlay(userID int) (*play.Play[T], bool) {
	model, err := s.store.Play().GetByUserID(userID)
	if err != nil {
		return nil, false
	}

	p := play.NewPlay[T](model, *s.scenes[model.SceneKey], s.store)

	return p, true
}
