package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andidroid/testgo/pkg/redis"
	"github.com/andidroid/testgo/pkg/redis/search"
	"github.com/andidroid/testgo/pkg/server/handler"
	"github.com/gin-contrib/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.6.1"

	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"

	// "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testHandler *handler.TestHandler
var streamHandler *handler.EventStreamHandler
var healthHandler *handler.HealthHandler

var truckHandler *handler.TruckHandler
var placeHandler *handler.PlaceHandler
var orderHandler *handler.OrderHandler
var placesearchHandler *handler.SearchPlaceHandler

func init() {

	logger := log.New(os.Stdout, "gin-server", log.LstdFlags)
	logger.Println("start server")

	ctx := context.Background()

	//

	// OTEL_EXPORTER_JAEGER_ENDPOINT default "http://localhost:14268/api/traces"
	// OTEL_EXPORTER_JAEGER_USER
	// OTEL_EXPORTER_JAEGER_PASSWORD
	//jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces"))
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		fmt.Println(err)
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("testgoservice"),
			attribute.String("environment", "dev"),
			attribute.Int64("ID", 42),
		)),
	)

	otel.SetTracerProvider(tp)
	otelmongo.WithTracerProvider(tp)

	//

	mongodbHost, ok := os.LookupEnv("MONGODB_HOST")
	if !ok {
		mongodbHost = "localhost"
		os.Setenv("MONGODB_HOST", "localhost")
	}
	fmt.Printf("MONGODB_HOST: %s\n", mongodbHost)

	mongodbPort, ok := os.LookupEnv("MONGODB_PORT")
	if !ok {
		mongodbPort = "27017"
		os.Setenv("MONGODB_PORT", "27017")
	}
	fmt.Printf("MONGODB_PORT: %s\n", mongodbPort)

	mongodb := os.ExpandEnv("mongodb://$MONGODB_HOST:$MONGODB_PORT/")
	fmt.Println("MongoDB URL: ", mongodb)

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	opts := options.Client()
	opts.Monitor = otelmongo.NewMonitor(otelmongo.WithTracerProvider(tp))
	opts.ApplyURI(mongodb)
	mongoClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	// err = mongoClient.Ping(ctx, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	mongoDatabase := mongoClient.Database("test")

	redisClient := redis.CreateClient()
	redisSearchClient := search.CreateClient()

	//Add distributed tracing to redis client
	//go get github.com/go-redis/redis/extra/redisotel/v8
	// redisClient.AddHook(redisotel.NewTracingHook())

	// status := redisClient.Ping(ctx)
	// fmt.Println(status)

	healthHandler = handler.NewHealthHandler(ctx, logger, mongoDatabase, redisClient)
	testHandler = handler.NewTestHandler(ctx, logger, mongoDatabase, redisClient)
	streamHandler = handler.NewEventStreamHandler()

	truckHandler = handler.NewTruckHandler(ctx, logger, mongoDatabase, redisClient)
	placeHandler = handler.NewPlaceHandler(ctx, logger, mongoDatabase, redisClient)
	orderHandler = handler.NewOrderHandler(ctx, logger, mongoDatabase, redisClient)
	placesearchHandler = handler.NewSearchPlaceHandler(ctx, logger, redisSearchClient)
}

func CreateRouter() *gin.Engine {

	fmt.Println("start server")

	router := gin.Default()

	// custom logger format
	// router.Use(gin.LoggerWithFormatter(func(
	// 	param gin.LogFormatterParams) string {
	// 	return fmt.Sprintf("[%s] %s %s %d %s\n",
	// 		param.TimeStamp.Format("2006-01-02T15:04:05"),
	// 		param.Method,
	// 		param.Path,
	// 		param.StatusCode,
	// 		param.Latency,
	// 	)
	// }))
	// configure file logger
	// 	gin.DisableConsoleColor()
	// f, _ := os.Create("debug.log")
	// gin.DefaultWriter = io.MultiWriter(f)

	// session store in redis, cockies
	// store, _ := redisStore.NewStore(10, "tcp",
	// 	"localhost:6379", "", []byte("secret"))
	// store := sessions.NewCookieStore([]byte("secret"))
	// router.Use(sessions.Sessions("testgo-sessions", store))

	// Add event-streaming headers
	//router.Use(handler.HeadersMiddleware())
	// Initialize new streaming server

	//router.Use(stream.ServeHTTP())

	// routinghandler := handler.NewRoutingHandler()

	// router.GET("/test", testHandler)

	router.Use(otelgin.Middleware("testgo"))

	router.Use(cors.Default())

	router.Use(Logger())
	router.Use(handler.AuthMiddleware())

	router.GET("/health", healthHandler.HandleGetRequest)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// router.POST("/tests", testHandler.HandlePostRequest)
	// router.GET("/tests", testHandler.HandleGetAllTestsRequest)
	// router.PUT("/tests/:id", testHandler.HandlePutRequest)
	// router.DELETE("/tests/:id", testHandler.HandleDeleteRequest)
	// router.GET("/tests/:id", testHandler.HandleGetTestByIdRequest)

	// // Simple group: v2
	// v2 := router.Group("/v2")
	// {
	// 	v2.POST("/login", loginEndpoint)
	// 	v2.POST("/submit", submitEndpoint)
	// 	v2.POST("/read", readEndpoint)
	// }

	router.StaticFile("/map", "./docs/map.html")
	router.StaticFile("/", "./docs/index.html")
	router.Static("/assets", "./docs/assets")
	return router
}

