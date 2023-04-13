package repository

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type KafkaWriterRepository struct {
	writer *kafka.Writer
	ctx    context.Context
}

func NewKafkaWriterRepository(ctx context.Context, brokers []string, topic string) *KafkaWriterRepository {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers[0]),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaWriterRepository{
		writer: writer,
		ctx:    ctx,
	}
}

func (kr *KafkaWriterRepository) WriteData(key []byte, data []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: data,
	}
	err := kr.writer.WriteMessages(kr.ctx, msg)
	if err != nil {
		return fmt.Errorf("error writing data to Kafka: %v", err)
	}

	return nil
}

func (kr *KafkaWriterRepository) Close() error {
	return kr.writer.Close()
}
