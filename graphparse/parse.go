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
	// "os"

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

func exprToObj(expr ast.Expr) (types.Object, error) {
	switch x := expr.(type) {
	case *ast.SelectorExpr:
		if sel := pkginfo.Selections[x]; sel != nil {
			return sel.Obj(), nil
		}
		// Probably fully-qualified
		return pkginfo.ObjectOf(x.Sel), nil
	
	case *ast.Ident:
		obj, err := getObjFromIdent(x)
		return obj, err
	
	default:
		ParserLog.Printf("missed type %T\n", x)
	}

	return nil, fmt.Errorf("couldn't get object for expression:", expr)
}

func objIsWorthy(obj types.Object) bool {
	if obj.Pkg() == nil {
		return false
	}
	return true
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

type parseEngine struct {
}


func (eng *parseEngine) parseRootPackage(f *ast.File) {
	if rootPackage == nil {
		rootPackage = LookupOrCreateCanonicalNode(f.Name.Name, RootPackage, f.Name.Name)
	}
}

func (eng *parseEngine) parseFile(f *ast.File) {
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
}

func (eng *parseEngine) parseTypeSpec(typeSpec *ast.TypeSpec) {
	switch x := typeSpec.Type.(type) {
	case *ast.StructType:
		obj, err := getObjFromIdent(typeSpec.Name)
		if err != nil {
			panic(err)
		}

		typeNode := LookupOrCreateNode(obj, Struct, obj.Name())
		Graph.AddEdge(rootPackage, typeNode)
	
	case *ast.Ident:
		obj, err := getObjFromIdent(typeSpec.Name)
		if err != nil {
			panic(err)
		}

		typeNode := LookupOrCreateNode(obj, Struct, obj.Name())
		Graph.AddEdge(rootPackage, typeNode)
	
	default:
		ParserLog.Printf("missed type %T\n", x)
	}
}

func (eng *parseEngine) parseFuncDecl(funcDecl *ast.FuncDecl) {
	obj, err := getObjFromIdent(funcDecl.Name)
	if err != nil {
		panic(err)
	}

	if !objIsWorthy(obj) {
		return
	}

	// TODO test code below for cases where it isn't a *types.Signature
	// then refactor using isMethod
	isMethod := funcDecl.Recv.NumFields() == 1
	variant := Func
	if isMethod {
		variant = Method
	}

	funcNode := LookupOrCreateNode(obj, variant, funcDecl.Name.Name)
	parentFunc = funcNode

	switch x := obj.Type().(type) {
	case *types.Signature:
		if isMethod {
			eng.parseMethod(funcNode, x.Recv())
		} else {
			Graph.AddEdge(rootPackage, funcNode)
		}

		eng.parseFuncDeclResults(funcNode, x.Results())

	default:
		ParserLog.Printf("missed type %T\n", x)
	}
}

func (eng *parseEngine) parseFuncDeclResults(funcNode Node, results *types.Tuple) {
	for i := 0; i < results.Len(); i++ {
		result := results.At(i)
		
		// parse result type
		obj := typeToObj(result.Type())
		if obj == nil {
			continue
		}
		if !objIsWorthy(obj) {
			continue
		}

		resultTypeNode := LookupOrCreateNode(obj, Func, obj.Name())
		Graph.AddEdge(funcNode, resultTypeNode)
	}
}

func (eng *parseEngine) parseMethod(funcNode Node, recv *types.Var) {
	obj := typeToObj(recv.Type())
	recvTypeNode := LookupOrCreateNode(obj, Struct, obj.Name())
	Graph.AddEdge(recvTypeNode, funcNode)
}

func (eng *parseEngine) parseCallExpr(callExp *ast.CallExpr) {
	obj, err := exprToObj(callExp.Fun)
	if err != nil || obj == nil {
		return
	}

	if !objIsWorthy(obj) {
		return
	}

	if obj.Pkg().Name() != thisPackage {
		eng.parseExternalFuncCall(callExp, obj)
		return
	}


	variant := Func
	switch x := obj.Type().(type) {
	case *types.Signature:
		if x.Recv() != nil {
			variant = Method
		}
	default:
		ParserLog.Printf("missed type %T\n", x)
	}
	
	funcNode := LookupOrCreateNode(obj, variant, obj.Name())
	
	// TEST
	if parentFunc != nil {
		Graph.AddEdge(parentFunc, funcNode)
	}
}

func (eng *parseEngine) parseExternalFuncCall(callExp *ast.CallExpr, obj types.Object) {
	importPath := obj.Pkg().Path()
	importNode := LookupOrCreateCanonicalNode(importPath, ImportedPackage, importPath)

	variant := Func
	switch x := obj.Type().(type) {
	case *types.Signature:
		if x.Recv() != nil {
			variant = Method
		}
	default:
		ParserLog.Printf("missed type %T\n", x)
	}
	funcNode := LookupOrCreateNode(obj, variant, obj.Name())

	Graph.AddEdge(importNode, funcNode)
	
	// TEST
	if parentFunc != nil {
		Graph.AddEdge(parentFunc, funcNode)
	}
	
	// package -> callnode <- parent
}


func (eng *parseEngine) parseImportSpec(importSpec *ast.ImportSpec) {
	importPath, err := strconv.Unquote(importSpec.Path.Value)
	if err != nil {
		ParserLog.Fatal(err)
	}
	
	LookupOrCreateCanonicalNode(importPath, ImportedPackage, importPath)	
	// importedPackage := LookupOrCreateCanonicalNode(importPath, ImportedPackage, importPath)
	// Graph.AddEdge(importedPackage, rootPackage)
}


var eng = parseEngine{}


func Visit(node ast.Node) bool {
	switch x := node.(type) {
	case *ast.ImportSpec:
		eng.parseImportSpec(x)
	case *ast.File:
		eng.parseRootPackage(x)
		eng.parseFile(x)
	case *ast.TypeSpec:
		eng.parseTypeSpec(x)
	case *ast.FuncDecl:
		eng.parseFuncDecl(x)
	case *ast.CallExpr:
		eng.parseCallExpr(x)
	case *ast.Ident, nil:
		break
	default:
		ParserLog.Printf("missed type %T\n", x)
		return true
	}
	return true
}