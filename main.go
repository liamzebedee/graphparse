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
	"strconv"
	"sort"
	// "github.com/liamzebedee/graphparse/graphparse"
	"github.com/dcadenas/pagerank"
)


// idea: actually model it like webpages
//       looking for a piece of code? use current scope as starting page
//       use the type system to autofill the vars? 

// other thing with VR:
// need a visual shape-based constraint/design language


// two approaches:
// - big global canonical lookup table of scope (bad, error-prone)
// - small recursive progressive graph AST approach (better)



type node struct {
	children []*node
	value interface{}
	label string
}
func (n *node) AddChild(val interface{}, label string) *node {
	child := &node{value: val, label: label}
	n.children = append(n.children, child)
	return child
}

// type Sortable interface {
// 	Len() int
// 	Less(i, j int) bool
// 	Swap(i, j int)
// }

// func sortMap(data valueSortedMap) valueSortedMap {
// 	sort.Sort(sort.Reverse(pl))
// 	return pl
// }

// type valueSortedMap 
type valueSortedMap map[nodeid]float64

type pair struct {
	k nodeid
	v float64
}

func sortMap(from map[nodeid]float64) map[nodeid]float64 {
	tmp := []pair
	for k, v := range from, i := 0 {
		tmp = append(tmp, pair{k, v})
	}
}
func (this valueSortedMap) Len() int {
	return len(this)
}
func (this valueSortedMap) Less(i, j int) bool {
	return this[i] < this[j]
}
func (this valueSortedMap) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}


// we don't have to deal with duplicate words just yet
// lets just use the addresses? 

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
	// maxNodeSize := 10.0 // inches
	// min node size is 1
	ranks := make(valueSortedMap)
	graph.Rank(probability_of_following_a_link, tolerance, func(identifier int, rank float64) {
		ranks[nodeid(identifier)] = rank
	})

	// normalise ranks to something that is nice to look at
	maxNodeSize := 3. // inches
	sort.Sort(ranks)

	scaleFactor := maxNodeSize / ranks[-1]
	for identifier, rank := range ranks {
		node := nodesLookup[nodeid(identifier)]
		rankStretched := rank * scaleFactor
		fmt.Fprintf(w, "%v [width=%v] [height=%v] [label=\"%v\"];\n", identifier, rankStretched, rankStretched, node.label)
	}


	// 2. Edges
	for _, edge := range edges {
		fmt.Fprintf(w, "\"%v\" -> \"%v\";\n", edge[0], edge[1])
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

func (this *node) Id() nodeid {
	if i, err := strconv.ParseInt(fmt.Sprintf("%p", this), 0, 64); err != nil {
		panic(err)
	} else { return nodeid(i) }
	
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
	parent ast.Node
	parentLabel string
	Graph *node
}

func NewVisitor(rootAstNode ast.Node) Visitor {
	v := Visitor{}
	// v.graph = graphparse.NewGraph()
	v.Graph = &node{value: rootAstNode, label: "/"}
	// v.parent = rootAstNode
	// v.parentLabel = ""
	return v
}

func (v Visitor) goDeeper(node ast.Node, label string) (w Visitor) {
	w = v
	w.Graph = w.Graph.AddChild(node, label)
	return w
}

// will every child only belong to one parent? 
// yeah, it's pagerank
// maybe in v2 we will have variables linked based on lexic similarity

func (v Visitor) registerNode(node ast.Node, label string) {
	// v.graph.Add(v.parent, node)
	v.Graph.AddChild(node, label)
}

func (v Visitor) Visit(node ast.Node) (w ast.Visitor) {
	switch x := node.(type) {
		case *ast.Package:
			for abspath, srcfile := range x.Files {
				path, err := filepath.Rel(dir, abspath)
				if err != nil {
					panic(err)
				}

				// pkgName := x.Name

				// label := "(" + path + ")"
				label := path

				ast.Walk(v.goDeeper(srcfile, label), srcfile)
				fmt.Println("Processing", path)
			}

		case *ast.TypeSpec:
			return v.goDeeper(x, x.Name.Name)
		
		case *ast.StructType:
			return v.goDeeper(x, x.Fields.List[0].Names[0].Name)

		case *ast.FuncDecl:
			if recv := x.Recv; recv != nil {
				// Method on struct
				structName := recv.List[0].Names[0].Name

				return v.goDeeper(x, structName)
			} else {
				// Just a function
				return v.goDeeper(x, x.Name.Name)
			}

		case *ast.ValueSpec:
			v.registerNode(x, x.Names[0].Name)


		// case *ast.Ident:
			// v.registerNode(x, x.Name)
			// fmt.Println(strings.Join(v.parents, ".") + "." + x.Name)

		default:
			return v
	}
	return w
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
		visitor.Graph.ToDot()
	}

}