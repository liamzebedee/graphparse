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
	"strings"
)


var nodeLookup = make(map[nodeid]Node)
var rootNode Node

func addNodeToLookup(n Node) {
	nodeLookup[n.Id()] = n
	if rootNode == nil && n.Variant() == RootPackage {
		rootNode = n
	}
}

type ranksMap = map[nodeid]float64

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
	String() string
}

type objNode struct {
	obj types.Object
	baseNode
}

type objlookup struct {
	id string
}
var objLookups = make(map[string]*objlookup)

// The Id of an object node is defined canonically 
// as the token.Pos of where the type is declared
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
	
	panic("can't get id")
}
func (n *objNode) Label() string {
	return n.baseNode.label
}
func (n *objNode) Variant() NodeType {
	return n.baseNode.variant
}
func (n *objNode) String() string {
	return fmt.Sprintf("%s <%s>", n.Label(), nodeTypes[n.Variant()])
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
func (n *canonicalNode) String() string {
	return fmt.Sprintf("%s <%s>", n.Label(), nodeTypes[n.Variant()])
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

func (e edge) String() string {
	return fmt.Sprintf("%s -> %s", e.from.String(), e.to.String())
}

type graph struct {
	edges []edge
}

func NewGraph() *graph {
	return &graph{
	}
}



func pathEnclosingNodes(a, b Node) []edge {

}

// Returns indices of edges that match cond
func filterEdges(edges []edge, cond func(edge) bool) []int {
	l := []int{}
	for i, e := range edges {
		if cond(e) {
			l = append(l, i)
		}
	}
	return l
}

func (g *graph) contractEdges(shouldContract func(Node) bool) []edge {
	edges := newArrayList()
	for _, v := range g.edges {
		fmt.Println(v.String())
		edges.append(v)
	}

	getUniqueNodes := func(edges *arraylist) []Node {
		unique := make(map[nodeid]Node)
		nodes := []Node{}
		for _, e := range edges.Map() {
			a := e.(edge).from
			b := e.(edge).to
			unique[a.Id()] = a
			unique[b.Id()] = b
		}
		for _, n := range unique {
			nodes = append(nodes, n)
		}
		return nodes
	}

	getNodesToContract := func(edges *arraylist) []Node {
		nodesToContract := []Node{}
		for _, n := range getUniqueNodes(edges) {
			if shouldContract(n) {
				nodesToContract = append(nodesToContract, n)
			}
		}
		return nodesToContract
	}

	// get all edges of this node
	filterEdges := func(edges *arraylist, cond func(e edge) bool) []int {
		l := []int{}
		for i, e := range edges.Map() {
			if cond(e.(edge)) {
				l = append(l, i)
			}
		}
		return l
	}


	nodesToContract := getNodesToContract(edges)
	fmt.Println( "contracting", len(nodesToContract), "nodes")

	for _, n := range nodesToContract {
		// fmt.Println(n.String())

		// for this node N
		// [A,B,C] -> N -> [E,F,G]
		//   ins             outs

		// do this
		// [A] -> [E,F,G]
		// [B] -> [E,F,G]
		// [C] -> [E,F,G]

		inEdgeIdxs := filterEdges(edges, func(e edge) bool {
			return e.to.Id() == n.Id()
		})
		outEdgeIdxs := filterEdges(edges, func(e edge) bool {
			return e.from.Id() == n.Id()
		})
		

		for _, in := range inEdgeIdxs {
			if len(outEdgeIdxs) > 0 {
				inEdge := edges.get(in).(edge)

				for _, out := range outEdgeIdxs {
					outEdge := edges.get(out).(edge)

					edges.append(edge{
						from: inEdge.from,
						to: outEdge.to,
					})
				}
			}

			edges.delete(in)
		}

		if len(outEdgeIdxs) > 0 {
			for _, out := range outEdgeIdxs {
				edges.delete(out)
			}
		}

	}

	edgesArr := []edge{}
	for _, v := range edges.Map() {
		edgesArr = append(edgesArr, v.(edge))
	}

	fmt.Println("Reduced from", len(g.edges), "to", len(edgesArr), "edges")

	return edgesArr
}


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


func (g *graph) computeNodeRanks(edges []edge) ranksMap {
	fmt.Println(len(edges), "edges and", len(nodeLookup), "nodes")

	// Compute PageRank distribution
	graph := pagerank.New()
	for _, edge := range edges {
		graph.Link(int(edge.from.Id()), int(edge.to.Id()))
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
	if len(ranks) == 0 {
		panic("no ranks")
	}
	min, max := ranks[0].Rank, ranks[len(ranks)-1].Rank

	scaleRank := func(rank float64) float64 {
		return (maxSize - minSize)  *  (rank - min)/(max - min) + minSize
	}
	fmt.Println("smallest node is", min, scaleRank(min))
	fmt.Println("biggest node is", max, scaleRank(max))

	rankMap := make(ranksMap)

	for _, rank := range ranks {
		rankScaled := scaleRank(rank.Rank)
		rankMap[rank.NodeId] = rankScaled
	}

	return rankMap
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

func (g *graph) mapEdges2(edges []edge, ranks ranksMap, fn func(e edge)) {
	for _, edge := range edges {
		fn(edge)
	}
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

	
	edges := g.contractEdges(func(n Node) bool {
		switch(n.Variant()) {
		case Struct, RootPackage:
			return false
		default:
			return true
		}
	})
	ranks := g.computeNodeRanks(edges)
	
	w.WriteString("digraph graphname {\n")
	
	g.mapRanks(ranks, func(node Node, rank float64) {
		fmt.Fprintf(w, "%v [width=%v] [height=%v] [label=\"%v\"];\n", node.Id(), rank, rank, node.Label())
	})

	g.mapEdges2(edges, ranks, func(edge edge) {
		fmt.Fprintf(w, "\"%v\" -> \"%v\";\n", edge.from.Id(), edge.to.Id())
	})

	w.WriteString("}")
}


type jsonNodeDef struct {
	Rank float64 `json:"rank"`
	Label string `json:"label"`
	Id nodeid `json:"id"`
	Variant NodeType `json:"variant"`
	Pos string
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
	edges := g.contractEdges(func(n Node) bool {
		switch(n.Variant()) {
		case Struct, RootPackage:
			return false
		default:
			if n.Variant() == Func && (strings.HasPrefix(n.Label(), "New") || strings.HasPrefix(n.Label(), "new")) {
				return false
			}
			return true
		}
	})
	ranks := g.computeNodeRanks(edges)

	g.mapRanks(ranks, func(node Node, rank float64) {
		n := jsonNodeDef{
			Id: node.Id(),
			Rank: rank,
			Label: node.Label(),
			Variant: node.Variant(),
		}
		jsonGraph.NodesLookup[node.Id()] = n
		jsonGraph.Nodes = append(jsonGraph.Nodes, n)
	})

	g.mapEdges(ranks, func(edge edge) {
		jsonGraph.Edges = append(jsonGraph.Edges, jsonNodeEdge{
			edge.from.Id(), 
			edge.to.Id(),
		})
	})

	json.NewEncoder(buf).Encode(jsonGraph)
	f.WriteString(buf.String())
}
