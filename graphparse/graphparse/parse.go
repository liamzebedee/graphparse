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


type packageFileInfo struct {
	Code string `json:"code"`
	Pos token.Pos `json:"pos"`
}


func GenerateCodeGraphFromProg(prog *loader.Program, pkgpath, pkgFilePath string) (*Graph, error) {
	eng := parseEngine{
		fileLookup: make(map[string]packageFileInfo),
	}
	eng.Graph = NewGraph()

	// TODO doesn't work with relative package imports
	if pkgFilePath == "" {
		pkgFilePath = build.Default.GOPATH + "/src/" + pkgpath
	}

	fmt.Println(pkgFilePath)
	eng.pkgPath = pkgFilePath

	eng.pkginfo = prog.Package(pkgpath)
	if eng.pkginfo == nil {
		for _, err := range eng.pkginfo.Errors {
			fmt.Printf(err.Error())
		}
		return nil, fmt.Errorf("pkg was not loaded")
	}
	eng.fset = prog.Fset

	eng.thisPackage = eng.pkginfo.Pkg.Name()
	fmt.Println("Generating graph for package", eng.thisPackage)

	// TODO process subpackages
	// eg github.com/twitchyliquid64/subnet/subnet/conn from subnet
	for _, f := range eng.pkginfo.Files {
		currentFilePath := eng.fset.File(f.Pos()).Name()
		fileName, _ := filepath.Rel(pkgFilePath, currentFilePath)

		code, err := ioutil.ReadFile(currentFilePath)
		if err != nil {
			return nil, fmt.Errorf("couldn't read source file:", err)
		}

		eng.fileLookup[fileName] = packageFileInfo{
			Code: string(code),
			Pos: f.Pos(),
		}

		fmt.Println("Processing", fileName)
		ast.Inspect(f, eng.Visit)
	}

	eng.processUnknownReferences()

	eng.Graph.markGenerated()

	return eng.Graph, nil
}


func GenerateCodeGraph(pkgpath string) (*Graph, error) {
	conf := loader.Config{ParserMode: parser.ParseComments}
	// when a package is already included in vendor/
	// then the call in GenerateCodeGraphFromProg to prog.Package(pkgpath) will return nil
	// FIX THIS BUG BUG BUG
	conf.Import(pkgpath)

	var err error
	prog, err := conf.Load()
	if err != nil {
		return nil, fmt.Errorf("couldn't load package:", err)
	}

	return GenerateCodeGraphFromProg(prog, pkgpath, "")
}

type parseEngine struct {
	pkgPath string
	rootPackage Node
	currentFile Node
	parentFunc *objNode
	Graph *Graph

	prog *loader.Program
	pkginfo *loader.PackageInfo 

	fset *token.FileSet
	thisPackage string

	fileLookup map[string]packageFileInfo
}

var optClusterFiles = true

func (eng *parseEngine) processUnknownReferences() {
	// var missing []edge
	// var edges []edge

	var i int
	for _, n := range eng.Graph.objNodeLookup {
		fmt.Println(n.obj, eng.DebugInfo(n))
		i++
	}

	fmt.Println(i, "missing nodes")
	// for _, e := range g.edges {
	// 	if _, ok := nodeLookup[e.from.Id()]; !ok {
	// 		missing = append(missing, e)
	// 		continue
	// 	}
	// 	if _, ok := nodeLookup[e.to.Id()]; !ok {
	// 		missing = append(missing, e)
	// 		continue
	// 	}
	// 	edges = append(edges, e)
	// }

	// for _, x := range missing {
	// 	fmt.Println(x)
	// }

	// g.edges = edges
}

func (eng *parseEngine) DebugInfo(n Node) string {
	return ""
}

func (eng *parseEngine) parseRootPackage(f *ast.File) {
	if eng.rootPackage == nil {
		eng.rootPackage = eng.Graph.LookupOrCreateCanonicalNode(f.Name.Name, RootPackage, f.Name.Name)
	}
}

