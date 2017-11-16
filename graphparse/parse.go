package graphparse

import (
	"strconv"
	"go/types"
	"go/token"
	"go/ast"
	"fmt"
	"errors"
	"os"
)




type Visitor struct {
	Graph *node
	Pkg *types.Package
	fset *token.FileSet
}

func NewVisitor(rootAstNode ast.Node, pkg *types.Package, fset *token.FileSet) Visitor {
	v := Visitor{Pkg: pkg, fset: fset}
	v.Graph = &node{parent: nil, _value: rootAstNode, label: "/"}
	return v
}

func (v Visitor) newVisitorWithParent(parent ast.Node, id nodeid) (w Visitor) {
	w = v
	w.Graph = w.Graph.AddChild(parent, "", id)
	return w
}

func (v Visitor) registerChild(child ast.Node, id nodeid) {
	v.Graph.AddChild(child, "", id)
}


var idcounter int = 0
func newNodeId() nodeid {
	idcounter++
	return nodeid(idcounter)
}


func pointerToNodeid(ptr interface{}) nodeid {
	if i, err := strconv.ParseInt(fmt.Sprintf("%p", &ptr), 0, 64); err != nil {
		panic(err)
	} else {
		return nodeid(i)
	}
}

func (vis Visitor) getIdFromPointer(node *ast.Ident) (nodeid, error) {
	// defer func() {
 //        if r := recover(); r != nil {
 //            fmt.Println("BAD BAD NODE without Obj")
 //            // ast.Fprint(os.Stderr, vis.fset, node, nil)
 //        }
 //    }() 
	
	if node.Obj != nil {
		return pointerToNodeid(node.Obj), nil
	} else {
		if obj := vis.Pkg.Scope().Lookup(node.Name); obj != nil {
			return pointerToNodeid(node.Obj), nil
		}

		return nodeid(-1), errors.New("couldn't get node Obj")
		// fmt.Println(node)
		// panic("no node.Obj found")
	}
}


type FindIdentVisitor struct {
	Ident *ast.Ident
}
func (this FindIdentVisitor) Visit(node ast.Node) ast.Visitor {
	switch x := node.(type) {
		case *ast.Ident:
			this.Ident = x
			return nil
		default:
			return this
	}

	return nil
}


func (vis Visitor) walkNodeFindIdent(parent ast.Node) (*ast.Ident, error) {
	// identVis := new(FindIdentVisitor)
	// ast.Walk(identVis, parent)
	var ident *ast.Ident
	ast.Inspect(parent, func(node ast.Node) bool {
		switch x := node.(type) {
			case *ast.Ident:
				ident = x
				return false
			default:
				return true
		}
		return true
	})
	// ident := identVis.Ident

	if ident == nil {
		return nil, errors.New("no ident found")
	}

	return ident, nil
}


