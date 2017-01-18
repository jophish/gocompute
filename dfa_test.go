package gocompute

import (
	//"errors"
	"github.com/jophish/golang-set"
	"strconv"
	"testing"
)

//things to test:
// constructing a dfa, catching errors:
//    -states of arbitrary type (string, int, struct)
//    - error catching:
//       - states, alphabet have elements of different type
//       - accept not in states
//       - transition invalid
//
// simulating a dfa:
//    -

//test constructor, simulate, union, intersection, complement, difference

var d1, d1err = makeEvenOnesDFAStringStates()
var d2, d2err = makeEvenOnesDFAIntStates()
var d3, d3err = makeEvenOnesDFAStructStates()
var d4, d4err = makeOddZerosDFAStringStates()
var d5, d5err = makeOddOnesDFAStringStates()
var d6, d6err = makeEvenZerosDFAStringStates()
var d7, d7err = makeAllZerosOnesStringStates()
var d8, d8err = makeNoneZerosOnesStringStates()

var simulateTests = []struct {
	d           *DFA
	err         error
	testStrings map[string]bool
	descriptor  string
}{
	{d1, d1err, map[string]bool{"0011": true, "": true, "0001011011": false, "100": false, "001001011001001011": true}, "DFA accepting strings with even number of 1s using strings for states"},
	{d2, d2err, map[string]bool{"0011": true, "": true, "0001011011": false, "100": false, "001001011001001011": true}, "DFA accepting strings with even number of 1s using ints for states"},
	{d3, d3err, map[string]bool{"0011": true, "": true, "0001011011": false, "100": false, "001001011001001011": true}, "DFA accepting strings with even number of 1s using structs for states"},
	{d4, d4err, map[string]bool{"00011": true, "": false, "0001011011": true, "100": false, "001001011001001011": false}, "DFA accepting strings with odd number of 0s using strings for states"},
	{d5, d5err, map[string]bool{"0001": true, "": false, "0001011011": true, "100": true, "001001011001001011": false}, "DFA accepting strings with odd number of 1s using strings for states"},
	{d6, d6err, map[string]bool{"0001": false, "": true, "0001011011": false, "100": true, "001001011001001011": true}, "DFA accepting strings with even number of 0s using strings for states"},
	{d7, d7err, map[string]bool{"0001": true, "": true, "0001011011": true, "100": true, "001001011001001011": true}, "DFA accepting all strings of 0s and 1s using strings for states"},
	{d8, d8err, map[string]bool{"0001": false, "": false, "0001011011": false, "100": false, "001001011001001011": false}, "DFA accepting no strings using strings for states"},
}

var uniond1, uniond1err = d1.Union(d6) // should accept strings which contain even number of 0s or 1s
var uniond2, uniond2err = d4.Union(d5) //should accept strings which contain odd number of 0s or 1s

//for union, check that union of same dfa is same as original
//maybe test union between different representations
var unionTests = []struct {
	d1          *DFA
	d1err       error
	d2          *DFA
	d2err       error
	testStrings map[string]bool
	descriptor  string
}{
	{d1, d1err, d1, d1err, map[string]bool{"0011": true, "": true, "0001011011": false, "100": false, "001001011001001011": true}, "Union of two DFAs accepting even number of 1s, both using strings as states"},
	{d1, d1err, d3, d3err, map[string]bool{"0011": true, "": true, "0001011011": false, "100": false, "001001011001001011": true}, "Union of two DFAs accepting even number of 1s, one using strings, other using structs as states"},
	{d1, d1err, d6, d6err, map[string]bool{"011": true, "": true, "0001011011": false, "10": false, "000101101": true}, "Union of DFAs, one accepting even number 1s, other accepting even number 0s"},
	{d4, d4err, d5, d5err, map[string]bool{"0011": false, "": false, "0001011011": true, "100": true, "001001011001001011": false}, "Union of DFAs, one accepting odd number 1s, other accepting odd number 0s"},
	{d1, d1err, d5, d5err, map[string]bool{"0011": true, "": true, "0001011011": true, "100": true, "001001011001001011": true}, "Union of DFAs, one accepting even number 1s, other accepting odd number 1s"},
	{uniond1, uniond1err, uniond2, uniond2err, map[string]bool{"0011": true, "": true, "0001011011": true, "100": true, "001001011001001011": true}, "Union of union of DFAs, should accept all strings of 0s and 1s"},
	{compld1, compld1err, d1, d1err, map[string]bool{"0011": true, "": true, "0001011011": true, "100": true, "001001011001001011": true}, "Union of a DFA accepting strings of even number of 1s and its complement"},
	{interd1, interd1err, d5, d5err, map[string]bool{"0011": true, "": true, "0001011011": true, "100": true, "00100101100100111": false}, "Union of an intersection and a DFA: should accept strings with even number of 0s and 1s OR strings with odd number 1s"},
}

