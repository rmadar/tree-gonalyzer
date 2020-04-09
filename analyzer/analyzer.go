// Package allowing to wrap all needed element of a TTree plotting analysis
package analyzer

import (

	"go-hep.org/x/hep/hbook"
	
	"tree-gonalyzer/sample"
	"tree-gonalyzer/variable"
)

// Starting to implement the Analyzer object
type Analyzer struct {
	Samples []sample.Spl      // sample on which to run
	SplGroup string           // specify how to group sample together
	Variables []*variable.Var // variables to plot
	Selections []string       // implement a type selection ?
	Histos []*hbook.H1D       // multiple map histo[var, sample, cut, syst]
}
