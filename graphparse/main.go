package graphparse

import (
	"fmt"
)

// Stuff does everything.
func Stuff() {
	GenerateCodeGraph("github.com/liamzebedee/graphparse/graphparse")
	// GenerateCodeGraph("github.com/twitchyliquid64/subnet/subnet")
	// GenerateCodeGraph("github.com/btcsuite/btcd/blockchain")
	
	Graph.WriteDotToFile("./www/graph.dot")
	Graph.WriteJsonToFile("./www/graph.json")
	
	fmt.Println("Running web api on port 8081...")
	WebAPI("8081")
}
