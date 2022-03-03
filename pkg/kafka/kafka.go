package kafka

import (
	"fmt"
)

func main() {
	fmt.Println("start kafka test")

	go Reader(nil)
	Writer(nil)
}