var interd1, interd1err = d1.Intersection(d6) // should accept strings with even number of 0s and 1s

var intersectionTests = []struct {
	d1          *DFA
	d1err       error
	d2          *DFA
	d2err       error
	testStrings map[string]bool
	descriptor  string
}{
	{d1, d1err, d1, d1err, map[string]bool{"0011": true, "": true, "0001011011": false, "100": false, "001001011001001011": true}, "Intersection of two DFAs accepting even nummber of 1s, both using strings as states"},
	{d1, d1err, d3, d3err, map[string]bool{"0011": true, "": true, "0001011011": false, "100": false, "001001011001001011": true}, "Intersection of two DFAs accepting even number of 1s, one using strings, other using structs as states"},
	{d1, d1err, d5, d5err, map[string]bool{"0011": false, "": false, "0001011011": false, "100": false, "001001011001001011": false}, "Intersection of DFAs, one accepting even number 1s, other accepting odd number 1s"},
	{d4, d4err, d5, d5err, map[string]bool{"0011": false, "": false, "0001011011": true, "100": false, "001001011001001011": false}, "Intersection of DFAs, one accepting odd number 1s, other accepting odd number 0s"},
	{d1, d1err, d7, d7err, map[string]bool{"0011": true, "": true, "0001011011": false, "100": false, "001001011001001011": true}, "Intersection of DFAs, one accepting all strings of 0s, 1s, other accepting even number 1s"},
	{compld1, compld1err, d1, d1err, map[string]bool{"0011": false, "": false, "0001011011": false, "100": false, "001001011001001011": false}, "Intersection of a DFA accepting strings of even number of 1s and its complement"},
	{uniond1, uniond1err, uniond2, uniond2err, map[string]bool{"0011": false, "": false, "000101101": true, "100": true, "001001011001001011": false}, "Intersection of union of DFAs, should accept strings with even number 1s and odd number 0s OR strings with odd number 1s and even number 0s"},
	{interd1, interd1err, d6, d6err, map[string]bool{"0011": true, "": true, "0001011011": false, "010": false, "001001011001001011": true}, "Intersection of an intersection and another DFA: should accept strings with even number of 0s and 1s"},
}

var compld1, compld1err = d1.Complement()

var complementTests = []struct {
	d1          *DFA
	d1err       error
	testStrings map[string]bool
	descriptor  string
}{
	{d1, d1err, map[string]bool{"0011": false, "": false, "0001011011": true, "100": true, "001001011001001011": false}, "Complement of a DFA recognizing an even number of 1s"},
	{d7, d7err, map[string]bool{"0011": false, "": false, "0001011011": false, "100": false, "001001011001001011": false}, "Complement of a DFA recognizing all strings of 0s and 1s"},
	{d8, d8err, map[string]bool{"0001": true, "": true, "0001011011": true, "100": true, "001001011001001011": true}, "Complement of a DFA recognizing no strings"},
	{uniond1, uniond1err, map[string]bool{"0001": true, "": false, "0001011011": true, "100": false, "001001011001001011": false}, "Complement of a union of DFAs which should recognize strings containing an even number of 0s or 1s"},
	{interd1, interd1err, map[string]bool{"0001": true, "": false, "0001011011": true, "100": true, "001001011001001011": false}, "Complement of an intersection of DFAs which should recognize strings containing an even number of 0s and 1s"},
	{compld1, compld1err, map[string]bool{"0001": false, "": true, "0001011011": false, "100": false, "001001011001001011": true}, "Complement of a complement of a DFA recognizing strings with even number of 1s"},
}

