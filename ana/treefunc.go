package ana

import (
	"log"

	"go-hep.org/x/hep/groot/rtree"
)

type TreeFunc struct {
	VarsName []string
	Fct      interface{}
}

// NewTreeFuncVarBool returns a TreeFunc to get
// a single boolean branch-based variable.
func NewTreeFuncVarBool(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x bool) bool { return x },
	}
}

// NewTreeFuncVarF64 returns a TreeFunc to get a single
// float64 branch-based variable.
func NewTreeFuncVarF64(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x float64) float64 { return x },
	}
}

// NewTreeFuncValF64 returns a TreeFunc to get float value,
// ie not a branch-based variable.
func NewTreeFuncValF64(v float64) TreeFunc {
	return TreeFunc{
		VarsName: []string{},
		Fct:      func() float64 { return v },
	}
}

// NewTreeFuncVarF32 returns a TreeFunc to get a float64 output
// from a single float32 branch-based variable.
func NewTreeFuncVarF32(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x float32) float64 { return float64(x) },
	}
}

// NewFormulaFunc returns the rtree.FormulaFunc function associated
// to the TreeFunc f, from a give rtree.Reader r.
func (f *TreeFunc) NewFormulaFunc(r *rtree.Reader) *rtree.FormulaFunc {
	ff, err := r.FormulaFunc(f.VarsName, f.Fct)
	if err != nil {
		log.Fatalf("could not create formulaFunc: %+v", err)
	}
	return ff
}

// GetFuncF64 returns a function to be called in the event loop to get
// the float64 value computed in f.Fct function.
func (f *TreeFunc) GetFuncF64(r *rtree.Reader) func() float64 {
	if len(f.VarsName) > 0 {
		return f.NewFormulaFunc(r).Func().(func() float64)
	} else {
		return f.Fct.(func() float64)
	}
}

// GetFuncBool returns the function to be called in the event loop to get
// the boolean value computed in f.Fct function.
func (f *TreeFunc) GetFuncBool(r *rtree.Reader) func() bool {
	if len(f.VarsName) > 0 {
		return f.NewFormulaFunc(r).Func().(func() bool)
	} else {
		return f.Fct.(func() bool)
	}
}