func (eng *parseEngine) parseFile(f *ast.File) {
	currentFilePath := eng.fset.File(f.Pos()).Name()
	fileName, _ := filepath.Rel(eng.pkgPath, currentFilePath)

	eng.currentFile = eng.Graph.LookupOrCreateCanonicalNode(currentFilePath, File, fileName)

	eng.Graph.AddEdge(eng.rootPackage, eng.currentFile, Def)
}

func (eng *parseEngine) parseTypeSpec(typeSpec *ast.TypeSpec) {
	fromNode := eng.rootPackage
	if optClusterFiles {
		fromNode = eng.currentFile
	}
	
	switch x := typeSpec.Type.(type) {
	case *ast.StructType:
		obj, err := eng.getObjFromIdent(typeSpec.Name)
		if err != nil {
			panic(err)
		}

		typeNode := eng.Graph.CreateNode(obj, Struct, obj.Name(), eng.getPos(obj.Pos()))
		eng.Graph.AddEdge(fromNode, typeNode, Def)
	
	case *ast.Ident:
		obj, err := eng.getObjFromIdent(typeSpec.Name)
		if err != nil {
			panic(err)
		}

		typeNode := eng.Graph.CreateNode(obj, Struct, obj.Name(), eng.getPos(obj.Pos()))
		eng.Graph.AddEdge(fromNode, typeNode, Def)
	
	default:
		ParserLog.Printf("missed type %T\n", x)
	}
}

func (eng *parseEngine) parseFuncDecl(funcDecl *ast.FuncDecl) {
	fromNode := eng.rootPackage
	if optClusterFiles {
		fromNode = eng.currentFile
	}


	obj, err := eng.getObjFromIdent(funcDecl.Name)
	if err != nil {
		panic(err)
	}

	if !eng.objIsWorthy(obj) {
		return
	}

	// TODO test code below for cases where it isn't a *types.Signature
	// then refactor using isMethod
	isMethod := funcDecl.Recv.NumFields() == 1
	variant := Func
	if isMethod {
		variant = Method
	}

	funcNode := eng.Graph.CreateNode(obj, variant, funcDecl.Name.Name, eng.getPos(obj.Pos()))
	eng.parentFunc = funcNode

	switch x := obj.Type().(type) {
	case *types.Signature:
		if isMethod {
			eng.parseMethod(funcNode, x.Recv())
		} else {
			eng.Graph.AddEdge(fromNode, funcNode, Def)
		}

		eng.parseFuncDeclResults(funcNode, x.Results())

	default:
		ParserLog.Printf("missed type %T\n", x)
	}

	eng.parseFuncBody(funcDecl.Body)
	eng.parentFunc = nil
}

func (eng *parseEngine) parseFuncBody(body *ast.BlockStmt) {
	ast.Inspect(body, eng.VisitBody)
}



func (eng *parseEngine) parseFuncDeclResults(funcNode Node, results *types.Tuple) {
	for i := 0; i < results.Len(); i++ {
		result := results.At(i)
		
		// parse result type
		obj := typeToObj(result.Type())
		if obj == nil {
			continue
		}
		if !eng.objIsWorthy(obj) {
			continue
		}

		// resultTypeNode := LookupOrCreateNode(obj, Func, obj.Name())
		resultTypeNode := eng.Graph.LookupNode(obj)	
		eng.Graph.AddEdge(funcNode, resultTypeNode, Use)
	}
}

func (eng *parseEngine) parseMethod(funcNode Node, recv *types.Var) {
	obj := typeToObj(recv.Type())
	// recvTypeNode := LookupOrCreateNode(obj, Struct, obj.Name())
	recvTypeNode := eng.Graph.LookupNode(obj)
	eng.Graph.AddEdge(recvTypeNode, funcNode, Def)
}

