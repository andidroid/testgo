package main

import (
	"fmt"
	"os"

	// "time"

	"github.com/andidroid/testgo/pkg/server"
	// "github.com/andidroid/testgo/pkg/kafka"
)

func main() {
	fmt.Println("Run Test Go")
	server.Start()

	// go kafka.Reader()

	// time.Sleep(2 * time.Second)

	// kafka.Writer()

	// time.Sleep(10 * time.Second)

}

func readEnvVars() {

	mongodbHost, ok := os.LookupEnv("MONGODB_HOST")
	if !ok {
		mongodbHost = "localhost"
	}
	fmt.Printf("MONGODB_HOST: %s\n", mongodbHost)

	mongodbPort, ok := os.LookupEnv("MONGODB_PORT")
	if !ok {
		mongodbPort = "27017"
	}
	fmt.Printf("MONGODB_PORT: %s\n", mongodbPort)

	kafkaHost, ok := os.LookupEnv("KAFKA_HOST")
	if !ok {
		kafkaHost = "localhost"
	}
	fmt.Printf("KAFKA_HOST: %s\n", kafkaHost)

	kafkaPort, ok := os.LookupEnv("KAFKA_PORT")
	if !ok {
		kafkaPort = "27017"
	}
	fmt.Printf("KAFKA_PORT: %s\n", kafkaPort)

	// dbURL := os.ExpandEnv("postgres://$DB_USERNAME:$DB_PASSWORD@DB_HOST:$DB_PORT/$DB_NAME")

}
