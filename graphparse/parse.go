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

func getIdOfObj(obj types.Object) (nodeid, error) {
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
	return nodeid(-1), fmt.Errorf("unexpected error", obj)
}


func getIdFromPointer(node *ast.Ident) (nodeid, error) {
	// https://golang.org/src/go/types/universe.go
	if obj := types.Universe.Lookup(node.Name); obj != nil {
		return nodeid(-1), fmt.Errorf("universe object ", node)
	}


	// I wrote this in a subconcious spree of, "I have a gut feeling that this will do it"
	obj := pkginfo.ObjectOf(node)
	return getIdOfObj(obj)
	
	// obj = pkginfo.Implicits[node]
	// if obj != nil {
	// 	return pointerToNodeid(obj), nil
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

		case *ast.ImportSpec:
			// pkgId, err := getIdFromPointer()
			obj := pkginfo.Implicits[x]
			importId, err := getIdOfObj(obj)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return true
			}

			importedPackageNode := NewNode(x, importId, obj.Name())
			Graph.AddEdge(importedPackageNode, pkgIdentNode)
			// THIS IS BAD BELOW:
			// Graph.AddEdge(pkgIdentNode, importedPackageNode)
			return true

		case *ast.TypeSpec:
			typeId, err := getIdFromPointer(x.Name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return true
			}

			typeNode := NewNode(x.Name, typeId, x.Name.Name)
			typeNode.extraAttrs = "[color=\"red\"]"
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
							// fmt.Fprintln(os.Stderr, "parsing type - missed Ident field", field)
							// ast.Print(fset, field.Names[0])
							identId, err := getIdFromPointer(field.Names[0])

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
								fieldTypeIdent = z.Sel
							default:
								fmt.Fprintln(os.Stderr, "parsing StarExpr field - missed StarExpr.X type", field)
								return true
							}

							fieldTypeId, err := getIdFromPointer(fieldTypeIdent)
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

							fieldTypeId, err := getIdFromPointer(fieldType)
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
			funcId, err := getIdFromPointer(x.Name)
			if err != nil {
				// Function is not referenced outside of this file
				// Thus does not have a .Obj
				fmt.Fprintln(os.Stderr, err.Error())
				return true
			}

			funcNode := NewNode(x, funcId, x.Name.Name)

			if x.Recv != nil && len(x.Recv.List) == 1 {
				recv := x.Recv.List[0]
				switch y := recv.Type.(type) {
				case *ast.StarExpr:
					recvId, err := getIdFromPointer(y.X.(*ast.Ident))

					if err != nil {
						fmt.Fprintln(os.Stderr, err.Error())
						return true
					}

					// recvVarName is the 'c' in (c *Client)
					// recvVarName := recv.Names[0].Name
					recvTypeName := y.X.(*ast.Ident).Name
					recvType := NewNode(y, recvId, recvVarName)

					Graph.AddEdge(funcNode, recvType)
				default:
					fmt.Fprintln(os.Stderr, "parsing receiver - missed type", recv)
				}

				// Graph.AddEdge(pkgIdentNode, funcNode)
			} else {
				Graph.AddEdge(pkgIdentNode, funcNode)
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
			// fmt.Fprintf(os.Stderr, "parsing - missed type %T\n", x)
			return true
	}
	return true
}