package ana_test

import (
	"github.com/rmadar/tree-gonalyzer/ana"
)

func ExampleNewVariable() {
	// Variable 'name' corresponding to the branch 'branchF64'
	// to be histogrammed with 100 bins between 0 and 1
	ana.NewVariable("name", ana.NewVarF64("branchF64"), 100, 0, 100)
}
