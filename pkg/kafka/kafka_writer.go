package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	//"github.com/hamba/avro"
	"github.com/segmentio/kafka-go"
)

// type AvroKafkaWriter struct {
// 	ctx context.Context
// 	writer *kafka.Writer
// 	avroschema avro.Schema
// }

// func NewAvroKafkaWriter(avroschema string, ctx context.Context) (*AvroKafkaWriter, error) {

// 	schema, err := avro.Parse(avroschema)
// 	if err != nil {
// 		return nil, err
// 	}
// ...
// 	return &AvroKafkaProducer{
// 		ctx:        ctx,
// 		writer:   w,
// 		avroschema: schema,
// 	}, nil
// }

func Writer() {

	fmt.Println("start kafka writer")
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "test-topic",
		Async:   true,
	})

	// for {
	// 	m, err := r.ReadMessage(context.Background())
	// 	if err != nil {
	// 		log.Fatal("error reading message:", err)
	// 		break
	// 	}
	// 	fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	// }

	ctx := context.Background()

	err := w.WriteMessages(ctx, kafka.Message{
		Key:   []byte("Key-A"),
		Value: []byte("Hello World!"),
	})

	if err != nil {
		log.Fatal("error writing message:", err)
	}

	done := make(chan bool)
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case t := <-ticker.C:
				fmt.Println("Tick")
				fmt.Println("Current time: ", t)

				err := w.WriteMessages(ctx, kafka.Message{
					Key:   []byte("Key-A"),
					Value: []byte("Hello World! " + t.String()),
				})

				if err != nil {
					log.Fatal("error writing message:", err)
				}
			}
		}
	}()

	// wait for 10 seconds
	time.Sleep(10 * time.Second)
	done <- true

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

// func main() {
// 	Writer()
// }
