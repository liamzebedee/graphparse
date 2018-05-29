package graphparse

import (
	"golang.org/x/tools/go/loader"
	"testing"
	"path/filepath"
	"os"
	"github.com/stretchr/testify/assert"
	"bytes"
	"fmt"
	"encoding/json"
)

func TestMain(m *testing.M) {
	setup()

	var buffer bytes.Buffer
	Graph.ToString(&buffer)
	fmt.Println(buffer.String())

	os.Exit(m.Run())
}

// func TestGraphdotOutput(t *testing.T) {
// 	var buffer bytes.Buffer
// 	Graph.ToDot(&buffer)
// 	assert.NotEmpty(t, buffer.String())
// 	fmt.Println(buffer.String())
// }

func TestJsonOutput(t *testing.T) {
	var buffer bytes.Buffer
	res := Graph.toJson()
	json.NewEncoder(&buffer).Encode(res)
	assert.NotEmpty(t, buffer.String())
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

	// Write test data to file for testing with JS.
	Graph.WriteJsonToFile("../test/graph.json")
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
	if(!optClusterFiles) {
		assert.True(t, Graph.HasEdgeBetween(findNode(RootPackage, "testpkg"), findNode(Struct, "Server")))
	}
}

func TestParseTypeAlias(t *testing.T) {
	assert.NotNil(t, findNode(Struct, "clientID"))
	if(!optClusterFiles) {
		assert.True(t, Graph.HasEdgeBetween(findNode(RootPackage, "testpkg"), findNode(Struct, "clientID")))
	}
}

func TestParseFuncDecl(t *testing.T) {
	assert.NotNil(t, findNode(Func, "NewServer"))
	if(!optClusterFiles) {
		assert.True(t, Graph.HasEdgeBetween(findNode(RootPackage, "testpkg"), findNode(Func, "NewServer")))
	}
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
	// assert.NotNil(t, findNode(ImportedPackage, "log"))

	// assert.NotNil(t, Graph.Edge(findNode(ImportedPackage, "log"), findNode(Func, "New")), "")
	// assert.NotNil(t, Graph.Edge(findNode(ImportedPackage, "log"), findNode(Func, "New")), "")
}

func TestAExternalPkgFuncCall(t *testing.T) {
	assert.NotNil(t, Graph.Edge(
		findNode(Func, "main"), 
		findNode(Func, "NewServer"),
	), "")
}

func TestParseValueSpec(t *testing.T) {
	assert.NotNil(t, findNode(Field, "logger"))

	assert.NotNil(t, Graph.Edge(
		findNode(File, "main.go"), 
		findNode(Field, "logger"),
	), "logger global is linked to defining file")

	assert.NotNil(t, Graph.Edge(
		findNode(Method, "Listen"), 
		findNode(Field, "logger"),
	), "usages of logger global are noted correctly")
	
	assert.Nil(t, Graph.Edge(
		findNode(File, "server.go"), 
		findNode(Field, "err"),
	), "value spec is parsed correctly with respect to parent funcs")
}

func TestParseFile(t *testing.T) {
	assert.NotNil(t, findNode(File, "main.go"))
	assert.NotNil(t, findNode(File, "server.go"))
	
	assert.NotNil(t, Graph.Edge(
		findNode(File, "server.go"), 
		findNode(Func, "NewServer"),
	), "")
}

func TestAllEdgesExist(t *testing.T) {
	for _, e := range Graph.edges {
		assert.True(t, Graph.nodeExists(e.From().ID()))
		assert.True(t, Graph.nodeExists(e.To().ID()))
	}
}

// func TestGenerates

// e.g. within Client.Close, conn.Close() is called to conn.
// func TestParseCallsToStructMembers(t *testing.T) {
// 	assert.NotNil(t, findNode(Method, "Listen"))
// 	assert.NotNil(t, findNode(Field, "listener"))

// 	assert.NotNil(t, Graph.Edge(
// 		findNode(Method, "Listen"),
// 		findNode(Field, "listener"),
// 	), "")
// 	assert.NotNil(t, Graph.Edge(
// 		findNode(Field, "listener"),
// 		findNode(Method, "Accept"),
// 	), "")
// }



