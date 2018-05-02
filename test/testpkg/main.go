package testpkg

import (
	"log"
	"os"
)

type mockTypeAlias = int64

var logger = log.New(os.Stdout, "", 0)

func main() {
	s, err := NewServer("localhost", 12345)
	// logger.Println("starting up...")

	if err != nil {
		panic(err)
	}
	go s.Listen()
}