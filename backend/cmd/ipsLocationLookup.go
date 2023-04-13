package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"ip-lookup-app/entity"
	"ip-lookup-app/repository"
	"ip-lookup-app/service"
	"ip-lookup-app/util"
	"strconv"
)

func RunIpsLocationLookup() {
	ctx := context.Background()

	db, _ := strconv.Atoi(util.GetEnv("REDIS_DB"))
	client := redis.NewClient(&redis.Options{
		Addr:     util.GetEnv("REDIS_HOST"),
		Password: util.GetEnv("REDIS_PASSWORD"),
		DB:       db,
	})
	redisRepository := repository.NewRedisRepository(client, ctx)
	geoIpService := service.NewGeoIpService()
	ipLocationService := service.NewIpLocationService(redisRepository, geoIpService)

	kafkaReaderRepo := repository.NewKafkaReaderRepository(
		ctx,
		[]string{util.GetEnv("KAFKA_HOST")},
		"ips-lookup-request",
	)

	for {
		// consume to ips-lookup-request
		messageKey, messageValue, _ := kafkaReaderRepo.ReadData()

		var ipsRequest entity.IpsRequest
		json.Unmarshal(messageValue, &ipsRequest)
		ipLocations, err := ipLocationService.GetBulkLocation(ipsRequest.Ip)
		if err != nil {
			fmt.Errorf("failed to lookup IPs location: %v", err)
		}

		//produce to ips-lookup-result
		kafkaWriterRepo := repository.NewKafkaWriterRepository(
			ctx,
			[]string{util.GetEnv("KAFKA_HOST")},
			"ips-lookup-result",
		)
		message, _ := json.Marshal(ipLocations)
		kafkaWriterRepo.WriteData([]byte(messageKey), message)
		kafkaWriterRepo.Close()
	}
}
