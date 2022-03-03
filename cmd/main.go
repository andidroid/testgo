package main

import (
	"fmt"
	"os"
	"sync"

	// "time"
	//"channel"

	"google.golang.org/grpc"

	// "github.com/andidroid/testgo/pkg/kafka"
	"github.com/andidroid/testgo/internal/channel"
	pb "github.com/andidroid/testgo/pkg/grpc"
	"github.com/andidroid/testgo/pkg/kafka"
	"github.com/andidroid/testgo/pkg/server"
)

var wg = sync.WaitGroup{}

func main() {
	fmt.Println("Run Test Go")

	apichannel := make(chan channel.ApiChannelEntity, 10)
	// messagechannel := make(chan channel.MessageChannelEntity, 10)

	wg.Add(1)
	go runServer(apichannel)
	// wg.Add(1)
	// go runKafkaProducer(apichannel)
	// wg.Add(1)
	// go runKafkaConsumer(messagechannel)
	// wg.Add(1)
	// go runMongoWriter(messagechannel)

	// go runGreeterServer(messagechannel)

	//

	// go kafka.Reader()

	// time.Sleep(2 * time.Second)

	// kafka.Writer()

	// time.Sleep(10 * time.Second)

	wg.Wait()

}

func runServer(apichannel chan<- channel.ApiChannelEntity) {
	// wg.Add(1)
	server.Start()
	wg.Done()
}

func runKafkaProducer(apichannel <-chan channel.ApiChannelEntity) {
	// wg.Add(1)
	kafka.Writer(apichannel)
	wg.Done()
}

func runKafkaConsumer(apichannel chan channel.MessageChannelEntity) {
	// wg.Add(1)
	kafka.Reader(apichannel)
	wg.Done()
}

func runMongoWriter(messagechannel <-chan channel.MessageChannelEntity) {
	// wg.Add(1)

	e := <-messagechannel
	fmt.Println(e.Message)

	wg.Done()
}

func runGreeterServer(messagechannel chan channel.MessageChannelEntity) {
	// wg.Add(1)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &greeterserver{})
	wg.Done()
}

type greeterserver struct {
	pb.UnimplementedGreeterServer
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
