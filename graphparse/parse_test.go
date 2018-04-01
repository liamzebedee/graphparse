package graphparse

import (
	"golang.org/x/tools/go/loader"
	"testing"
	"path/filepath"
)

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



func TestCreatesPackageNode(t *testing.T) {
	setup()

	for _, node := range nodeLookup {
		if(node.Variant() == RootPackage && node.Label() == "testpkg") {
			return
		}
	}
	
	t.Fail()
}

func TestParseStruct(t *testing.T) {
	setup()

	var structNode Node

	for _, node := range nodeLookup {
		if(node.Variant() == Struct && node.Label() == "Server") {
			structNode = node
		}
	}

	if structNode == nil {
		t.Fail()
	}


	
	var constructorNode Node
	// var methodNode Node

	for _, child := range Graph.children(structNode) {
		if child.Label() == "NewServer" {
			constructorNode = child
		}
	}

	if constructorNode == nil {
		t.Fail()
	}
}