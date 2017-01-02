package gocompute

import (
	"gopkg.in/fatih/set.v0"
)

type PDA struct {
	states     *set.Set
	slphabet   *set.Set
	transition func(state, input, stackSymbol string) (nextStates *set.Set)
	start      string
	stackStart string
	accept     *set.Set
}

func newPDA(states,
	alphabet *set.Set,
	transition func(state, input, stackSymbol string) (nextStates *set.Set),
	start,
	stackStart string,
	accept *set.Set) *PDA {

	return &PDA{states, alphabet, transition, start, stackStart, accept}
}

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
