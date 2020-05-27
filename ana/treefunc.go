package ana

import (
	"log"

	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/groot/rtree/rfunc"
)

// TreeFunc is a wrapper to use rtree.FormulaFunc in an easy way.
// It provides a set of functions to ease the simple cases
// of boolean, float32 and float64 branches. Once the slice of variable
// (branch) names and the function are given, one can either access
// the rtree.FormulaFunc or directly the GO function to be called
// in the event loop for boolean, float32 and float64.
type TreeFunc struct {
	VarsName []string
	Fct      interface{}
}

// NewTreeFuncVarBool returns a TreeFunc to get
// a single boolean branch-based variable.
// The output value is a boolean.
func NewVarBool(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x bool) bool { return x },
	}
}

// NewTreeFuncVarF64 returns a TreeFunc to get a single
// float64 branch-based variable. The output value is a float64.
func NewVarF64(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x float64) float64 { return x },
	}
}

// NewTreeFuncVarF32 returns a TreeFunc to get a single
// float32 branch-based variable. The output  value is a float64.
func NewVarF32(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x float32) float64 { return float64(x) },
	}
}

// NewTreeFuncVarI64 returns a TreeFunc to get a single
// int64 branch-based variable. The output value is a float64.
func NewVarI64(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x int64) float64 { return float64(x) },
	}
}

// NewTreeFuncVarI32 returns a TreeFunc to get a single
// int32 branch-based variable. The output value is a float64.
func NewVarI32(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x int32) float64 { return float64(x) },
	}
}

// NewTreeFuncValF64 returns a TreeFunc to get float value,
// ie not a branch-based variable.
func NewValF64(v float64) TreeFunc {
	return TreeFunc{
		VarsName: []string{},
		Fct:      func() float64 { return v },
	}
}

// FormulaFuncFromReader returns the rtree.FormulaFunc function associated
// to the TreeFunc f, from a give rtree.Reader r.
func (f *TreeFunc) FormulaFuncFromReader(r *rtree.Reader) rfunc.Formula {
	ff, err := r.FormulaFunc(f.VarsName, f.Fct)
	if err != nil {
		log.Fatalf("could not create formulaFunc: %+v", err)
	}
	return ff
}

// GetFuncF64 returns a function to be called in the event loop to get
// the float64 value computed in f.Fct function.
func (f *TreeFunc) GetFuncF64(r *rtree.Reader) func() float64 {
	return f.FormulaFuncFromReader(r).Func().(func() float64)
}

// GetFuncBool returns the function to be called in the event loop to get
// the boolean value computed in f.Fct function.
func (f *TreeFunc) GetFuncBool(r *rtree.Reader) func() bool {
	return f.FormulaFuncFromReader(r).Func().(func() bool)
}
