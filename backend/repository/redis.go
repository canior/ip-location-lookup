package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepository(client *redis.Client, ctx context.Context) *RedisRepository {
	return &RedisRepository{
		client: client,
		ctx:    ctx,
	}
}

func (r *RedisRepository) WriteData(key string, data interface{}) error {
	// Encode data as JSON
	encoded, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error encoding data as JSON: %v", err)
	}

	// Write encoded data to Redis
	err = r.client.Set(r.ctx, key, encoded, 0).Err()
	if err != nil {
		return fmt.Errorf("error writing data to Redis: %v", err)
	}

	return nil
}

func (r *RedisRepository) GetData(key string, result interface{}) error {
	val, err := r.client.Get(r.ctx, key).Result()

	exists, _ := r.client.Exists(r.ctx, key).Result()
	if exists != 1 {
		return nil
	}

	err = json.Unmarshal([]byte(val), result)
	if err != nil {
		return fmt.Errorf("error decoding data: %v", err)
	}

	return nil
}
