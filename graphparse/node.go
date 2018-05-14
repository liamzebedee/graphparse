package graphparse

import (
	"go/types"
	"fmt"
	// "errors"
	"github.com/jimlawless/whereami"
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
	DebugInfo() string
}

type objNode struct {
	obj types.Object
	baseNode
}

type objlookup struct {
	id string
}
var objLookups = make(map[string]*objlookup)

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
	panic("can't get id")
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

func (n *objNode) DebugInfo() string {
	pos := n.obj.Pos()
	return fset.Position(pos).String()
}


var objNodeLookup = make(map[nodeid]*objNode)

func LookupNode(obj types.Object) *objNode {
	if obj == nil {
		panic("obj must be non-nil")
	}

	id := objToId(obj)
	
	if node, ok := nodeLookup[id]; ok {
		return node.(*objNode)
	} else if node, ok := objNodeLookup[id]; ok {
		return node
	} else {
		// create template of node for later.
		template := &objNode{
			obj,
			baseNode{Template, obj.Name()},
		}
		objNodeLookup[id] = template
		return template
	}
}

func CreateNode(obj types.Object, variant NodeType, label string) *objNode {
	if obj == nil {
		panic("obj must be non-nil")
	}
	if label == "" {
		panic("label must be something")
	}

	id := objToId(obj)

	if node, ok := nodeLookup[id]; ok {
		ParserLog.Fatalln("already created: ", node)
		return node.(*objNode)
	}

	node := objNodeLookup[id]
	if node == nil {
		node = &objNode{
			obj,
			baseNode{},
		}
	} else {
		delete(objNodeLookup, id)
	}

	node.baseNode = baseNode{variant, label}

	// check for template node (uncompleted references)
	// template, _ := objNodeLookup[id]
	// if template != nil {
	// 	*template = *node
	// 	delete(objNodeLookup, id)
	// }

	ParserLog.Println("created node: ", node.String())
	addNodeToLookup(node)
	
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

func LookupOrCreateCanonicalNode(key string, variant NodeType, label string) *canonicalNode {
	node, ok := canonicalNodeLookup[key]

	if !ok {
		node = &canonicalNode{
			baseNode{variant, label},
		}
		canonicalNodeLookup[key] = node
		addNodeToLookup(node)
	}

	return node
} 

type nodeid int64