// Doc for cutflow package
package cflow

import (
	"fmt"
)

// Event model interface
type Evt interface{
	Vars()   []Var
	Weight() float64
}

// TreeVar groups the branch name and a value
// of the proper type.
type Var struct {
	Name string
	Value interface{}
}

// Event yields type with both raw and
// weighted yields, and a name for a cut stage.
type yields struct {
	Name   string  // Name of the cut stage.
	Nraw   float64 // Raw yields.
	Nwgt   float64 // Weighted yields, weight is defined by Event.weight()
}

// cutFlow is a slice of Yields, once per cut.
type cutFlow []yields

// Cut contains the needed information
type Cut struct {
	Name string
	Sel  func(e Evt) bool
}

// newCutFlow creates a CutFlow object corresponding
// to a given cut sequence.
func newCutFlow(cuts []Cut) cutFlow {
	cf := make([]yields, len(cuts))
	for i, cut := range cuts {
		cf[i].Name = cut.Name
	}
	return cf
}

// Cutflow string formater
func (cf cutFlow) String() string {
	var str string
	for _, y := range cf {
		str += fmt.Sprintf("%v %v %v\n", y.Name, y.Nraw, y.Nwgt)
	}
	return str
}
