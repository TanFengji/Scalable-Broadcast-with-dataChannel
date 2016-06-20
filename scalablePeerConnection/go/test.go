package main

import (
    "fmt"
    "github.com/guanyilun/go-sampling/sampling"
)

func main () {
    s := sampling.NewSampling() 
    s.AddValue(1, 1)
    s.AddValue(0, 1)
    s.Normalize()
    
    total := 0
    for i:=0; i<1000; i++ {
	val := s.Sample()
	total += val
    }
    fmt.Println(total)
}