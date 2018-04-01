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
	"go/parser"

	"golang.org/x/tools/go/loader"
)

// GO GURU
// https://golang.org/lib/godoc/analysis/help.html
// https://github.com/golang/tools/blob/master/go/callgraph/cha/cha.go#L38:46
// https://github.com/dominikh/implements/blob/master/main.go#L103:16
// https://github.com/golang/tools/blob/master/cmd/guru/implements.go#L47:16

var optIncludeFilesAsNodes = false

var prog *loader.Program
var pkginfo *loader.PackageInfo 
var Graph *graph
var fset *token.FileSet
var thisPackage string

var fileLookup = make(map[string]packageFileInfo)

func GenerateCodeGraphFromProg(prog *loader.Program, pkgpath, pkgFilePath string) {
	// TODO doesn't work with relative package imports
	if pkgFilePath == "" {
		pkgFilePath = build.Default.GOPATH + "/src/" + pkgpath
	}

	fmt.Println(pkgFilePath)

	pkginfo = prog.Package(pkgpath)
	if pkginfo == nil {
		panic("pkg was not loaded")
	}
	fset = prog.Fset
	Graph = NewGraph()

	thisPackage = pkginfo.Pkg.Name()
	fmt.Println("Generating graph for package", thisPackage)

	// TODO process subpackages
	// eg github.com/twitchyliquid64/subnet/subnet/conn from subnet
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

	Graph.markGenerated()
}

func GenerateCodeGraph(pkgpath string) {
	// We use loader instead of Importer here deliberately
	// Since we identify objects by their declaring position obj.Pos()
	// TODO alternatively we could simply use the pointer to the type
	// but this would involve parsing the underlying type (in the case of pointers)
	// Could be a potential avenue.
	conf := loader.Config{ParserMode: parser.ParseComments}
	conf.Import(pkgpath)

	var err error
	prog, err = conf.Load()
	if err != nil {
		panic(err)
	}

	GenerateCodeGraphFromProg(prog, pkgpath, "")
}




func getObjFromIdent(ident *ast.Ident) (types.Object, error) {
	// I wrote this in a subconcious spree of, "I have a gut feeling that this will do it"
	obj := pkginfo.ObjectOf(ident)

	// we want to have a canonical obj so we can make an id out of it
	// in this case, we use the def obj
	// obj := pkginfo.Defs[ident]
	
	// if obj == nil {
	// 	return nil, fmt.Errorf("obj not found for ident", ident)
	// }

	// Universe scope
	// if obj.Pkg() == nil {
	// 	return nil, fmt.Errorf("universe object ", obj.Type(), ident)
	// }

	if obj != nil {
		return obj, nil
	}

	return nil, fmt.Errorf("unexpected error", obj)
}

func exprToObj(expr ast.Expr) (error, types.Object) {
	switch y := expr.(type) {
	case *ast.SelectorExpr:
		if sel := pkginfo.Selections[y]; sel != nil {
			return nil, sel.Obj()
		}

		// Probably fully-qualified
		return nil, pkginfo.ObjectOf(y.Sel)
	case *ast.Ident:
		obj, err := getObjFromIdent(y)
		return err, obj
	}

	return fmt.Errorf("couldnt get obj for expr", expr), nil
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
		return true
	
	case *ast.TypeSpec:
		obj, err := getObjFromIdent(x.Name)
		if err != nil {
			panic(err)
		}

		if !typeIsNamed(obj.Type()) {
			return true
		}

		typeNode := LookupOrCreateNode(obj, Struct, obj.Name())
		Graph.AddEdge(rootPackage, typeNode)

	case *ast.FuncDecl:
		obj, err := getObjFromIdent(x.Name)
		if err != nil {
			panic(err)
		}

		if !typeIsNamed(obj.Type()) {
			return true
		}

		funcNode := LookupOrCreateNode(obj, Func, x.Name.Name)
		parentFunc = funcNode

		// 1. Link receiver
		if x.Recv != nil && len(x.Recv.List) > 0 {
			recvTypeObj := obj.(*types.Func).Type().(*types.Signature).Recv()

			structName := ""

			switch typ := recvTypeObj.Type().(type) {
			case *types.Pointer:
				structName = typ.Elem().(*types.Named).Obj().Name()
			case *types.Named:
				structName = typ.Obj().Name()
			default:
				fmt.Printf("%T\n", typ)
				panic(typ)
			}

			if err != nil {
				panic(err)
			}
			
			// varName := recvTypeObj.Name()

			// Type of the receiver
			// ie the type that this method operates on
			typeNode := LookupOrCreateNode(recvTypeObj, Struct, structName)
			// Graph.AddEdge(funcNode, typeNode)
			Graph.AddEdge(typeNode, funcNode)
		} else {
			Graph.AddEdge(rootPackage, funcNode)
		}

		// Link params
		if params := x.Type.Params.List; len(params) > 0 {
			for _, y := range params {
				obj, err := getObjFromIdent(y.Names[len(y.Names) - 1])
				if err != nil {
					panic(err)
				}

				// Ignore unnamed params
				if !typeIsNamed(obj.Type()) {
					continue
				}

				paramTypeNode := LookupOrCreateNode(obj, Struct, obj.Name())
				Graph.AddEdge(paramTypeNode, funcNode)
			}
		}
	
	case *ast.CallExpr:		
		typ := pkginfo.TypeOf(x)
		// if !typeIsNamed(typ) {
		// 	fmt.Println("ignoring call for unnamed type", typ)
		// }
		obj := typeToObj(typ)
		// fmt.Println(obj)
		// return false

		// err, obj := exprToObj(x.Fun)

		// if err != nil {
		// 	fmt.Println("error parsing callexpr", x, err)
		// 	return true
		// }

		// if objIsBuiltin(obj) {
		// 	return true
		// }

		return true

		if !typeIsNamed(obj.Type()) {
			return true
		}

		
		if obj.Pkg() != nil && obj.Pkg().Name() == thisPackage {
			funcCall := LookupOrCreateNode(obj, FuncCall, obj.Name())
			Graph.AddEdge(parentFunc, funcCall)
		}
	
	default:
		// fmt.Fprintf(os.Stderr, "parsing - missed type %T\n", x)
		return true
	}
	return true
}

// func objIsBuiltin(obj types.Object) bool {
// 	fmt.Println(obj.Pkg())
// 	return true
// }