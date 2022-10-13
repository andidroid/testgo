package main

import (
	"fmt"

	"github.com/andidroid/testgo/pkg/server"
)

func main() {
	fmt.Println("Run RoutingService Go")

	router := server.CreateRouter()
	server.AddRoutingRoutes(router)

	router.Run(":80")
}
