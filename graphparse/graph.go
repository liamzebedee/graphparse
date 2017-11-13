package graphparse

import (
	"fmt"
	"sort"
	"bufio"
	"os"

	
	"github.com/dcadenas/pagerank"
)

type node struct {
	children []*node
	value interface{}
	label string
	parent *node
	id nodeid
}

type nodeid int64

type edge []nodeid

func (n *node) AddChild(val interface{}, label string, id nodeid) *node {
	child := node{}
	child.value = val
	child.label = label
	child.parent = n
	child.id = nodeid
	n.children = append(n.children, &child)
	return &child
}

func (this *node) Id() nodeid {
	// if i, err := strconv.ParseInt(fmt.Sprintf("%p", this), 0, 64); err != nil {
		// panic(err)
	// } else { return nodeid(i) }
	

	// if this.label == "c" {
	// 	fmt.Println(this.parent.label)
	// }

	// fmt.Println(fqn)

	// return getIdForName(this.label)
	// return getIdForName(this.FQN())
	return this.nodeid
}

func (this *node) ToDot() {
	printToStdout := false


	f, err := os.OpenFile("/Users/liamz/parser/src/github.com/liamzebedee/graphparse/graph.dot", os.O_WRONLY, 0600)
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
	var edges []edge
	nodesLookup := make(map[nodeid]*node)

	this.walk(makeEdges(&edges))
	// this.walk(makeNodeLookup(nodesLookup))

	fmt.Println(len(edges), "edges and", len(nodesLookup), "nodes")
	
	// Compute PageRank distribution
	graph := pagerank.New()
	for _, edge := range edges {
		graph.Link(int(edge[0]), int(edge[1]))
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
	minAllowed, maxAllowed := 1.0, 5.0
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
	fmt.Println("smallest node is", scaleRank(min))
	fmt.Println("biggest node is", scaleRank(max))

	for _, rank := range ranks {
		node := nodesLookup[rank.NodeId]
		rankStretched := scaleRank(rank.Rank)

		fmt.Fprintf(w, "%v [width=%v] [height=%v] [label=\"%v\"];\n", rank.NodeId, rankStretched, rankStretched, node.label)
	}

	// 2. Edges
	for _, edge := range edges {
		fmt.Fprintf(w, "\"%v\" -> \"%v\";\n", edge[0], edge[1])
	}


	// for _, node := range nodesLookup {
	// 	fmt.Println(node.FQN())
	// }

	w.WriteString("}\n")
}

type walkfn func(n *node)

func (n *node) walk(fn walkfn) {
	for _, child := range n.children {
		fn(child)
		child.walk(fn)
	}	
}




// var idmap map[string]nodeid = make(map[string]nodeid)
// var idmap_i int = 0

// func getIdForName(id string) nodeid {
// 	if _, ok := idmap[id]; !ok {
// 		idmap_i++
// 		idmap[id] = nodeid(idmap_i)
// 	}
// 	return idmap[id]
// }



// // canonical format is as follows:
// // PACKAGE FILE DECLARATIONline
// func (this *node) FQN() string {
// 	fqn := ""

// 	p := this
// 	for p != nil && p.label != "/" {
// 		fqn = p.label + "/" + fqn
// 		p = p.parent
// 	}

// 	return fqn
// }