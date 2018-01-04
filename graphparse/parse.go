package graphparse

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io/ioutil"
	"strconv"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/loader"
)

var packageFilePath = "/Users/liamz/go/src/github.com/twitchyliquid64/subnet/"
var optIncludeFilesAsNodes = false

var prog *loader.Program
var pkginfo *loader.PackageInfo 
var Graph *graph
var fset *token.FileSet
var currentFile *ast.File
var fileLookup = make(map[string]packageFileInfo)

type packageFileInfo struct {
	file *ast.File
	Code string `json:"code"`
	Pos token.Pos `json:"pos"`
	// Nodes map[postuple] `json:"nodes"`
}

var canonicalRefsToNodes map[string]nodeid
var objsIdentified map[string]nodeid
var importedPackages map[string]*node

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
	canonicalRefsToNodes = make(map[string]nodeid)
	objsIdentified = make(map[string]nodeid)
	importedPackages = make(map[string]*node)

	for _, f := range pkginfo.Files {
		currentFilePath := fset.File(f.Package).Name()
		fileName, _ := filepath.Rel(packageFilePath, currentFilePath)

		code, err := ioutil.ReadFile(currentFilePath)
		if err != nil {
			panic(err)
		}

		fileLookup[fileName] = packageFileInfo{
			file: currentFile,
			Code: string(code),
			Pos: f.Pos(),
		}

		fmt.Println("Processing", fileName)
		ast.Inspect(f, Visit)
		// ast.Print(fset, f)
	}
}

type CanonicalNode struct {
	val interface{}
}

func GetCanonicalNodeId(ref string) (nodeid) {
	if _, ok := canonicalRefsToNodes[ref]; !ok {
		canonicalRefsToNodes[ref] = pointerToId(&CanonicalNode{})
	}
	return canonicalRefsToNodes[ref]
}

func pointerToStr(ptr interface{}) string {
	return fmt.Sprintf("%p", &ptr)
}

func pointerToId(ptr interface{}) nodeid {
	if i, err := strconv.ParseInt(pointerToStr(ptr), 0, 64); err != nil {
		panic(err)
	} else {
		return nodeid(i)
	}
}


func getIdOfIdent(node *ast.Ident) (nodeid, error) {
	// https://golang.org/src/go/types/universe.go
	if obj := types.Universe.Lookup(node.Name); obj != nil {
		return nodeid(-1), fmt.Errorf("universe object ", node)
	}

	// I wrote this in a subconcious spree of, "I have a gut feeling that this will do it"
	obj := pkginfo.ObjectOf(node)

	objId := string(obj.Pos()) + obj.Id()
	// fmt.Println(obj.Name(), obj.Pos(), objId)

	if obj != nil {
		if id, ok := objsIdentified[objId]; ok {
			return id, nil
		} else {
			id := pointerToId(obj)
			objsIdentified[objId] = id
			return id, nil
		}
	}

	return nodeid(-1), fmt.Errorf("unexpected error", obj)
}

var pkgIdentNode *node
var currentFileNode *node

var parentNode *node


func getIdOfNode(node ast.Node) nodeid {
	switch x := node.(type) {
	case *ast.ImportSpec:
		importName, err := strconv.Unquote(x.Path.Value)
		if err != nil {
			panic(err)
		}
		id := GetCanonicalNodeId(importName)
		return id
	}
	panic(1)
}


