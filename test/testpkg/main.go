package testpkg

import (
	"log"
	"os"
)

type mockTypeAlias = int64

var logger = log.New(os.Stdout, "", 0)

func main() {
	c, err := NewServer("localhost", 12345)
	// logger.Println("starting up...")

	if err != nil {
		panic(err)
	}
	go c.Listen()
}