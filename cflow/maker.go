package cflow

import (
	"log"
	
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

type Maker struct {

	// Event model implementing the Event interface.
	// It requires two functions, namely Event.rvars()
	// and Event.weight(), to be defined.
	Event *Event

	// Selection applied before checking individual cuts
	// of the cut sequence.
	Preselection func(e *Event) bool

	// Cut sequence defining each stage of the cut flow.
	Cuts CutSequence

	// List of the name of files to be analyzed.
	FilesName []string

	// Name of the TTree to be analyzed.
	TreeName string
}

func (ana *Maker) CutFlow() CutFlow {

	// Full rtree
	var tree rtree.Tree

	// Loop over files to get the full tree
	for iFile, fName := range ana.FilesName {

		// Open the file
		f, err := groot.Open(fName)
		if err != nil {
			log.Fatal(err)
		}
		
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
	
	// Event type
	var e Event

	// Tree reader
        r, err := rtree.NewReader(tree, e.rvars())  
        if err != nil {                                               
                log.Fatalf("could not create tree reader: %+v", err)  
        }                                                             
        defer r.Close()
	
	// Cut sequence and associated cutflow
	cutSeq  := NewCutSeq(ana.Cuts...)
	cutFlow := From(cutSeq)

	// Loop over events
        err = r.Read(func(ctx rtree.RCtx) error {  

		// Apply preselection if any
		if ana.Preselection != nil {
			if !ana.Preselection(&e) {
				return nil
			}
		}
		
		// Loop over the cuts
		for ic, cut := range ana.Cuts {
			if cut.Pass(e) {
				cutFlow[ic].Nraw += 1
				cutFlow[ic].Nwgt += e.weight()
			}
		}
		
		return nil
	})

	// return the result
	return cutFlow
}