func (vis Visitor) Visit(node ast.Node) (ast.Visitor) {
	// goDeeperWithIdent := func(parent *ast.Node, x *ast.Ident) Visitor {
	// 	var id nodeid
	// 	if x.Obj != nil {
	// 		if i, err := strconv.ParseInt(fmt.Sprintf("%p", x.Obj), 0, 64); err != nil {
	// 			panic(err)
	// 		} else {
	// 			id = nodeid(i)
	// 		}				
	// 	} else {
	// 		id = newNodeId()
	// 	}

	// 	// // v.Pkg.Scope()
	// 	// if x.Obj != nil {
	// 	// 	// fmt.Println(x.Obj)
	// 	// } else {
	// 	// 	fmt.Println(parent, x.Name)
	// 	// }

	// 	// fmt.Println(1)

	// 	return vis.goDeeper(x, x.Name, id)
	// }

	switch x := node.(type) {
		case *ast.FuncDecl:

			// defer func() {
		 //        if r := recover(); r != nil {
		 //            fmt.Println("BAD BAD NODE without Type.Results.List")
		 //            ast.Fprint(os.Stderr, vis.fset, x, nil)
		 //        }
		 //    }() 

			// ast.Fprint(os.Stderr, vis.fset, x, nil)
			fmt.Println("FUNC:", x.Name.Name)

			// Function as parent
			id, err := vis.getIdFromPointer(x.Name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Can't get id for node", x.Name)
				return vis
			}
			newVis := vis.newVisitorWithParent(x, id)

			// Edges between (Func => ReturnType)
			if x.Type.Results != nil {
				for _, funcResult := range x.Type.Results.List {
					// walk get Ident.Obj
					// /*ident := */funcResult)
					_, err := vis.walkNodeFindIdent(funcResult.Type)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Can't find Ident within node", funcResult.Type)
						return newVis
					}
					// newVis.registerChild(ident, vis.getIdFromPointer(ident.Name))
				}				
			}

			return newVis
			// return vis


		// case *ast.TypeSpec:
			// return vis.newVisitorWithParent(x, vis.getIdFromPointer(x.Name))

		// case *ast.AssignStmt:
			// Link as such:
			// obj.field = thing.thingo
			
			// LHS ONLY:
			// Link (field => obj)

			// LHS TO RHS
			// Link (field => thingo)

			// RHS only
			// Link thingo => thing

			// So it looks like:
			// obj => field => thingo
			// thing => thingo

			// lhs := x.Lhs
			// rhs := x.Rhs

/*
	Something like:
	tlsConf, err := conn.TLSConfig(certPemPath, keyPemPath, caCertPath)

	(conn, TLSConfig)
	(tlsConf, TLSConfig)

						  Lhs: []ast.Expr (len = 2) {
   672  .  .  .  .  .  .  .  0: *ast.Ident {
   673  .  .  .  .  .  .  .  .  NamePos: /Users/liamz/parser/src/github.com/twitchyliquid64/subnet/subnet/client.go:50:2
   674  .  .  .  .  .  .  .  .  Name: "tlsConf"
   675  .  .  .  .  .  .  .  .  Obj: *ast.Object {
   676  .  .  .  .  .  .  .  .  .  Kind: var
   677  .  .  .  .  .  .  .  .  .  Name: "tlsConf"
   678  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 670)
   679  .  .  .  .  .  .  .  .  }
   680  .  .  .  .  .  .  .  }
   681  .  .  .  .  .  .  .  1: *ast.Ident {
   682  .  .  .  .  .  .  .  .  NamePos: /Users/liamz/parser/src/github.com/twitchyliquid64/subnet/subnet/client.go:50:11
   683  .  .  .  .  .  .  .  .  Name: "err"
   684  .  .  .  .  .  .  .  .  Obj: *ast.Object {
   685  .  .  .  .  .  .  .  .  .  Kind: var
   686  .  .  .  .  .  .  .  .  .  Name: "err"
   687  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 670)
   688  .  .  .  .  .  .  .  .  }
   689  .  .  .  .  .  .  .  }
   690  .  .  .  .  .  .  }
   691  .  .  .  .  .  .  TokPos: /Users/liamz/parser/src/github.com/twitchyliquid64/subnet/subnet/client.go:50:15
   692  .  .  .  .  .  .  Tok: :=
   693  .  .  .  .  .  .  Rhs: []ast.Expr (len = 1) {
   694  .  .  .  .  .  .  .  0: *ast.CallExpr {
   695  .  .  .  .  .  .  .  .  Fun: *ast.SelectorExpr {
   696  .  .  .  .  .  .  .  .  .  X: *ast.Ident {
   697  .  .  .  .  .  .  .  .  .  .  NamePos: /Users/liamz/parser/src/github.com/twitchyliquid64/subnet/subnet/client.go:50:18
   698  .  .  .  .  .  .  .  .  .  .  Name: "conn"
   699  .  .  .  .  .  .  .  .  .  }
   700  .  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
   701  .  .  .  .  .  .  .  .  .  .  NamePos: /Users/liamz/parser/src/github.com/twitchyliquid64/subnet/subnet/client.go:50:23
   702  .  .  .  .  .  .  .  .  .  .  Name: "TLSConfig"
   703  .  .  .  .  .  .  .  .  .  }
   704  .  .  .  .  .  .  .  .  }
   705  .  .  .  .  .  .  .  .  Lparen: /Users/liamz/parser/src/github.com/twitchyliquid64/subnet/subnet/client.go:50:32
   706  .  .  .  .  .  .  .  .  Args: []ast.Expr (len = 3) {
   707  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
   708  .  .  .  .  .  .  .  .  .  .  NamePos: /Users/liamz/parser/src/github.com/twitchyliquid64/subnet/subnet/client.go:50:33
   709  .  .  .  .  .  .  .  .  .  .  Name: "certPemPath"
   710  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 611)
   711  .  .  .  .  .  .  .  .  .  }
   712  .  .  .  .  .  .  .  .  .  1: *ast.Ident {
   713  .  .  .  .  .  .  .  .  .  .  NamePos: /Users/liamz/parser/src/github.com/twitchyliquid64/subnet/subnet/client.go:50:46
   714  .  .  .  .  .  .  .  .  .  .  Name: "keyPemPath"
   715  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 620)
   716  .  .  .  .  .  .  .  .  .  }
   717  .  .  .  .  .  .  .  .  .  2: *ast.Ident {
   718  .  .  .  .  .  .  .  .  .  .  NamePos: /Users/liamz/parser/src/github.com/twitchyliquid64/subnet/subnet/client.go:50:58
   719  .  .  .  .  .  .  .  .  .  .  Name: "caCertPath"
   720  .  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 629)
   721  .  .  .  .  .  .  .  .  .  }
   722  .  .  .  .  .  .  .  .  }
   723  .  .  .  .  .  .  .  .  Ellipsis: -
   724  .  .  .  .  .  .  .  .  Rparen: /Users/liamz/parser/src/github.com/twitchyliquid64/subnet/subnet/client.go:50:68
   725  .  .  .  .  .  .  .  }
   726  .  .  .  .  .  .  }
   727  .  .  .  .  .  }
*/
   			// return vis
			// return goDeeperWithIdent(x, x.Name)

		default:
			return vis
	}
	return vis
}

// case *ast.Package:
// 	fmt.Println("#pkg")
// 	for abspath, srcfile := range x.Files {
// 		path, err := filepath.Rel(dir, abspath)
// 		if err != nil {
// 			panic(err)
// 		}

// 		ast.Walk(v.goDeeper(srcfile, path), srcfile)
// 		fmt.Println("Processing", path)
// 	}