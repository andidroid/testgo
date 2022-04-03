package main

import (
	"fmt"

	"github.com/andidroid/testgo/pkg/server"
)

func main() {
	fmt.Println("Run MessagingService Go")

	router := server.CreateRouter()
	server.AddStreamingRoutes(router)

	router.Run(":80")
}
