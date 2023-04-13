package repository_test

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"ip-lookup-app/entity"
	"ip-lookup-app/repository"
	"ip-lookup-app/util"
	"os"
	"strconv"
	"testing"
)

func TestRedisRepository(t *testing.T) {
	os.Setenv("APP_ENV", "testing")

	db, _ := strconv.Atoi(util.GetEnv("REDIS_DB"))
	client := redis.NewClient(&redis.Options{
		Addr:     util.GetEnv("REDIS_HOST"),
		Password: util.GetEnv("REDIS_PASSWORD"),
		DB:       db,
	})

	ctx := context.Background()
	client.FlushAll(ctx)

	redisRepository := repository.NewRedisRepository(client, ctx)

	var actualEmpty entity.IpLocation
	err := redisRepository.GetData("0.0.0.1", &actualEmpty)
	assert.Equal(t, "", actualEmpty.Ip)
	assert.Empty(t, err)

	expect := entity.IpLocation{
		Ip:   "0.0.0.0",
		City: "San Francisco",
	}
	redisRepository.WriteData("0.0.0.0", expect)

	var actual entity.IpLocation
	redisRepository.GetData("0.0.0.0", &actual)
	assert.Equal(t, expect, actual)

	client.FlushAll(ctx)
}
