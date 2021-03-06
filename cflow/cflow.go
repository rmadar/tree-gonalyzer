// cflow package exposes types allowing to
// perform a cut flow on a single file.
package cflow

import (
	"fmt"
)

// Event model interface
type Evt interface{
	Vars()   []Var   // Return the list of needed variables.
	Weight() float64 // Define the weight (from available variables).
}

// TreeVar groups the branch name and a value
// of the proper type, as needed by rtree.ReadVar.
type Var struct {
	Name string        // Name of the branch
	Value interface{}  // Pointer of the same type of the stored branch.
}

// Event yields type with both raw and
// weighted yields, and a name for a cut stage.
type yields struct {
	Name  string  // Name of the cut stage.
	Raw   float64 // Raw yields.
	Wgt   float64 // Weighted yields, as defined by Evt.weight()
}

// cutFlow is a slice of Yields, once per cut.
type cutFlow []yields

// Cut contains the needed information
type Cut struct {
	Name string           // Name of the cut.
	Sel  func(e Evt) bool // Function defining the cut.
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

// Print outputs nicely the result
func (cf cutFlow) Print() {

	// Table header
	ul25 := "--------------------------"
	ul30 := "-------------------------------"
	fmt.Printf("\n| %-25s| %-30s| %-30s|\n", "Cut name", "Raw Yields", "Weighted Yields")
	fmt.Printf(  "| %-25s| %22s %6s | %22s %6s |\n", "", "Abs", "Rel", "Abs", "Rel")
	fmt.Printf("|%s|%s|%s|\n", ul25, ul30, ul30)

	// Print each cut yields
	for i, y := range cf {
		var yref yields
		switch i {
		case 0 : yref = cf[0]
		default: yref = cf[i-1]
		}
		absEff := efficiency(y, cf[0])
		relEff := efficiency(y, yref )
		fmt.Printf("| %-25s|%16.0f  %4.0f%%  %4.0f%% |%16.2f  %4.0f%%  %4.0f%% |\n",
			y.Name,
			y.Raw, absEff.Raw, relEff.Raw, 
			y.Wgt, absEff.Wgt, relEff.Wgt, 
		)
	}
	fmt.Printf("\n")
}

func efficiency(y, yref yields) yields {
	return yields{
		Name: y.Name,
		Raw:  y.Raw / yref.Raw * 100,
		Wgt:  y.Wgt / yref.Wgt * 100,
	}
}
