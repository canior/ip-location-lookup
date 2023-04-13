package api

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"ip-lookup-app/entity"
	"ip-lookup-app/repository"
	"ip-lookup-app/util"
	"net"
	"net/http"
)

func IpsLookupHandler(w http.ResponseWriter, r *http.Request, clientChannelMap map[string]chan string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	if r.Method == "OPTIONS" {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var ipsRequest entity.IpsRequest
	err := decoder.Decode(&ipsRequest)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for _, ip := range ipsRequest.Ip {
		if net.ParseIP(ip) == nil {
			http.Error(w, "Invalid Ips found in request body", http.StatusBadRequest)
			return
		}
	}

	clientId := mux.Vars(r)["clientId"]
	messageKey := uuid.New().String()

	if clientChannelMap[clientId] != nil {
		ctx := context.Background()

		// Create Kafka writer repository
		kafkaWriterRepo := repository.NewKafkaWriterRepository(
			ctx,
			[]string{util.GetEnv("KAFKA_HOST")},
			"ips-lookup-request",
		)
		defer kafkaWriterRepo.Close()

		encoded, _ := json.Marshal(ipsRequest)
		kafkaWriterRepo.WriteData([]byte(messageKey), encoded)

		kafkaReaderRepo := repository.NewKafkaReaderRepository(
			ctx,
			[]string{util.GetEnv("KAFKA_HOST")},
			"ips-lookup-result",
		)
		defer kafkaReaderRepo.Close()

		for {
			// Read a message from the Kafka topic
			k, v, _ := kafkaReaderRepo.ReadData()
			if k == messageKey {
				clientChannelMap[clientId] <- string(v)
				break
			}
		}

	}
}
