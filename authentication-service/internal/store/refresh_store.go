package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RefreshStore 用redis管理 jti 的存储
type RefreshStore struct {
	rdb *redis.Client
}

func NewRefreshStore(rdb *redis.Client) *RefreshStore {
	return &RefreshStore{
		rdb: rdb,
	}
}

// Save 保存 refresh token 的 jti
func (s *RefreshStore) Save(ctx context.Context, jti string, ttl time.Duration) error {
	return s.rdb.Set(ctx, jti, "1", ttl).Err()
}

func (s *RefreshStore) SavePair(ctx context.Context, key string, value string, ttl time.Duration) error {
	return s.rdb.Set(ctx, key, value, ttl).Err()
}

// Exists 判断 jti 是否存在
func (s *RefreshStore) Exists(ctx context.Context, jti string) (bool, error) {
	_, err := s.rdb.Get(ctx, jti).Result()
	if err == redis.Nil {
		return false, nil
	}
	return err == nil, err
}

func (s *RefreshStore) GetValue(ctx context.Context, key string) (string, error) {
	value, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return value, err
}

// Delete 删除 jti（登出 / rotation）
func (s *RefreshStore) Delete(ctx context.Context, jti string) error {
	return s.rdb.Del(ctx, jti).Err()
}