func Visit(node ast.Node) bool {
	switch x := node.(type) {
	case *ast.File:
		if pkgIdentNode == nil {
			pkgIdent := x.Name
			pkgIdentId := pointerToId(pkgIdent)
			pkgIdentNode = NewNode(pkgIdent, pkgIdentId, x.Name.Name)
		}

		if optIncludeFilesAsNodes {
			fileName, err := filepath.Rel(packageFilePath, fset.File(x.Package).Name())
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return false
			}

			fileNodeId := GetCanonicalNodeId(fileName)
			currentFileNode = NewNode(x, fileNodeId, fileName)
			currentFileNode.extraAttrs = "[color=\"red\"]"
			Graph.AddEdge(pkgIdentNode, currentFileNode)
		}

	case *ast.ImportSpec:
		importName, _ := strconv.Unquote(x.Path.Value)
		importedPackageNode := NewNode(nil, getIdOfNode(x), importName)
		Graph.AddEdge(importedPackageNode, pkgIdentNode)
		// currentFileNode.extraAttrs = "[color=\"red\"]"

		// THIS IS BAD BELOW:
		// Graph.AddEdge(pkgIdentNode, importedPackageNode)
		return true

	case *ast.TypeSpec:
		typeId, err := getIdOfIdent(x.Name)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return true
		}

		typeNode := NewNode(x.Name, typeId, x.Name.Name)

		if optIncludeFilesAsNodes {
			Graph.AddEdge(currentFileNode, typeNode)
		} else {
			Graph.AddEdge(pkgIdentNode, typeNode)
		}

		// If struct, loop over fields
		switch y := x.Type.(type) {
		case *ast.StructType:
			if y.Fields != nil {
				for _, field := range y.Fields.List {
					switch y := field.Type.(type) {
					case *ast.SelectorExpr:
						// selection := pkginfo.Selections[y]
						// TODO
						// ast.Print(fset, y)
						// fmt.Println(y, selection)

						fieldTypeId, err := getIdOfIdent(y.Sel)
						if err != nil {
							fmt.Fprintln(os.Stderr, err.Error())
							return true
						}

						fieldType := NewNode(y, fieldTypeId, y.Sel.Name)

						Graph.AddEdge(typeNode, fieldType)
					case *ast.Ident:
						// fmt.Fprintln(os.Stderr, "parsing type - missed Ident field", field)
						// ast.Print(fset, field.Names[0])
						identId, err := getIdOfIdent(field.Names[0])

						if err != nil {
							fmt.Fprintln(os.Stderr, err.Error())
							return true
						}

						// fieldType := NewNode(field.Names[0], identId, y.Name)
						fieldNode := NewNode(field.Names[0], identId, field.Names[0].Name)

						// Graph.AddEdge(typeNode, fieldType)
						Graph.AddEdge(typeNode, fieldNode)

					case *ast.StarExpr:
						var fieldTypeIdent *ast.Ident

						switch z := y.X.(type) {
						case *ast.Ident:
							fieldTypeIdent = z
						case *ast.SelectorExpr:
							// selection := pkginfo.Selections[z]
							// fmt.Println(pkginfo.Selections)
							fieldTypeIdent = z.Sel
						default:
							fmt.Fprintln(os.Stderr, "parsing StarExpr field - missed StarExpr.X type", field)
							return true
						}

						fieldTypeId, err := getIdOfIdent(fieldTypeIdent)
						if err != nil {
							fmt.Fprintln(os.Stderr, err.Error())
							return true
						}

						fieldType := NewNode(fieldTypeIdent, fieldTypeId, fieldTypeIdent.Name)
						Graph.AddEdge(typeNode, fieldType)

					case *ast.ChanType:
						var fieldType *ast.Ident
						switch z := y.Value.(type) {
						case *ast.StarExpr:
							fieldType = z.X.(*ast.Ident)
						default:
							fmt.Fprintln(os.Stderr, "parsing ChanType field - missed Value type", field)
						}

						fieldTypeId, err := getIdOfIdent(fieldType)
						if err != nil {
							fmt.Fprintln(os.Stderr, err.Error())
							return true
						}

						fieldTypeNode := NewNode(fieldType, fieldTypeId, fieldType.Name)
						Graph.AddEdge(typeNode, fieldTypeNode)
						// fmt.Fprintln(os.Stderr, "parsing type - missed ChanType field", field)
						// ast.Print(fset, field)
					}

				}
			}
			break
			
		default:
			fmt.Fprintln(os.Stderr, "parsing type - missed type", y)
		}

	case *ast.FuncDecl:
		// Function as parent
		funcId, err := getIdOfIdent(x.Name)
		if err != nil {
			// Function is not referenced outside of this file
			// Thus does not have a .Obj
			fmt.Fprintln(os.Stderr, err.Error())
			return true
		}

		funcNode := NewNode(x, funcId, x.Name.Name)
		parentNode = funcNode

		if x.Recv != nil && len(x.Recv.List) == 1 {
			recv := x.Recv.List[0]
			switch y := recv.Type.(type) {
			case *ast.StarExpr:
				recvId, err := getIdOfIdent(y.X.(*ast.Ident))

				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					return true
				}

				// recvVarName is the 'c' in (c *Client)
				// recvVarName := recv.Names[0].Name
				recvTypeName := y.X.(*ast.Ident).Name
				recvType := NewNode(y, recvId, recvTypeName)

				Graph.AddEdge(funcNode, recvType)
			default:
				fmt.Fprintln(os.Stderr, "parsing receiver - missed type", recv)
			}

			// Graph.AddEdge(pkgIdentNode, funcNode)
		} else {
			if optIncludeFilesAsNodes {
				Graph.AddEdge(currentFileNode, funcNode)
			} else {
				Graph.AddEdge(pkgIdentNode, funcNode)
			}

		}

		// Loop over return values
		if x.Type.Results != nil {
			for _, funcResult := range x.Type.Results.List {
				// each *ast.Field
				switch y := funcResult.Type.(type) {
				case *ast.StarExpr:
					starExprIdent := y.X.(*ast.Ident)

					funcResultId, err := getIdOfIdent(starExprIdent)
					if err != nil {
						fmt.Fprintln(os.Stderr, err.Error())
						return true
					}
					funcResultNode := NewNode(funcResult, funcResultId, starExprIdent.Name)
					Graph.AddEdge(funcNode, funcResultNode)
					break

				case *ast.Ident:
					funcResultId, err := getIdOfIdent(y)
					if err != nil {
						// Function is not referenced outside of this file
						// Thus does not have a .Obj
						fmt.Fprintln(os.Stderr, err.Error())
						return true
					}

					funcResultNode := NewNode(funcResult, funcResultId, y.Name)
					Graph.AddEdge(funcNode, funcResultNode)
					break

				default:
					fmt.Fprintln(os.Stderr, "parsing result - missed type", funcResult)
				}
			}
		}

		return true
	
	// 
	
	// func (s *Scope) Innermost(pos token.Pos) *Scope
	case *ast.CallExpr:
		// ftype, fval := pkginfo.Types[x]
		switch y := x.Fun.(type) {
		case *ast.SelectorExpr:
			obj := pkginfo.ObjectOf(y.Sel)
			pkg := obj.Pkg()
			if pkg != nil {
				if pkg.Name() == "subnet" {
					id, err := getIdOfIdent(y.Sel)
					if err != nil {
						fmt.Fprintln(os.Stderr, err.Error())
						return true
					}
		
					callNode := NewNode(y, id, y.Sel.Name)
					// Graph.AddEdge(callNode, parentNode)
					Graph.AddEdge(parentNode, callNode)
				} else {
					importName := pkg.Path()
					pkgId := GetCanonicalNodeId(importName)
					importPkgNode := NewNode(nil, pkgId, pkg.Name())


					id, err := getIdOfIdent(y.Sel)
					if err != nil {
						fmt.Fprintln(os.Stderr, err.Error())
						return true
					}
		
					callNode := NewNode(y, id, y.Sel.Name)

					Graph.AddEdge(importPkgNode, callNode)
					Graph.AddEdge(callNode, parentNode)
					// importPkgNode := 
					// package -> thing -> otherthing <- subnet
				}
			}
			
	
		case *ast.Ident:
			// id, _ := getIdOfIdent(y)
			// fmt.Println("unreg'd ident", y.Name, id)
		default:
			fmt.Fprintln(os.Stderr, "parsing call - missed type", y)
		}
	default:
		// fmt.Fprintf(os.Stderr, "parsing - missed type %T\n", x)
		return true
	}
	return true
}
