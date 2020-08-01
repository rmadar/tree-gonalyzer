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
	}
}

// Definition of the event weight
func (e *usrEvt) Weight() float64 {
	return float64(e.pt / 1000.)
}

// Definition of the Cuts 
var (
	Cut0 = func(e cflow.Event) bool {
		evt := e.(*usrEvt)
		return evt.pt > 10
	}
	Cut1 = func(e cflow.Event) bool {
		evt := e.(*usrEvt)
		return evt.eta > 0.5
	}
	Cut2 = func(e cflow.Event) bool {
		evt := e.(*usrEvt)
		return evt.phi < 2.0
	}
)

func Example_Analysis() {

	// Input files
	files := []string{
		"../testdata/file1.root",
		"../testdata/file2.root",
	}

	// User-defined event model, based on cflow.Event
	var e cflow.Event; e = &usrEvt{}
	
	// Cut sequence
	cutSeq := cflow.NewCutSeq(
		cflow.Cut{Name: "CUT0", Sel: Cut0},
		cflow.Cut{Name: "CUT1", Sel: Cut1},
		cflow.Cut{Name: "CUT2", Sel: Cut2},
	)

	// Define the cutflow analyzer
	ana := cflow.Analysis{
		Event:     &e,
		Cuts:      cutSeq,
		FilesName: files,
		TreeName: "truth",
	}

	// Run the cutflow
	ana.Run()

	// Output:
	// CUT0 28897 1551.0207463856786
	// CUT1 10645 564.4076810218394
	// CUT2 8701 460.4137644460425
}


