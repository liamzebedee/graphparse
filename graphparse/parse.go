package graphparse

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"go/build"

	"golang.org/x/tools/go/loader"
)

var optIncludeFilesAsNodes = false

var prog *loader.Program
var pkginfo *loader.PackageInfo 
var Graph *graph
var fset *token.FileSet
var thisPackage string

var fileLookup = make(map[string]packageFileInfo)

type packageFileInfo struct {
	// file *ast.File
	Code string `json:"code"`
	Pos token.Pos `json:"pos"`
}


func GenerateCodeGraph(pkgpath string, pkgFilePath string) {
	// We use loader instead of Importer here deliberately
	// Since we identify objects by their declaring position obj.Pos()
	// TODO alternatively we could simply use the pointer to the type
	// but this would involve parsing the underlying type (in the case of pointers)
	// Could be a potential avenue.
	conf := loader.Config{ParserMode: 0}
	conf.Import(pkgpath)

	var err error
	prog, err = conf.Load()
	if err != nil {
		panic(err)
	}

	// TODO doesn't work with relative package imports
	if pkgFilePath == "" {
		pkgFilePath = build.Default.GOPATH + "/src/" + pkgpath
	}

	pkginfo = prog.Package(pkgpath)
	if pkginfo == nil {
		panic("pkg was not loaded")
	}
	fset = prog.Fset
	Graph = NewGraph()

	thisPackage = pkginfo.Pkg.Name()
	fmt.Println("Generating graph for package", thisPackage)

	for _, f := range pkginfo.Files {
		currentFilePath := fset.File(f.Pos()).Name()
		fileName, _ := filepath.Rel(pkgFilePath, currentFilePath)

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
	// obj := pkginfo.ObjectOf(ident)

	// we want to have a canonical obj so we can make an id out of it
	// in this case, we use the def obj
	// obj := pkginfo.Defs[ident]
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


func exprToObj(expr ast.Expr) types.Object {
	var obj types.Object
	switch y := expr.(type) {
	// case *ast.StarExpr:
		// starExp := y.X
	case *ast.SelectorExpr:
		if sel := pkginfo.Selections[y]; sel != nil {
			obj = sel.Obj()
		} else {
			// Probably fully-qualified
			obj = pkginfo.ObjectOf(y.Sel)
			// obj = pkginfo.Defs[y.Sel]
		}
	case *ast.Ident:
		obj = pkginfo.ObjectOf(y)
	default:
		fmt.Println(expr)
	}
	return obj
}


func Visit(node ast.Node) bool {
	switch x := node.(type) {
	case *ast.File:
		if rootPackage == nil {
			rootPackage = LookupOrCreateCanonicalNode(x.Name.Name, RootPackage, x.Name.Name)
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
		importName, _ := strconv.Unquote(x.Path.Value)
		importedPackage := LookupOrCreateCanonicalNode(importToCanonicalKey(x), ImportedPackage, importName)
		Graph.AddEdge(importedPackage, rootPackage)
		
		// currentFileNode.extraAttrs = "[color=\"red\"]"
		return true
	
	case *ast.TypeSpec:
		obj, err := getObjFromIdent(x.Name)
		if err != nil {
			panic(err)
		}

		typeNode := LookupOrCreateNode(obj, Struct, obj.Type().String())
		Graph.AddEdge(rootPackage, typeNode)

	case *ast.FuncDecl:
		obj, err := getObjFromIdent(x.Name)

		if err != nil {
			panic(err)
		}

		funcNode := LookupOrCreateNode(obj, Func, x.Name.Name)

		if x.Recv != nil && len(x.Recv.List) > 0 {
			recvTypeObj := obj.(*types.Func).Type().(*types.Signature).Recv()

			structName := ""

			switch typ := recvTypeObj.Type().(type) {
			case *types.Pointer:
				structName = typ.Elem().(*types.Named).Obj().Name()
			default:
				panic(typ)
			}

			if err != nil {
				panic(err)
			}
			
			// varName := recvTypeObj.Name()
			recvTypeNode := LookupOrCreateNode(recvTypeObj, Struct, structName)
			Graph.AddEdge(funcNode, recvTypeNode)
		} else {
			Graph.AddEdge(rootPackage, funcNode)
		}
		
		parentFunc = funcNode
	
	case *ast.CallExpr:
		var obj types.Object
		
		obj = exprToObj(x.Fun)

		if obj == nil {
			panic(x.Fun)
		}
		
		if obj.Pkg() != nil && obj.Pkg().Name() == thisPackage {
			funcCall := LookupOrCreateNode(obj, FuncCall, obj.Name())
			Graph.AddEdge(parentFunc, funcCall)
		}

		
		// switch y := x.Fun.(type) {
		// case *ast.SelectorExpr:
		// 	obj := pkginfo.ObjectOf(y.Sel)
		// 	pkg := obj.Pkg()

		// 	if pkg != nil {
		// 		if pkg.Name() == "subnet" {
		// 			id, err := getIdOfIdent(y.Sel)
		// 			if err != nil {
		// 				fmt.Fprintln(os.Stderr, err.Error())
		// 				return true
		// 			}
		
		// 			callNode := NewNode(y, id, y.Sel.Name, ShouldAlreadyExist)
		// 			Graph.AddEdge(callNode, parentNode)
		// 			// Graph.AddEdge(parentNode, callNode)
					
		// 		} else {
		// 			importName := pkg.Path()
		// 			pkgId := GetCanonicalNodeId(importName)
		// 			importPkgNode := NewNode(nil, pkgId, pkg.Name(), ImportedPackage)

		// 			id, err := getIdOfIdent(y.Sel)
		// 			if err != nil {
		// 				fmt.Fprintln(os.Stderr, err.Error())
		// 				return true
		// 			}
		
		// 			callNode := NewNode(y, id, y.Sel.Name, ImportedFunc)

		// 			Graph.AddEdge(importPkgNode, callNode)
		// 			Graph.AddEdge(parentNode, callNode)
		// 			// package -> callnode <- parent
		// 		}
		// 	}
			
		
		// 	case *ast.Ident:
		// 		id, err := getIdOfIdent(y)
		// 		if err != nil {
		// 			fmt.Fprintln(os.Stderr, err.Error())
		// 			return true
		// 		}
				
		// 		callNode := NewNode(y, id, y.Name, ShouldAlreadyExist)
		// 		Graph.AddEdge(callNode, parentNode)
		// 	default:
		// 		fmt.Fprintln(os.Stderr, "parsing call - missed type", y)
		// 	}
	
	default:
		// fmt.Fprintf(os.Stderr, "parsing - missed type %T\n", x)
		return true
	}
	return true
}
