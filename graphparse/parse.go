package graphparse

import (
	
	"go/types"
	"go/ast"
)

func makeEdges(edges *[]edge) walkfn {
	return func(node *node) {
		for _, child := range node.children {
			*edges = append(*edges, edge{node.Id(), child.Id()})
		}
	}
}

type Visitor struct {
	Graph *node
	Pkg *types.Package
}

func NewVisitor(rootAstNode ast.Node, pkg *types.Package) Visitor {
	v := Visitor{Pkg: pkg}
	v.Graph = &node{parent: nil, value: rootAstNode, label: "/"}
	return v
}

func (v Visitor) goDeeper(node ast.Node, label string) (w Visitor) {
	w = v
	w.Graph = w.Graph.AddChild(node, label)
	return w
}

func (v Visitor) registerNode(node ast.Node, label string) {
	v.Graph.AddChild(node, label)
}



func (v Visitor) Visit(node ast.Node) (w ast.Visitor) {
	switch x := node.(type) {
		case *ast.Ident:
			// go deeper
			return v.goDeeper(x, x.Name)

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