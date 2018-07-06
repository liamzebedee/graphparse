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
	// "gonum.org/v1/gonum/graph/path"
	// "gonum.org/v1/gonum/graph/simple"
	// "github.com/gonum/gonum/tree/master/graph/path"

	"go/types"
	"time"
)



func (g *Graph) addNodeToLookup(n Node) {
	g.nodeLookup[n.Id()] = n
	if g.rootNode == nil && n.Variant() == RootPackage {
		g.rootNode = n
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


type Graph struct {
	edges []edge
	generatedAt string
	objNodeLookup map[nodeid]*objNode
	nodeLookup map[nodeid]Node
	rootNode Node
}

func NewGraph() *Graph {
	return &Graph{
		objNodeLookup: make(map[nodeid]*objNode),
		nodeLookup: make(map[nodeid]Node),
	}
}

func (g *Graph) nodeExists(id int64) bool {
	for _, n := range g.Nodes() {
		if n.ID() == id {
			return true
		}
	}
	return false
}

func (g *Graph) ToString(w io.Writer) {
	for _, n := range g.nodeLookup {
		w.Write([]byte(n.String() + "\n"))
	}
	w.Write([]byte("\n"))
	for _, e := range g.edges {
		w.Write([]byte(e.String() + "\n"))
	}
}

func (g *Graph) lookupObjectNode(obj types.Object) (*objNode) {
	for _, n := range g.nodeLookup {
		if objNode, ok := n.(*objNode); ok {
			if objNode.obj == obj {
				return objNode
			}
		}
	}

	return nil
}

// func pathEnclosingNodes(g *Graph,  a, b Node) ([]edge) {
// 	// g2 := simple.NewDirectedGraph()
// 	g2 := simple.NewUndirectedGraph()
// 	for _, n := range g.nodeLookup {
// 		g2.AddNode(n)
// 	}
// 	for _, e := range g.edges {
// 		g2.SetEdge(e)
// 	}


// 	pathsTree := path.DijkstraAllPaths(g2)
// 	paths, _ := pathsTree.AllBetween(a, b)

// 	edges := []edge{}

// 	for _, path := range paths {
// 		for i, n := range path {
// 			if i == 0 || i == len(path) {
// 				continue
// 			}
// 			e := edge{
// 				from: path[i - 1].(Node),
// 				to: n.(Node),
// 			}
// 			edges = append(edges, e)
// 		}
// 	}
	
// 	return edges
// }

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

func (g *Graph) contractEdges(shouldContract func(Node) bool) []edge {
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


func (g *Graph) AddEdge(from, to Node, variant EdgeType) {
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

func (g *Graph) markGenerated() {
	g.generatedAt = time.Now().String()
}

func (g *Graph) parents(n Node) (l []Node) {
	for _, edge := range g.edges {
		if edge.to == n {
			l = append(l, edge.to)
		}
	}
	return l
}

func (g *Graph) children(n Node) (l []Node) {
	for _, edge := range g.edges {
		if edge.from == n {
			l = append(l, edge.from)
		}
	}
	return l
}


func (g *Graph) computeNodeRanks(edges []edge) ranksMap {
	if len(edges) == 0 {
		panic("no edges given, can't compute ranks")
	}
	fmt.Println("computing ranks for", len(g.nodeLookup), "nodes (", len(edges), " edges)")

	// Compute PageRank distribution
	Graph := pagerank.New()
	for _, edge := range edges {
		Graph.Link(int(edge.from.Id()), int(edge.to.Id()))
	}

	probability_of_following_a_link := 0.2
	tolerance := 0.0001

	var ranks rankPairList
	Graph.Rank(probability_of_following_a_link, tolerance, func(identifier int, rank float64) {
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
func (g *Graph) mapRanks(ranks ranksMap, fn func(n Node, rank float64)) {
	for id, rank := range ranks {
		node, ok := g.nodeLookup[id]
		if !ok {
			// panic(id)
			fmt.Println("Couldn't find node for id during ranks:", id)
			continue
		}
		fn(node, rank)
	}
}


func (g *Graph) mapEdges(edges []edge, fn func(e edge)) {
	for _, edge := range edges {
		fn(edge)
	}
}

func (g *Graph) WriteDotToFile(path string) {
	dotfilePath, _ := filepath.Abs(path)
	f, err := os.Create(dotfilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	g.ToDot(f)
}

func (g *Graph) ToDot(w io.Writer) {
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
	// DebugInfo string `json:"debugInfo"`
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
	RootNode nodeid `json:"rootNode"`
}
func newJsonGraph() jsonGraph {
	return jsonGraph{
		NodesLookup: make(map[nodeid]jsonNodeDef),
		NodeTypes: nodeTypes,
		Edges: []jsonNodeEdge{},
		Nodes: []jsonNodeDef{},
	}
}

func (g *Graph) _toJson(edges []edge) jsonGraph {
	jsonGraph := newJsonGraph()
	jsonGraph.RootNode = g.rootNode.Id()
	
	ranks := g.computeNodeRanks(edges)

	g.mapRanks(ranks, func(node Node, rank float64) {
		n := jsonNodeDef{
			Id: node.Id(),
			Rank: rank,
			Label: node.Label(),
			Variant: node.Variant(),
			// DebugInfo: node.DebugInfo(),
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

func (g *Graph) ToJson() jsonGraph {
	return g._toJson(g.edges)
}

func (g *Graph) WriteJsonToFile(relpath string) {
	path, _ := filepath.Abs(relpath)
	f, err := os.Create(path)
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(f)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	defer w.Flush()

	
	jsonGraph := g.ToJson()

	json.NewEncoder(buf).Encode(jsonGraph)
	f.WriteString(buf.String())
}
