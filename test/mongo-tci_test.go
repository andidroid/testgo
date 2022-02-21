package test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type appContainer struct {
	testcontainers.Container
	URI string
}

type mongoContainer struct {
	testcontainers.Container
	URI string
}

const defaultAppPort = 8090
const defaultMongoDBPort = 27017

func setupApp(ctx context.Context) (*appContainer, error) {

	fmt.Println("Start setup testcontainers app")
	timeout := 5 * time.Minute // Default timeout

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: "/",
		},

		ExposedPorts: []string{"8090"},
		WaitingFor:   wait.ForHTTP("/").WithStartupTimeout(timeout),
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

	mappedPort, err := container.MappedPort(ctx, "8090")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &appContainer{Container: container, URI: uri}, nil
}

func setupMongo(ctx context.Context) (*mongoContainer, error) {

	fmt.Println("Start setup testcontainers mongodb")
	timeout := 5 * time.Minute // Default timeout

	req := testcontainers.ContainerRequest{
		Image:        "mongo",
		ExposedPorts: []string{"27017"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(timeout),
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

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &mongoContainer{Container: container, URI: uri}, nil
}

func TestIntegrationMonggitoLatestReturn(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	mongoC, err := setupMongo(ctx)
	if err != nil {
		t.Fatal(err)
	}

	appC, err := setupApp(ctx)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Running Tests...")

	// Clean up the container after the test is complete
	defer appC.Terminate(ctx)
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

	mongoC, err := setupMongo(ctx)
	if err != nil {
		fmt.Println(err)
	}

	appC, err := setupApp(ctx)
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
