package cmd_test

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"ip-lookup-app/cmd"
	"ip-lookup-app/entity"
	"ip-lookup-app/repository"
	"ip-lookup-app/util"
	"os"
	"testing"
)

func TestRunIpsLocationLookup(t *testing.T) {
	os.Setenv("APP_ENV", "testing")

	// create a producer ips-lookup-request
	ctx := context.Background()

	// Create Kafka writer repository
	kafkaWriterRepo := repository.NewKafkaWriterRepository(
		ctx,
		[]string{util.GetEnv("KAFKA_HOST")},
		"ips-lookup-request",
	)

	// Create Kafka reader repository
	kafkaReaderRepo := repository.NewKafkaReaderRepository(
		ctx,
		[]string{util.GetEnv("KAFKA_HOST")},
		"ips-lookup-result",
	)

	ip := "216.209.131.154"
	randomKey := uuid.New().String()
	ipRequest := entity.IpsRequest{
		Ip: []string{ip},
	}
	encoded, _ := json.Marshal(ipRequest)
	kafkaWriterRepo.WriteData([]byte(randomKey), encoded)

	go cmd.RunIpsLocationLookup()

	// create a consumer to wait for ips-lookup-result
	for {
		// Read a message from the Kafka topic
		messageKey, messageValue, _ := kafkaReaderRepo.ReadData()
		if messageKey == randomKey {
			assert.Equal(t, "[{\"ip\":\"216.209.131.154\",\"city\":\"Toronto\",\"timezone\":\"America/Toronto\",\"accuracy_radius\":5,\"postal_code\":\"M4M\"}]", string(messageValue))
			break
		}
	}

}
