package main

type Node struct {
    Value string
    edges []Edge
}

type Edge struct {
    Parent, Child Node
}

func (n *Node) AddEdge(e Edge) {
    n.edges = append(n.edges, e)
    // TODO: Add error handling
}

func (n *Node) RemoveEdge(e Edge) {
    // TODO: Implement
}

