package store

import "github.com/liminin/goscene/model"

type storeKey string

const (
	currentIndexKey storeKey = "current"
	firstTimeKey    storeKey = "first_time"
)

type Store interface {
	State() StateRepository
	Play() PlayRepository
}

type StateRepository interface {
	Get(playID int, key string) (i *model.Item, err error)
	Set(playID int, key string, value any) (err error)
}

type Upd func(*model.Play)

type PlayRepository interface {
	// New creates new play
	New(key string, userID int) (err error)
	// Update updates play's fields
	Update(playID int, upd ...Upd) (err error)
	// GetByUserID returns the play model by user id.
	// if the user has more than 1 active Play, then the earliest play is returned
	GetByUserID(userID int) (s *model.Play, err error)
	// End deletes play
	End(playID int) (err error)

	PlayExist(playID int) bool
}

func UpdateCurrentStepIndex(index int) Upd {
	return func(p *model.Play) {
		p.CurrentIndex = index
	}
}

func ToggleFirstTime() Upd {
	return func(p *model.Play) {
		p.FirstTime = !p.FirstTime
	}
}

func UpdateFirstTime(firstTime bool) Upd {
	return func(p *model.Play) {
		p.FirstTime = firstTime
	}
}
