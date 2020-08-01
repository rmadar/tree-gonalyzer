package cflow_test

import (
	"github.com/rmadar/tree-gonalyzer/cflow"
)

// User defined event model
type usrEvt struct {
	pt  float32
	eta float32
	phi float32
	pid int32
}

// How to connect the event variables to the original tree
func (e *usrEvt) Vars() []cflow.Var {
	return []cflow.Var{
		{Name: "l_pt" , Value: &e.pt },
		{Name: "l_eta", Value: &e.eta},
		{Name: "l_phi", Value: &e.phi},
		{Name: "l_pid", Value: &e.pid},
	}
}

// Definition of the event weight
func (e *usrEvt) Weight() float64 {
	return float64(e.pt / 1000.)
}

// Definition of the Cuts 
var (
	presel = func(e cflow.Evt) bool {
		evt := e.(*usrEvt)
		return evt.pid == 11
	}
	cut0 = func(e cflow.Evt) bool {
		evt := e.(*usrEvt)
		return evt.pt > 10
	}
	cut1 = func(e cflow.Evt) bool {
		evt := e.(*usrEvt)
		return evt.eta > 0.5
	}
	cut2 = func(e cflow.Evt) bool {
		evt := e.(*usrEvt)
		return evt.phi < 2.0
	}
)

func ExampleAnalysis_basicCutFlow() {

	// List of input files
	files := []string{
		"../testdata/file1.root",
		"../testdata/file2.root",
	}

	// User-defined event model, based on cflow.Evt interface.
	var e cflow.Evt;
	e = &usrEvt{}
	
	// Cut sequence - they are cumulated.
	cutSeq := []cflow.Cut{
		{Name: "CUT0", Sel: cut0},
		{Name: "CUT1", Sel: cut1},
		{Name: "CUT2", Sel: cut2},
	}

	// Define the cutflow analyzer
	ana := cflow.Analysis{
		EventModel:   &e,
		Preselection: presel,
		Cuts:         cutSeq,
		FilesName:    files,
		TreeName:     "truth",
	}

	// Run the cutflow
	ana.Run()

	// Output:
	// Cut name     Raw Yields   Weighted Yields
	// --------     ----------   ---------------
	// CUT0         14516        778.10         
	// CUT1         5281         280.66         
	// CUT2         4312         228.75
}


