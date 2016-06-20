package main

import "fmt"

func main () {
    graph := NewGraph()
    graph.AddNode("a")
    graph.AddNode("b")
    graph.AddNode("c")
    graph.AddUniEdge("a", "b", 1)
    graph.AddBiEdge("a", "c", 1)
    
    // Test RemoveEdge()
    a := graph.GetNode("a")
    b := graph.GetNode("b")
    c := graph.GetNode("c")
    ae := a.GetEdges()
    be := b.GetEdges()
    ce := c.GetEdges()
    fmt.Println(len(ae)) // 2
    fmt.Println(len(be)) // 0
    fmt.Println(len(ce)) // 1
    
    graph.RemoveUniEdge("a", "c")
    ae = a.GetEdges()
    be = b.GetEdges()
    ce = c.GetEdges()
    fmt.Println(len(ae)) // 1
    fmt.Println(len(be)) // 0
    fmt.Println(len(ce)) // 1
    
    e := graph.GetEdge("c", "a")
    fmt.Println(e.Parent.Value) // c
    fmt.Println(e.Child.Value) // a
    
    graph.removeEdge(e)
    ae = a.GetEdges()
    be = b.GetEdges()
    ce = c.GetEdges()
    fmt.Println(len(ae)) // 1
    fmt.Println(len(be)) // 0
    fmt.Println(len(ce)) // 0
    
    graph.AddNode("d")
    graph.AddNode("e")    
    fmt.Println(graph.GetTotalNodes()) // 5 
    
    graph.AddUniEdge("d", "e", 1)
    d := graph.GetNode("d")
    de := d.GetEdges()
    fmt.Println(len(de)) // 1
    fmt.Println(graph.HasUniEdge("d", "e")) // true
    fmt.Println(graph.HasBiEdge("d", "e")) // false
    
    graph.AddUniEdge("c", "d", 1)
    graph.AddUniEdge("b", "d", 1)
    graph.AddUniEdge("d", "a", 1)
    edges := graph.GetInEdges("d")
    for _, e := range edges {
	fmt.Println(e.Parent.Value) // c, b
    }
    
    edges = graph.GetOutEdges("d")
    for _, e := range edges {
	fmt.Println(e.Child.Value) // e, a
    }
    
    graph.RemoveNode("d")
    fmt.Println(graph.GetTotalNodes()) // 4 
    fmt.Println(graph.HasUniEdge("d", "e")) // false 
    fmt.Println(graph.HasUniEdge("c", "d")) // false
    fmt.Println(graph.HasUniEdge("b", "d")) // false
    fmt.Println(graph.HasUniEdge("d", "a")) // false ERR
    
    
} 