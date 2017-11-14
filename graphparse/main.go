package graphparse

import (
	// "io"
	// "strings"
	// "path/filepath"
	// "io/ioutil"
	// "strconv"
	// godsUtils "github.com/emirpasic/gods/utils" 
	// "github.com/emirpasic/gods/trees/btree"
	// "github.com/emirpasic/gods/maps/treemap"
	"golang.org/x/tools/go/loader"
	"go/parser"
	"go/ast"
	// "go/token"
)

const dir string = "/Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/"

func Stuff() {
	pkgpath := "github.com/twitchyliquid64/subnet/subnet"

	conf := loader.Config{ParserMode: parser.ParseComments}
	conf.Import(pkgpath)
	prog, err := conf.Load()
	if err != nil {
		panic(err)
	}

	pkginfo := prog.Package(pkgpath)
	// fset := prog.Fset
	// ast.Print(fset, rootAst)
	
	visitor := NewVisitor(nil, pkginfo.Pkg)

	for i, f := range pkginfo.Files {
		if i < 3 { continue }
		ast.Walk(visitor, f)
		if i == 3 { break }
		
		

		// ast.Print(fset, f)
	}
	visitor.Graph.ToDot()
	visitor.Graph.String()
}