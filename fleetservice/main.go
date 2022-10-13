package main

import (
	"fmt"

	"github.com/andidroid/testgo/pkg/server"
)

func main() {
	fmt.Println("Run FleetService Go")

	router := server.CreateRouter()
	server.AddFleetRoutes(router)

	router.Run(":80")
}
