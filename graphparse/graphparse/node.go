package graphparse

import (
	"go/types"
	"go/token"
	"fmt"
	// "errors"
	// "github.com/jimlawless/whereami"
)

type NodeType int
const (
	Struct NodeType = iota
	Method
	Func
	Field
	RootPackage
	File
	ImportedPackage
	ImportedFunc
	FuncCall
	Template
)
var nodeTypes = []string{
	"Struct",
	"Method",
	"Func",
	"Field",
	"RootPackage",
	"File",
	"ImportedPackage",
	"ImportedFunc",
	"FuncCall",
	"___template___",
}

type baseNode struct {
	variant NodeType
	label string
}

type Node interface {
	Id() nodeid
	ID() int64
	Label() string
	Variant() NodeType
	String() string
	PosInfo() posInfo
}

type posInfo struct {
	Filename string // filename, if any
	Offset   int    // offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (byte count)
}

type objNode struct {
	obj types.Object
	baseNode
}

func (n *objNode) PosInfo() posInfo {
	return posInfo{}
}

/*

	func (t *Basic) Underlying() Type     { return t }
	func (t *Array) Underlying() Type     { return t }
	func (t *Slice) Underlying() Type     { return t }
	func (t *Struct) Underlying() Type    { return t }
	func (t *Pointer) Underlying() Type   { return t }
	func (t *Tuple) Underlying() Type     { return t }
	func (t *Signature) Underlying() Type { return t }
	func (t *Interface) Underlying() Type { return t }
	func (t *Map) Underlying() Type       { return t }
	func (t *Chan) Underlying() Type      { return t }
	func (t *Named) Underlying() Type     { return t.underlying }

*/

func typeIsNamed(typ types.Type) bool {
	switch typ := typ.(type) {
	// End recursive case.
	case *types.Basic:
		return false
	case *types.Named:
		return true
	
	// Recursive cases
	case *types.Array:
	case *types.Slice:
	case *types.Chan:
	case *types.Pointer:
		return typeIsNamed(typ.Elem())
	
	// Error cases
	case *types.Struct:
	case *types.Tuple:
	case *types.Signature:
		return false
		// fmt.Println(typ)
		// panic("this type doesn't have underlying type, we shouldn't encounter it")
	}
	return false
}

func typeToObj(typ types.Type) (types.Object) {
	switch typ := typ.(type) {
	// End recursive case.
	case *types.Basic:
		// panic("basic type has no obj")
		return nil
	case *types.Named:
		return typ.Obj()

	// Recursive cases
	case *types.Array:
	case *types.Slice:
	case *types.Chan:
	case *types.Pointer:
		return typeToObj(typ.Elem())
	
	// Error cases
	case *types.Struct:
	case *types.Tuple:
	case *types.Signature:
		panic("types have no obj")
		// fmt.Println(typ)
		// panic("this type doesn't have underlying type, we shouldn't encounter it")
	}
	return nil
}

// TODO refactor
func objToId(obj types.Object) nodeid {
	return nodeid(obj.Pos())
}

// The Id of an object node is defined canonically 
// as the token.Pos of where the type is declared
func (n *objNode) Id() nodeid {
	return objToId(n.obj)
}
func (n *objNode) Label() string {
	return n.baseNode.label
}
func (n *objNode) Variant() NodeType {
	return n.baseNode.variant
}
func (n *objNode) String() string {
	return fmt.Sprintf("%s <%s>", n.Label(), nodeTypes[n.Variant()])
}

// func (n *objNode) Pos() posInfo {
// 	pos := n.obj.Pos()
// 	return fset.Position(pos).String()
// }




func (g *Graph) LookupNode(obj types.Object) *objNode {
	if obj == nil {
		panic("obj must be non-nil")
	}

	id := objToId(obj)
	
	if node, ok := g.nodeLookup[id]; ok {
		return node.(*objNode)
	} else if node, ok := g.objNodeLookup[id]; ok {
		return node
	} else {
		// create template of node for later.
		template := &objNode{
			obj,
			baseNode{Template, obj.Name()},
		}
		g.objNodeLookup[id] = template
		return template
	}
}

func (g *Graph) CreateNode(obj types.Object, variant NodeType, label string, pos token.Position) *objNode {
	if obj == nil {
		panic("obj must be non-nil")
	}
	if label == "" {
		panic("label must be something")
	}

	id := objToId(obj)

	if node, ok := g.nodeLookup[id]; ok {
		ParserLog.Fatalln("already created: ", node)
		return node.(*objNode)
	}

	node := g.objNodeLookup[id]
	if node == nil {
		node = &objNode{
			obj,
			baseNode{},
		}
	} else {
		delete(g.objNodeLookup, id)
	}

	node.baseNode = baseNode{variant, label}

	// check for template node (uncompleted references)
	// template, _ := objNodeLookup[id]
	// if template != nil {
	// 	*template = *node
	// 	delete(objNodeLookup, id)
	// }

	ParserLog.Println("created node: ", node.String())
	g.addNodeToLookup(node)
	
	return node
}

var canonicalNodeLookup = make(map[string]*canonicalNode)

type canonicalNode struct {
	baseNode
}

func (n *canonicalNode) Id() nodeid {
	return pointerToId(n)
}
func (n *canonicalNode) Label() string {
	return n.baseNode.label
}
func (n *canonicalNode) Variant() NodeType {
	return n.baseNode.variant
}
func (n *canonicalNode) String() string {
	return fmt.Sprintf("%s <%s>", n.Label(), nodeTypes[n.Variant()])
}
func (n *canonicalNode) DebugInfo() string {
	return ""
}

func (n *canonicalNode) PosInfo() posInfo {
	return posInfo{}
}

func (g *Graph) LookupOrCreateCanonicalNode(key string, variant NodeType, label string) *canonicalNode {
	node, ok := canonicalNodeLookup[key]

	if !ok {
		node = &canonicalNode{
			baseNode{variant, label},
		}
		canonicalNodeLookup[key] = node
		g.addNodeToLookup(node)
	}

	return node
} 

type nodeid int64