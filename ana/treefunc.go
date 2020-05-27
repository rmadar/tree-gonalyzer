package ana

import (
	"fmt"
	"log"

	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/groot/rtree/rfunc"
)

// TreeFunc is a wrapper to use rtree.Formula in an easy way.
// It provides a set of functions to ease the simple cases
// of boolean, float32 and float64 branches. Once the slice of variable
// (branch) names and the function are given, one can either access
// the rtree.Formula or directly the GO function to be called
// in the event loop for boolean, float32 and float64.
type TreeFunc struct {
	VarsName []string
	Fct      interface{}
}

// NewVarBool returns a TreeFunc to get
// a single boolean branch-based variable.
// The output value is a boolean.
func NewVarBool(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x bool) bool { return x },
	}
}

// NewVarF64 returns a TreeFunc to get a single
// float64 branch-based variable. The output value is a float64.
func NewVarF64(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x float64) float64 { return x },
	}
}

// NewVarF32 returns a TreeFunc to get a single
// float32 branch-based variable. The output  value is a float64.
func NewVarF32(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x float32) float64 { return float64(x) },
	}
}

// NewVarI64 returns a TreeFunc to get a single
// int64 branch-based variable. The output value is a float64.
func NewVarI64(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x int64) float64 { return float64(x) },
	}
}

// NewVarI32 returns a TreeFunc to get a single
// int32 branch-based variable. The output value is a float64.
func NewVarI32(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x int32) float64 { return float64(x) },
	}
}

// NewValF64 returns a TreeFunc to get float value,
// ie not a branch-based variable.
func NewValF64(v float64) TreeFunc {
	return TreeFunc{
		VarsName: []string{},
		Fct:      func() float64 { return v },
	}
}

// NewVarF64s returns a TreeFunc to get a slice of
// float64 branch-based variable. The output value is a float64.
func NewVarF64s(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct: func(xs []float64) []float64 { return xs },
	}
}

// NewVarF32s returns a TreeFunc to get a slice of
// float32 branch-based variable. The output  value is a float64.
func NewVarF32s(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct: func(xs []float32) []float64 {
			res := make([]float64, len(xs))
			for i, x := range xs {
				fmt.Println(i, x)
				res[i] = float64(x)
			}
			return res
		},
	}
}

// NewVarI64s returns a TreeFunc to get a slice of
// int64 branch-based variable. The output value is a float64.
func NewVarI64s(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct: func(xs []int64) []float64 {
			res := make([]float64, len(xs))
			for i, x := range xs {
				res[i] = float64(x)
			}
			return res
		},
	}
}

// NewVarI32s returns a TreeFunc to get a slice of
// int32 branch-based variable. The output value is a float64.
func NewVarI32s(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct: func(xs []int32) []float64 {
			res := make([]float64, len(xs))
			for i, x := range xs {
				res[i] = float64(x)
			}
			return res
		},
	}
}

// FormulaFrom returns the rtree.FormulaFunc function associated
// to the TreeFunc f, from a give rtree.Reader r.
func (f *TreeFunc) FormulaFrom(r *rtree.Reader) rfunc.Formula {
	ff, err := r.FormulaFunc(f.VarsName, f.Fct)
	if err != nil {
		log.Fatalf("could not create formulaFunc: %+v", err)
	}
	return ff
}

// GetFuncF64 returns a function to be called in the event loop to get
// the float64 value computed in f.Fct function.
func (f *TreeFunc) GetFuncF64(r *rtree.Reader) func() float64 {
	return f.FormulaFrom(r).Func().(func() float64)
}

// GetFuncF64s returns a function to be called in the event loop to get
// a slice []float64 values computed in f.Fct function.
func (f *TreeFunc) GetFuncF64s(r *rtree.Reader) func() []float64 {
	return f.FormulaFrom(r).Func().(func() []float64)
}

// GetFuncBool returns the function to be called in the event loop to get
// the boolean value computed in f.Fct function.
func (f *TreeFunc) GetFuncBool(r *rtree.Reader) func() bool {
	return f.FormulaFrom(r).Func().(func() bool)
}
