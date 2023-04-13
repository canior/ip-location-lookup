package repository

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"ip-lookup-app/entity"
	"ip-lookup-app/repository"
	"ip-lookup-app/util"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestKafkaRepository(t *testing.T) {
	os.Setenv("APP_ENV", "testing")
	// Set up context
	ctx := context.Background()

	// Create Kafka writer repository
	kafkaWriterRepo := repository.NewKafkaWriterRepository(
		ctx,
		[]string{util.GetEnv("KAFKA_HOST")},
		"test-topic",
	)

	// Create Kafka reader repository
	kafkaReaderRepo := repository.NewKafkaReaderRepository(
		ctx,
		[]string{util.GetEnv("KAFKA_HOST")},
		"test-topic",
	)

	randomKey := uuid.New().String()
	randomIp := generateRandomIP()
	ipRequest := entity.IpsRequest{
		Ip: []string{randomIp},
	}
	encoded, _ := json.Marshal(ipRequest)
	go kafkaWriterRepo.WriteData([]byte(randomKey), encoded)

	// Continuously read messages
	for {
		// Read a message from the Kafka topic
		messageKey, messageValue, _ := kafkaReaderRepo.ReadData()

		if messageKey == randomKey {
			var actual entity.IpsRequest
			json.Unmarshal(messageValue, &actual)
			assert.Equal(t, randomIp, actual.Ip[0])
			break
		}
	}
}

func generateRandomIP() string {
	// Generate four random integers between 0 and 255
	octet1 := rand.Intn(256)
	octet2 := rand.Intn(256)
	octet3 := rand.Intn(256)
	octet4 := rand.Intn(256)

	// Combine the four integers into an IP address string
	ip := strconv.Itoa(octet1) + "." + strconv.Itoa(octet2) + "." + strconv.Itoa(octet3) + "." + strconv.Itoa(octet4)

	return ip
}
