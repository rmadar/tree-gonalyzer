package ana_test

import (
	"github.com/rmadar/tree-gonalyzer/ana"
)

func ExampleNewVariable() {
	// Variable 'name' corresponding to the branch 'branchF64'
	// to be histogrammed with 100 bins between 0 and 1
	varFunc := ana.NewTreeFuncVarF64("branchF64")
	ana.NewVariable("name", varFunc, 100, 0, 100)
}
