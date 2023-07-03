package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v9"
	"github.com/liminin/goscene/model"
	"github.com/mitchellh/mapstructure"
)

var (
	ctx = context.Background()
)

type RedisStore struct {
	redis           *redis.Client
	redisKey        func(a ...any) string
	playRepository  *RedisPlayRepository
	stateRepository *RedisStateRepository
}

func NewRedisStore(r *redis.Client, rk func(a ...any) string) *RedisStore {
	return &RedisStore{
		redis:    r,
		redisKey: rk,
	}
}

func (s *RedisStore) Play() PlayRepository {
	if s.playRepository != nil {
		return s.playRepository
	}

	s.playRepository = &RedisPlayRepository{
		redis:    s.redis,
		redisKey: s.redisKey,
	}

	return s.playRepository
}

func (s *RedisStore) State() StateRepository {
	if s.stateRepository != nil {
		return s.stateRepository
	}

	s.stateRepository = &RedisStateRepository{
		redis:    s.redis,
		redisKey: s.redisKey,
	}

	return s.stateRepository
}

type RedisPlayRepository struct {
	redis    *redis.Client
	redisKey func(a ...any) string
}

// Update updates play's fields
func (r *RedisPlayRepository) Update(playID int, upd ...Upd) (err error) {
	play, err := r.get(playID)

	for _, u := range upd {
		u(play)
	}

	return r.redis.HSet(ctx, r.redisKey("play", playID), map[string]any{
		"current":    play.CurrentIndex,
		"first_time": play.FirstTime,
	}).Err()
}

// GetByUserID returns the play model by user id.
// if the user has more than 1 active Play, then the earliest play is returned
func (r *RedisPlayRepository) GetByUserID(userID int) (s *model.Play, err error) {
	playID, err := r.GetIDByUserID(userID)

	if err != nil {
		return
	}

	return r.get(playID)
}

func (r *RedisPlayRepository) get(playID int) (s *model.Play, err error) {
	v := r.redis.HGetAll(ctx, r.redisKey("play", playID)).Val()

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           s,
	})

	err = decoder.Decode(v)

	return
}

func (r *RedisPlayRepository) GetIDByUserID(userID int) (id int, err error) {
	v := r.redis.LRange(ctx, r.redisKey("user", userID), 0, -1).Val()

	if len(v) == 0 {
		err = errUserHasNotActivePlay
		return
	}

	i, err := strconv.ParseInt(v[0], 10, 64)

	return int(i), err
}

func (r *RedisPlayRepository) PlayExist(playID int) bool {
	result := r.redis.Exists(ctx, r.redisKey("play", playID)).Val()

	return result != 0
}

func (r *RedisPlayRepository) New(key string, userID int) (err error) {
	id, _ := r.redis.Incr(ctx, r.redisKey("next_id")).Uint64()

	r.redis.LPush(ctx, r.redisKey("user", userID), id)

	err = r.redis.HSet(ctx, r.redisKey("play", id), map[string]any{
		"id":         id,
		"user_id":    userID,
		"scene_key":  key,
		"current":    0,
		"first_time": true,
	}).Err()

	return
}

func (r *RedisPlayRepository) End(playID int) (err error) {
	uID, _ := r.redis.HGet(ctx, r.redisKey("play", playID), "user_id").Int()

	r.redis.LRem(ctx, r.redisKey("user", uID), 0, playID)

	err = r.redis.Unlink(
		ctx,
		r.redisKey("play", playID),
		r.redisKey("play", playID, "state"),
	).Err()

	return
}

type RedisStateRepository struct {
	redis    *redis.Client
	redisKey func(a ...any) string
}

func (r *RedisStateRepository) Get(playID int, key string) (i *model.Item, err error) {
	v, err := r.redis.HGet(ctx, r.redisKey("play", playID, "state"), key).Result()

	i = &model.Item{
		Key:    key,
		PlayID: playID,
		Value:  v,
	}

	return
}

func (r *RedisStateRepository) Set(playID int, key string, value any) (err error) {
	data := []byte{}

	switch value.(type) {
	case string:
		data = []byte(fmt.Sprint(value))
	default:
		data, err = json.Marshal(value)
		if err != nil {
			return err
		}

	}

	return r.redis.HSet(ctx, r.redisKey("play", playID, "state"), key, data).Err()
}