func AddAllRoutes(router *gin.Engine) {
	AddRoutingRoutes(router)

	//fleet service
	// search service
	AddFleetRoutes(router)
	AddSearchRoutes(router)
	AddStreamingRoutes(router)
}

func AddRoutingRoutes(router *gin.Engine) {
	// Simple group: v1
	v1 := router.Group("/routing")
	{
		v1.GET("/tsp", handler.GetTSP)
		v1.GET("/list", handler.GetRouteAsList)
		v1.GET("/geometry", handler.GetRouteAsGeometry)

		v1.GET("/info", handler.GetRouteInformation)
		v1.GET("/poi", handler.GetPOIs)
		v1.GET("/poi/:osm_id/node", handler.GetNearestNode)

	}

	v2 := router.Group("/node")
	{
		v2.GET("/source", handler.GetNodeSearchSource)
		v2.GET("/target", handler.GetNodeSearchTarget)
		v2.GET("/unknown", handler.GetNodeSearchQuery)
		v2.GET("/:id", handler.GetNodeById)
	}
}

func AddFleetRoutes(router *gin.Engine) {
	v3 := router.Group("/fleet")
	{
		v3.GET("/start", handler.HandlePostStartOrderRequest)

		v3.GET("/truck", truckHandler.HandleGetAllTrucksRequest)
		v3.GET("/truck/:id", truckHandler.HandleGetTruckByIdRequest)
		v3.POST("/truck", truckHandler.HandlePostTruckRequest)
		v3.PUT("/truck/:id", truckHandler.HandlePutTruckRequest)
		v3.DELETE("/truck/:id", truckHandler.HandleDeleteTruckRequest)

		v3.GET("/order", orderHandler.HandleGetAllOrdersRequest)
		v3.GET("/order/:id", orderHandler.HandleGetOrderByIdRequest)
		v3.POST("/order", orderHandler.HandlePostOrderRequest)
		v3.PUT("/order/:id", orderHandler.HandlePutOrderRequest)
		v3.DELETE("/order/:id", orderHandler.HandleDeleteOrderRequest)

		v3.GET("/place", placeHandler.HandleGetAllPlacesRequest)
		v3.GET("/place/:id", placeHandler.HandleGetPlaceByIdRequest)
		v3.POST("/place", placeHandler.HandlePostPlaceRequest)
		v3.PUT("/place/:id", placeHandler.HandlePutPlaceRequest)
		v3.DELETE("/place/:id", placeHandler.HandleDeletePlaceRequest)
	}
}

func AddStreamingRoutes(router *gin.Engine) {
	// messaging service
	router.GET("/stream", handler.HeadersMiddleware(), streamHandler.ServeHTTP(), streamHandler.GetPositionStream)

}

func AddSearchRoutes(router *gin.Engine) {
	v3 := router.Group("/search")
	{
		v3.GET("/place", placesearchHandler.HandleGetSearchRequest)
	}
}

func Start() {

	router := CreateRouter()
	AddAllRoutes(router)
	router.Run(":80")
	//router.RunTLS(":443", "certs/localhost.crt", "certs/localhost.key")

}

// router.Use(cors.New(cors.Config{
// 	AllowOrigins:     []string{"http://localhost:80"},
// 	AllowMethods:     []string{"GET","PUT","POST","DELETE", "OPTIONS"},
// 	AllowHeaders:     []string{"Origin"},
// 	ExposeHeaders:    []string{"Content-Length"},
// 	AllowCredentials: true,
// 	MaxAge: 12 * time.Hour,
//  }))

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// before request
		c.Next()
		// after request
		latency := time.Since(t)
		// access the status we are sending
		status := c.Writer.Status()
		log.Println("response:", status, latency)
	}
}
