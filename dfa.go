package gocompute

import (
	"errors"
	"gopkg.in/fatih/set.v0"
	"reflect"
)

type DFA struct {
	states     *set.Set
	alphabet   *set.Set
	transition func(state, input string) (transition string)
	start      string
	accept     *set.Set
}

//preconditions:
// states, transitions, are finite and all elements have type string
// accept is subset of states, start is in states
// for every pair of states and alphabet characters, (state, char),
// transition(state, char) is in states

func NewDFA(states,
	alphabet *set.Set,
	transition func(state, input string) (nextState string),
	start string,
	accept *set.Set) (*DFA, error) {

	d := &DFA{states, alphabet, transition, start, accept}
	ans, err := d.CheckDFA()
	if ans != true && err != nil {
		return nil, err
	}
	return d, nil
}

//preconditions:
//d must be a valid DFA,
//each character of w must be in d.alphabet
func (d DFA) Simulate(w string) (bool, error) {
	ans, err := d.CheckDFA()
	if ans != true && err != nil {
		return false, errors.New("gocompute/dfa: invalid DFA")
	}

	currentState := d.start
	for _, r := range w {
		if !d.alphabet.Has(string(r)) {
			return false, errors.New("gocompute/dfa: string to test not in alphabet of DFA")
		}
		currentState = d.transition(currentState, string(r))
	}
	if d.accept.Has(currentState) {
		return true, nil
	}
	return false, nil
}

func (d DFA) CheckDFA() (bool, error) {
	//check that all elements of states are strings
	stateItems := d.states.List()
	for i := 0; i < d.states.Size(); i++ {
		if reflect.TypeOf(stateItems[i]).Kind() != reflect.String {
			return false, errors.New("gocompute/dfa: set of states contains non-string type")
		}
	}

	//check that all elements of alphabet are strings
	alphItems := d.alphabet.List()
	for i := 0; i < d.alphabet.Size(); i++ {
		if reflect.TypeOf(alphItems[i]).Kind() != reflect.String {
			return false, errors.New("gocompute/dfa: alphabet contains non-string type")
		}
	}

	//check that start state is a member of states
	if !d.states.Has(d.start) {
		return false, errors.New("gocompute/dfa: start state not in set of states")
	}

	//check that set of accept states is subset of states
	if !d.states.IsSubset(d.accept) {
		return false, errors.New("gocompute/dfa: set of accept states not a subset of set of all states")
	}

	//check that for every pair q,a with q in states and a in alphabet,
	//transition(q,a) is in states
	for i := 0; i < d.states.Size(); i++ {
		for j := 0; j < d.alphabet.Size(); j++ {
			if !d.states.Has(d.transition(stateItems[i].(string), alphItems[j].(string))) {
				return false, errors.New("gocompute/dfa: incomplete or invalid transition function")
			}
		}
	}
	return true, nil
}
