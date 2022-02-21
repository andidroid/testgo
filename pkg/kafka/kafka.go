package kafka

import (
	"fmt"
)

func main() {
	fmt.Println("start kafka test")
	go Reader()
	Writer()
}