func (eng *parseEngine) parseCallExpr(callExp *ast.CallExpr) {
	obj, err := eng.exprToObj(callExp.Fun)

	switch x := callExp.Fun.(type) {
	case *ast.SelectorExpr:
		sel := eng.pkginfo.Selections[x]
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

	if !eng.objIsWorthy(obj) {
		return
	}

	if obj.Pkg().Name() != eng.thisPackage {
		// eng.parseExternalFuncCall(callExp, obj)
		return
	}

	funcNode := eng.funcObjToNode(obj)
	
	// TEST
	if eng.parentFunc != nil {
		eng.Graph.AddEdge(eng.parentFunc, funcNode, Use)
	}
}

func (eng *parseEngine) parseSelectorCallExpr2(sel *types.Selection, selAst *ast.SelectorExpr) {
	// Selection X.Y
	
	// Parse the X
	// recv := sel.Recv()
	// recvObj := typeToObj(recv)
	recvObj, err := eng.exprToObj(selAst.X)
	if err == nil {
		if eng.objIsWorthy(recvObj) {
			recvNode := eng.Graph.LookupNode(recvObj)
			eng.Graph.AddEdge(eng.parentFunc, recvNode, Use)
		}
	} else {
		ParserLog.Println("")
	}

	

	// Parse the Y
	selObj := sel.Obj()
	if eng.objIsWorthy(selObj) {
		objNode := eng.Graph.LookupNode(selObj)
		eng.Graph.AddEdge(eng.parentFunc, objNode, Use)
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
			currNode := eng.Graph.LookupNode(obj)
			eng.Graph.AddEdge(prev, currNode, Use)
			prev = currNode
		}
	}
}

func (eng *parseEngine) parseExternalFuncCall(callExp *ast.CallExpr, obj types.Object) {
	importPath := obj.Pkg().Path()
	importNode := eng.Graph.LookupOrCreateCanonicalNode(importPath, ImportedPackage, importPath)

	variant := Func
	switch x := obj.Type().(type) {
	case *types.Signature:
		if x.Recv() != nil {
			variant = Method
		}
	default:
		ParserLog.Printf("missed type %T\n", x)
	}
	funcNode := eng.Graph.CreateNode(obj, variant, obj.Name(), eng.getPos(obj.Pos()))
	// funcNode := LookupOrCreateNode(obj, variant, obj.Name())

	eng.Graph.AddEdge(importNode, funcNode, Use)	
	
	// TEST
	if eng.parentFunc != nil {
		eng.Graph.AddEdge(eng.parentFunc, funcNode, Use)
		eng.Graph.AddEdge(funcNode, eng.parentFunc, Use)
	}
	
	// package -> callnode <- parent
}


func (eng *parseEngine) parseImportSpec(importSpec *ast.ImportSpec) {
	importPath, err := strconv.Unquote(importSpec.Path.Value)
	if err != nil {
		ParserLog.Fatal(err)
	}
	
	eng.Graph.LookupOrCreateCanonicalNode(importPath, ImportedPackage, importPath)	
	importedPackage := eng.Graph.LookupOrCreateCanonicalNode(importPath, ImportedPackage, importPath)
	eng.Graph.AddEdge(importedPackage, eng.rootPackage, Def)
}

// Parses constant and variable declarations
func (eng *parseEngine) parseValueSpec(x *ast.ValueSpec) {
	// trying first with one value
	ident := x.Names[0]
	if len(x.Names) > 1 {
		ParserLog.Println("value spec has more than one ident %s", x.Names)
	}

	obj, err := eng.getObjFromIdent(ident)
	if err != nil {
		ParserLog.Fatal("couldn't parse valuespec obj:", err)
	}

	// node := LookupOrCreateNode(obj, Field, obj.Name())
	node := eng.Graph.CreateNode(obj, Field, obj.Name(), eng.getPos(obj.Pos()))
	if eng.parentFunc != nil {
		eng.Graph.AddEdge(eng.parentFunc, node, Def)
	} else {
		eng.Graph.AddEdge(eng.currentFile, node, Def)
	}
}


func (eng *parseEngine) Visit(node ast.Node) bool {
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


func (eng *parseEngine) VisitBody(node ast.Node) bool {
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