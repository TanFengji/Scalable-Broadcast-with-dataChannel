package main

import (
    "fmt"
    "github.com/gonum/floats" // TODO: Use go get to obtain the repository
)

// Graph struct represents a graph that is either undirectional or directional.
// It may or may not have a head
type Graph struct {
    edges map[string]map[string]Edge
    nodes map[string]Node
    head Node
}

// NewGraph creates an empty new graph
func NewGraph() (*Graph) {
    var g Graph
    g.edges = make(map[string]map[string]Edge)
    g.nodes = make(map[string]Node)
    return &g
}

// AddNode adds a node to a graph and refer to the node with a unique name
func (g *Graph) AddNode(s string) {
    n := Node{Value: s}
    if _, exist := g.nodes[s]; !exist {
	g.nodes[s] = n
    } 
    //TODO: Add error handling
}

// GetNode returns a the node with associated string s. The second argument
// denotes the error and becomes nil if things go correctly`
func (g *Graph) GetNode(s string) Node, string {
    //NOTE: Make sure the node returned isn't a copy but the actual node
    if n, exist := g.nodes[s]; exist {
	return n, nil
    }
    return Node{}, "Error getting node"
}

// RemoveNode removes a node associated with string s from the graph and 
// also removes the edges both linking in it and linking out of it`
func (g *Graph) RemoveNode(s string) {
    if n, exist := g.nodes[s]; exist {
	delete(g.nodes, s)
	
	//Remove edges linked to it
	es := n.GetInEdges(s)
	for _, e := range es {
	    g.removeEdge(e)
	}
	
	//Remove edges linking out of it
	es := n.GetOutEdge(p)
	for _, e := range es {
	    g.removeEdge(e)
	}
    } 
    //TODO: Add error handling
}

func (g *Graph) GetChildren(s string) []Node {
    // Assuming s node exists
    children := make([]Node, 0)
    if n, exist := g.nodes[s]; exist {
	val := n.Value
	for _, e := range n.edges {
	    children = append(children, e.Child)
	}
    }
    return children
}

func (g *Graph) GetParent(s string) []Node {
    parents := make([]Node, 0)
    edges := g.GetInEdges(s)
    for _, e := range edges {
	parents = append(parents, e.Parent)
    }
    return parents
}

func (g *Graph) GetInEdges(c string) []Edge {
    edges := make([]Edge, 0)
    for _, n := range g.nodes {
	p := n.Value
	if g.HasUniEdge(p, c) {
	    e := g.GetUniEdge(p, c) 
	    edges = append(edges, e) 
	}
    }
}

func (g *Graph) GetOutEdges(p string) []Edge {
    return g.nodes[p].edges
}

func (g *Graph) HasUniEdge(parent, child string) bool {
    _, exist := g.edges[parent][child]
    return exist
}

func (g *Graph) HasBiEdge(parent, child string) bool {
    exist1 := g.HasUniEdge(parent, child)
    exist2 := g.HasUniEdge(child, parent)
    return exist1 && exist2
}

func (g *Graph) GetEdge(parent, child string) {
    e := g.edges[parent][child]
    return e
}

func (g *Graph) SetHead(s string) {
    if n, exist := g.nodes[s]; exist {
	g.head = n
    }
    //TODO: Add error handling
}

func (g *Graph) AddUniEdge(parent, child string, weight int) {
    p := g.GetNode(parent) //NOTE: Make sure it allows changing original value
    c := g.GetNode(Child) 
    //TODO: Add error handling
    ef := Edge{Parent: p, Child: c, Weight: weight}
    p.AddEdge(ef)
    g.edges[parent][child] = ef
}

// AddBiEdge adds a bidirectional edge between parent and child with the 
// same weight associated with it. 
func (g *Graph) AddBiEdge(parent, child string, weight int) {
    g.AddUniEdge(parent, child, weight)
    g.AddUniEdge(child, parent, weight)
}

func (g *Graph) RemoveUniEdge(parent, child string) {
    e := g.edges[parent][child]
    delete(g.edges[parent], child)
    
    //TODO: Use a better encapsulation, now removeEdge (private) is calling
    //RemoveUniEdge (public)
}

func (g *Graph) removeEdge(e Edge) {
    parent := e.Parent
    child := e.Child
    g.RemoveUniEdge(parent, child)
    //TODO: Use a better encapsulation, now removeEdge (private) is calling
    //RemoveUniEdge (public)
    
}

// GetDCMST function finds a Degree-Constrained Maximum Spanning Tree for a 
// given graph and return it as a subgraph. The current implementation refers
// to the learning automata method (JA 2013)
func (g *Graph) GetDCMST(deg int) {
    // Starting by making an empty map from node to automata 
    var wt float64 = 0
    var nc int = 0
    autos := make(map[string]Automata)
    
    
    // Generate an isomorphic graph to g
    for k := range g.nodes {
	nc := len(g.GetChildren(k))
	autos[k] = NewAutomata(k, deg)
    }
    
    // Start to generate a MST
    mst := NewGraph()
    
    // Start from the head -> assuming head exists
    head := g.head.Value
    
    mst.AddNode(head)
    mst.SetHead(head)
    parent := head
    
    for {
	// Loop begins
	a := autos[parent]
	
	if a.IsActive() {
	    children := g.GetChildren(parent)
	    
	    // check if all of its children is inactive, if all of the children
	    // automata are inactive then hasActive flag is false, otherwise true
	    var hasActive bool = false
	    for i, n := range children {
		if !hasActive {
		    hasActive = autos[n.Value].IsActive()
		} else {
		    break //NOTE: Make sure break doesn't refer to the if statement
		}
	    }
	    
	    // If automata has active children, then start to enumerate
	    if hasActive {
		i := a.Enum()
		child := children[i].Value
		
		// Re-enumerate until an active node is found
		// May need a more efficient algorithm
		for !autos[child].IsActive() {
		    i = a.ReEnum()
		    child := children[i].Value
		}
		
		e := g.edges[parent][child]
		
		if e.weight < a.delta {
		    a.delta = e.weight
		    a.Reward(i)
		} else {
		    a.Penalize(i)
		}
		
		mst.AddNode(child)
		mst.AddUniEdge(parent, child, e.weight)
		
	    } else { // It means that a node is either at the tail or there is no
		     // active children
		if mst.GetTotalNodes() == g.GetTotalNodes() {
		    // It means that all the nodes are covered and we have found
		    // a mst
		} else {
		    // It means that we have run out of choices but there are more 
		    // nodes to include
		}
	    }
	}
    }
}