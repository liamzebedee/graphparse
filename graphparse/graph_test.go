package graphparse

import (
	// "fmt"
	// "github.com/liamzebedee/graphparse/graphparse"
	// "testing"
	"golang.org/x/tools/go/loader"
	"os"
	"go/ast"

)

// type mockNode struct {
// }
// func (n *mockNode) Id() graphparse.nodeid {
// }

func parseFile(conf loader.Config, path string) *ast.File {
	conf.CreateFromFilenames("testpkg", path)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	astF, err := conf.ParseFile("main.go", f)
	if err != nil {
		panic(err)
	}
	return astF
}


// const pkgPath = "github.com/liamzebedee/graphparse/test/testpkg"
const pkgPath = "github.com/liamzebedee/graphparse/graphparse"


// func TestShortestPath(t *testing.T) {
// 	// create graph
// 	var conf loader.Config

// 	// Create an ad hoc package with path "foo" from
// 	// the specified already-parsed files.
// 	// All ASTs must have the same 'package' declaration.
// 	// conf.CreateFromFiles("foo", parseFile(conf, "./testpkg/main.go"))
// 	conf.Import(pkgPath)

//  	prog, err := conf.Load()

// 	if err != nil {
// 		panic(err)
// 	}

// 	GenerateCodeGraphFromProg(prog, pkgPath)

// 	nodeA := lookupObjectNode(prog.Package(pkgPath).Pkg.Scope().Lookup("graph"))
// 	nodeB := lookupObjectNode(prog.Package(pkgPath).Pkg.Scope().Lookup("pathEnclosingNodes"))		
	
// 	fmt.Println(pathEnclosingNodes(Graph, nodeA, nodeB))

// 	Graph.ToDot(os.Stdout)
	
// 	// create two nodes
// }