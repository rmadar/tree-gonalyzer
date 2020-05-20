package ana_test

import (
	"log"
	"fmt"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
	
	"github.com/rmadar/tree-gonalyzer/ana"
)

func Example_NewTreeFuncVarBool() {

	// Get a reader for the example
	r := getReader(5)
	
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
		fmt.Println(ctx.Entry, vTreeFunc, vormFunc)
		return nil
	})

	// Output:
	// 0 false false
	// 1 false false
	// 2 false false
	// 3 false false
	// 4 false false
}


// Helper function get a reader for the examples
func getReader(nmax int64) *rtree.Reader {

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

	return r
}
