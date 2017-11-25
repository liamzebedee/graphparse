package graphparse

import (
	"strconv"
	"go/types"
	"go/token"
	"go/ast"
	"fmt"

	"golang.org/x/tools/go/loader"
	"os"
)

var pkgpath string
func DoMainStuff() {
	pkgpath := "github.com/twitchyliquid64/subnet/subnet"
	// pkgpath := "github.com/liamzebedee/graphparse/graphparse"

	conf := loader.Config{ParserMode: 0}
	conf.Import(pkgpath)
	prog, err := conf.Load()
	if err != nil {
		panic(err)
	}

	pkginfo = prog.Package(pkgpath)

	Pkg = pkginfo.Pkg
	fset = prog.Fset
	
	Graph = NewGraph()
	unresolvedIdentToId = make(map[*ast.Ident]nodeid)
	objsIdentified = make(map[string]nodeid)

	for i, f := range pkginfo.Files {
		currentFile = f

		fmt.Println("Processing", fset.File(f.Package).Name())
		ast.Inspect(f, Visit)
		if i == 12*1000 { break }
		// ast.Print(fset, f)
	}
	Graph.ToDot()
}

var pkginfo *loader.PackageInfo
var Graph *graph
var Pkg *types.Package
var fset *token.FileSet
var currentFile *ast.File
var unresolvedIdentToId map[*ast.Ident]nodeid
var pkgIdentNode *node

var objsIdentified map[string]nodeid

type Visitor struct {
}


func pointerToNodeid(ptr interface{}) nodeid {
	if i, err := strconv.ParseInt(fmt.Sprintf("%p", &ptr), 0, 64); err != nil {
		panic(err)
	} else {
		return nodeid(i)
	}
}

type CanonicalUnresolvedIdent struct {
	ident *ast.Ident
}


func getIdFromPointer(node *ast.Ident) (nodeid, error) {
	// https://golang.org/src/go/types/universe.go
	if obj := types.Universe.Lookup(node.Name); obj != nil {
		return nodeid(-1), fmt.Errorf("universe object ", node)
	}


	// I wrote this in a subconcious spree of, "I have a gut feeling that this will do it"
	obj := pkginfo.ObjectOf(node)
	objId := obj.Id()
	if obj != nil {
		if id, ok := objsIdentified[objId]; ok {
			return id, nil
		} else {
			id := pointerToNodeid(obj)
			objsIdentified[objId] = id
			return id, nil
		}
	}

	obj = pkginfo.Implicits[node]
	if obj != nil {
		return pointerToNodeid(obj), nil
	}

	
	// if node.Obj != nil {
	// 	return pointerToNodeid(node.Obj), nil
	// } else {
	// 	if _, obj := Pkg.Scope().LookupParent(node.Name, token.NoPos); obj != nil {
	// 		return pointerToNodeid(obj), nil
	// 	}

	// 	if id, ok := unresolvedIdentToId[node]; ok {
	// 		return id, nil
	// 	}
	// 	for _, ident := range currentFile.Unresolved {
	// 		if ident == node {
	// 			x := &CanonicalUnresolvedIdent{node}

	// 			id := pointerToNodeid(x)

	// 			unresolvedIdentToId[node] = id
	// 			return id, nil
	// 		}
	// 	}

	// 	x := &CanonicalUnresolvedIdent{node}
	// 	id := pointerToNodeid(x)
	// 	unresolvedIdentToId[node] = id
	// 	return id, nil

	return nodeid(-1), fmt.Errorf("couldn't get node Obj -", node)
	// }
}

// TODO  Names: []*ast.Ident (len = 2)


func Visit(node ast.Node) bool {
	switch x := node.(type) {
		case *ast.File:
			if pkgIdentNode == nil {
				pkgIdent := x.Name
				pkgIdentId := pointerToNodeid(pkgIdent)
				pkgIdentNode = NewNode(pkgIdent, pkgIdentId, x.Name.Name)
			}

		case *ast.TypeSpec:
			typeId, err := getIdFromPointer(x.Name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return true
			}

			typeNode := NewNode(x, typeId, x.Name.Name)
			Graph.AddEdge(pkgIdentNode, typeNode)

			// If struct, loop over fields
			switch y := x.Type.(type) {
			case *ast.StructType:
				if y.Fields != nil {
					for _, field := range y.Fields.List {
						switch y := field.Type.(type) {
						case *ast.SelectorExpr:
							fieldTypeId, err := getIdFromPointer(y.Sel)
							if err != nil {
								fmt.Fprintln(os.Stderr, err.Error())
								return true
							}
							
							fieldType := NewNode(y, fieldTypeId, y.Sel.Name)

							Graph.AddEdge(typeNode, fieldType)
						case *ast.Ident:
							fmt.Fprintln(os.Stderr, "parsing type - missed Ident field", field)
						case *ast.StarExpr:
							fmt.Fprintln(os.Stderr, "parsing type - missed StarExpr field", field)
						case *ast.ChanType:
							// fieldTypeId, err := getIdFromPointer(y.Sel)
							// if err != nil {
							// 	fmt.Fprintln(os.Stderr, err.Error())
							// 	return true
							// }
							
							// fieldType := NewNode(y, fieldTypeId, y.Sel.Name)

							// Graph.AddEdge(typeNode, fieldType)
							fmt.Fprintln(os.Stderr, "parsing type - missed ChanType field", field)
							break
						}

					}
				}
				break
			default:
				fmt.Fprintln(os.Stderr, "parsing type - missed type", y)
			}

		case *ast.FuncDecl:
			// Function as parent
			funcId, err := getIdFromPointer(x.Name)
			if err != nil {
				// Function is not referenced outside of this file
				// Thus does not have a .Obj
				fmt.Fprintln(os.Stderr, err.Error())
				return true
			}

			funcNode := NewNode(x, funcId, x.Name.Name)
			Graph.AddEdge(pkgIdentNode, funcNode)

			if x.Recv != nil && len(x.Recv.List) == 1 {
				recv := x.Recv.List[0]
				switch y := recv.Type.(type) {
				case *ast.StarExpr:
					recvId, err := getIdFromPointer(y.X.(*ast.Ident))

					if err != nil {
						fmt.Fprintln(os.Stderr, err.Error())
						return true
					}

					// recvName is 'c' in (c *Client)
					// recvVar := NewNode(y, recvId, recv.Names[0].Name)
					recvVar := NewNode(y, recvId, y.X.(*ast.Ident).Name)

					// recvType := NewNode(y, )
					Graph.AddEdge(funcNode, recvVar)
					break
				default:
					fmt.Fprintln(os.Stderr, "parsing receiver - missed type", recv)
				}
				

			}

			// Loop over return values
			if x.Type.Results != nil {
				for _, funcResult := range x.Type.Results.List {
					// each *ast.Field
					switch y := funcResult.Type.(type) {
					case *ast.StarExpr:
						starExprIdent := y.X.(*ast.Ident)

						funcResultId, err := getIdFromPointer(starExprIdent)
						if err != nil {
							fmt.Fprintln(os.Stderr, err.Error())
							return true
						}
						funcResultNode := NewNode(funcResult, funcResultId, starExprIdent.Name)
						Graph.AddEdge(funcNode, funcResultNode)
						break

					case *ast.Ident:
						funcResultId, err := getIdFromPointer(y)
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

		default:
			return true
	}
	return true
}