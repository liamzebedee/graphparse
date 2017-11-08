package main

import (
	"fmt"
	"go/parser"
	"go/ast"
	"go/token"
	"strings"
	"path/filepath"
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




type Visitor struct {
	parents stack
}

func NewVisitor() Visitor {
	return Visitor{}
}

func (v Visitor) deeperStack(parentToAdd string) (w Visitor) {
	w = Visitor{}
	w.parents = v.parents.Push(parentToAdd)
	return w
}

func (v Visitor) Visit(node ast.Node) (w ast.Visitor) {
	switch x := node.(type) {
		case *ast.Package:
			for path, f := range x.Files {
				relpath, err := filepath.Rel(dir, path)
				if err != nil {
					panic(err)
				}

				pkgName := x.Name
	  			
				// v.deeperStack("(" + relpath + ")").Visit(f)
				ast.Walk(v.deeperStack(pkgName + "(" + relpath + ")"), f)

				fmt.Println(relpath)
	  		}

		case *ast.TypeSpec:
			return v.deeperStack(x.Name.Name)
		
		case *ast.FuncDecl:
			// Functions have one receiver
			// This is the parent
			if recv := x.Recv; recv != nil {
				// Method on struct
				// ast.Print(nil,x)
				structName := recv.List[0].Names[0].Name
				// fmt.Println(recv.List[0].Names[0].Name)
				// fmt.Println(structName)

				return v.deeperStack(structName)	
			} else {
				// Just a function
				return v.deeperStack(x.Name.Name)
			}

		case *ast.Ident:
			fmt.Println(strings.Join(v.parents, ".") + "." + x.Name)

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

		visitor := NewVisitor()
		ast.Walk(visitor, pkg)

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