package graphparse

import (
	// "flag"
	"fmt"
)

// Stuff does everything.
func Stuff() {
	// runApi := flag.Bool("api", false, "run API server for web experiment")
	// flag.Parse()

	GenerateCodeGraph("github.com/liamzebedee/graphparse/graphparse")
	// GenerateCodeGraph("github.com/twitchyliquid64/subnet/subnet")
	// GenerateCodeGraph("github.com/btcsuite/btcd/blockchain")
	Graph.ToJson()
	// Graph.ToDot()

	fmt.Println("Running web api on port 8081...")
	WebAPI("8081")
}
