package graphparse

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"golang.org/x/tools/go/loader"
)

var packageFilePath = "/Users/liamz/go/src/github.com/twitchyliquid64/subnet/"
var optIncludeFilesAsNodes = false

var prog *loader.Program
var pkginfo *loader.PackageInfo 
var Graph *graph
var fset *token.FileSet

var fileLookup = make(map[string]packageFileInfo)

type packageFileInfo struct {
	// file *ast.File
	Code string `json:"code"`
	Pos token.Pos `json:"pos"`
}


func GenerateCodeGraph() {
	var err error

	pkgpath := "github.com/twitchyliquid64/subnet/subnet"
	// pkgpath := "github.com/liamzebedee/graphparse/graphparse"

	conf := loader.Config{ParserMode: 0}
	conf.Import(pkgpath)
	prog, err = conf.Load()
	if err != nil {
		panic(err)
	}

	pkginfo = prog.Package(pkgpath)
	fset = prog.Fset
	Graph = NewGraph()

	for _, f := range pkginfo.Files {
		currentFilePath := fset.File(f.Package).Name()
		fileName, _ := filepath.Rel(packageFilePath, currentFilePath)

		code, err := ioutil.ReadFile(currentFilePath)
		if err != nil {
			panic(err)
		}

		fileLookup[fileName] = packageFileInfo{
			Code: string(code),
			Pos: f.Pos(),
		}

		fmt.Println("Processing", fileName)
		ast.Inspect(f, Visit)
		// ast.Print(fset, f)
	}
}




func getObjFromIdent(ident *ast.Ident) (types.Object, error) {
	// I wrote this in a subconcious spree of, "I have a gut feeling that this will do it"
	obj := pkginfo.ObjectOf(ident)
	
	if obj == nil {
		return nil, fmt.Errorf("obj not found for ident", ident)
	}

	// Universe scope
	if obj.Pkg() == nil {
		return nil, fmt.Errorf("universe object ", obj.Type(), ident)
	}

	if obj != nil {
		return obj, nil
	}

	// if obj != nil {
	// 	if id, ok := objsIdentified[objId]; ok {
	// 		return id, nil
	// 	} else {
	// 		id := pointerToId(obj)
	// 		objsIdentified[objId] = id
	// 		return id, nil
	// 	}
	// }

	return nil, fmt.Errorf("unexpected error", obj)
}


var rootPackage Node
var currentFile Node
var parentFunc Node



func importToCanonicalKey(importSpec *ast.ImportSpec) string {
	importName, err := strconv.Unquote(importSpec.Path.Value)
	if err == nil {
		return importName
	} else {
		panic(err)
	}
}


func Visit(node ast.Node) bool {
	switch x := node.(type) {
	case *ast.File:
		if rootPackage == nil {
			pkgIdent := x.Name.Name
			rootPackage = LookupOrCreateCanonicalNode(pkgIdent, RootPackage)
			// LABEL x.Name.Name
		}

		// TODO REFACTOR
		// if optIncludeFilesAsNodes {
		// 	fileName, err := filepath.Rel(packageFilePath, fset.File(x.Package).Name())
		// 	if err != nil {
		// 		fmt.Fprintln(os.Stderr, err.Error())
		// 		return false
		// 	}

		// 	fileNodeId := GetCanonicalNodeId(fileName)
		// 	currentFileNode = NewNode(x, fileNodeId, fileName, File)
		// 	// currentFileNode.extraAttrs = "[color=\"red\"]"
		// 	Graph.AddEdge(pkgIdentNode, currentFileNode)
		// }

	case *ast.ImportSpec:
		// LABEL importName, _ := strconv.Unquote(x.Path.Value)
		importedPackage := LookupOrCreateCanonicalNode(importToCanonicalKey(x), ImportedPackage)
		Graph.AddEdge(importedPackage, rootPackage)
		
		// currentFileNode.extraAttrs = "[color=\"red\"]"
		return true
	
	default:
		// fmt.Fprintf(os.Stderr, "parsing - missed type %T\n", x)
		return true
	}
	return true
}
