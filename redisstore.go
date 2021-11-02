package goscene

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
)

var (
	ctx = context.Background()
)

type RedisStore struct {
	redis           *redis.Client
	redisKey        func(a ...interface{}) string
	playRepository  *RedisPlayRepository
	stateRepository *RedisStateRepository
	infoRepository  *RedisPlayInfoRepository
}

func NewRedisStore(r *redis.Client, rk func(a ...interface{}) string) *RedisStore {
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

func (s *RedisStore) Info() PlayInfoRepository {
	if s.infoRepository != nil {
		return s.infoRepository
	}

	s.infoRepository = &RedisPlayInfoRepository{
		redis:    s.redis,
		redisKey: s.redisKey,
	}

	return s.infoRepository
}

type RedisPlayRepository struct {
	redis    *redis.Client
	redisKey func(a ...interface{}) string
}

func (r *RedisPlayRepository) New(key string, userID int) (err error) {
	id, _ := r.redis.Incr(ctx, r.redisKey("next_id")).Uint64()

	r.redis.LPush(ctx, r.redisKey("user", userID), id)

	err = r.redis.HSet(ctx, r.redisKey("play", id), map[string]interface{}{
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

type RedisPlayInfoRepository struct {
	redis    *redis.Client
	redisKey func(a ...interface{}) string
}

func (r *RedisPlayInfoRepository) Set(playID int, key storeKey, value interface{}) {
	r.redis.HSet(ctx, r.redisKey("play", playID), fmt.Sprint(key), value)
}

func (r *RedisPlayInfoRepository) GetIDByUserID(userID int) (id int64, err error) {
	v := r.redis.LRange(ctx, r.redisKey("user", userID), -1, 1).Val()

	if len(v) == 0 {
		err = errUserHasNotActivePlay
		return
	}

	id, err = strconv.ParseInt(v[0], 10, 64)
	return
}

func (r *RedisPlayInfoRepository) GetByUserID(userID int) (spi *ScenePlayInfo, err error) {
	playID, err := r.GetIDByUserID(userID)

	if err != nil {
		return
	}

	v := r.redis.HGetAll(ctx, r.redisKey("play", playID)).Val()

	spi = &ScenePlayInfo{}

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &spi,
	})

	err = decoder.Decode(v)

	return
}

type RedisStateRepository struct {
	redis    *redis.Client
	redisKey func(a ...interface{}) string
}

func (r *RedisStateRepository) Get(playID int, key string) (v interface{}, err error) {
	v, err = r.redis.HGet(ctx, r.redisKey("play", playID, "state"), key).Result()

	return
}

func (r *RedisStateRepository) Set(playID int, key string, value interface{}) (err error) {
	err = r.redis.HSet(ctx, r.redisKey("play", playID, "state"), key, value).Err()

	return
}
