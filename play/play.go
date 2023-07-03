package play

import (
	"errors"

	"github.com/liminin/goscene/model"
	"github.com/liminin/goscene/scene"
	"github.com/liminin/goscene/state"
	"github.com/liminin/goscene/store"
)

var (
	errInvalidCurrentIndex = errors.New("invalid current index")
)

const (
	firstTimeKey    = "first_time"
	currentIndexKey = "current_index"
)

type Play[T any] struct {
	ID           int
	UserID       int
	firstTime    bool
	store        store.Store
	state        *state.State
	scene        scene.Scene[T]
	currentIndex int
	data         T
}

func NewPlay[T any](m *model.Play, scene scene.Scene[T], store store.Store) *Play[T] {
	return &Play[T]{
		ID:           m.ID,
		UserID:       m.UserID,
		firstTime:    m.FirstTime,
		currentIndex: m.CurrentIndex,
		scene:        scene,
		store:        store,
		state:        state.NewState(m.ID, store),
	}
}

func (s *Play[T]) FirstTime() bool {
	return s.firstTime
}

func (s *Play[T]) State() *state.State {
	return s.state
}

func (s *Play[T]) SetData(data T) {
	s.data = data
}

func (s *Play[T]) Execute() error {
	if s.FirstTime() {
		s.store.Play().Update(s.ID, store.ToggleFirstTime())
	}

	step, err := s.scene.GetStep(s.currentIndex)
	if err != nil {
		return err
	}

	return step.Run(s, s.data)
}

func (s *Play[T]) Next() error {
	if s.currentIndex >= s.scene.LastStepIndex() {
		s.Exit()

		return nil
	}

	return s.Go(s.currentIndex + 1)
}

func (s *Play[T]) Back() error {
	if s.currentIndex == 0 {
		s.Exit()

		return nil
	}

	return s.Go(s.currentIndex - 1)
}

func (s *Play[T]) Go(i int) error {
	if i < 0 || i > s.scene.LastStepIndex() {
		return errInvalidCurrentIndex
	}

	s.currentIndex = i
	s.firstTime = true

	s.store.Play().Update(
		s.ID,
		store.UpdateCurrentStepIndex(i),
		store.UpdateFirstTime(true),
	)

	return s.Execute()
}

func (s *Play[T]) Exit() error {
	s.store.Play().End(s.ID)
	return nil
}
