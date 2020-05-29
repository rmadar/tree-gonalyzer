package ana_test

import (
	"fmt"
	"log"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"

	"github.com/rmadar/tree-gonalyzer/ana"
)

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

	// rtree.Formula object
	formula := treeFunc.FormulaFrom(r)

	// Go function to be called in the event loop
	getValue := formula.Func().(func() float64)

	// Event loop
	r.Read(func(ctx rtree.RCtx) error {
		fmt.Printf(" pt*eta=%.2f\n", getValue())
		return nil
	})

	// Output:
	// pt=145.13, eta=-2.08, pt*eta=-302.15
	// pt=13.85, eta=-1.69, pt*eta=-23.35
	// pt=44.03, eta=-3.93, pt*eta=-173.06
	// pt=136.98, eta=0.64, pt*eta=87.88
	// pt=77.47, eta=2.93, pt*eta=226.79
}

func ExampleTreeFunc_withBranchBoolForPlot() {
	// Get a reader for the example
	f, r := getReader(5)
	defer f.Close()
	defer r.Close()

	// TreeFunc object from a boolean branch name in the TTree
	treeFunc := ana.NewVarBool("init_qq")

	// Go function to be called in the event loop
	// The return type is float64, since it's for plotting.
	getValue, ok := treeFunc.GetFuncF64(r)
	if !ok {
		log.Fatal("type assertion failed: expect float64")
	}

	// rtree.Formula object
	formula := treeFunc.FormulaFrom(r)

	// Event loop
	r.Read(func(ctx rtree.RCtx) error {
		vTreeFunc := getValue()
		vFormula := formula.Func().(func() float64)()
		fmt.Printf("%v %v %v\n", ctx.Entry, vTreeFunc, vFormula)
		return nil
	})

	// Output:
	// 0 0 0
	// 1 0 0
	// 2 1 1
	// 3 0 0
	// 4 0 0
}

func ExampleTreeFunc_withBranchBoolForCut() {
	// Get a reader for the example
	f, r := getReader(5)
	defer f.Close()
	defer r.Close()

	// TreeFunc object from a boolean branch name in the TTree
	// The return type is boolean, since it's for a cut.
	treeFunc := ana.NewCutBool("init_qq")

	// Go function to be called in the event loop
	getValue, ok := treeFunc.GetFuncBool(r)
	if !ok {
		log.Fatal("type assertion failed: expect bool")
	}

	// rtree.Formula object
	formula := treeFunc.FormulaFrom(r)

	// Event loop
	r.Read(func(ctx rtree.RCtx) error {
		vTreeFunc := getValue()
		vFormula := formula.Func().(func() bool)()
		fmt.Printf("%v %v %v\n", ctx.Entry, vTreeFunc, vFormula)
		return nil
	})

	// Output:
	// 0 false false
	// 1 false false
	// 2 true true
	// 3 false false
	// 4 false false
}

func ExampleTreeFunc_withBranchF64() {
	// Get a reader for the example
	f, r := getReader(5)
	defer f.Close()
	defer r.Close()

	// TreeFunc object from a float64 branch name in the TTree.
	// The return type is []float64.
	treeFunc := ana.NewVarF64("truth_dphi_ll")

	// Go function to be called in the event loop
	getValue, ok := treeFunc.GetFuncF64(r)
	if !ok {
		log.Fatal("type assertion failed: expect float64")
	}

	// rtree.Formula object
	formula := treeFunc.FormulaFrom(r)

	// Event loop
	r.Read(func(ctx rtree.RCtx) error {
		vTreeFunc := getValue()
		vFormula := formula.Func().(func() float64)()
		fmt.Printf("%v %.2f %.2f\n", ctx.Entry, vTreeFunc, vFormula)
		return nil
	})

	// Output:
	// 0 2.99 2.99
	// 1 1.07 1.07
	// 2 3.03 3.03
	// 3 0.07 0.07
	// 4 2.35 2.35
}

func ExampleTreeFunc_withBranchF32s() {
	// Get a reader for the example
	f, r := getReaderWithSlices(5)
	defer f.Close()
	defer r.Close()

	// TreeFunc object from a []float32 branch name in the TTree.
	// The return type is []float64.
	treeFunc := ana.NewVarF32s("hits_time_mc")

	// Go function to be called in the event loop
	getVal, ok := treeFunc.GetFuncF64s(r)
	if !ok {
		log.Fatal("type assertion failed: expect []float64")
	}

	// rtree.Formula object
	formula := treeFunc.FormulaFrom(r)
	getValForm := formula.Func().(func() []float64)

	// Event loop
	r.Read(func(ctx rtree.RCtx) error {
		vTreeFunc := getVal()
		vFormula := getValForm()
		fmt.Printf("Evt[%v] %v %v\n", ctx.Entry, vTreeFunc, vFormula)
		return nil
	})

	// Output:
	// 0 2.99 2.99
	// 1 1.07 1.07
	// 2 3.03 3.03
	// 3 0.07 0.07
	// 4 2.35 2.35
}

func ExampleTreeFunc_debugSlices() {

	f, err := groot.Open("../testdata/fileSlices.root")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	o, err := f.Get("modules")
	if err != nil {
		panic(err)
	}
	t := o.(rtree.Tree)

	r, err := rtree.NewReader(t, []rtree.ReadVar{}, rtree.WithRange(0, 5))
	if err != nil {
		panic(err)
	}
	defer r.Close()

	form, err := r.FormulaFunc(
		[]string{"hits_time_mc"},
		func(xs []float32) []float64 {
			o := make([]float64, len(xs))
			for i, v := range xs {
				o[i] = float64(2 * v)
			}
			return o
		},
	)
	if err != nil {
		panic(err)
	}

	fct := form.Func().(func() []float64)

	err = r.Read(func(ctx rtree.RCtx) error {
		fmt.Printf("evt[%d]: %v\n", ctx.Entry, fct())
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Output:
	// blabla
}

// Example showing how to load a numerical value in a TreeFunc.
// The reason why this approach exists is to be able
// to pass a simple constant to a sample, using the
// same API  ana.With.Weight(f TreeFunc).
func ExampleTreeFunc_withNumericalValue() {
	// Get a reader for the example
	f, r := getReader(5)
	defer f.Close()
	defer r.Close()

	// TreeFunc object from a float64
	treeFunc := ana.NewValF64(0.33)

	// Go function to be called in the event loop
	getValue, ok := treeFunc.GetFuncF64(r)
	if !ok {
		log.Fatal("type assertion failed: expect float64")
	}

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

// Helper function get a reader (w/o slices) for the examples
func getReader(nmax int64) (*groot.File, *rtree.Reader) {
	return getReaderFile("../testdata/file1.root", "truth", nmax)
}

// Helper function get a reader with slices for the examples
func getReaderWithSlices(nmax int64) (*groot.File, *rtree.Reader) {
	return getReaderFile("../testdata/fileSlices.root", "modules", nmax)
}

// Helper function to get a tree tname from a file fname.
func getReaderFile(fname, tname string, nmax int64) (*groot.File, *rtree.Reader) {

	// Get the file
	f, err := groot.Open(fname)
	if err != nil {
		log.Fatal("example_treefunc_test.go: could not open "+fname+": %w", err)
	}

	// Get the tree
	obj, err := f.Get(tname)
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
