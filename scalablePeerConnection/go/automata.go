package graph

import (
    "github.com/guanyilun/go-sampling/samping"
)

type Automata struct {
    var limit int
    var probs []float64
    var active bool
    var sampling sample.Sampling // May not be necessary
    var counter int
    var delta int
}

func NewAutomata(actions, limit int) *Automata {
    var a Automata
    a.limit = limit
    a.probs = make([]float64, actions)
    a.active = true
    a.counter = 0
    a.sampling = sampling.NewSampling()
    a.sampling.AddBundleProbs(a.probs)
    a.delta = 100000 // A large number 
    return &a
}

func (a *Automata) Enum int, string {
    if a.counter < a.limit {
	a.counter++
	return a.sampling.Sample(), nil
    } else {
	a.active = false
	return 0, "ERROR - Enum limit is reached"
    }
}

func (a *Automata) ReEnum int, string {
    return a.sampling.Sample()
}

func (a *Automata) IsActive bool {
    return a.active
}