var differenceTests = []struct {
	d1          *DFA
	d1err       error
	d2          *DFA
	d2err       error
	testStrings map[string]bool
	descriptor  string
}{
	{d1, d1err, d1, d1err, map[string]bool{"0011": false, "": false, "0001011011": false, "100": false, "001001011001001011": false}, "Difference of two identical DFAs"},
	{d3, d3err, d4, d4err, map[string]bool{"00011": false, "": true, "0001011011": false, "100": false, "001001011001001011": true}, "Difference between DFAs using different representations for states: should accept strings with even number of 0s and 1s"},
	{d1, d1err, d5, d5err, map[string]bool{"0011": true, "": true, "0001011011": false, "100": false, "001001011001001011": true}, "Difference of DFAs, first accepting even number 1s, other accepting odd number 1s"},
	{d4, d4err, d5, d5err, map[string]bool{"0011": false, "": false, "00010110111": true, "100": false, "01001011001001011": true}, "Difference of DFAs, first accepting odd number 0s, other accepting odd number 1s: should accept odd number 0s even number 1s"},
	{d7, d7err, d1, d1err, map[string]bool{"0011": false, "": false, "0001011011": true, "100": true, "001001011001001011": false}, "Difference of DFAs, first accepting all strings of 0s, 1s, other accepting even number 1s"},
	{d1, d1err, compld1, compld1err, map[string]bool{"0011": true, "": true, "0001011011": false, "100": false, "001001011001001011": true}, "Difference of a DFA accepting strings of even number of 1s and its complement"},
	{uniond1, uniond1err, uniond2, uniond2err, map[string]bool{"0011": true, "": true, "000101101": false, "100": false, "001001011001001011": true}, "Difference of union of DFAs, should accept strings with an even number of 0s and 1s"},
	{interd1, interd1err, d6, d6err, map[string]bool{"0011": false, "": false, "0001011011": false, "010": false, "001001011001001011": false}, "Difference of an intersection and another DFA: should not accept any strings"},
}

func TestDFASimulate(t *testing.T) {
	for _, test := range simulateTests {
		if test.err != nil {
			t.Error("On test: " + test.descriptor + ", error: " + test.err.Error())
			t.FailNow()
		}
		for k, v := range test.testStrings {
			ans, err := test.d.Simulate(k)
			if ans != v {
				t.Error("On test: " + test.descriptor + ", error: DFA should have answered " + strconv.FormatBool(v) + " to string " + k)
			}
			if err != nil {
				t.Error("On test: " + test.descriptor + ", while testing string " + k + ", error: " + err.Error())
			}
		}
	}
}

func TestDFAUnion(t *testing.T) {
	for _, test := range unionTests {
		if test.d1err != nil {
			t.Error("On test: " + test.descriptor + ", error: " + test.d1err.Error())
			t.FailNow()
		}
		if test.d2err != nil {
			t.Error("On test: " + test.descriptor + ", error: " + test.d2err.Error())
			t.FailNow()
		}
		testd, err := test.d1.Union(test.d2)
		if err != nil {
			t.Error("On test: " + test.descriptor + ", error: " + err.Error())
			t.FailNow()
		}
		for k, v := range test.testStrings {
			ans, err := testd.Simulate(k)
			if ans != v {
				t.Error("On test: " + test.descriptor + ", error: DFA should have answered " + strconv.FormatBool(v) + " to string " + k)
			}
			if err != nil {
				t.Error("On test: " + test.descriptor + ", while testing string " + k + ", error: " + err.Error())
			}
		}
	}
}

