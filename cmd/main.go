package main

import (
	"fmt"
	"sync"

	// "time"
	//"channel"

	// "github.com/andidroid/testgo/pkg/kafka"

	"github.com/andidroid/testgo/pkg/server"
)

var wg = sync.WaitGroup{}

func main() {
	fmt.Println("Run Test Go")
	server.Start()
}
