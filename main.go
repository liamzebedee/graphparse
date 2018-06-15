package main

import (
	"fmt"
	"github.com/liamzebedee/graphparse/graphparse"
	"github.com/liamzebedee/graphparse/graphparse/webapi"
)

var Graph1 *graphparse.Graph

// Stuff does everything.
func main() {
	// Graph1, err := graphparse.GenerateCodeGraph("github.com/liamzebedee/graphparse/graphparse")
	// if err != nil {
	// 	panic(err)
	// }
	// GenerateCodeGraph("github.com/twitchyliquid64/subnet/subnet")
	// graph := GenerateCodeGraph("github.com/btcsuite/btcd/blockchain")
	// Graph1.WriteDotToFile("./www/graph.dot")
	// Graph1.WriteJsonToFile("./www/graph.json")
	fmt.Println("Running web api on port 8082...")
	webapi.WebAPI("8082")
}