func TestDFAIntersection(t *testing.T) {
	for _, test := range intersectionTests {
		if test.d1err != nil {
			t.Error("On test: " + test.descriptor + ", error: " + test.d1err.Error())
			t.FailNow()
		}
		if test.d2err != nil {
			t.Error("On test: " + test.descriptor + ", error: " + test.d2err.Error())
			t.FailNow()
		}
		testd, err := test.d1.Intersection(test.d2)
		if err != nil {
			t.Error("On test: " + test.descriptor + ", error: " + err.Error())
			t.FailNow()
		}
		for k, v := range test.testStrings {
			ans, err := testd.Simulate(k)
			if ans != v {
				t.Error("On test: " + test.descriptor + ", error: DFA should have answered " + strconv.FormatBool(v) + " to string " + k)
			}
			if err != nil {
				t.Error("On test: " + test.descriptor + ", while testing string " + k + ", error: " + err.Error())
			}
		}
	}
}

func TestDFAComplement(t *testing.T) {
	for _, test := range complementTests {
		if test.d1err != nil {
			t.Error("On test: " + test.descriptor + ", error: " + test.d1err.Error())
			t.FailNow()
		}
		testd, err := test.d1.Complement()
		if err != nil {
			t.Error("On test: " + test.descriptor + ", error: " + err.Error())
			t.FailNow()
		}
		for k, v := range test.testStrings {
			ans, err := testd.Simulate(k)
			if ans != v {
				t.Error("On test: " + test.descriptor + ", error: DFA should have answered " + strconv.FormatBool(v) + " to string " + k)
			}
			if err != nil {
				t.Error("On test: " + test.descriptor + ", while testing string " + k + ", error: " + err.Error())
			}
		}
	}
}

func TestDFADifference(t *testing.T) {
	for _, test := range differenceTests {
		if test.d1err != nil {
			t.Error("On test: " + test.descriptor + ", error: " + test.d1err.Error())
			t.FailNow()
		}
		if test.d2err != nil {
			t.Error("On test: " + test.descriptor + ", error: " + test.d2err.Error())
			t.FailNow()
		}
		testd, err := test.d1.Difference(test.d2)
		if err != nil {
			t.Error("On test: " + test.descriptor + ", error: " + err.Error())
			t.FailNow()
		}
		for k, v := range test.testStrings {
			ans, err := testd.Simulate(k)
			if ans != v {
				t.Error("On test: " + test.descriptor + ", error: DFA should have answered " + strconv.FormatBool(v) + " to string " + k)
			}
			if err != nil {
				t.Error("On test: " + test.descriptor + ", while testing string " + k + ", error: " + err.Error())
			}
		}
	}
}

//the three methods below produce equivalent DFAs, with different representations. all should accept
//strings of 0s and 1s with an even number of 1s.
func makeEvenOnesDFAStringStates() (*DFA, error) {
	states := mapset.NewSet("q0", "q1")
	alphabet := mapset.NewSet("0", "1")
	transition := func(state interface{}, input string) (nextState interface{}) {
		q0map := map[string]interface{}{"0": "q0", "1": "q1"}
		q1map := map[string]interface{}{"0": "q1", "1": "q0"}
		fullmap := map[interface{}](map[string]interface{}){"q0": q0map, "q1": q1map}
		return fullmap[state][input]
	}
	start := "q0"
	accept := mapset.NewSet("q0")
	return NewDFA(states, alphabet, transition, start, accept)
}

