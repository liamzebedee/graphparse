package graphparse

import (
	"fmt"
	"sort"
	"bufio"
	"os"
	"github.com/dcadenas/pagerank"

)



type node struct {
	value interface{}
	id nodeid
	label string
	extraAttrs string
}

var clientseen bool = false
func NewNode(value interface{}, id nodeid, label string) *node {
	return &node{value, id, label, ""}
}

type nodeid int64

type edge []*node

type graph struct {
	edges []edge
}

func NewGraph() *graph {
	return &graph{}
}

func (this *graph) AddEdge(from, to *node) {
	e := edge{from, to}
	this.edges = append(this.edges, e)
}

func (this *graph) String() string {
	return ""
}

func (this *graph) ToDot() {
	printToStdout := false

	f, err := os.Create("/Users/liamz/parser/src/github.com/liamzebedee/graphparse/graph.dot")
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



	// Compute edges from pointers
	nodesLookup := make(map[nodeid]*node)
	for _, edge := range this.edges {
		node1 := edge[0]
		node2 := edge[1]
		nodesLookup[node1.id] = node1
		nodesLookup[node2.id] = node2
	}
	fmt.Println(len(this.edges), "edges and", len(nodesLookup), "nodes")
	
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
		return (maxAllowed - minAllowed) * (rank - min) / (max - min) + minAllowed;
	}
	// fmt.Println("smallest node is", scaleRank(min))
	// fmt.Println("biggest node is", scaleRank(max))

	for _, rank := range ranks {
		node := nodesLookup[rank.NodeId]
		rankStretched := scaleRank(rank.Rank)
		fmt.Fprintf(w, "%v [width=%v] [height=%v] [label=\"%v\"] %v;\n", rank.NodeId, rankStretched, rankStretched, node.label, node.extraAttrs)
	}

	// 2. Edges
	for _, edge := range this.edges {
		fmt.Fprintf(w, "\"%v\" -> \"%v\";\n", edge[1].id, edge[0].id)
	}


	w.WriteString("}\n")
}






