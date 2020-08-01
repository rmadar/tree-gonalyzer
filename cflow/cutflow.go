// Doc for cutflow package
package cflow

import (
	"fmt"
	
	"go-hep.org/x/hep/groot/rtree"
)

// Event model interface
type Event interface{
	rvars() []rtree.ReadVar
	weight() float64
}

// Event yields type with both raw and
// weighted yields, and a name for a cut stage.
type Yields struct {
	Name   string  // Name of the cut stage.
	Nraw   float64 // Raw yields.
	Nwgt   float64 // Weighted yields, weight is defined by Event.weight()
}

// CutFlow is a slice of Yields, once per cut.
type CutFlow []Yields

// Cut contains the needed information
type Cut struct {
	Name string
	Pass func(e Event) bool
}

// Serie of several cuts.
type CutSequence []Cut

// NewCutSequence create a cut sequence from
// a list of cuts.
func NewCutSeq(cuts ...Cut) CutSequence {
	cs := make([]Cut, len(cuts))
	for i, c := range cuts {
		cs[i] = c
	}
	return cs
}

// From() creates a CutFlow object corresponding
// to a given cut sequence.
func From(cuts CutSequence) CutFlow {
	cf := make([]Yields, len(cuts))
	for i, cut := range cuts {
		cf[i].Name = cut.Name
	}
	return cf
}

// Cutflow string formater
func (cf CutFlow) String() string {
	var str string
	for _, y := range cf {
		str += fmt.Sprintf("%v %v %v %v\n", y.Name, y.Nraw, y.Nwgt)
	}
	return str
}