func makeEvenOnesDFAIntStates() (*DFA, error) {
	states := mapset.NewSet(0, 1)
	alphabet := mapset.NewSet("0", "1")
	transition := func(state interface{}, input string) (nextState interface{}) {
		q0map := map[string]interface{}{"0": 0, "1": 1}
		q1map := map[string]interface{}{"0": 1, "1": 0}
		fullmap := map[interface{}](map[string]interface{}){0: q0map, 1: q1map}
		return fullmap[state][input]
	}
	start := 0
	accept := mapset.NewSet(0)
	return NewDFA(states, alphabet, transition, start, accept)
}

func makeEvenOnesDFAStructStates() (*DFA, error) {
	type state struct {
		str string
		num int
	}
	q0 := state{"q0", 0}
	q1 := state{"q1", 1}
	states := mapset.NewSet(q0, q1)
	alphabet := mapset.NewSet("0", "1")
	transition := func(state interface{}, input string) (nextState interface{}) {
		q0map := map[string]interface{}{"0": q0, "1": q1}
		q1map := map[string]interface{}{"0": q1, "1": q0}
		fullmap := map[interface{}](map[string]interface{}){q0: q0map, q1: q1map}
		return fullmap[state][input]
	}
	start := q0
	accept := mapset.NewSet(q0)
	return NewDFA(states, alphabet, transition, start, accept)
}

/////

func makeOddZerosDFAStringStates() (*DFA, error) {
	states := mapset.NewSet("q0", "q1")
	alphabet := mapset.NewSet("0", "1")
	transition := func(state interface{}, input string) (nextState interface{}) {
		q0map := map[string]interface{}{"0": "q1", "1": "q0"}
		q1map := map[string]interface{}{"0": "q0", "1": "q1"}
		fullmap := map[interface{}](map[string]interface{}){"q0": q0map, "q1": q1map}
		return fullmap[state][input]
	}
	start := "q0"
	accept := mapset.NewSet("q1")
	return NewDFA(states, alphabet, transition, start, accept)
}

func makeOddOnesDFAStringStates() (*DFA, error) {
	states := mapset.NewSet("q0", "q1")
	alphabet := mapset.NewSet("0", "1")
	transition := func(state interface{}, input string) (nextState interface{}) {
		q0map := map[string]interface{}{"0": "q0", "1": "q1"}
		q1map := map[string]interface{}{"0": "q1", "1": "q0"}
		fullmap := map[interface{}](map[string]interface{}){"q0": q0map, "q1": q1map}
		return fullmap[state][input]
	}
	start := "q0"
	accept := mapset.NewSet("q1")
	return NewDFA(states, alphabet, transition, start, accept)
}

func makeEvenZerosDFAStringStates() (*DFA, error) {
	states := mapset.NewSet("q0", "q1")
	alphabet := mapset.NewSet("0", "1")
	transition := func(state interface{}, input string) (nextState interface{}) {
		q0map := map[string]interface{}{"0": "q1", "1": "q0"}
		q1map := map[string]interface{}{"0": "q0", "1": "q1"}
		fullmap := map[interface{}](map[string]interface{}){"q0": q0map, "q1": q1map}
		return fullmap[state][input]
	}
	start := "q0"
	accept := mapset.NewSet("q0")
	return NewDFA(states, alphabet, transition, start, accept)
}

func makeAllZerosOnesStringStates() (*DFA, error) {
	states := mapset.NewSet("q0")
	alphabet := mapset.NewSet("0", "1")
	transition := func(state interface{}, input string) (nextState interface{}) {
		return "q0"
	}
	start := "q0"
	accept := mapset.NewSet("q0")
	return NewDFA(states, alphabet, transition, start, accept)
}

func makeNoneZerosOnesStringStates() (*DFA, error) {
	states := mapset.NewSet("q0", "q1")
	alphabet := mapset.NewSet("0", "1")
	transition := func(state interface{}, input string) (nextState interface{}) {
		return "q0"
	}
	start := "q0"
	accept := mapset.NewSet("q1")
	return NewDFA(states, alphabet, transition, start, accept)
}
