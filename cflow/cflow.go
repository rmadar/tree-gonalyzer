// Doc for cutflow package
package cflow

import (
	"os"
	"fmt"
	"text/tabwriter"
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
	Name   string  // Name of the cut stage.
	Nraw   float64 // Raw yields.
	Nwgt   float64 // Weighted yields, as defined by Evt.weight()
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

	// minwidth, tabwidth, padding, padchar, flags
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 13, 5, 0, ' ', tabwriter.TabIndent)
	defer w.Flush()

	// Headers
	fmt.Fprintf(w, "\n%s\t%s\t%s\t", "Cut name", "Raw Yields", "Weighted Yields")
	fmt.Fprintf(w, "\n%s\t%s\t%s\t", "--------", "----------", "---------------")
	
	// Loop over yields (cuts)
	for _, y := range cf {
		fmt.Fprintf(w, "\n%s\t%.0f\t%.2f\t", y.Name, y.Nraw, y.Nwgt)
	}
}
