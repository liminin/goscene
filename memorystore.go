package goscene

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"
)

type MemoryStore struct {
	play  *MemoryPlayRepository
	state *MemoryStateRepository
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) State() StateRepository {
	if s.state != nil {
		return s.state
	}

	s.state = &MemoryStateRepository{
		items: map[int]state{},
	}

	return s.state
}

func (s *MemoryStore) Play() PlayRepository {
	if s.play != nil {
		return s.play
	}

	s.play = &MemoryPlayRepository{
		items: map[int]*ScenePlayInfo{},
	}

	return s.play
}

func (s *MemoryStore) Info() PlayInfoRepository {
	if s.play != nil {
		return s.play
	}

	s.play = &MemoryPlayRepository{
		items: map[int]*ScenePlayInfo{},
	}

	return s.play
}

type state map[string]string

type MemoryStateRepository struct {
	items map[int]state

	playRepository *MemoryPlayRepository

	mx sync.RWMutex
}

func (r *MemoryStateRepository) Get(playID int, key string) (v string, err error) {
	if !r.playRepository.IsPlayExist(playID) {
		return "", errUserHasNotActivePlay
	}

	state := r.items[playID]

	v, ok := state[key]
	if !ok {
		return "", errors.New("not found")
	}

	return
}

func (r *MemoryStateRepository) Set(playID int, key string, value any) (err error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if !r.playRepository.IsPlayExist(playID) {
		return errUserHasNotActivePlay
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	state := r.items[playID]
	state[key] = string(data)
	r.items[playID] = state

	return nil
}

type MemoryPlayRepository struct {
	items map[int]*ScenePlayInfo

	nextID int

	mx sync.RWMutex
}

func (r *MemoryPlayRepository) IsPlayExist(playID int) bool {
	_, ok := r.items[playID]
	return ok
}

func (r *MemoryPlayRepository) Set(playID int, key storeKey, value interface{}) (err error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	play := r.items[playID]

	switch key {
	case firstTimeKey:
		v, ok := value.(bool)
		if !ok {
			return errors.New("invalid value")
		}

		play.FirstTime = v
	case currentIndexKey:
		v, ok := value.(int)
		if !ok {
			return errors.New("invalid value")
		}

		play.CurrentIndex = v
	}

	return nil
}

func (r *MemoryPlayRepository) GetByUserID(userID int) (s *ScenePlayInfo, err error) {
	id, err := r.GetIDByUserID(userID)
	if err != nil {
		return nil, err
	}

	play := *r.items[int(id)]

	return &play, nil
}

func (r *MemoryPlayRepository) GetIDByUserID(userID int) (int64, error) {
	ids := []int{}

	for _, play := range r.items {
		if play.UserID == userID {
			ids = append(ids, play.ID)
		}
	}

	if len(ids) == 0 {
		return 0, errUserHasNotActivePlay
	}

	sort.Ints(ids)
	id := ids[len(ids)-1]

	return int64(id), nil
}

func (r *MemoryPlayRepository) New(key string, userID int) (err error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.nextID++

	play := &ScenePlayInfo{
		ID:       r.nextID,
		SceneKey: key,
	}

	r.items[play.ID] = play

	return nil
}

func (r *MemoryPlayRepository) End(playID int) (err error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	delete(r.items, playID)

	return nil
}
