package main

import (
	"fmt"
	"go/parser"
	"go/ast"
	"go/token"
	// "io"
	// "strings"
	"path/filepath"
	// "io/ioutil"
	"bufio"
	"os"
	// "strconv"
	"sort"
	// "github.com/liamzebedee/graphparse/graphparse"
	"github.com/dcadenas/pagerank"
	// godsUtils "github.com/emirpasic/gods/utils" 
	// "github.com/emirpasic/gods/trees/btree"
	// "github.com/emirpasic/gods/maps/treemap"
)


// idea: actually model it like webpages
//       looking for a piece of code? use current scope as starting page
//       use the type system to autofill the vars? 

// other thing with VR:
// need a visual shape-based constraint/design language


// two approaches:
// - big global canonical lookup table of scope (bad, error-prone)
// - small recursive progressive graph AST approach (better)



// next goal: get methods on structs to be linked back to the struct type
// really what this entails is changing the canonicalisation of the id


var idmap map[string]nodeid = make(map[string]nodeid)
var idmap_i int = 0

func getIdForName(id string) nodeid {
	if _, ok := idmap[id]; !ok {
		idmap_i++
		idmap[id] = nodeid(idmap_i)
	}
	return idmap[id]
}

type node struct {
	children []*node
	value interface{}
	label string
	parent *node
}
func (n *node) AddChild(val interface{}, label string) *node {
	child := &node{value: val, label: label, parent: n}
	n.children = append(n.children, child)
	return child
}

// we don't have to deal with duplicate words just yet
// lets just use the addresses? 

type rankPair struct {
	NodeId nodeid
	Rank float64
}

// A slice of pairs that implements sort.Interface to sort by values
type rankPairList []rankPair

func (p rankPairList) Len() int           { return len(p) }
func (p rankPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p rankPairList) Less(i, j int) bool { return p[i].Rank < p[j].Rank }



// now we need a canonical id
// it's much better to use the type system to help with this
// rather than defining a lookup map which is simply difficult

// so we can use each file's scope

// scope = file.Scope
// ident = x.Name // or ident, idk
// for ; scope != nil; scope = scope.Outer {
// 	if obj := scope.Lookup(ident.Name); obj != nil {
// 		ident.Obj = obj
// 		return true
// 	}
// }


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
	this.walk(makeNodeLookup(nodesLookup))

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


	for _, node := range nodesLookup {
		fmt.Println(node.FQN())
	}

	w.WriteString("}\n")
}

type walkfn func(n *node)

func (n *node) walk(fn walkfn) {
	for _, child := range n.children {
		fn(child)
		child.walk(fn)
	}	
}

type nodeid int64
type edge []nodeid

// canonical format is as follows:
// PACKAGE FILE DECLARATIONline

func (this *node) FQN() string {
	fqn := ""

	p := this
	for p != nil && p.label != "/" {
		fqn = p.label + "/" + fqn
		p = p.parent
	}

	return fqn
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
	return getIdForName(this.FQN())
}

// each child of children
func makeEdges(edges *[]edge) walkfn {
	return func(node *node) {
		for _, child := range node.children {
			*edges = append(*edges, edge{node.Id(), child.Id()})
		}
	}
}

func makeNodeLookup(nodeLookup map[nodeid]*node) walkfn {
	return func(node *node) {
		nodeLookup[node.Id()] = node
	}
}



type Visitor struct {
	Graph *node
}

func NewVisitor(rootAstNode ast.Node) Visitor {
	v := Visitor{}
	v.Graph = &node{parent: nil, value: rootAstNode, label: "/"}
	return v
}

func (v Visitor) goDeeper(node ast.Node, label string) (w Visitor) {
	w = v
	w.Graph = w.Graph.AddChild(node, label)
	return w
}

func (v Visitor) registerNode(node ast.Node, label string) {
	v.Graph.AddChild(node, label)
}

// will every child only belong to one parent? 
// yeah, it's pagerank
// maybe in v2 we will have variables linked based on lexic similarity


// next step: fix entries that have no .go file as canonical path ?? (wtf?!)


func (v Visitor) Visit(node ast.Node) (w ast.Visitor) {
	switch x := node.(type) {
		case *ast.Package:
			for abspath, srcfile := range x.Files {
				path, err := filepath.Rel(dir, abspath)
				if err != nil {
					panic(err)
				}

				ast.Walk(v.goDeeper(srcfile, path), srcfile)
				fmt.Println("Processing", path)
			}

		default:
			return w
	}
	return v
}

const dir string = "/Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/"

func main() {
	fset := token.NewFileSet()
	// dir := "/Users/liamz/parser/src/github.com/liamzebedee/graphparse/testsrc/"

	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		fmt.Println(err)
	}

	for name, pkg := range pkgs {
		fmt.Println("Package:", name)

		visitor := NewVisitor(pkg)
		ast.Walk(visitor, pkg)
		// ast.Print(fset, pkg)
		visitor.Graph.ToDot()
	}

}