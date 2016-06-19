package main

type Node struct {
    Value string
    edges []Edge
}

type Edge struct {
    Parent, Child Node
}

func (n *Node) AddEdge(e Edge) 

func (n *Node) RemoveEdge(e Edge)

