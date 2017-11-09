package main

import (
	"fmt"
	"go/parser"
	"go/ast"
	"go/token"
	// "io"
	// "strings"
	"path/filepath"
	// "io/ioutil"
	"bufio"
	"os"
	// "github.com/liamzebedee/graphparse/graphparse"
	// "github.com/gonum/graph"
)


type stack []string

func (s stack) Push(v string) stack {
    return append(s, v)
}

func (s stack) Pop() (stack, string) {
    l := len(s)
    return  s[:l-1], s[l-1]
}


// two approaches:
// - big global canonical lookup table of scope (bad, error-prone)
// - small recursive progressive graph AST approach (better)



type node struct {
	children []*node
	value interface{}
	label string
}
func (n *node) AddChild(val interface{}, label string) *node {
	child := &node{value: val, label: label}
	n.children = append(n.children, child)
	return child
}



func (this *node) ToDot() {
	f, err := os.OpenFile("/Users/liamz/parser/src/github.com/liamzebedee/graphparse/graph.dot", os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	if true {
		// f = os.Stdout
		defer f.Close()
	} else {
		defer f.Close()
	}
	w := bufio.NewWriter(f)
	defer w.Flush()


	// root
	w.WriteString("digraph graphname {\n")


	addChildren(this, w)

	w.WriteString("}\n")
}

// each child of children
func addChildren(n *node, w *bufio.Writer) {
	nodeA := n.label

	for _, child := range n.children {
		nodeB := child.label
		w.WriteString("\""+nodeA+"\"" + " -> " + "\"" + nodeB + "\"" + ";\n")

		if len(child.children) > 0 {
			addChildren(child, w)
		}
	}
}



type Visitor struct {
	parent ast.Node
	parentLabel string
	Graph *node
}

func NewVisitor(rootAstNode ast.Node) Visitor {
	v := Visitor{}
	// v.graph = graphparse.NewGraph()
	v.Graph = &node{value: rootAstNode, label: "/"}
	// v.parent = rootAstNode
	// v.parentLabel = ""
	return v
}

func (v Visitor) goDeeper(node ast.Node, label string) (w Visitor) {
	w = v
	w.Graph = w.Graph.AddChild(node, label)
	return w
}

// will every child only belong to one parent? 
// yeah, it's pagerank
// maybe in v2 we will have variables linked based on lexic similarity

func (v Visitor) registerNode(node ast.Node, label string) {
	// v.graph.Add(v.parent, node)
	v.Graph.AddChild(node, label)
}

func (v Visitor) Visit(node ast.Node) (w ast.Visitor) {
	switch x := node.(type) {
		case *ast.Package:
			for abspath, srcfile := range x.Files {
				path, err := filepath.Rel(dir, abspath)
				if err != nil {
					panic(err)
				}

				// pkgName := x.Name

				// label := "(" + path + ")"
				label := path

				ast.Walk(v.goDeeper(srcfile, label), srcfile)
				fmt.Println("Processing", path)
	  		}

		case *ast.TypeSpec:
			return v.goDeeper(x, x.Name.Name)
		
		case *ast.StructType:
			return v.goDeeper(x, x.Fields.List[0].Names[0].Name)

		case *ast.FuncDecl:
			if recv := x.Recv; recv != nil {
				// Method on struct
				structName := recv.List[0].Names[0].Name

				return v.goDeeper(x, structName)
			} else {
				// Just a function
				return v.goDeeper(x, x.Name.Name)
			}

		// case *ast.Ident:
			// v.registerNode(x, x.Name)
			// fmt.Println(strings.Join(v.parents, ".") + "." + x.Name)

		default:
			return v
	}
	return w
}

const dir string = "/Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/"

func main() {
	fset := token.NewFileSet()
	// dir := "/Users/liamz/parser/src/github.com/liamzebedee/graphparse/testsrc/"

	pkgs, err := parser.ParseDir(fset, dir, nil, 0)

	if err != nil {
		fmt.Println(err)
	}

	for name, pkg := range pkgs {
		fmt.Println("Package:", name)
		// fmt.Println(name, pkg)

		visitor := NewVisitor(pkg)
		ast.Walk(visitor, pkg)
		visitor.Graph.ToDot()

		// for path, file := range pkg.Files {
		// 	fmt.Println("File:", path)
		// 	// ast.Print(fset, file)
			
		// 	// for _, decl := range file.Decls {
		// 	// 	// var parent ast.Node = ast.Node(decl)
		// 	// 	// var parent string = ""


		// 	// 	// actually need to recurse and be able to keep track of the parent hierarchy 
		// 	// 	// otherwise we would never know when we go up level, because parent wouldn't change
		// 	// 	// H: does parent change suitably

		// 	// 	// we need to know depth to make a stack

				

		// 	// 	// ast.Inspect(decl, func(node ast.Node) bool {
					
					
		// 	// 	// 	return true
		// 	// 	// })
		// 	// }

		// 	fmt.Println("")

		// 	// fmt.Println(file)
		// }

		// fmt.Println("Merging package files...")
		// pkgFile := ast.MergePackageFiles(pkg, ast.FilterFuncDuplicates & ast.FilterUnassociatedComments & ast.FilterImportDuplicates)

		// // fmt.Println(pkgFile)
		// for _, dec := range pkgFile.Decls {
		// 	fmt.Println(dec)
		// }
	}



	// ast.Inspect(f, func(n ast.Node) bool {
	// 	var s string
	// 	switch x := n.(type) {
	// 	case *ast.BasicLit:
	// 		s = x.Value
	// 	case *ast.Ident:
	// 		s = x.Name
	// 	}
	// 	if s != "" {
	// 		fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
	// 	}
	// 	return true
	// })

}