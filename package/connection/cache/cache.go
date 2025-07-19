package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	Set(ctx context.Context, key string, val interface{}, exp time.Duration) error
	SetLRU(ctx context.Context, key string, identifier, val interface{}) error
	GetLRU(ctx context.Context, key string, identifier, dest interface{}) error
	MSet(ctx context.Context, key string, val interface{}) error
	MGet(ctx context.Context, keys []string) ([]interface{}, error)
	Del(ctx context.Context, key []string) error
	DelWithoutPrefix(ctx context.Context, key []string) error
	LPush(ctx context.Context, key string, value interface{}) error
	RPush(ctx context.Context, key string, value interface{}) error
	Remember(ctx context.Context, key string, dest interface{}, expiry time.Duration, fn func() error) error
	RememberWithLock(ctx context.Context, key string, dest interface{}, expiry time.Duration, fn func() error) error
	RememberWithLockLocal(ctx context.Context, key string, dest interface{}, expiry time.Duration, fn func() error) error
	SIsMember(ctx context.Context, key string, members interface{}) (bool, error)

	SMembers(ctx context.Context, key string) ([]string, error)

	SRem(ctx context.Context, key string, members interface{}) (int64, error)

	SAdd(ctx context.Context, key string, val interface{}) (int64, error)
	ZAdd(ctx context.Context, key string, args redis.ZAddArgs) (int64, error)
	ZRemByScore(ctx context.Context, key, min, max string) (int64, error)
	ZGetByScore(ctx context.Context, key, min, max string) ([]string, error)
	Scan(ctx context.Context, cursor, count int64, key string) ([]string, uint64, error)
}
