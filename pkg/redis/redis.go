package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	//docker run -d -p 6379:6379 redis:7.0-rc1-alpine

	ctx := context.Background()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	status := redisClient.Ping(ctx)
	fmt.Println(status)
	// mongodbHost, ok := os.LookupEnv("MONGODB_HOST")
	// 	if !ok {
	// 		mongodbHost = "localhost"
	// 		os.Setenv("MONGODB_HOST", "localhost")
	// 	}
	// 	fmt.Printf("MONGODB_HOST: %s\n", mongodbHost)

	// 	mongodbPort, ok := os.LookupEnv("MONGODB_PORT")
	// 	if !ok {
	// 		mongodbPort = "27017"
	// 		os.Setenv("MONGODB_PORT", "27017")
	// 	}
	// 	fmt.Printf("MONGODB_PORT: %s\n", mongodbPort)

	// 	mongodb := os.ExpandEnv("mongodb://$MONGODB_HOST:$MONGODB_PORT")
	// 	fmt.Println("MongoDB URL: ", mongodb)

	// opt, err := redis.ParseURL("redis://<user>:<pass>@localhost:6379/<db>")
	// if err != nil {
	// 	panic(err)
	// }

	// redisClient := redis.NewClient(opt)

	err := redisClient.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := redisClient.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

}
