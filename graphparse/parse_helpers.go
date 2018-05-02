package graphparse

import (
	"go/ast"
	"go/types"
	"fmt"
	"strconv"
)

func getObjFromIdent(ident *ast.Ident) (types.Object, error) {
	// I wrote this in a subconcious spree of, "I have a gut feeling that this will do it"
	obj := pkginfo.ObjectOf(ident)

	// we want to have a canonical obj so we can make an id out of it
	// in this case, we use the def obj
	// obj := pkginfo.Defs[ident]
	
	// if obj == nil {
	// 	return nil, fmt.Errorf("obj not found for ident", ident)
	// }

	// Universe scope
	// if obj.Pkg() == nil {
	// 	return nil, fmt.Errorf("universe object ", obj.Type(), ident)
	// }

	if obj != nil {
		return obj, nil
	}

	return nil, fmt.Errorf("unexpected error", obj)
}

func exprToObj(expr ast.Expr) (types.Object, error) {
	switch x := expr.(type) {
	case *ast.SelectorExpr:
		if sel := pkginfo.Selections[x]; sel != nil {
			return sel.Obj(), nil
		}
		// Probably fully-qualified
		return pkginfo.ObjectOf(x.Sel), nil
	
	case *ast.Ident:
		obj, err := getObjFromIdent(x)
		return obj, err
	
	default:
		ParserLog.Printf("missed type %T\n", x)
	}

	return nil, fmt.Errorf("couldn't get object for expression:", expr)
}

func objIsWorthy(obj types.Object) bool {
	if obj.Pkg() == nil {
		return false
	}
	if obj.Pkg().Name() != thisPackage {
		return false
	}
	return true
}



func importToCanonicalKey(importSpec *ast.ImportSpec) string {
	importName, err := strconv.Unquote(importSpec.Path.Value)
	if err == nil {
		return importName
	} else {
		panic(err)
	}
}

/*
FieldVal, MethodVal - something.X()
MethodExpr - fn := something.X; // fn(something, yada)
*/
func parseXOfSelectorExpr(selX ast.Expr) string {
	switch y := selX.(type) {
	case *ast.SelectorExpr:
		return parseXOfSelectorExpr(y.X) + "." + y.Sel.Name
	case *ast.Ident:
		return y.Name
	default:
		ParserLog.Printf("didn't understand X of selector %T", y)
	}
	return ""
}

type annotatedSelectorObject struct {
	kind types.SelectionKind
	types.Object
}


func getObjectsFromSelector(sel ast.Expr) (objs []annotatedSelectorObject) {
	switch x := sel.(type) {
	case *ast.SelectorExpr:
		sel := pkginfo.Selections[x]
		
		if sel == nil {
			ParserLog.Printf("skipping selector expr (likely qualified identifier) %T\n", x)
			break
		} else {
			annSelObj := annotatedSelectorObject{
				sel.Kind(),
				sel.Obj(),
			}
			objs = append(objs, annSelObj)
		}


		if ident, ok := x.X.(*ast.Ident); ok {
			obj := pkginfo.ObjectOf(ident)
			if !objIsWorthy(obj) {
				ParserLog.Printf("encountered unexpected bad obj in selector expression %T - %T", obj, x)
			} else {
				annSelObj := annotatedSelectorObject{
					sel.Kind(),
					obj,
				}
				objs = append(objs, annSelObj)
			}
			return objs
		}

		return append(objs, getObjectsFromSelector(x.X)...)
		
	default:
		ParserLog.Printf("didn't understand X of selector %T", x)
	}
	return objs
}

func funcObjToNode(obj types.Object) *objNode {
	// variant := Func
	// switch x := obj.Type().(type) {
	// case *types.Signature:
	// 	if x.Recv() != nil {
	// 		variant = Method
	// 	}
	// default:
	// 	ParserLog.Printf("missed type %T\n", x)
	// }
	
	// funcNode := LookupOrCreateNode(obj, variant, obj.Name())
	funcNode := LookupNode(obj)
	return funcNode
}

func objToVariant(obj types.Object) NodeType {
	switch x := obj.Type().(type) {
	case *types.Signature:
		if x.Recv() != nil {
			return Method
		}
		return Func
	
	// // Recursive cases
	// case *types.Array:
	// case *types.Slice:
	// case *types.Chan:
	// case *types.Pointer:
	// 	return typeToObj(typ.Elem())

	// case *types.Struct:
	// 	return Struct
		
	default:
		return Field
	}
}