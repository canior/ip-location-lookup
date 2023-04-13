package service_test

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"ip-lookup-app/entity"
	"ip-lookup-app/repository"
	"ip-lookup-app/service"
	"ip-lookup-app/util"
	"os"
	"strconv"
	"testing"
)

func TestIpLocationService(t *testing.T) {
	os.Setenv("APP_ENV", "testing")

	db, _ := strconv.Atoi(util.GetEnv("REDIS_DB"))
	client := redis.NewClient(&redis.Options{
		Addr:     util.GetEnv("REDIS_HOST"),
		Password: util.GetEnv("REDIS_PASSWORD"),
		DB:       db,
	})

	redisRepository := repository.NewRedisRepository(client, context.Background())
	geoIpService := service.NewGeoIpService()
	ipLocationService := service.NewIpLocationService(redisRepository, geoIpService)

	ip := "216.209.131.154"
	redisRepository.WriteData(ip, nil)

	ipLocation, _ := ipLocationService.GetLocation(ip)
	assert.Equal(t, ip, ipLocation.Ip)

	var ipLocationFromRedis entity.IpLocation
	redisRepository.GetData(ip, &ipLocationFromRedis)
	assert.Equal(t, ipLocation, ipLocationFromRedis)
}
