package gocompute

import (
	"github.com/jophish/golang-set"
)
// Internal representation of a PDA
type PDA struct {
	states     *mapset.Set
	slphabet   *mapset.Set
	transition func(state, input, stackSymbol string) (nextStates *mapset.Set)
	start      string
	stackStart string
	accept     *mapset.Set
}

// Constructor for a new PDA.
func newPDA(states,
	alphabet *mapset.Set,
	transition func(state, input, stackSymbol string) (nextStates *mapset.Set),
	start,
	stackStart string,
	accept *mapset.Set) *PDA {

	return &PDA{states, alphabet, transition, start, stackStart, accept}
}

// Simulates a PDA. Unimplemented.
func (p PDA) Simulate(w string) bool {
	//each configuration (state, input, stack head) is a node in a graph, with directed edges
	//going out as specified by the transition function. do a BFS, and return TRUE if any nodes
	//reachable from the initial configuration are ones representing configs containing an
	//accept state. if no such state is found, return FALSE
	/*	currentState := p.start
		stackHead := p.stackStart
		for _, r := range w {
			currentState = p.transition(currentState, string(r), stackHead)
		}
	*/
	return true
}
