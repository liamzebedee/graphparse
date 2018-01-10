package testing

import (
	"testing"
  "path/filepath"
  "go/token"
  "go/parser"
  "go/ast"
  
  // "github.com/liamzebedee/graphparse/graphparse"
)

var testfile, _ = filepath.Abs("./testpkg/method.go")

func TestParseMethodReceiver(t *testing.T) {
  // graphparse.GenerateCodeGraph("./testpkg", "")
  fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, testfile, nil, 0)
  ast.Print(fset, file)
  // panic(err)
}


