package graphparse

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"path/filepath"

	"github.com/dcadenas/pagerank"
)

type node struct {
	value      interface{}
	id         nodeid
	label      string
	extraAttrs string
}

var nodeLookup = make(map[nodeid]node)

func NewNode(value interface{}, id nodeid, label string) *node {
	newNode, ok := nodeLookup[id]

	if !ok {
		newNode = node{value, id, label, ""}
		nodeLookup[id] = newNode
	} else {
		// fmt.Println("Reusing node")
	}

	return &newNode
}

type nodeid int64

type edge []*node

type graph struct {
	edges []edge
}

func NewGraph() *graph {
	return &graph{
	}
}

func (graph *graph) AddEdge(from, to *node) {
	e := edge{from, to}
	graph.edges = append(graph.edges, e)
}

func (this *graph) String() string {
	return ""
}

// Writes the graph into a .dot graph format for web viz
func (this *graph) WriteDotFile() {
	printToStdout := false
	dotfilePath, _ := filepath.Abs("./www/graph.dot")
	f, err := os.Create(dotfilePath)

	if err != nil {
		panic(err)
	}
	if printToStdout {
		f = os.Stdout
	} else {
		defer f.Close()
	}
	w := bufio.NewWriter(f)
	defer w.Flush()

	fmt.Println(len(this.edges), "edges and", len(nodeLookup), "nodes")

	// Compute PageRank distribution
	graph := pagerank.New()
	for _, edge := range this.edges {
		graph.Link(int(edge[0].id), int(edge[1].id))
	}

	probability_of_following_a_link := 0.85
	tolerance := 0.0001

	// Generate .dot file for graphviz
	// ------
	w.WriteString("digraph graphname {\n")

	// 1. Node definitions
	var ranks rankPairList
	graph.Rank(probability_of_following_a_link, tolerance, func(identifier int, rank float64) {
		ranks = append(ranks, rankPair{nodeid(identifier), rank})
	})

	// normalise ranks to something that is nice to look at
	sort.Sort(ranks)
	minAllowed, maxAllowed := 1.0, 6.0
	var min float64
	var max float64
	if len(ranks) == 0 {
		min, max = 0., 1.
	} else {
		min, max = ranks[0].Rank, ranks[len(ranks)-1].Rank
	}

	scaleRank := func(rank float64) float64 {
		return (maxAllowed-minAllowed)*(rank-min)/(max-min) + minAllowed
	}
	// fmt.Println("smallest node is", scaleRank(min))
	// fmt.Println("biggest node is", scaleRank(max))

	for _, rank := range ranks {
		node := nodeLookup[rank.NodeId]
		rankStretched := scaleRank(rank.Rank)
		fmt.Fprintf(w, "%v [width=%v] [height=%v] [label=\"%v\"] %v;\n", rank.NodeId, rankStretched, rankStretched, node.label, node.extraAttrs)
	}

	// 2. Edges
	for _, edge := range this.edges {
		fmt.Fprintf(w, "\"%v\" -> \"%v\";\n", edge[1].id, edge[0].id)
	}

	w.WriteString("}\n")
}
