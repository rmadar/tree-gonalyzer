package ana_test

import (
	"github.com/rmadar/tree-gonalyzer/ana"
)

func ExampleNewVariable() {
	// Variable 'name' corresponding to the branch 'branchF64'
	// to be histogrammed with 100 bins between 0 and 100.
	ana.NewVariable("name", "branchF64", new(float64), 100, 0, 100)
}

func ExampleNewVariableFromTreeFunc() {
	// Define the TreeFunc
	fct := ana.TreeFunc{
		VarsName: []string{"br1", "br2", "br3"},
		Fct:      func(x1, x2, x3 float64) float64 { return x1 + x2/x3 },
	}

	// Variable 'name' corresponding to the branch 'br1+br2/br3'
	// to be histogrammed with 100 bins between 0 and 100.
	ana.NewVariableFromTreeFunc("name", fct, 100, 0, 100)
}
