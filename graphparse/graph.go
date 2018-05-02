package graphparse

import (
	"io"
	"bufio"
	"fmt"
	"os"
	"sort"
	"path/filepath"
	"encoding/json"
	"bytes"
	"github.com/dcadenas/pagerank"
	// "strings"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"go/types"
	"time"
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

type EdgeType int
const (
	Use EdgeType = iota
	Def
)

type edge struct {
	from Node
	to Node
	variant EdgeType
	id int
}

var lastEdgeId = 0

func newEdge(from Node, to Node, variant EdgeType) edge {
	lastEdgeId++
	return edge{
		from,
		to,
		variant,
		lastEdgeId,
	}
}

func (e edge) String() string {
	return fmt.Sprintf("%s -> %s", e.from.String(), e.to.String())
}

type graph struct {
	edges []edge
	generatedAt string
}

func NewGraph() *graph {
	return &graph{
	}
}

func (g *graph) ToString(w io.Writer) {
	for _, n := range nodeLookup {
		w.Write([]byte(n.String() + "\n"))
	}
	w.Write([]byte("\n"))
	for _, e := range g.edges {
		w.Write([]byte(e.String() + "\n"))
	}
}

func (g *graph) processUnknownReferences() {
	var missing []Node
	var edges []edge

	for _, e := range g.edges {
		if node, ok := nodeLookup[e.from.Id()]; !ok {
			missing = append(missing, node)
			continue
		}
		if node, ok := nodeLookup[e.to.Id()]; !ok {
			missing = append(missing, node)
			continue
		}
		edges = append(edges, e)
	}

	fmt.Println(len(missing), "missing nodes")
	fmt.Println(missing)

	g.edges = edges
}

func lookupObjectNode(obj types.Object) (*objNode) {
	for _, n := range nodeLookup {
		if objNode, ok := n.(*objNode); ok {
			if objNode.obj == obj {
				return objNode
			}
		}
	}

	return nil
}

func pathEnclosingNodes(g *graph,  a, b Node) ([]edge) {
	// g2 := simple.NewDirectedGraph()
	g2 := simple.NewUndirectedGraph()
	for _, n := range nodeLookup {
		g2.AddNode(n)
	}
	for _, e := range g.edges {
		g2.SetEdge(e)
	}


	pathsTree := path.DijkstraAllPaths(g2)
	paths, _ := pathsTree.AllBetween(a, b)

	edges := []edge{}

	for _, path := range paths {
		for i, n := range path {
			if i == 0 || i == len(path) {
				continue
			}
			e := edge{
				from: path[i - 1].(Node),
				to: n.(Node),
			}
			edges = append(edges, e)
		}
	}
	
	return edges
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

	fmt.Println("Contracted", len(g.edges), "to", len(edgesArr), "edges")

	return edgesArr
}


func (g *graph) AddEdge(from, to Node, variant EdgeType) {
	if from == nil {
		panic("from node must be non-nil")
	}
	if to == nil {
		panic("to node must be non-nil")
	}
	// if _, ok := nodeLookup[from.Id()]; !ok {
	// 	ParserLog.Fatalln("from node doesn't exist, cannot add edge", from)
	// }
	// if _, ok := nodeLookup[to.Id()]; !ok {
	// 	ParserLog.Fatalln("to node doesn't exist, cannot add edge", to)
	// }
	e := newEdge(from, to, variant)
	g.edges = append(g.edges, e)
}

func (g *graph) markGenerated() {
	g.generatedAt = time.Now().String()
}

func (g *graph) parents(n Node) (l []Node) {
	for _, edge := range g.edges {
		if edge.to == n {
			l = append(l, edge.to)
		}
	}
	return l
}

func (g *graph) children(n Node) (l []Node) {
	for _, edge := range g.edges {
		if edge.from == n {
			l = append(l, edge.from)
		}
	}
	return l
}


func (g *graph) computeNodeRanks(edges []edge) ranksMap {
	if len(edges) == 0 {
		panic("no edges given, can't compute ranks")
	}
	fmt.Println("computing ranks for", len(nodeLookup), "nodes (", len(edges), ")")

	// Compute PageRank distribution
	graph := pagerank.New()
	for _, edge := range edges {
		graph.Link(int(edge.from.Id()), int(edge.to.Id()))
	}

	probability_of_following_a_link := 0.2
	tolerance := 0.0001

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
		node, ok := nodeLookup[id]
		if !ok {
			panic(id)
		}
		fn(node, rank)
	}
}


