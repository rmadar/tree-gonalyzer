package ana

import (
	"go-hep.org/x/hep/groot/rtree"
	"log"
)

type TreeFunc struct {
	VarsName []string
	Fct      interface{}
}

// NewVarF64 return a TreeFunc to get a single float64 variable
func NewFuncVarF64(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x float64) float64 { return x },
	}
}

// Return a
func NewFuncF64(v float64) TreeFunc {
	return TreeFunc{
		VarsName: []string{},
		Fct:      func() float64 { return v },
	}
}

// NewVarF64 return a TreeFunc to get a single float32 variable
func NewFuncVarF32(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x float32) float64 { return float64(x) },
	}
}

// Get the rtree.FormulaFunc function from the reader
func (f *TreeFunc) NewFormulaFunc(r *rtree.Reader) *rtree.FormulaFunc {
	ff, err := r.FormulaFunc(f.VarsName, f.Fct)
	if err != nil {
		log.Fatalf("could not create formulaFunc: %+v", err)
	}
	return ff
}

// Get the function to be called in the event loop to get
// 'Var' float64 value, from the reader
func (f *TreeFunc) GetVarFunc(r *rtree.Reader) func() float64 {
	return f.NewFormulaFunc(r).Func().(func() float64)
}

// Get the function to be called in the event loop to get
// 'Cut' boolean value, from the reader
func (f *TreeFunc) GetCutFunc(r *rtree.Reader) func() bool {
	return f.NewFormulaFunc(r).Func().(func() bool)
}
