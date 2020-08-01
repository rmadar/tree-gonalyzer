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
	w.Init(os.Stdout, 1, 5, 7, ' ', tabwriter.TabIndent)
	defer w.Flush()

	// Headers
	fmt.Fprintf(w, "\n%s\t%s\t%s\t", "Cut name", "  Raw Yields  ", "Weighted Yields")
	fmt.Fprintf(w, "\n%s\t%s\t%s\t", "--------", "---------------", "---------------")
	
	// Loop over yields (cuts)
	
	for i, y := range cf {
		absEff := efficiency(y, cf[0])
		relEff := efficiency(y, cf[0])
				
		if i>0 {
			relEff = efficiency(y, cf[i-1])
			fmt.Fprintf(w, "\n%s\t%.0f %0.1f %0.1f\t%.2f %0.1f %0.1f\t",
				y.Name,
				y.Nraw, absEff.Nraw, relEff.Nraw,
				y.Nwgt, absEff.Nwgt, relEff.Nwgt,
			)	
		} else {
			
			fmt.Fprintf(w, "\n%s\t%.0f %%abs %%rel\t%.2f %%abs %%rel\t",
				y.Name,
				y.Nraw, 
				y.Nwgt, 
			)
			
		}
	}
}

func efficiency(y, yref yields) yields {
	return yields{
		Name: y.Name,
		Nraw: y.Nraw / yref.Nraw * 100,
		Nwgt: y.Nwgt / yref.Nwgt * 100,
	}
}
