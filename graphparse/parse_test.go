package graphparse

import (
	"golang.org/x/tools/go/loader"
	"testing"
	"path/filepath"
	"os"
	"github.com/stretchr/testify/assert"
	"bytes"
	"fmt"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func setup() {
	// create graph
	var conf loader.Config


	pkgPath := "github.com/liamzebedee/graphparse/test/testpkg"
	pkgFilePath, _ := filepath.Abs("../test/testpkg")

	conf.Import(pkgPath)

 	prog, err := conf.Load()

	if err != nil {
		panic(err)
	}

	GenerateCodeGraphFromProg(prog, pkgPath, pkgFilePath)
}

func findNode(variant NodeType, label string) Node {
	var found Node

	for _, node := range nodeLookup {
		if(node.Variant() == variant && node.Label() == label) {
			if found != nil {
				// same income, same super.
				panic("two nodes, same variant, same label. seems fishy")
			}
			found = node
		}
	}

	return found
}


func TestParseRootPackage(t *testing.T) {
	assert.NotNil(t, findNode(RootPackage, "testpkg"))
}

func TestParseStruct(t *testing.T) {
	assert.NotNil(t, findNode(Struct, "Server"))
	assert.True(t, Graph.HasEdgeBetween(findNode(RootPackage, "testpkg"), findNode(Struct, "Server")))
}

func TestParseTypeAlias(t *testing.T) {
	assert.NotNil(t, findNode(Struct, "clientID"))
	assert.True(t, Graph.HasEdgeBetween(findNode(RootPackage, "testpkg"), findNode(Struct, "clientID")))	
}

func TestParseFuncDecl(t *testing.T) {
	assert.NotNil(t, findNode(Func, "NewServer"))
	assert.True(t, Graph.HasEdgeBetween(findNode(RootPackage, "testpkg"), findNode(Func, "NewServer")))
}

func TestParseFuncDeclResults(t *testing.T) {
	assert.NotNil(t, Graph.Edge(findNode(Func, "NewServer"), findNode(Struct, "Server")), "links to return result type")
}

func TestParseMethod(t *testing.T) {
	assert.NotNil(t, Graph.Edge(findNode(Struct, "Server"), findNode(Method, "Listen")), "links from type to method")
}

func TestDoesntLinkMethodAsChildOfPackage(t *testing.T) {
	assert.Nil(t, Graph.Edge(findNode(RootPackage, "testpkg"), findNode(Method, "Listen")), "should only link methods as child of their struct")
}

func TestParseFuncCall(t *testing.T) {
	assert.NotNil(t, findNode(Func, "main"))
	assert.NotNil(t, Graph.Edge(findNode(Func, "main"), findNode(Func, "NewServer")), "")
}



// func TestParseMethodCall(t *testing.T) {
// 	assert.NotNil(t, Graph.Edge(findNode(Func, "main"), findNode(Func, "NewServer")), "")
// }

func TestHandlesBuiltins(t *testing.T) {
	for _, node := range nodeLookup {
		if(node.Label() == "panic") {
			assert.Failf(t, "", "builtin functions aren't added: ", node)
		}

		if(node.Label() == "error") {
			assert.Failf(t, "", "builtin types aren't added: ", node)
		}
	}
}

func TestParseImports(t *testing.T) {
	assert.NotNil(t, findNode(ImportedPackage, "log"))
	// assert.NotNil(t, Graph.Edge(findNode(ImportedPackage, "log"), findNode(Func, "New")), "")
	// assert.NotNil(t, Graph.Edge(findNode(ImportedPackage, "log"), findNode(Func, "New")), "")
}
func TestExternalPkgFuncCall(t *testing.T) {
	// assert.Nil(t, Graph.Edge(findNode(Func, "main"), findNode(Func, "NewServer")), ""))
}


// func TestParseFuncCallToVars(t *testing.T) {
// }



func TestGraphdotOutput(t *testing.T) {
	var buffer bytes.Buffer
	Graph.ToDot(&buffer)
	assert.NotEmpty(t, buffer.String())
	fmt.Println(buffer.String())
}