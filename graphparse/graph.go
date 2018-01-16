package graphparse

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"path/filepath"
	"go/types"
	"encoding/json"
	"bytes"
	"github.com/dcadenas/pagerank"
)


var nodeLookup = make(map[nodeid]Node)

func addNodeToLookup(n Node) {
	nodeLookup[n.Id()] = n
}

// how to prune intermediates?
// a -> b -> c
// prune b
// a -> c

// e -> 
// a -> b -> c
//		  -> d -> a
// a -> c

// then

// e -> a
// a -> c
// a -> d

type NodeType int
const (
	Struct NodeType = iota
	Method
	Func
	Field
	RootPackage
	File
	ImportedPackage
	ImportedFunc
	FuncCall
)
var nodeTypes = []string{
	"Struct",
	"Method",
	"Func",
	"Field",
	"RootPackage",
	"File",
	"ImportedPackage",
	"ImportedFunc",
	"FuncCall",
}

type baseNode struct {
	variant NodeType
	label string
}

type Node interface {
	Id() nodeid
	Label() string
	Variant() NodeType
}

type objNode struct {
	obj types.Object
	baseNode
}

type objlookup struct {
	id string
}
var objLookups = make(map[string]*objlookup)
func (n *objNode) Id() nodeid {
	switch n.variant {
	case Struct:
		switch typ := n.obj.Type().(type) {
		case *types.Named:
			id := nodeid(typ.Obj().Pos())
			return id
		case *types.Pointer:
			id := nodeid(typ.Elem().(*types.Named).Obj().Pos())
			return id
		default:
			panic("cant find type of struct Object")
		}
	default:
		objId := n.obj.String()
		x, ok := objLookups[objId]
		if !ok {
			x = &objlookup{objId}
			objLookups[objId] = x
		}
		return pointerToId(x)
	}

	// if objId == "Run" {
	// 	fmt.Println(n.obj.String())
	// }
	// objId := string(n.obj.Pos()) + n.obj.Id()
	
	panic("cant get id")
}
func (n *objNode) Label() string {
	return n.baseNode.label
}
func (n *objNode) Variant() NodeType {
	return n.baseNode.variant
}


func LookupOrCreateNode(obj types.Object, variant NodeType, label string) *objNode {
	if obj == nil {
		panic("obj must be non-nil")
	}
	id := pointerToId(obj)

	node, ok := nodeLookup[id]

	if !ok {
		node = &objNode{
			obj,
			baseNode{variant, label},
		}
		addNodeToLookup(node)
	}

	return node.(*objNode)
}

func LookupNode(obj types.Object) *objNode {
	if obj == nil {
		panic("obj must be non-nil")
	}
	id := pointerToId(obj)

	node, ok := nodeLookup[id]
	if !ok {
		panic("node not found")
	}
	return node.(*objNode)
}

var canonicalNodeLookup = make(map[string]*canonicalNode)

type canonicalNode struct {
	baseNode
}

func (n *canonicalNode) Id() nodeid {
	return pointerToId(n)
}
func (n *canonicalNode) Label() string {
	return n.baseNode.label
}
func (n *canonicalNode) Variant() NodeType {
	return n.baseNode.variant
}

func LookupOrCreateCanonicalNode(key string, variant NodeType, label string) *canonicalNode {
	node, ok := canonicalNodeLookup[key]

	if !ok {
		node = &canonicalNode{
			baseNode{variant, label},
		}
		canonicalNodeLookup[key] = node
		addNodeToLookup(node)
	}

	return node
} 

type nodeid int64

type edge struct {
	from Node
	to Node
}

type graph struct {
	edges []edge
}

func NewGraph() *graph {
	return &graph{
	}
}



// type pointerGraph struct {
// }
// func (g *graph) toPointerGraph() *pointerGraph {

// }

func (g *graph) AddEdge(from, to Node) {
	if from == nil {
		panic("from node must be non-nil")
	}
	if to == nil {
		panic("to node must be non-nil")
	}
	if _, ok := nodeLookup[from.Id()]; !ok {
		panic("from node doesn't exist, cannot add edge")
	}
	if _, ok := nodeLookup[to.Id()]; !ok {
		panic("to node doesn't exist, cannot add edge")
	}
	e := edge{from, to}
	g.edges = append(g.edges, e)
}


func (g *graph) computeNodeRanks() map[nodeid]float64 {
	fmt.Println(len(g.edges), "edges and", len(nodeLookup), "nodes")

	// Compute PageRank distribution
	graph := pagerank.New()
	for _, edge := range g.edges {
		if edge.from.Variant() == RootPackage {	
			graph.Link(int(edge.from.Id()), int(edge.to.Id()))
			// continue
		} else {
			graph.Link(int(edge.from.Id()), int(edge.to.Id()))
		}
	}

	probability_of_following_a_link := 0.9
	tolerance := 0.05

	var ranks rankPairList
	graph.Rank(probability_of_following_a_link, tolerance, func(identifier int, rank float64) {
		ranks = append(ranks, rankPair{nodeid(identifier), rank})
	})

	// normalise ranks to something that is nice to look at
	sort.Sort(ranks)

	minSize, maxSize := 1.0, 6.0
	min, max := ranks[0].Rank, ranks[len(ranks)-1].Rank

	scaleRank := func(rank float64) float64 {
		return (maxSize - minSize)  *  (rank - min)/(max - min) + minSize
	}
	fmt.Println("smallest node is", min, scaleRank(min))
	fmt.Println("biggest node is", max, scaleRank(max))

	rankMap := make(map[nodeid]float64)

	for _, rank := range ranks {
		rankScaled := scaleRank(rank.Rank)
		rankMap[rank.NodeId] = rankScaled
	}

	return rankMap
}


