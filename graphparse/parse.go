package graphparse

import (
	"strconv"
	"go/types"
	"go/ast"
	"fmt"
)



type Visitor struct {
	Graph *node
	Pkg *types.Package
}

func NewVisitor(rootAstNode ast.Node, pkg *types.Package) Visitor {
	v := Visitor{Pkg: pkg}
	v.Graph = &node{parent: nil, _value: rootAstNode, label: "/"}
	return v
}

func (v Visitor) goDeeper(node ast.Node, label string, id nodeid) (w Visitor) {
	w = v
	w.Graph = w.Graph.AddChild(node, label, id)
	return w
}

var idcounter int = 0
func newNodeId() nodeid {
	idcounter++
	return nodeid(idcounter)
}
var parent ast.Node
func (v Visitor) Visit(node ast.Node) (ast.Visitor) {
	// defer func() { parent = node }()
	if(node == nil) {
		fmt.Println("0")
	}

	switch x := node.(type) {
		case *ast.Ident:
			var id nodeid
			if x.Obj != nil {
				if i, err := strconv.ParseInt(fmt.Sprintf("%p", x.Obj), 0, 64); err != nil {
					panic(err)
				} else {
					id = nodeid(i)
				}				
			} else {
				id = newNodeId()
			}

			// v.Pkg.Scope()
			if x.Obj != nil {
				// fmt.Println(x.Obj)
			} else {
				fmt.Println(parent, x.Name)
			}

			fmt.Println(1)

			return v.goDeeper(x, x.Name, id)

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

		// case *ast.StructType:
			// fmt.Println("#struct")
			// v.goDeeper(x, x.Fields.List[0].Names[0].Name)

		// case *ast.Ident:
			// v.registerNode(x, x.Name)
			// fmt.Println(strings.Join(v.parents, ".") + "." + x.Name)

		default:
			return v
	}
	return v
}