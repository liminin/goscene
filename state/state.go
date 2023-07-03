package state

import "github.com/liminin/goscene/store"

type State struct {
	playID int
	store  store.Store
}

func NewState(playID int, store store.Store) *State {
	return &State{
		playID: playID,
		store:  store,
	}
}

func (s *State) Get(key string) (*Item, error) {
	m, err := s.store.State().Get(s.playID, key)

	return NewItem(m, err), err
}

func (s *State) Set(key string, value any) error {
	return s.store.State().Set(s.playID, key, value)
}
