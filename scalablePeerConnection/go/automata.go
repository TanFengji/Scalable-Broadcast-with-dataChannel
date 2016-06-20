package graph

import (
    "github.com/guanyilun/go-sampling/samping"
)

type Automata struct {
    var limit int
    var actions int
    var probs []float64
    var active bool
    var sampling sample.Sampling // May not be necessary
    var counter int
    var delta int
    var reward float64
    var penalize float64
    var threshold float64
}

func NewAutomata(actions, limit int) *Automata {
    var a Automata
    a.limit = limit
    a.probs = make([]float64, actions)
    
    // Initialize probabilities for automata
    for i := range a.probs {
	a.probs[i] = 1/actions
    }
    a.active = true
    a.counter = 0
    a.actions = actions
    a.sampling = sampling.NewSampling()
    a.sampling.AddBundleProbs(a.probs)
    a.delta = 100000 // A large number 
    
    // By default we adopt L_R-I model following JA (2013)
    a.reward = 0.09
    a.penalize = 0
    a.threshold = 0.9
    return &a
}

func (a *Automata) Enum() int, string {
    if a.counter < a.limit {
	a.counter++
	return a.sampling.Sample(), nil
    } else {
	a.active = false
	return 0, "ERROR - Enum limit is reached"
    }
}

func (a *Automata) ReEnum() int, string {
    return a.sampling.Sample()
}

func (a *Automata) IsActive() bool {
    return a.active
}

func (a *Automata) Reward(j int) {
    // Assuming learning reward-penalty (L_R-I) algorithm
    var sum float64 = 0
    r := a.reward
    for i := range a.probs {
	if i == j {
	    a.probs[i] = a.probs[i] + r * (1 - a.probs[i])
	} else {
	    a.probs[i] = (1 - r) * a.probs[i] 
	}
	sum += a.probs[i]
    }
    
    // Normalize the probabilities after modifying
    a.Normalize()
}

func (a *Automata) Penalize(j int) {
    // Assuming learning reward-penalty (L_R-I) algorithm, r = 0, and 
    // it can be seen that that Penalize() does nothing in such case,
    // hence we will comment this part unless required in the future.
    
    /*
    r := a.penalize
    for i := range a.probs {
	if i == j {
	    a.probs[i] = (1 - r) * a.probs[i] 
	} else {
	    a.probs[i] = r / (a.actions - 1) + (1 - r) * a.probs[i] 
	}
    }
    // Normalize the probabilities after modifying
    a.Normalize()
    */
}


func (a *Automata) Normalize() {
    var norm float64 = 0
    for _, v := range a.probs {
	norm += v
    }
    
    for i := range a.probs {
	a.probs[i] = a.probs[i] / norm
    }
}