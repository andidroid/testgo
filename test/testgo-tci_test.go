package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var itest string

func init() {
	flag.StringVar(&itest, "itest", "", "the foo bar bang")
}

type AppContainer struct {
	testcontainers.Container
	URI  string
	Host string
	Port string
}

type MongoContainer struct {
	testcontainers.Container
	URI  string
	Host string
	Port string
}

type RedisContainer struct {
	testcontainers.Container
	URI  string
	Host string
	Port string
}

const defaultAppPort = 80
const defaultMongoDBPort = 27017
const defaultredisPort = 6379

//, redisC RedisContainer
func setupApp(ctx context.Context, mongoC MongoContainer) (*AppContainer, error) {

	fmt.Println("Start setup testcontainers app")
	timeout := 5 * time.Minute // Default timeout

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       "D:\\Programmierung\\Git\\testgo",
			Dockerfile:    "Dockerfile",
			PrintBuildLog: true,
		},

		ExposedPorts: []string{"80"},
		// TODO use /health
		WaitingFor: wait.ForHTTP("/metrics").WithStartupTimeout(timeout),
		Env: map[string]string{
			// "REDIS_HOST": redisC.Host,
			// "REDIS_PORT": redisC.Port,
			"MONGODB_HOST": "network", //mongoC.Host,
			"MONGODB_PORT": mongoC.Port,
		},
		Networks: []string{"network"},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "80")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &AppContainer{Container: container, URI: uri, Host: ip, Port: mappedPort.Port()}, nil
}

func setupMongo(ctx context.Context) (*MongoContainer, error) {

	fmt.Println("Start setup testcontainers mongodb")
	timeout := 5 * time.Minute // Default timeout

	req := testcontainers.ContainerRequest{
		Image:        "mongo",
		ExposedPorts: []string{"27017"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(timeout),
		Env: map[string]string{
			"MONGO_INITDB_DATABASE": "test",
			// "MONGO_INITDB_ROOT_USERNAME": "test",
			// "MONGO_INITDB_ROOT_PASSWORD": "test",
		},
		Networks: []string{"network"},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("mongodb://%s:%s", ip, mappedPort.Port())

	return &MongoContainer{Container: container, URI: uri, Host: ip, Port: mappedPort.Port()}, nil
}

func setupRedis(ctx context.Context) (*RedisContainer, error) {

	fmt.Println("Start setup testcontainers redis")
	timeout := 5 * time.Minute // Default timeout

	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections").WithStartupTimeout(timeout),
		Networks:     []string{"network"},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("redis started: ", container)
	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &RedisContainer{Container: container, URI: uri, Host: ip, Port: mappedPort.Port()}, nil
}

func TestTC(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	if itest != "tcit" {
		fmt.Print("skip itest")
		t.Skip("Skipping for itest != tcit")
	}

	ctx := context.Background()

	networkRequest := testcontainers.NetworkRequest{
		Driver:     "bridge",
		Name:       "network",
		Attachable: true,
	}

	net, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: networkRequest,
	})

	if err != nil {
		t.Fatalf("cannot create network: %s", err)
	}

	defer net.Remove(ctx)

	fmt.Printf("created network: %s", net)
	mongoC, err := setupMongo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("mongo db container started at: %s", mongoC.URI)
	fmt.Println()
	// redisC, err := setupRedis(ctx)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Printf("redis container started at: %s", redisC.URI)
	// fmt.Println()
	//, *redisC
	appC, err := setupApp(ctx, *mongoC)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println()
	fmt.Println("Running Tests...")

	// Clean up the container after the test is complete

	defer appC.Terminate(ctx)
	// defer redisC.Terminate(ctx)
	defer mongoC.Terminate(ctx)

	resp, err := http.Get(appC.URI + "/test")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("calling /test: %s", resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d. Got %d.", http.StatusOK, resp.StatusCode)
	}
}

func main() {
	ctx := context.Background()

	// redisC, err := setupRedis(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	mongoC, err := setupMongo(ctx)
	if err != nil {
		fmt.Println(err)
	}
	//, *redisC
	appC, err := setupApp(ctx, *mongoC)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Running Tests...")

	// Clean up the container after the test is complete
	defer appC.Terminate(ctx)
	defer mongoC.Terminate(ctx)

	resp, err := http.Get(appC.URI + "/test")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("calling /test: %s", resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Expected status code %d. Got %d.", http.StatusOK, resp.StatusCode)
	}
}
