package cache

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/encoding/json"
	"github.com/spf13/cast"

	"news-service/internal/constant"
	"news-service/package/config"
)

const (
	WebRedis = "web"
)

type Redis struct {
	conn    *redis.Client
	conf    *config.Cache
	lock    *sync.Map
	prefix  string
	lruSize int64
}

func withCtxNR(ctx context.Context) context.Context {
	trx := ctx.Value(constant.NewRelicTransactionCtx)
	if trx == nil {
		return ctx
	}
	return newrelic.NewContext(ctx, trx.(*newrelic.Transaction))
}

func (r Redis) SIsMember(ctx context.Context, key string, members interface{}) (bool, error) {
	ctx = withCtxNR(ctx)
	key = r.prefixKey(key)
	return r.conn.SIsMember(ctx, key, members).Result()
}

func (r Redis) SMembers(ctx context.Context, key string) ([]string, error) {
	ctx = withCtxNR(ctx)
	key = r.prefixKey(key)
	return r.conn.SMembers(ctx, key).Result()
}

func (r Redis) SRem(ctx context.Context, key string, members interface{}) (int64, error) {
	ctx = withCtxNR(ctx)
	key = r.prefixKey(key)
	return r.conn.SRem(ctx, key, members).Result()
}

// SAdd implements Cache.
func (r Redis) SAdd(ctx context.Context, key string, members interface{}) (int64, error) {
	ctx = withCtxNR(ctx)
	key = r.prefixKey(key)
	return r.conn.SAdd(ctx, key, members).Result()
}

// Del implements Cache.
func (r Redis) Del(ctx context.Context, key []string) error {
	ctx = withCtxNR(ctx)
	keys := []string{}
	for _, v := range key {
		keys = append(keys, r.prefixKey(v))
	}
	return r.conn.Del(ctx, keys...).Err()
}

// Scan implements Cache.
func (r Redis) Scan(ctx context.Context, cursor, count int64, key string) ([]string, uint64, error) {
	ctx = withCtxNR(ctx)
	return r.conn.Scan(ctx, uint64(cursor), key, count).Result()
}

func (r Redis) ZGetByScore(ctx context.Context, key, min, max string) ([]string, error) {
	ctx = withCtxNR(ctx)
	return r.conn.ZRangeByScore(ctx, key, &redis.ZRangeBy{Min: min, Max: max}).Result()
}

func (r Redis) ZRemByScore(ctx context.Context, key, min, max string) (int64, error) {
	ctx = withCtxNR(ctx)
	return r.conn.ZRemRangeByScore(ctx, key, min, max).Result()
}

func (r Redis) ZAdd(ctx context.Context, key string, args redis.ZAddArgs) (int64, error) {
	ctx = withCtxNR(ctx)
	return r.conn.ZAddArgs(ctx, key, args).Result()
}

func (r Redis) TTL(ctx context.Context, key string) (time.Duration, error) {
	ctx = withCtxNR(ctx)
	key = r.prefixKey(key)
	return r.conn.TTL(ctx, key).Result()
}

