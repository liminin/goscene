package goscene

type storeKey string

const (
	currentIndexKey storeKey = "current"
	firstTimeKey    storeKey = "first_time"
)

type Store interface {
	State() StateRepository
	Play() PlayRepository
	Info() PlayInfoRepository
}

type StateRepository interface {
	Get(playID int, key string) (v interface{}, err error)
	Set(playID int, key string, value interface{}) (err error)
}

type PlayRepository interface {
	New(key string, userID int) (err error)
	End(playID int) (err error)
}

type PlayInfoRepository interface {
	Set(playID int, key storeKey, value interface{})
	GetByUserID(userID int) (s *ScenePlayInfo, err error)
	GetIDByUserID(userID int) (id int64, err error)
}
