package cflow_test

import (
	"go-hep.org/x/hep/groot/rtree"	

	"github.com/rmadar/tree-gonalyzer/cflow"
)

type usrEvent struct {
	pt  float32
	eta float32
	phi float32
	pid int32
}

func (e usrEvent) rvars() []rtree.ReadVar {
	return []rtree.ReadVar{
		{Name: "l_pt" , Value: &e.pt },
		{Name: "l_eta", Value: &e.eta},
		{Name: "l_phi", Value: &e.phi},
		{Name: "l_pid", Value: &e.pid},
	}
}

func (e usrEvent) weight() float64 {
	return float64(e.pt * (e.eta + e.phi) / 100.)
}

func Example_aSimpleCutFlow() {

	// Input files
	files := []string{
		"../../testdata/file1.root",
		"../../testdata/file2.root",
	}
	
	// User-defined event model
	var e cflow.Event 
	e = usrEvent{}
	
	// Cuts
	ptCut  := func(e cflow.Event) bool {return e.pt >50}
	etaCut := func(e cflow.Event) bool {return ptCut(e)  && e.eta>0.5}
	phiCut := func(e cflow.Event) bool {return etaCut(e) && e.phi<2}

	// Cut sequence
	cutSeq := cflow.NewCutSeq(
		cflow.Cut{Name: "CUT0", Pass: ptCut },
		cflow.Cut{Name: "CUT1", Pass: etaCut},
		cflow.Cut{Name: "CUT2", Pass: phiCut},
	)

	// Define the cutflow analyzer
	ana := cflow.Maker{
		Event:     eCF,
		Cuts:      cutSeq,
		FilesName: files,
		TreeName: "truth",
	}

	// Run the cutflow
	cutFlow := ana.CutFlow()
}
