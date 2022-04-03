package main

import (
	"fmt"

	"github.com/andidroid/testgo/pkg/server"
)

func main() {
	fmt.Println("Run SearchService Go")

	router := server.CreateRouter()
	server.AddSearchRoutes(router)

	router.Run(":80")
}
