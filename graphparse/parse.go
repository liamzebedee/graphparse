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


var prog *loader.Program
var pkginfo *loader.PackageInfo 

var fset *token.FileSet
var thisPackage string

var fileLookup = make(map[string]packageFileInfo)

var eng = parseEngine{}

var Graph *graph

type packageFileInfo struct {
	Code string `json:"code"`
	Pos token.Pos `json:"pos"`
}


func GenerateCodeGraphFromProg(prog *loader.Program, pkgpath, pkgFilePath string) {
	Graph = NewGraph()

	// TODO doesn't work with relative package imports
	if pkgFilePath == "" {
		pkgFilePath = build.Default.GOPATH + "/src/" + pkgpath
	}

	fmt.Println(pkgFilePath)
	eng.pkgPath = pkgFilePath

	pkginfo = prog.Package(pkgpath)
	if pkginfo == nil {
		panic("pkg was not loaded")
	}
	fset = prog.Fset

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

	Graph.processUnknownReferences()

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

type parseEngine struct {
	pkgPath string
	rootPackage Node
	currentFile Node
	parentFunc *objNode
}

var optClusterFiles = true



func (eng *parseEngine) parseRootPackage(f *ast.File) {
	if eng.rootPackage == nil {
		eng.rootPackage = LookupOrCreateCanonicalNode(f.Name.Name, RootPackage, f.Name.Name)
	}
}

func (eng *parseEngine) parseFile(f *ast.File) {
	currentFilePath := fset.File(f.Pos()).Name()
	fileName, _ := filepath.Rel(eng.pkgPath, currentFilePath)

	eng.currentFile = LookupOrCreateCanonicalNode(currentFilePath, File, fileName)

	Graph.AddEdge(eng.rootPackage, eng.currentFile, Def)
}

func (eng *parseEngine) parseTypeSpec(typeSpec *ast.TypeSpec) {
	fromNode := eng.rootPackage
	if optClusterFiles {
		fromNode = eng.currentFile
	}
	
	switch x := typeSpec.Type.(type) {
	case *ast.StructType:
		obj, err := getObjFromIdent(typeSpec.Name)
		if err != nil {
			panic(err)
		}

		typeNode := CreateNode(obj, Struct, obj.Name())
		Graph.AddEdge(fromNode, typeNode, Def)
	
	case *ast.Ident:
		obj, err := getObjFromIdent(typeSpec.Name)
		if err != nil {
			panic(err)
		}

		typeNode := CreateNode(obj, Struct, obj.Name())
		Graph.AddEdge(fromNode, typeNode, Def)
	
	default:
		ParserLog.Printf("missed type %T\n", x)
	}
}

func (eng *parseEngine) parseFuncDecl(funcDecl *ast.FuncDecl) {
	fromNode := eng.rootPackage
	if optClusterFiles {
		fromNode = eng.currentFile
	}


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

	funcNode := CreateNode(obj, variant, funcDecl.Name.Name)
	eng.parentFunc = funcNode

	switch x := obj.Type().(type) {
	case *types.Signature:
		if isMethod {
			eng.parseMethod(funcNode, x.Recv())
		} else {
			Graph.AddEdge(fromNode, funcNode, Def)
		}

		eng.parseFuncDeclResults(funcNode, x.Results())

	default:
		ParserLog.Printf("missed type %T\n", x)
	}

	eng.parseFuncBody(funcDecl.Body)
	eng.parentFunc = nil
}

func (eng *parseEngine) parseFuncBody(body *ast.BlockStmt) {
	ast.Inspect(body, VisitBody)
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

		// resultTypeNode := LookupOrCreateNode(obj, Func, obj.Name())
		resultTypeNode := LookupNode(obj)	
		Graph.AddEdge(funcNode, resultTypeNode, Use)
	}
}

func (eng *parseEngine) parseMethod(funcNode Node, recv *types.Var) {
	obj := typeToObj(recv.Type())
	// recvTypeNode := LookupOrCreateNode(obj, Struct, obj.Name())
	recvTypeNode := LookupNode(obj)
	Graph.AddEdge(recvTypeNode, funcNode, Def)
}

