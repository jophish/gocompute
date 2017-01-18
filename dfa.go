// gocompute is a package for experimenting with structures and automata from the field of computability theory.
// This package contains (or, will contain) implementations of DFAs, NFAs, CFLs, PDAs, Turing machines, and the like,
// as well as ways to simulate their operation on inputs and manipulate the machines themselves. Some features include
// generating new automaton using the closure properties of certain classes of languages.
package gocompute

import (
	"errors"
	"github.com/jophish/golang-set"
	"reflect"
)

// Internal representation of a DFA
type DFA struct {
	states     mapset.Set
	alphabet   mapset.Set
	transition func(state interface{}, input string) (transition interface{})
	start      interface{}
	accept     mapset.Set
	simulate   func(w string, d *DFA) (output bool, err error)
}

//preconditions:
// states, transitions, are finite and all elements have same type
// accept is subset of states, start is in states
// for every pair of states and alphabet characters, (state, char),
// transition(state, char) is in state

// Constructor method for creating a DFA. Takes as input a set of states of the same type, a set of strings representing the
// input alphabet, a pointer to a transition function mapping states and strings to states, a start state, and a set of accept
// states. Returns a pointer to the newly created DFA and an error, which is non-nil if the input was improperly formatted.
//
// States can be any type, but they must all be the same type. The start state must be a member of the set of all states, and
// the set of accept states must be a subset of the set of all states. The transition function must return a valid state for
// each possible input combination of state and alphabet symbol.
func NewDFA(states,
	alphabet mapset.Set,
	transition func(state interface{}, input string) (nextState interface{}),
	start interface{},
	accept mapset.Set) (*DFA, error) {

	d := &DFA{states, alphabet, transition, start, accept, nil}
	simulate := simDFA
	d.simulate = simulate

	ans, err := d.CheckDFA()
	if ans != true && err != nil {
		return nil, err
	}
	return d, nil
}

//all the closure methods require that d1 and d2 share the same alphabet. this doesn't have to be, but it is the current implementation

//given two DFAs d1, d2, returns pointer to new DFA which recognizes the
//language union(L(d1), L(d2)). does not modify d1, d2.

// Given DFAs d1 and d2, which recognize languages L(d1) and L(d2) respectively, d1.Union(d2) returns
// a pointer to a new DFA which recognizes both L(d1) and L(d2).
func (d1 DFA) Union(d2 *DFA) (*DFA, error) {
	if !d1.alphabet.Equal(d2.alphabet) {
		return nil, errors.New("gocompute/dfa: alphabets of input DFAs must be equal")
	}
	//states are represented as ordered pairs
	states := d1.states.CartesianProduct(d2.states)
	alphabet := d1.alphabet

	transition := func(state interface{}, input string) (nextState interface{}) {
		st := state.(mapset.OrderedPair)
		first := d1.transition(st.First.(interface{}), input)
		second := d2.transition(st.Second.(interface{}), input)
		return mapset.OrderedPair{first, second}
	}

	start := mapset.OrderedPair{d1.start, d2.start}
	accept := (d1.accept.CartesianProduct(d2.states)).Union(d1.states.CartesianProduct(d2.accept))
	d3 := &DFA{states, alphabet, transition, start, accept, nil}
	simulate := simDFA
	d3.simulate = simulate
	return d3, nil

}

//make sure d1.alphabet == d2.alphabet

// Given DFAs d1 and d2, which recognize languages L(d1) and L(d2) respectively, d1.Intersection(d2) returns
// a pointer to a new DFA which recognizes the intersection of L(d1) and L(d2).
func (d1 DFA) Intersection(d2 *DFA) (*DFA, error) {
	if !d1.alphabet.Equal(d2.alphabet) {
		return nil, errors.New("gocompute/dfa: alphabets of input DFAs must be equal")
	}
	ans, err := d1.CheckDFA()
	if ans == false && err != nil {
		return nil, errors.New("gocompute/dfa: invalid DFA: " + err.Error())
	}
	ans, err = d2.CheckDFA()
	if ans == false && err != nil {
		return nil, errors.New("gocompute/dfa: invalid DFA: " + err.Error())
	}
	//states are represented as ordered pairs
	states := d1.states.CartesianProduct(d2.states)
	alphabet := d1.alphabet
	transition := func(state interface{}, input string) (nextState interface{}) {
		st := state.(mapset.OrderedPair)
		first := d1.transition(st.First.(interface{}), input)
		second := d2.transition(st.Second.(interface{}), input)
		return mapset.OrderedPair{first, second}
	}

	start := mapset.OrderedPair{d1.start, d2.start}
	accept := d1.accept.CartesianProduct(d2.accept)
	d3 := &DFA{states, alphabet, transition, start, accept, nil}
	simulate := simDFA
	d3.simulate = simulate
	return d3, nil

}

// Given a DFA d1 with alphabet E, recognizing language L(d1), d1.Complement() returns a pointer to a new DFA recognizing
// the set E* - L(d1), where E* is the set of all strings that can be created from symbols in the alphabet E.
func (d1 DFA) Complement() (*DFA, error) {
	if !d1.alphabet.Equal(d2.alphabet) {
		return nil, errors.New("gocompute/dfa: alphabets of input DFAs must be equal")
	}
	ans, err := d1.CheckDFA()
	if ans == false && err != nil {
		return nil, errors.New("gocompute/dfa: invalid DFA: " + err.Error())
	}
	states := d1.states
	alphabet := d1.alphabet
	transition := d1.transition
	start := d1.start
	accept := d1.states.Difference(d1.accept)
	d2 := &DFA{states, alphabet, transition, start, accept, simDFA}
	return d2, nil
}

