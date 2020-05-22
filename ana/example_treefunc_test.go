package ana_test

import (
	"fmt"
	"log"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"

	"github.com/rmadar/tree-gonalyzer/ana"
)

// Example showing how a general TreeFunc object works.
func ExampleTreeFunc_general() {
	// Get a reader for the example
	f, r := getReader(5)
	defer f.Close()
	defer r.Close()

	// TreeFunc object computing t_pt*t_eta
	treeFunc := ana.TreeFunc{
		VarsName: []string{"t_pt", "t_eta"},
		Fct: func(pt, eta float32) float64 {
			fmt.Printf("pt=%.2f, eta=%.2f,", pt, eta)
			return float64(pt * eta)
		},
	}

	// rtree.FormulaFunc object
	formFunc := treeFunc.FormulaFuncFromReader(r)

	// Go function to be called in the event loop
	getValue := formFunc.Func().(func() float64)

	// Event loop
	r.Read(func(ctx rtree.RCtx) error {
		fmt.Printf(" pt*eta=%.2f\n", getValue())
		return nil
	})

	// Output:
	// pt=128.41, eta=-1.17, pt*eta=-149.63
	// pt=149.89, eta=0.70, pt*eta=104.62
	// pt=212.82, eta=1.16, pt*eta=245.99
	// pt=108.26, eta=-1.63, pt*eta=-176.88
	// pt=133.40, eta=-2.63, pt*eta=-350.77
}

// Example showing how NewTreeFuncVarBool() works and compares
// to the rtree.FormulaFunc.
func ExampleNewTreeFuncVarBool() {
	// Get a reader for the example
	f, r := getReader(5)
	defer f.Close()
	defer r.Close()

	// branch name of a boolean variable in the TTree
	varName := "init_qq"

	// TreeFunc object
	treeFunc := ana.NewTreeFuncVarBool(varName)

	// rtree.FormulaFunc object
	formFunc := treeFunc.FormulaFuncFromReader(r)

	// Go function to be called in the event loop
	getValue := treeFunc.GetFuncBool(r)

	// Event loop
	r.Read(func(ctx rtree.RCtx) error {
		vTreeFunc := getValue()
		vormFunc := formFunc.Func().(func() bool)()
		fmt.Printf("%v %v %v\n", ctx.Entry, vTreeFunc, vormFunc)
		return nil
	})

	// Output:
	// 0 false false
	// 1 false false
	// 2 false false
	// 3 false false
	// 4 false false
}

// Example showing how NewTreeFuncVarF64() works and compares
// to the rtree.FormulaFunc.
func ExampleNewTreeFuncVarF64() {
	// Get a reader for the example
	f, r := getReader(5)
	defer f.Close()
	defer r.Close()

	// branch name of a float64 variable in the TTree
	varName := "truth_dphi_ll"

	// TreeFunc object
	treeFunc := ana.NewTreeFuncVarF64(varName)

	// rtree.FormulaFunc object
	formFunc := treeFunc.FormulaFuncFromReader(r)

	// Go function to be called in the event loop
	getValue := treeFunc.GetFuncF64(r)

	// Event loop
	r.Read(func(ctx rtree.RCtx) error {
		vTreeFunc := getValue()
		vormFunc := formFunc.Func().(func() float64)()
		fmt.Printf("%v %.2f %.2f\n", ctx.Entry, vTreeFunc, vormFunc)
		return nil
	})

	// Output:
	// 0 0.14 0.14
	// 1 2.17 2.17
	// 2 2.23 2.23
	// 3 1.78 1.78
	// 4 1.65 1.65
}

// Example showing how NewTreeFuncValF64() works.
// The reason why this approach exists is to be able
// to pass a simple constant to a sample, using the
// same API  ana.With.Weight(f TreeFunc).
func ExampleNewTreeFuncValF64() {
	// Get a reader for the example
	f, r := getReader(5)
	defer f.Close()
	defer r.Close()

	// TreeFunc object
	treeFunc := ana.NewTreeFuncValF64(0.33)

	// Go function to be called in the event loop
	getValue := treeFunc.GetFuncF64(r)

	// Event loop
	r.Read(func(ctx rtree.RCtx) error {
		vTreeFunc := getValue()
		fmt.Printf("%v %.2f\n", ctx.Entry, vTreeFunc)
		return nil
	})

	// Output:
	// 0 0.33
	// 1 0.33
	// 2 0.33
	// 3 0.33
	// 4 0.33
}

// Helper function get a reader for the examples
func getReader(nmax int64) (*groot.File, *rtree.Reader) {

	// Get the file
	f, err := groot.Open("../testdata/ttbar_ME.root")
	if err != nil {
		log.Fatal("could not open ROOT file ../testdata/ttbar_ME.root: %w", err)
	}

	// Get the tree
	obj, err := f.Get("truth")
	if err != nil {
		log.Fatal("could not retrieve object: %w", err)
	}
	t := obj.(rtree.Tree)

	// Get Reader
	r, err := rtree.NewReader(t, []rtree.ReadVar{}, rtree.WithRange(0, nmax))
	if err != nil {
		log.Fatal("could not create tree reader: %w", err)
	}

	return f, r
}
