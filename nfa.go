package gocompute

import (
	//	"fmt"
	"golang-set"
	//	"reflect"
)

type NFA struct {
	states     mapset.Set
	alphabet   mapset.Set
	transition func(state interface{}, input string) (transitionSet mapset.Set)
	start      interface{}
	accept     mapset.Set
	simulate   func(w string, d *DFA) (output bool, err error)
}