func (g *graph) mapEdges(edges []edge, fn func(e edge)) {
	for _, edge := range edges {
		fn(edge)
	}
}

func (g *graph) WriteDotToFile() {
	dotfilePath, _ := filepath.Abs("./www/graph.dot")
	f, err := os.Create(dotfilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	g.ToDot(f)
}

func (g *graph) ToDot(w io.Writer) {
	buf := bufio.NewWriter(w)
	defer buf.Flush()
	
	// edges := g.contractEdges(func(n Node) bool {
	// 	switch(n.Variant()) {
	// 	case Struct, RootPackage:
	// 		return false
	// 	default:
	// 		return true
	// 	}
	// })
	edges := g.edges
	ranks := g.computeNodeRanks(edges)
	
	fmt.Fprintf(w, "digraph graphname {\n")
	
	g.mapRanks(ranks, func(node Node, rank float64) {
		fmt.Fprintf(w, "%v [width=%v] [height=%v] [label=\"%v\"];\n", node.Id(), rank, rank, node.Label())
	})

	g.mapEdges(edges, func(edge edge) {
		fmt.Fprintf(w, "\"%v\" -> \"%v\";\n", edge.from.Id(), edge.to.Id())
	})

	fmt.Fprintf(w, "}\n")
}


type jsonNodeDef struct {
	Rank float64 `json:"rank"`
	Label string `json:"label"`
	Id nodeid `json:"id"`
	Variant NodeType `json:"variant"`
	Pos string
	DebugInfo string `json:"debugInfo"`
}
type jsonNodeEdge struct {
	From nodeid      `json:"source"`
	To nodeid        `json:"target"`
	Variant EdgeType `json:"variant"`
	Id int           `json:"id"`
}
type jsonGraph struct {
	NodesLookup map[nodeid]jsonNodeDef `json:"nodeLookup"`
	Nodes []jsonNodeDef `json:"nodes"`
	Edges []jsonNodeEdge `json:"edges"`
	NodeTypes []string `json:"nodeTypes"`
	AdjList map[nodeid][]nodeid `json:"adjList"`
}
func newJsonGraph() jsonGraph {
	return jsonGraph{
		NodesLookup: make(map[nodeid]jsonNodeDef),
		NodeTypes: nodeTypes,
		Edges: []jsonNodeEdge{},
		Nodes: []jsonNodeDef{},
	}
}

func (g *graph) _toJson(edges []edge) jsonGraph {
	jsonGraph := newJsonGraph()
	
	ranks := g.computeNodeRanks(edges)

	g.mapRanks(ranks, func(node Node, rank float64) {
		n := jsonNodeDef{
			Id: node.Id(),
			Rank: rank,
			Label: node.Label(),
			Variant: node.Variant(),
			DebugInfo: node.DebugInfo(),
		}
		jsonGraph.NodesLookup[node.Id()] = n
		jsonGraph.Nodes = append(jsonGraph.Nodes, n)
	})

	g.mapEdges(edges, func(edge edge) {
		jsonGraph.Edges = append(jsonGraph.Edges, jsonNodeEdge{
			edge.from.Id(), 
			edge.to.Id(),
			edge.variant,
			edge.id,
		})
	})

	jsonGraph.AdjList = make(map[nodeid][]nodeid)
	g.mapEdges(edges, func(e edge) {
		jsonGraph.AdjList[e.from.Id()] = append(jsonGraph.AdjList[e.from.Id()], e.to.Id())
	})

	return jsonGraph
}

func (g *graph) toJson() jsonGraph {
	// edges := g.contractEdges(func(n Node) bool {
	// 	// return false
		
	// 	switch(n.Variant()) {
	// 	case Struct, RootPackage:
	// 		return false
	// 	default:
	// 		if n.Variant() == Func && (strings.HasPrefix(n.Label(), "New") || strings.HasPrefix(n.Label(), "new")) {
	// 			return false
	// 		}
	// 		return true
	// 	}
	// })
	return g._toJson(g.edges)
}

func (g *graph) WriteJsonToFile() {
	path, _ := filepath.Abs("./www/graph.json")
	f, err := os.Create(path)
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(f)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	defer w.Flush()

	
	jsonGraph := g.toJson()

	json.NewEncoder(buf).Encode(jsonGraph)
	f.WriteString(buf.String())
}