func (g *graph) computeAggregateNodeRanks(ranks map[nodeid]float64) map[nodeid]float64 {
	getNodesToRemove := func(ranks map[nodeid]float64, shouldRemove func(n Node) bool) []nodeid {
		ids := []nodeid{}
		for id, rank := range ranks {
			if rank > 0 && shouldRemove(nodeLookup[id]) {
				ids = append(ids, id)
			}
		}
		return ids
	}

	outEdges := make(map[nodeid][]Node)
	for _, edge := range g.edges {
		if l, ok := outEdges[edge.from.Id()]; !ok {
			outEdges[edge.from.Id()] = []Node{edge.to}
		} else {
			outEdges[edge.from.Id()] = append(l, edge.to)
		}
	}

	removeCond := func(n Node) bool {
		return n.Variant() != Struct 
	}

	iters := 0

	for nodes := getNodesToRemove(ranks, removeCond); len(nodes) > 0; nodes = getNodesToRemove(ranks, removeCond) {
		fmt.Println("aggregating nodes, iteration", iters)
		for _, id := range nodes {
			outs := outEdges[id]

			rankDist := ranks[id] / float64(len(outs))
			for _, out := range outs {
				ranks[out.Id()] += rankDist
			}
			// Technically equivalent to deletion, 
			// since PageRank will never issue a rank of 0.
			ranks[id] = 0
		}
		iters++
	}

	for id, r := range ranks {
		if r == 0 {
			delete(ranks, id)
		}
	}

	return ranks
}


func (g *graph) ToDot() {
	dotfilePath, _ := filepath.Abs("./www/graph.dot")
	f, err := os.Create(dotfilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	
	w := bufio.NewWriter(f)
	defer w.Flush()

	// Generate .dot file for graphviz
	// ------
	w.WriteString("digraph graphname {\n")
	
	ranks := g.computeNodeRanks()
	aggRanks := g.computeAggregateNodeRanks(ranks)
	
	g.mapRanks(aggRanks, func(node Node, rank float64) {
		fmt.Fprintf(w, "%v [width=%v] [height=%v] [label=\"%v\"];\n", node.Id(), rank, rank, node.Label())
	})

	g.mapEdges(aggRanks, func(edge edge) {
		fmt.Fprintf(w, "\"%v\" -> \"%v\";\n", edge.to.Id(), edge.from.Id())
	})

	w.WriteString("}\n")
}


type jsonNodeDef struct {
	Rank float64 `json:"rank"`
	Label string `json:"label"`
	Id nodeid `json:"id"`
	Variant NodeType `json:"variant"`
}
type jsonNodeEdge struct {
	From nodeid `json:"source"`
	To nodeid   `json:"target"`
}
type jsonGraph struct {
	NodesLookup map[nodeid]jsonNodeDef `json:"nodesLookup"`
	Nodes []jsonNodeDef `json:"nodes"`
	Edges []jsonNodeEdge `json:"edges"`
	NodeTypes []string `json:"nodeTypes"`
}
func newJsonGraph() jsonGraph {
	return jsonGraph{
		NodesLookup: make(map[nodeid]jsonNodeDef),
		NodeTypes: nodeTypes,
		Edges: []jsonNodeEdge{},
		Nodes: []jsonNodeDef{},
	}
}

type ranksMap = map[nodeid]float64

func (g *graph) thing() {
	// get file
	// defer
	// create graph serial structure
	// make ranks, aggRanks
	// output node defs
	// output edges
	// write
}

// Maps ranks with node lookup
func (g *graph) mapRanks(ranks ranksMap, fn func(n Node, rank float64)) {
	for id, rank := range ranks {
		fn(nodeLookup[id], rank)
	}
}

// Maps all edges, ignoring those that aren't in ranks.
func (g *graph) mapEdges(ranks ranksMap, fn func(e edge)) {
	for _, edge := range g.edges {
		// Ignore edges that aren't in ranks
		if _, ok := ranks[edge.from.Id()]; !ok {
			continue
		}
		if _, ok := ranks[edge.to.Id()]; !ok {
			continue
		}

		fn(edge)
	}
}

func (g *graph) ToJson() {
	path, _ := filepath.Abs("./www/graph.json")
	f, err := os.Create(path)
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(f)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	defer w.Flush()

	jsonGraph := newJsonGraph()
	ranks := g.computeNodeRanks()
	aggRanks := g.computeAggregateNodeRanks(ranks)

	g.mapRanks(aggRanks, func(node Node, rank float64) {
		n := jsonNodeDef{
			Id: node.Id(),
			Rank: rank,
			Label: node.Label(),
			Variant: node.Variant(),
		}
		jsonGraph.NodesLookup[node.Id()] = n
		jsonGraph.Nodes = append(jsonGraph.Nodes, n)
	})

	g.mapEdges(aggRanks, func(edge edge) {
		jsonGraph.Edges = append(jsonGraph.Edges, jsonNodeEdge{
			edge.to.Id(), 
			edge.from.Id(),
		})
	})

	json.NewEncoder(buf).Encode(jsonGraph)
	f.WriteString(buf.String())
}
