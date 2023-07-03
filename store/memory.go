package store

import (
	"errors"
	"sort"
	"sync"

	"github.com/liminin/goscene/model"
)

var (
	errUserHasNotActivePlay = errors.New("user does not have active play")
	errPlayNotFound         = errors.New("play not found")
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
		items:          map[int]state{},
		playRepository: s.Play(),
	}

	return s.state
}

func (s *MemoryStore) Play() PlayRepository {
	if s.play != nil {
		return s.play
	}

	s.play = &MemoryPlayRepository{
		items: map[int]model.Play{},
	}

	return s.play
}

type state map[string]model.Item

type MemoryStateRepository struct {
	items map[int]state

	playRepository PlayRepository

	mx sync.RWMutex
}

func (r *MemoryStateRepository) Get(playID int, key string) (m *model.Item, err error) {
	if !r.playRepository.PlayExist(playID) {
		return nil, errUserHasNotActivePlay
	}

	state := r.items[playID]

	item, ok := state[key]
	if !ok {
		return nil, errors.New("not found")
	}

	return &item, nil
}

func (r *MemoryStateRepository) Set(playID int, key string, value any) (err error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if !r.playRepository.PlayExist(playID) {
		return errUserHasNotActivePlay
	}

	s, ok := r.items[playID]
	if !ok {
		s = state{}
	}

	item := model.Item{
		Key:    key,
		Value:  value,
		PlayID: playID,
	}

	s[key] = item
	r.items[playID] = s

	return nil
}

type MemoryPlayRepository struct {
	items map[int]model.Play

	nextID int

	mx sync.RWMutex
}

func (r *MemoryPlayRepository) PlayExist(playID int) bool {
	_, ok := r.items[playID]
	return ok
}

func (r *MemoryPlayRepository) Set(playID int, key storeKey, value any) (err error) {
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

// Update updates play's fields
func (r *MemoryPlayRepository) Update(playID int, upd ...Upd) (err error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	play, ok := r.items[playID]
	if !ok {
		return errPlayNotFound
	}

	for _, u := range upd {
		u(&play)
	}

	r.items[playID] = play

	return nil
}

func (r *MemoryPlayRepository) GetByUserID(userID int) (s *model.Play, err error) {
	id, err := r.GetIDByUserID(userID)
	if err != nil {
		return nil, err
	}

	play := r.items[int(id)]

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

	play := model.Play{
		ID:        r.nextID,
		SceneKey:  key,
		UserID:    userID,
		FirstTime: true,
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
