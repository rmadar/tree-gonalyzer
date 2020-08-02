package cflow

import (
	"log"
	
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

// Analysis type groups together all needed
// information to perform a cutflow analysis
// on a single sample.
type Analysis struct {

	// EventMoel implements the Event interface.
	// It requires two functions, namely Event.Vars()
	// and Event.Weight(), to be defined.
	EventModel *Evt

	// Selection applied before applying individual cuts
	// of the cut sequence.
	Preselection func(e Evt) bool

	// Slice of cuts defining each stage of the cut flow.
	// The cuts are cumulated: if an event passes n-th cut,
	// it means it passes cut[0] && cut[1] && ... && cut[n].
	Cuts []Cut

	// List of the name of files to be analyzed.
	FilesName []string

	// Name of the TTree to be analyzed.
	TreeName string
}

// Run executes the event loop in order to count raw
// and weighted yields. The final cutflow is printed
// in this function, after the event loop. The typical
// output of this function is:
//  | Cut name                 | Raw Yields                    | Weighted Yields               |
//  |                          |                    Abs    Rel |                    Abs    Rel |
//  |--------------------------|-------------------------------|-------------------------------|
//  | Electron channel         |            5526   100%   100% |        28230.59   100%   100% |
//  | pT > 10 GeV              |            5281    96%    96% |        28065.97    99%    99% |
//  | Phi < 2.0 rad            |            4312    78%    82% |        22874.73    81%    82% |
//
func (ana *Analysis) Run() {
	
	// Full rtree
	var tree rtree.Tree

	// Loop over files to get the full tree
	for iFile, fName := range ana.FilesName {

		// Open the file
		f, err := groot.Open(fName)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		
		// Get the tree
		obj, err := f.Get(ana.TreeName)
		if err != nil {
			log.Fatal(err)
		}
		t := obj.(rtree.Tree)

		// Chain to the full tree
		switch iFile {
		case 0:
			tree = t
		default:
			tree = rtree.Chain(tree, t) 
		}
	}
	
	// Event 
	evt := *ana.EventModel

	// Variables to read.
	vars  := evt.Vars()
	rvars := make([]rtree.ReadVar, len(vars))
	for i, v := range vars {
		rvars[i] = rtree.ReadVar{Name: v.Name, Value: v.Value}
	}
	
	// Tree reader
        r, err := rtree.NewReader(tree, rvars)  
        if err != nil {                                               
                log.Fatalf("could not create tree reader: %+v", err)  
        }                                                             
        defer r.Close()
	
	// Cutflow corresponding to the slice of cuts.
	cutFlow := newCutFlow(ana.Cuts)

	// Loop over events
        err = r.Read(func(ctx rtree.RCtx) error {  

		// Apply preselection if any
		if ana.Preselection != nil {
			if !ana.Preselection(evt) {
				return nil
			}
		}
		
		// Loop over the cuts and cumulate them.
		pass := true
		for ic, cut := range ana.Cuts {
			pass = pass && cut.Sel(evt)
			if pass {
				cutFlow[ic].Raw += 1
				cutFlow[ic].Wgt += evt.Weight()
			}
		}
		
		return nil
	})
	
	// Print the result
        cutFlow.Print()
}
