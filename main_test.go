package testing

import (
	"testing"
	"go/ast"
	"go/token"
	"go/parser"
)


func TestAST(t *testing.T) {
	fset := token.NewFileSet()
	const dir string = "/Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/client.go"
	f, _ := parser.ParseFile(fset, dir, nil, 0)
	ast.Print(fset, f)
}


/*

Function calls

Rhs: []ast.Expr (len = 1) {
  186											0: *ast.CallExpr {
  186												Fun: *ast.Ident {
  186													NamePos: /Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/client.go:106:29
  186													Name: "GetNetGateway"
  186												}
  186												Lparen: /Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/client.go:106:42
  187												Ellipsis: -
  187												Rparen: /Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/client.go:106:43
  187											}


Selectors on structs (something.subthing)

If: /Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/client.go:104:2
  180							Cond: *ast.BinaryExpr {
  180								X: *ast.SelectorExpr {
  181									X: *ast.Ident {
  181										NamePos: /Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/client.go:104:5
  181										Name: "c"
  181										Obj: *(obj @ 1344)
  181									}
  181									Sel: *ast.Ident {
  181										NamePos: /Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/client.go:104:7
  181										Name: "newGateway"
  181									}
  181								}


Methods on structs

  226			4: *ast.FuncDecl {
  226				Recv: *ast.FieldList {
  226					Opening: /Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/client.go:128:6
  226					List: []*ast.Field (len = 1) {
  226						0: *ast.Field {
  226							Names: []*ast.Ident (len = 1) {
  226								0: *ast.Ident {
  226									NamePos: /Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/client.go:128:7
  227									Name: "c"
  227									Obj: *ast.Object {
  227										Kind: var
  227										Name: "c"
  227										Decl: *(obj @ 2266)
  227									}
  227								}
  227							}
  227							Type: *ast.StarExpr {
  227								Star: /Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/client.go:128:9
  228								X: *ast.Ident {
  228									NamePos: /Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/client.go:128:10
  228									Name: "Client"
  228									Obj: *(obj @ 104)
  228								}
  228							}
  228						}
  228					}
  228					Closing: /Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/client.go:128:16
  228				}
  229				Name: *ast.Ident {
  229					NamePos: /Users/liamz/parser/src/github.com/liamzebedee/graphparse/subnet/subnet/client.go:128:18
  229					Name: "Run"
  229				}


*/