//returns DFA recognizing L(d1)\L(d2) = L(d1) - L(d2) = L(d1) intersect L(d2)^c
//this closure property is sometimes known as relative complement

// Given DFAs d1 and d2, which recognize languages L(d1) and L(d2) respectively, d1.Difference(d2) returns a pointer
// to a new DFA which recognizes the language L(d1) - L(d2), or the set of all strings recognized by d1 but not d2.
func (d1 DFA) Difference(d2 *DFA) (*DFA, error) {
	if !d1.alphabet.Equal(d2.alphabet) {
		return nil, errors.New("gocompute/dfa: alphabets of input DFAs must be equal")
	}
	ans, err := d1.CheckDFA()
	if ans == false && err != nil {
		return nil, errors.New("gocompute/dfa: invalid DFA: " + err.Error())
	}
	ans, err = d2.CheckDFA()
	if ans == false && err != nil {
		return nil, errors.New("gocompute/dfa: invalid DFA: " + err.Error())
	}
	d2c, err := d2.Complement()
	if err != nil {
		return nil, errors.New("gocompute/dfa: invalid DFA: " + err.Error())
	}
	d3, err := d1.Intersection(d2c)
	if err != nil {
		return nil, errors.New("gocompute/dfa: invalid DFA: " + err.Error())
	}
	return d3, nil
}

//note: the following transformation functions generally require the use of
//the equivalence of NFAs and DFAs, so they might be tough.

//this one's tricky. if x in L(d1), y in L(d2), then d3
//recognizes the language with elements of the form xy

// Given DFAs d1 and d2, which recognize languages L(d1) and L(d2) respectively, d1.Concatenation(d2) returns a pointer
// to a new DFA which recognizes strings of the form xy, where x is in L(d1) and y is in L(d2). Unimplemented.
func (d1 DFA) Concatenation(d2 *DFA) (*DFA, error) {
	return nil, nil
}

//recognizes language of strings that are concatenations of strings in
//L(d1)

// Given a DFA d1 which recognizes the language L(d1), d1.Star() returns a pointer
// to a new DFA which recognizes strings which are the concatenation of any number of strings in L(d1). Unimplemented.
func (d1 DFA) Star() (*DFA, error) {
	return nil, nil
}

//outputs a DFA such that if x in L(d1), the DFA recognizes reversed(x)

// Given a DFA d1 recognizing language L(d1), returns a pointer to a new DFA which recognizes strings which are reversed
// versions of those in L(d1). Unimplemented.
func (d1 DFA) Reverse() (*DFA, error) {
	return nil, nil
}

//we use this abstraction of the simulate funciton in order to
//construct union, intersection, etc. of DFAs more easily

// Given a DFA d and a string w, d.Simulate(w) simulates the automaton operation of d on w, and returns true if
// d recognizes w, that is, if the DFA ends in an accept state. If the DFA does not end in an accept state, we return false,
// to represent the DFA not recognizing w.
func (d DFA) Simulate(w string) (bool, error) {
	return d.simulate(w, &d)
}

//preconditions:
//d must be a valid DFA,
//each character of w must be in d.alphabet
func simDFA(w string, d *DFA) (bool, error) {
	ans, err := d.CheckDFA()
	if ans != true && err != nil {
		return false, errors.New("gocompute/dfa: invalid DFA: " + err.Error())

	}

	currentState := d.start
	for _, r := range w {
		if !d.alphabet.Contains(string(r)) {
			return false, errors.New("gocompute/dfa: string to test not in alphabet of DFA")
		}
		currentState = d.transition(currentState, string(r))
	}
	if d.accept.Contains(currentState) {
		return true, nil
	}
	return false, nil
}
// Checks to make sure a given DFA d is properly formatted with correct input data.
func (d DFA) CheckDFA() (bool, error) {
	//check that all elements of states are of same type (interface{})
	stateItems := d.states.Iter()
	var t reflect.Type
	for elem := range stateItems {
		if t == nil {
			t = reflect.TypeOf(elem)
		}
		if reflect.TypeOf(elem) != t {
			return false, errors.New("gocompute/dfa: set of states contains non-string type")
		}
	}

	//check that all elements of alphabet are strings
	alphItems := d.alphabet.Iter()
	for elem := range alphItems {
		if reflect.TypeOf(elem).Kind() != reflect.String {
			return false, errors.New("gocompute/dfa: alphabet contains non-string type")
		}
	}

	//fmt.Println(reflect.TypeOf(d.start))
	//check that start state is a member of states
	if !d.states.Contains(d.start) {
		return false, errors.New("gocompute/dfa: start state not in set of states")
	}

	//check that set of accept states is subset of states
	if !d.states.IsSuperset(d.accept) {
		return false, errors.New("gocompute/dfa: set of accept states not a subset of set of all states")
	}

	//check that for every pair q,a with q in states and a in alphabet,
	//transition(q,a) is in states
	for stateElem := range stateItems {
		for alphElem := range alphItems {
			if !d.states.Contains(d.transition(stateElem.(interface{}), alphElem.(string))) {
				return false, errors.New("gocompute/dfa: incomplete or invalid transition function")
			}
		}
	}
	return true, nil
}
