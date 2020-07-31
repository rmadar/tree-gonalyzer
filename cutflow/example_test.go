package cutflow_test

import (
	_ "fmt"

	"go-hep.org/x/hep/groot/rtree"
	
	"github.com/rmadar/tree-gonalyzer/cutflow"
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
	var e usrEvent
	e = &cutflow.Event{}
	
	// Cuts
	ptCut  := func(e usrEvent) bool {return e.pt >50}
	etaCut := func(e usrEvent) bool {return ptCut(e)  && e.eta>0.5}
	phiCut := func(e usrEvent) bool {return etaCut(e) && e.phi<2}

	// Cut sequence
	cutSeq := cutflow.NewCutSeq(
		cutflow.Cut{Name: "CUT0", Pass: ptCut },
		cutflow.Cut{Name: "CUT1", Pass: etaCut},
		cutflow.Cut{Name: "CUT2", Pass: phiCut},
	)

	// Define the cutflow analyzer
	ana := cutflow.Maker{
		Event:     &e.(cutflow.Event),
		Cuts:      cutSeq,
		FilesName: files,
		TreeName: "truth",
	}

	// Run the cutflow
	cutFlow := ana.CutFlow()
}