func (r Redis) SetLRU(ctx context.Context, key string, identifier, val interface{}) error {
	ctx = withCtxNR(ctx)
	keyScore := r.prefixKeyLRU(key)
	key = fmt.Sprintf("%v:%v", key, identifier)
	if err := r.Set(ctx, key, val, redis.KeepTTL); err != nil {
		return err
	}
	if err := r.conn.ZAddNX(ctx, keyScore, redis.Z{Score: 1, Member: identifier}).Err(); err != nil {
		return err
	}
	length, err := r.conn.ZCard(ctx, keyScore).Result()
	if err != nil {
		return err
	}
	if length > r.lruSize {
		overflow := length - r.lruSize
		if err = r.conn.ZPopMin(ctx, keyScore, overflow).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r Redis) GetLRU(ctx context.Context, key string, identifier, dest interface{}) error {
	ctx = withCtxNR(ctx)
	keyScore := r.prefixKeyLRU(key)
	key = fmt.Sprintf("%v:%v", key, identifier)
	if err := r.conn.ZScore(ctx, keyScore, cast.ToString(identifier)).Err(); err != nil {
		return err
	}
	if err := r.conn.IncrBy(ctx, keyScore, 1).Err(); err != nil {
		return err
	}
	return r.Get(ctx, key, dest)
}

func (r Redis) RememberWithLockLocal(ctx context.Context, key string, dest interface{}, expiry time.Duration, fn func() error) error {
	ctx = withCtxNR(ctx)
	count, _ := r.Exists(ctx, []string{key})
	if count > 0 {
		return r.Get(ctx, key, dest)
	}
	lockKey := fmt.Sprintf("lock:%v", key)
	_, loaded := r.lock.LoadOrStore(lockKey, 1)
	log.Debugf("LOCK %v key: %v", !loaded, lockKey)
	if !loaded {
		log.Debug("LOCAL LOCKING: CALL SERVICE/FETCH FROM DB")
		if err := fn(); err != nil {
			r.lock.Delete(lockKey)
			return err
		}
		log.Debug("LOCAL LOCKING: CALL SERVICE/FETCH FROM DB SAVE")
		err := r.Set(ctx, key, dest, expiry)
		r.lock.Delete(lockKey)
		return err
	} else {
		for i := 0; i < 15; i++ {
			count, _ := r.Exists(ctx, []string{key})
			if count > 0 {
				return r.Get(ctx, key, dest)
			}
			min := 100
			max := 1000
			ran := rand.Intn(max-min) + min
			time.Sleep(time.Millisecond * time.Duration(ran))
		}
		if err := fn(); err != nil {
			r.lock.Delete(lockKey)
			return err
		}
		log.Debug("LOCAL LOCKING: FALLBACK CALL SERVICE INSTEAD")
		err := r.Set(ctx, key, dest, expiry)
		r.lock.Delete(lockKey)
		return err
	}
}

func (r Redis) RememberWithLock(ctx context.Context, key string, dest interface{}, expiry time.Duration, fn func() error) error {
	ctx = withCtxNR(ctx)
	count, _ := r.Exists(ctx, []string{key})
	if count > 0 {
		return r.Get(ctx, key, dest)
	}
	lockKey := fmt.Sprintf("lock:%v", key)

	if locking, err := r.conn.Conn().SetNX(ctx, lockKey, 1, time.Second*30).Result(); err != nil {
		return err
	} else if locking {
		log.Print("LOCKING: CALL SERVICE/FETCH FROM DB")
		if err := fn(); err != nil {
			return err
		}
		r.Set(ctx, key, dest, expiry)
		return r.Del(ctx, []string{lockKey})
	} else {
		for i := 0; i < 10; i++ {
			count, _ := r.Exists(ctx, []string{key})
			if count > 0 {
				log.Printf("LOCKING: WAITING FOUND: iteration %v", i)
				return r.Get(ctx, key, dest)
			}
			min := 100
			max := 1000
			ran := rand.Intn(max-min) + min
			time.Sleep(time.Millisecond * time.Duration(ran))
		}
		if err := fn(); err != nil {
			return err
		}
		log.Print("LOCKING: FALLBACK CALL SERVICE INSTEAD")
		r.Set(ctx, key, dest, expiry)
		return r.Del(ctx, []string{lockKey})
	}
}

func (r Redis) Remember(ctx context.Context, key string, dest interface{}, expiry time.Duration, fn func() error) error {
	ctx = withCtxNR(ctx)
	count, _ := r.Exists(ctx, []string{key})
	if count > 0 {
		return r.Get(ctx, key, dest)
	}
	if err := fn(); err != nil {
		return err
	}
	return r.Set(ctx, key, dest, expiry)
}

func (r Redis) LPush(ctx context.Context, key string, value interface{}) error {
	ctx = withCtxNR(ctx)
	key = r.prefixKey(key)
	return r.conn.LPush(ctx, key, value).Err()
}

func (r Redis) RPush(ctx context.Context, key string, value interface{}) error {
	ctx = withCtxNR(ctx)
	key = r.prefixKey(key)
	return r.conn.RPush(ctx, key, value).Err()
}

func (r Redis) DelWithoutPrefix(ctx context.Context, key []string) error {
	ctx = withCtxNR(ctx)
	return r.conn.Del(ctx, key...).Err()
}

func (r Redis) Get(ctx context.Context, key string, dest interface{}) error {
	ctx = withCtxNR(ctx)
	key = r.prefixKey(key)
	jsonByte, err := r.conn.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsonByte, dest); err != nil {
		log.Debug(err, "ERR UNMARSHAL")
		return err
	}
	return nil
}

func (r Redis) Set(ctx context.Context, key string, val interface{}, exp time.Duration) error {
	ctx = withCtxNR(ctx)
	key = r.prefixKey(key)
	jsonByte, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if exp == 0 {
		exp = 5 * time.Minute
	}
	return r.conn.Set(ctx, key, string(jsonByte), exp).Err()
}

func (r Redis) MSet(ctx context.Context, key string, val interface{}) error {
	ctx = withCtxNR(ctx)
	key = r.prefixKey(key)
	return r.conn.HSet(ctx, key, val).Err()
}

func (r Redis) MGet(ctx context.Context, keys []string) ([]interface{}, error) {
	ctx = withCtxNR(ctx)
	for i := range keys {
		keys[i] = r.prefixKey(keys[i])
	}
	return r.conn.MGet(ctx, keys...).Result()
}

func (r Redis) Exists(ctx context.Context, keys []string) (int64, error) {
	ctx = withCtxNR(ctx)
	for i := range keys {
		keys[i] = r.prefixKey(keys[i])
	}
	return r.conn.Exists(ctx, keys...).Result()
}

func (r Redis) prefixKey(key string) string {
	return fmt.Sprintf("%v%v", r.prefix, key)
}

func (r Redis) prefixKeyLRU(key string) string {
	return fmt.Sprintf("lru:%v%v", r.prefix, key)
}

func NewRedis(connType string, conf *config.Cache) Cache {
	r := Redis{
		conf: conf,
		lock: &sync.Map{},
	}
	r.prefix = conf.RedisPrefix
	r.lruSize = cast.ToInt64(conf.LRUSize)
	switch connType {
	case WebRedis:
		r.conn = redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%v:%v", conf.RedisHost, conf.RedisPort),
			PoolTimeout:  1 * time.Minute,
			PoolSize:     cast.ToInt(conf.PoolSize),
			MinIdleConns: cast.ToInt(conf.MinIdleConn),
			MaxIdleConns: cast.ToInt(conf.MaxIdleConn),
		})
	default:
		panic("unknown connection type")
	}
	r.conn.AddHook(nrredis.NewHook(r.conn.Options()))
	return r
}
