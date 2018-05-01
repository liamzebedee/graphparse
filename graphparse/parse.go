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
var Graph *graph
var fset *token.FileSet
var thisPackage string

var fileLookup = make(map[string]packageFileInfo)

var eng = parseEngine{}

type packageFileInfo struct {
	Code string `json:"code"`
	Pos token.Pos `json:"pos"`
}

func GenerateCodeGraphFromProg(prog *loader.Program, pkgpath, pkgFilePath string) {
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
	if obj.Pkg().Name() != thisPackage {
		return false
	}
	return true
}



func importToCanonicalKey(importSpec *ast.ImportSpec) string {
	importName, err := strconv.Unquote(importSpec.Path.Value)
	if err == nil {
		return importName
	} else {
		panic(err)
	}
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

	Graph.AddEdge(eng.rootPackage, eng.currentFile)
	
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

		typeNode := LookupOrCreateNode(obj, Struct, obj.Name())
		Graph.AddEdge(fromNode, typeNode)
	
	case *ast.Ident:
		obj, err := getObjFromIdent(typeSpec.Name)
		if err != nil {
			panic(err)
		}

		typeNode := LookupOrCreateNode(obj, Struct, obj.Name())
		Graph.AddEdge(fromNode, typeNode)
	
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

	funcNode := LookupOrCreateNode(obj, variant, funcDecl.Name.Name)
	eng.parentFunc = funcNode

	switch x := obj.Type().(type) {
	case *types.Signature:
		if isMethod {
			eng.parseMethod(funcNode, x.Recv())
		} else {
			Graph.AddEdge(fromNode, funcNode)
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

/*
FieldVal, MethodVal - something.X()
MethodExpr - fn := something.X; // fn(something, yada)
*/
func parseXOfSelectorExpr(selX ast.Expr) string {
	switch y := selX.(type) {
	case *ast.SelectorExpr:
		return parseXOfSelectorExpr(y.X) + "." + y.Sel.Name
	case *ast.Ident:
		return y.Name
	default:
		ParserLog.Printf("didn't understand X of selector %T", y)
	}
	return ""
}

type annotatedSelectorObject struct {
	kind types.SelectionKind
	types.Object
}

func getObjectsFromSelector(sel ast.Expr) (objs []annotatedSelectorObject) {
	switch x := sel.(type) {
	case *ast.SelectorExpr:
		sel := pkginfo.Selections[x]
		
		if sel == nil {
			ParserLog.Printf("skipping selector expr (likely qualified identifier) %T\n", x)
			break
		} else {
			annSelObj := annotatedSelectorObject{
				sel.Kind(),
				sel.Obj(),
			}
			objs = append(objs, annSelObj)
		}


		if ident, ok := x.X.(*ast.Ident); ok {
			obj := pkginfo.ObjectOf(ident)
			if !objIsWorthy(obj) {
				ParserLog.Printf("encountered unexpected bad obj in selector expression %T - %T", obj, x)
			} else {
				annSelObj := annotatedSelectorObject{
					sel.Kind(),
					obj,
				}
				objs = append(objs, annSelObj)
			}
			return objs
		}

		return append(objs, getObjectsFromSelector(x.X)...)
		
	default:
		ParserLog.Printf("didn't understand X of selector %T", x)
	}
	return objs
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
		eng.parseSelectorCallExpr2(sel)
		return
		
	case *ast.Ident:
		// continue as usual.
	default:
		ParserLog.Printf("missed type of call expr %T\n", x)
	}

	if err != nil || obj == nil {
		return
	}

	if !objIsWorthy(obj) {
		return
	}

	if obj.Pkg().Name() != thisPackage {
		// eng.parseExternalFuncCall(callExp, obj)
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
	if eng.parentFunc != nil {
		Graph.AddEdge(eng.parentFunc, funcNode)
	}
}

func (eng *parseEngine) parseSelectorCallExpr2(sel *types.Selection) {
	// recvTyp := sel.Recv()
	// recv := typeToObj(recvTyp)
	obj := sel.Obj()
	
	
	if objIsWorthy(obj) {
		objNode := LookupOrCreateNode(obj, FuncCall, obj.Name())
		Graph.AddEdge(eng.parentFunc, objNode)
	}

	// recvNode := LookupOrCreateNode(recv, Func, objA.Name())

	// nxt := eng.parentFunc
	
	// if objIsWorthy(objA) {
		// n1 := LookupOrCreateNode(objA, Func, objA.Name())
	// 	// Graph.AddEdge(eng.parentFunc, n1)
	// 	// Graph.AddEdge(nxt, n1)
	// 	nxt = n1
	// }

	// if objIsWorthy(b) {
	// 	n2 := LookupOrCreateNode(b, FuncCall, b.Name())
	// 	Graph.AddEdge(nxt, n2)
	// }
}

func (eng *parseEngine) parseSelectorCallExpr(objs []annotatedSelectorObject) {
	prev := eng.parentFunc

	for i := len(objs) - 1; i != -1; i-- {
		obj := objs[i]

		nodeTyp := FuncCall
		switch obj.kind {
		case types.FieldVal:
			nodeTyp = Field
		case types.MethodExpr, types.MethodVal:
			nodeTyp = Method
		}

		if eng.parentFunc.obj != obj.Object {
			currNode := LookupOrCreateNode(obj, nodeTyp, obj.Name())
			Graph.AddEdge(prev, currNode)
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
	funcNode := LookupOrCreateNode(obj, variant, obj.Name())

	Graph.AddEdge(importNode, funcNode)	
	
	// TEST
	if eng.parentFunc != nil {
		Graph.AddEdge(eng.parentFunc, funcNode)
		Graph.AddEdge(funcNode, eng.parentFunc)
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
	Graph.AddEdge(importedPackage, eng.rootPackage)
}




func Visit(node ast.Node) bool {
	switch x := node.(type) {
	case *ast.ImportSpec:
		// eng.parseImportSpec(x)
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