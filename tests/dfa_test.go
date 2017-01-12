package gocompute

import (
	"fmt"
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

func TestDFAConstructorStringStates(t *testing.T) {

}

//////////////// helper fn's /////////////////////

//returns pointer to different DFAs for testing purposes based on the given int
func getDFA(int num) *DFA {

}
