package repository

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type KafkaReaderRepository struct {
	reader *kafka.Reader
	ctx    context.Context
}

func NewKafkaReaderRepository(ctx context.Context, brokers []string, topic string) *KafkaReaderRepository {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
	})

	return &KafkaReaderRepository{
		reader: reader,
		ctx:    ctx,
	}
}

func (kr *KafkaReaderRepository) ReadData() (string, []byte, error) {
	msg, err := kr.reader.ReadMessage(kr.ctx)

	if err != nil {
		return "", nil, fmt.Errorf("error reading data from Kafka: %v", err)
	}

	return string(msg.Key), msg.Value, nil
}

func (kr *KafkaReaderRepository) Close() error {
	return kr.reader.Close()
}
