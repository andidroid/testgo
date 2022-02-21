package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	//"github.com/hamba/avro"
	"github.com/segmentio/kafka-go"
)

// type AvroKafkaReader struct {
// 	ctx context.Context
// 	reader *kafka.Reader
// 	avroschema avro.Schema
// }

// func NewAvroKafkaReader(avroschema string, ctx context.Context) (*AvroKafkaReader, error) {

// 	schema, err := avro.Parse(avroschema)
// 	if err != nil {
// 		return nil, err
// 	}
// ...
// 	return &AvroKafkaReader{
// 		ctx:        ctx,
// 		reader:   r,
// 		avroschema: schema,
// 	}, nil
// }

func Reader() {

	fmt.Println("start kafka reader")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		Topic:          "test-topic",
		GroupID:        "test-group",
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
		Partition:      0,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
	})
	r.SetOffset(42)

	// for {
	// 	m, err := r.ReadMessage(context.Background())
	// 	if err != nil {
	// 		log.Fatal("error reading message:", err)
	// 		break
	// 	}
	// 	fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	// }

	ctx := context.Background()
	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			log.Fatal("error reading message:", err)
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
