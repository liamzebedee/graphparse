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
	"fmt"
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
	fset := prog.Fset
	
	visitor := NewVisitor(nil, pkginfo.Pkg, fset)

	for i, f := range pkginfo.Files {
		fmt.Println("Processing", fset.File(f.Package).Name())
		ast.Walk(visitor, f)
		if i == 1 { break }

		ast.Print(fset, f)
	}
	// visitor.Graph.ToDot()
	visitor.Graph.String()
}