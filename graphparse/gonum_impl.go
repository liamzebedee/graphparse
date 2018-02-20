package graphparse

import (
	graphGn "gonum.org/v1/gonum/graph"
)

// Implementation of the Gonum Graph interface for graphparse.Graph


// Graph
// -----

// Has returns whether the node exists within the graph.
func (g *graph) Has(n graphGn.Node) bool {
	for _, e := range g.edges {
		// if n == graphGn.Node(e.from)  { return true }
		// if n == graphGn.Node(e.to)    { return true }
		if n.ID() == e.to.ID()   { return true }
		if n.ID() == e.from.ID() { return true }
	}
	return false
}

// Nodes returns all the nodes in the graph.
func (g *graph) Nodes() []graphGn.Node {
	unique := make(map[nodeid]Node)
	nodes := []graphGn.Node{}
	for _, e := range g.edges {
		unique[e.from.Id()] = e.from
		unique[e.to.Id()] = e.to
	}
	for _, n := range unique {
		nodes = append(nodes, n)
	}
	return nodes
}

// From returns all nodes that can be reached directly
// from the given node.
func (g *graph) From(n graphGn.Node) []graphGn.Node {
	nodes := []graphGn.Node{}
	for _, e := range g.edges {
		if e.to.ID() == n.ID() {
			nodes = append(nodes, e.to)
		}
		if e.from.ID() == n.ID() {
			nodes = append(nodes, e.from)
		}
	}
	return nodes
}

// HasEdgeBeteen returns whether an edge exists between
// nodes x and y without considering direction.
func (g *graph) HasEdgeBetween(x, y graphGn.Node) bool {
	return g.Edge(x, y) != nil
	// for _, e := range g.edges {
	// 	if e.from == x && e.to == y {
	// 		return true
	// 	}
	// 	if e.from == y && e.to == x {
	// 		return true
	// 	}
	// }
	// return false
}

// Edge returns the edge from u to v if such an edge
// exists and nil otherwise. The node v must be directly
// reachable from u as defined by the From method.
func (g *graph) Edge(u, v graphGn.Node) graphGn.Edge {
	for _, e := range g.edges {
		if e.from.ID() == u.ID() && 
		   e.to.ID()   == v.ID() {
			return e
		}
		// We ignore directness?
		if e.from.ID() == v.ID() && 
		   e.to.ID()   == u.ID() {
			return e
		}
	}
	return nil
}


// Node
// ----

func (n *objNode) ID() int64 {
	return int64(n.Id())
}

func (n *canonicalNode) ID() int64 {
	return int64(n.Id())
}


// Edge
// ----

func (e edge) From() graphGn.Node {
	return e.from
}

func (e edge) To() graphGn.Node {
	return e.to
}

func (e edge) Weight() float64 {
	return 1.
}