func (eng *parseEngine) parseCallExpr(callExp *ast.CallExpr) {
	obj, err := exprToObj(callExp.Fun)

	switch x := callExp.Fun.(type) {
	case *ast.SelectorExpr:
		sel := pkginfo.Selections[x]
		if sel == nil {
			ParserLog.Printf("skipping selector expr (likely qualified identifier) %T\n", x)
			break
		}
		// objs := getObjectsFromSelector(x)
		eng.parseSelectorCallExpr2(sel, x)
		return
		
	case *ast.Ident:
		// continue as usual.

	default:
		ParserLog.Printf("missed type of call expr %T\n", x)
	}

	if err != nil || obj == nil {
		ParserLog.Println("ignoring call expr for obj - ", obj, ", error: ", err)
		return
	}

	if !objIsWorthy(obj) {
		return
	}

	if obj.Pkg().Name() != thisPackage {
		// eng.parseExternalFuncCall(callExp, obj)
		return
	}

	funcNode := funcObjToNode(obj)
	
	// TEST
	if eng.parentFunc != nil {
		Graph.AddEdge(eng.parentFunc, funcNode, Use)
	}
}

func (eng *parseEngine) parseSelectorCallExpr2(sel *types.Selection, selAst *ast.SelectorExpr) {
	// Selection X.Y
	
	// Parse the X
	// recv := sel.Recv()
	// recvObj := typeToObj(recv)
	recvObj, err := exprToObj(selAst.X)
	if err == nil {
		if objIsWorthy(recvObj) {
			recvNode := LookupNode(recvObj)
			Graph.AddEdge(eng.parentFunc, recvNode, Use)
		}
	} else {
		ParserLog.Println("")
	}

	

	// Parse the Y
	selObj := sel.Obj()
	if objIsWorthy(selObj) {
		objNode := LookupNode(selObj)
		Graph.AddEdge(eng.parentFunc, objNode, Use)
	}
}

func (eng *parseEngine) parseSelectorCallExpr(objs []annotatedSelectorObject) {
	prev := eng.parentFunc

	for i := len(objs) - 1; i != -1; i-- {
		obj := objs[i]

		// nodeTyp := FuncCall
		// switch obj.kind {
		// case types.FieldVal:
		// 	nodeTyp = Field
		// case types.MethodExpr, types.MethodVal:
		// 	nodeTyp = Method
		// }

		if eng.parentFunc.obj != obj.Object {
			// currNode := LookupOrCreateNode(obj, nodeTyp, obj.Name())
			currNode := LookupNode(obj)
			Graph.AddEdge(prev, currNode, Use)
			prev = currNode
		}
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
	funcNode := CreateNode(obj, variant, obj.Name())
	// funcNode := LookupOrCreateNode(obj, variant, obj.Name())

	Graph.AddEdge(importNode, funcNode, Use)	
	
	// TEST
	if eng.parentFunc != nil {
		Graph.AddEdge(eng.parentFunc, funcNode, Use)
		Graph.AddEdge(funcNode, eng.parentFunc, Use)
	}
	
	// package -> callnode <- parent
}


func (eng *parseEngine) parseImportSpec(importSpec *ast.ImportSpec) {
	importPath, err := strconv.Unquote(importSpec.Path.Value)
	if err != nil {
		ParserLog.Fatal(err)
	}
	
	LookupOrCreateCanonicalNode(importPath, ImportedPackage, importPath)	
	importedPackage := LookupOrCreateCanonicalNode(importPath, ImportedPackage, importPath)
	Graph.AddEdge(importedPackage, eng.rootPackage, Def)
}

// Parses constant and variable declarations
func (eng *parseEngine) parseValueSpec(x *ast.ValueSpec) {
	// trying first with one value
	ident := x.Names[0]
	if len(x.Names) > 1 {
		ParserLog.Println("value spec has more than one ident %s", x.Names)
	}

	obj, err := getObjFromIdent(ident)
	if err != nil {
		ParserLog.Fatal("couldn't parse valuespec obj:", err)
	}

	// node := LookupOrCreateNode(obj, Field, obj.Name())
	node := CreateNode(obj, Field, obj.Name())
	if eng.parentFunc != nil {
		Graph.AddEdge(eng.parentFunc, node, Def)
	} else {
		Graph.AddEdge(eng.currentFile, node, Def)
	}
}


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
		return false
	case *ast.ValueSpec:
		eng.parseValueSpec(x)
	case *ast.Ident, nil:
		break
	default:
		ParserLog.Printf("missed type %T\n", x)
		return true
	}
	return true
}


func VisitBody(node ast.Node) bool {
	switch x := node.(type) {
	case *ast.CallExpr:
		// ast.Print(fset, x)
		eng.parseCallExpr(x)
	case *ast.ValueSpec:
		eng.parseValueSpec(x)
	case *ast.Ident, nil:
		break
	default:
		ParserLog.Printf("missed type %T\n", x)
		return true
	}
	return true
}