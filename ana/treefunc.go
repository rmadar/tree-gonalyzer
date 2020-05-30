package ana

import (
	"fmt"
	"log"

	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/groot/rtree/rfunc"
)

// TreeFunc is a wrapper to use rtree.Formula in an easy way.
// It provides a set of functions to ease the simple cases
// of boolean, float32 and float64 branches.
//
// Once the slice of variable (branch) names and the function
// are given, one can either access the rtree.Formula or
// directly the GO function to be called in the event loop
// for boolean, float32 and float64.
//
// Except `NewCutBool()`, all `NewXXX()` functions lead to a treeFunc.Fct
// returning a float64 a or []float64 in order to be accpeted
// by hplot.h1d.Fill(v, w) method. Instead `NewCutBool()` lead to
// a treeFunc.Fct returning a bool, specifically for cuts.
type TreeFunc struct {
	VarsName []string
	Fct      interface{}
	formula  rfunc.Formula
}

// NewCutBool returns a TreeFunc to get
// a single boolean branch-based variable for cuts.
// The output value is a boolean and cannot be used
// to be plotted. To plot a boolean, use NewVarBool(v).
func NewCutBool(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct:      func(x bool) bool { return x },
	}
}

// NewVarBool returns a TreeFunc to get
// a single boolean branch-based variable to plot.
// The output value is a float64 and cannot be
// used for selection. For cuts, use NewCutBool(v).
func NewVarBool(v string) TreeFunc {
	return TreeFunc{
		VarsName: []string{v},
		Fct: func(x bool) float64 {
			if x {
				return 1
			}
			return 0
		},
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
		Fct:      func(xs []float64) []float64 { return xs },
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
func (f *TreeFunc) rfuncFormula() rfunc.Formula {

	if f.formula != nil {
		return f.formula

	} else {
		var ff rfunc.Formula
		var err error

		switch fct := f.Fct.(type) {
		case func() bool:
			ff = rfunc.NewFuncToBool(f.VarsName, fct)
		case func() float64:
			ff = rfunc.NewFuncToF64(f.VarsName, fct)
		case func(float64) float64:
			ff = rfunc.NewFuncF64ToF64(f.VarsName, fct)
		case func(float32) float64:
			ff = rfunc.NewFuncF32ToF64(f.VarsName, fct)
		case func(float32, float32) float64:
			ff = newUsrFuncF32F32ToF64(f.VarsName, fct)
		case func(float32) bool:
			ff = rfunc.NewFuncF32ToBool(f.VarsName, fct)
		case func(float64) bool:
			ff = rfunc.NewFuncF64ToBool(f.VarsName, fct)
		case func(bool) float64:
			ff = newUsrFuncBoolToF64(f.VarsName, fct)
		default:
			ff, err = rfunc.NewGenericFormula(f.VarsName, f.Fct)
			if err != nil {
				log.Fatalf("could not create formula func: %+v", err)
			}
		}
		return ff
	}
}

func (f *TreeFunc) treeFormulaFrom(r *rtree.Reader) rfunc.Formula {
	tf, err := r.Formula(f.rfuncFormula())
	if err != nil {
		log.Fatalf("could not create formulaFunc: %+v", err)
	}
	return tf
}

// GetFuncF64 returns a function to be called in the event loop to get
// the float64 value computed in f.Fct function.
func (f *TreeFunc) GetFuncF64(r *rtree.Reader) (func() float64, bool) {
	fct, ok := f.treeFormulaFrom(r).Func().(func() float64)
	return fct, ok
}

// GetFuncF64s returns a function to be called in the event loop to get
// a slice []float64 values computed in f.Fct function.
func (f *TreeFunc) GetFuncF64s(r *rtree.Reader) (func() []float64, bool) {
	fct, ok := f.treeFormulaFrom(r).Func().(func() []float64)
	return fct, ok
}

// GetFuncBool returns the function to be called in the event loop to get
// the boolean value computed in f.Fct function.
func (f *TreeFunc) GetFuncBool(r *rtree.Reader) (func() bool, bool) {
	fct, ok := f.treeFormulaFrom(r).Func().(func() bool)
	return fct, ok
}

func newUsrFuncBoolToF64(varsName []string, fct func(bool) float64) *userFuncBoolToF64 {
	return &userFuncBoolToF64{
		rvars: varsName,
		fct:   fct,
	}
}

type userFuncBoolToF64 struct {
	rvars []string
	v1    *bool
	fct   func(bool) float64
}

func (usr *userFuncBoolToF64) RVars() []string { return usr.rvars }

func (usr *userFuncBoolToF64) Bind(args []interface{}) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*bool)
	return nil
}

func (usr *userFuncBoolToF64) Func() interface{} {
	return func() float64 {
		return usr.fct(*usr.v1)
	}
}

func newUsrFuncF32F32ToF64(varsName []string, fct func(float32, float32) float64) *userFuncF32F32ToF64 {
	return &userFuncF32F32ToF64{
		rvars: varsName,
		fct:   fct,
	}
}

type userFuncF32F32ToF64 struct {
	rvars []string
	v1    *float32
	v2    *float32
	fct   func(float32, float32) float64
}

func (usr *userFuncF32F32ToF64) RVars() []string { return usr.rvars }

func (usr *userFuncF32F32ToF64) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*float32)
	usr.v2 = args[0].(*float32)
	return nil
}

func (usr *userFuncF32F32ToF64) Func() interface{} {
	return func() float64 {
		return usr.fct(*usr.v1, *usr.v2)
	}
}
