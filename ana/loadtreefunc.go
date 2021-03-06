package ana

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rtree/rfunc"
)

// (bool) -> float64
func newFuncBoolToF64(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncBoolToF64{
		rvars: varsName,
		fct:   fct.(func(bool) float64),
	}, nil
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

// ([]float32, float64) -> []float64
func newFuncF32sF64ToF64s(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncF32sF64ToF64s{
		rvars: varsName,
		fct:   fct.(func([]float32, float64) []float64),
	}, nil
}

type userFuncF32sF64ToF64s struct {
	rvars []string
	v1    *[]float32
	v2    *float64
	fct   func([]float32, float64) []float64
}

func (usr *userFuncF32sF64ToF64s) RVars() []string { return usr.rvars }

func (usr *userFuncF32sF64ToF64s) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*[]float32)
	usr.v2 = args[1].(*float64)
	return nil
}

func (usr *userFuncF32sF64ToF64s) Func() interface{} {
	return func() []float64 {
		return usr.fct(*usr.v1, *usr.v2)
	}
}

// ([]float32, float32) -> []float64
func newFuncF32sF32ToF64s(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncF32sF32ToF64s{
		rvars: varsName,
		fct:   fct.(func([]float32, float32) []float64),
	}, nil
}

type userFuncF32sF32ToF64s struct {
	rvars []string
	v1    *[]float32
	v2    *float32
	fct   func([]float32, float32) []float64
}

func (usr *userFuncF32sF32ToF64s) RVars() []string { return usr.rvars }

func (usr *userFuncF32sF32ToF64s) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*[]float32)
	usr.v2 = args[1].(*float32)
	return nil
}

func (usr *userFuncF32sF32ToF64s) Func() interface{} {
	return func() []float64 {
		return usr.fct(*usr.v1, *usr.v2)
	}
}

// (float32, float32) -> float64
func newFuncF32F32ToF64(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncF32F32ToF64{
		rvars: varsName,
		fct:   fct.(func(float32, float32) float64),
	}, nil
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
	usr.v2 = args[1].(*float32)
	return nil
}

func (usr *userFuncF32F32ToF64) Func() interface{} {
	return func() float64 {
		return usr.fct(*usr.v1, *usr.v2)
	}
}

// (bool) -> bool
func newFuncBoolToBool(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncBoolToBool{
		rvars: varsName,
		fct:   fct.(func(bool) bool),
	}, nil
}

type userFuncBoolToBool struct {
	rvars []string
	v1    *bool
	fct   func(bool) bool
}

func (usr *userFuncBoolToBool) RVars() []string { return usr.rvars }

func (usr *userFuncBoolToBool) Bind(args []interface{}) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*bool)
	return nil
}

func (usr *userFuncBoolToBool) Func() interface{} {
	return func() bool {
		return usr.fct(*usr.v1)
	}
}

// (I32) -> bool
func newFuncI32ToBool(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncI32ToBool{
		rvars: varsName,
		fct:   fct.(func(int32) bool),
	}, nil
}

type userFuncI32ToBool struct {
	rvars []string
	v1    *int32
	fct   func(int32) bool
}

func (usr *userFuncI32ToBool) RVars() []string { return usr.rvars }

func (usr *userFuncI32ToBool) Bind(args []interface{}) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*int32)
	return nil
}

func (usr *userFuncI32ToBool) Func() interface{} {
	return func() bool {
		return usr.fct(*usr.v1)
	}
}

// (bool, float32) -> bool
func newFuncBoolF32ToBool(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncBoolF32ToBool{
		rvars: varsName,
		fct:   fct.(func(bool, float32) bool),
	}, nil
}

type userFuncBoolF32ToBool struct {
	rvars []string
	v1    *bool
	v2    *float32
	fct   func(bool, float32) bool
}

func (usr *userFuncBoolF32ToBool) RVars() []string { return usr.rvars }

func (usr *userFuncBoolF32ToBool) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*bool)
	usr.v2 = args[1].(*float32)
	return nil
}

func (usr *userFuncBoolF32ToBool) Func() interface{} {
	return func() bool {
		return usr.fct(*usr.v1, *usr.v2)
	}
}

// (bool, int32) -> bool
func newFuncBoolI32ToBool(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncBoolI32ToBool{
		rvars: varsName,
		fct:   fct.(func(bool, int32) bool),
	}, nil
}

type userFuncBoolI32ToBool struct {
	rvars []string
	v1    *bool
	v2    *int32
	fct   func(bool, int32) bool
}

func (usr *userFuncBoolI32ToBool) RVars() []string { return usr.rvars }

func (usr *userFuncBoolI32ToBool) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*bool)
	usr.v2 = args[1].(*int32)
	return nil
}

func (usr *userFuncBoolI32ToBool) Func() interface{} {
	return func() bool {
		return usr.fct(*usr.v1, *usr.v2)
	}
}

// ([]float32) -> []float64
func newFuncF32sToF64s(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncF32sToF64s{
		rvars: varsName,
		fct:   fct.(func([]float32) []float64),
	}, nil
}

type userFuncF32sToF64s struct {
	rvars []string
	v1    *[]float32
	fct   func([]float32) []float64
}

func (usr *userFuncF32sToF64s) RVars() []string { return usr.rvars }

func (usr *userFuncF32sToF64s) Bind(args []interface{}) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*[]float32)
	return nil
}

func (usr *userFuncF32sToF64s) Func() interface{} {
	return func() []float64 {
		return usr.fct(*usr.v1)
	}
}

// ([]float32, []float32, []float32) -> float64
func newFuncF32sF32sF32sToF64(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncF32sF32sF32sToF64{
		rvars: varsName,
		fct:   fct.(func([]float32, []float32, []float32) float64),
	}, nil
}

type userFuncF32sF32sF32sToF64 struct {
	rvars []string
	v1    *[]float32
	v2    *[]float32
	v3    *[]float32
	fct   func([]float32, []float32, []float32) float64
}

func (usr *userFuncF32sF32sF32sToF64) RVars() []string { return usr.rvars }

func (usr *userFuncF32sF32sF32sToF64) Bind(args []interface{}) error {
	if got, want := len(args), 3; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*[]float32)
	usr.v2 = args[1].(*[]float32)
	usr.v3 = args[2].(*[]float32)
	return nil
}

func (usr *userFuncF32sF32sF32sToF64) Func() interface{} {
	return func() float64 {
		return usr.fct(*usr.v1, *usr.v2, *usr.v3)
	}
}

// ([]int32, []float32) -> []float64
func newFuncI32sF32sToF64s(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncI32sF32sToF64s{
		rvars: varsName,
		fct:   fct.(func([]int32, []float32) []float64),
	}, nil
}

type userFuncI32sF32sToF64s struct {
	rvars []string
	v1    *[]int32
	v2    *[]float32
	fct   func([]int32, []float32) []float64
}

func (usr *userFuncI32sF32sToF64s) RVars() []string { return usr.rvars }

func (usr *userFuncI32sF32sToF64s) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*[]int32)
	usr.v2 = args[1].(*[]float32)
	return nil
}

func (usr *userFuncI32sF32sToF64s) Func() interface{} {
	return func() []float64 {
		return usr.fct(*usr.v1, *usr.v2)
	}
}

// (int32, float64) -> float64
func newFuncI32F64ToF64(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncI32F64ToF64{
		rvars: varsName,
		fct:   fct.(func(int32, float64) float64),
	}, nil
}

type userFuncI32F64ToF64 struct {
	rvars []string
	v1    *int32
	v2    *float64
	fct   func(int32, float64) float64
}

func (usr *userFuncI32F64ToF64) RVars() []string { return usr.rvars }

func (usr *userFuncI32F64ToF64) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*int32)
	usr.v2 = args[1].(*float64)
	return nil
}

func (usr *userFuncI32F64ToF64) Func() interface{} {
	return func() float64 {
		return usr.fct(*usr.v1, *usr.v2)
	}
}

// (int32, float32) -> float64
func newFuncI32F32ToF64(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncI32F32ToF64{
		rvars: varsName,
		fct:   fct.(func(int32, float32) float64),
	}, nil
}

type userFuncI32F32ToF64 struct {
	rvars []string
	v1    *int32
	v2    *float32
	fct   func(int32, float32) float64
}

func (usr *userFuncI32F32ToF64) RVars() []string { return usr.rvars }

func (usr *userFuncI32F32ToF64) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*int32)
	usr.v2 = args[1].(*float32)
	return nil
}

func (usr *userFuncI32F32ToF64) Func() interface{} {
	return func() float64 {
		return usr.fct(*usr.v1, *usr.v2)
	}
}

// (int32, int32) -> float64
func newFuncI32I32ToF64(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncI32I32ToF64{
		rvars: varsName,
		fct:   fct.(func(int32, int32) float64),
	}, nil
}

type userFuncI32I32ToF64 struct {
	rvars []string
	v1    *int32
	v2    *int32
	fct   func(int32, int32) float64
}

func (usr *userFuncI32I32ToF64) RVars() []string { return usr.rvars }

func (usr *userFuncI32I32ToF64) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*int32)
	usr.v2 = args[1].(*int32)
	return nil
}

func (usr *userFuncI32I32ToF64) Func() interface{} {
	return func() float64 {
		return usr.fct(*usr.v1, *usr.v2)
	}
}

// (float32, float32, float32, float32) -> float64
func newFuncF32F32F32F32ToF64(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncF32F32F32F32ToF64{
		rvars: varsName,
		fct:   fct.(func(float32, float32, float32, float32) float64),
	}, nil
}

type userFuncF32F32F32F32ToF64 struct {
	rvars []string
	v1    *float32
	v2    *float32
	v3    *float32
	v4    *float32
	fct   func(float32, float32, float32, float32) float64
}

func (usr *userFuncF32F32F32F32ToF64) RVars() []string { return usr.rvars }

func (usr *userFuncF32F32F32F32ToF64) Bind(args []interface{}) error {
	if got, want := len(args), 4; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*float32)
	usr.v2 = args[1].(*float32)
	usr.v3 = args[1].(*float32)
	usr.v4 = args[1].(*float32)
	return nil
}

func (usr *userFuncF32F32F32F32ToF64) Func() interface{} {
	return func() float64 {
		return usr.fct(*usr.v1, *usr.v2, *usr.v3, *usr.v4)
	}
}






// ([]float32) -> float64
func newFuncF32sToF64(varsName []string, fct interface{}) (rfunc.Formula, error) {
	return &userFuncF32sToF64{
		rvars: varsName,
		fct:   fct.(func([]float32) float64),
	}, nil
}

type userFuncF32sToF64 struct {
	rvars []string
	v1    *[]float32
	fct   func([]float32) float64
}

func (usr *userFuncF32sToF64) RVars() []string { return usr.rvars }

func (usr *userFuncF32sToF64) Bind(args []interface{}) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	usr.v1 = args[0].(*[]float32)
	return nil
}

func (usr *userFuncF32sToF64) Func() interface{} {
	return func() float64 {
		return usr.fct(*usr.v1)
	}
}





// Maps of all pre-defined function types.
var funcs = make(map[reflect.Type]func(rvars []string, fct interface{}) (rfunc.Formula, error))

// Load rfunc.Formula functions which are pre-defined in rtree/rfunc
// or user-defined rfunc.
func init() {

	// () -> bool
	funcs[reflect.TypeOf((func() bool)(nil))] = func(rvars []string, fct interface{}) (rfunc.Formula, error) {
		return rfunc.NewFuncToBool(rvars, fct.(func() bool)), nil
	}

	// () -> float64
	funcs[reflect.TypeOf((func() float64)(nil))] = func(rvars []string, fct interface{}) (rfunc.Formula, error) {
		return rfunc.NewFuncToF64(rvars, fct.(func() float64)), nil
	}

	// (float64) -> float64
	funcs[reflect.TypeOf((func(float64) float64)(nil))] = func(rvars []string, fct interface{}) (rfunc.Formula, error) {
		return rfunc.NewFuncF64ToF64(rvars, fct.(func(float64) float64)), nil
	}

	// (float32) -> float64
	funcs[reflect.TypeOf((func(float32) float64)(nil))] = func(rvars []string, fct interface{}) (rfunc.Formula, error) {
		return rfunc.NewFuncF32ToF64(rvars, fct.(func(float32) float64)), nil
	}

	// (float32) -> bool
	funcs[reflect.TypeOf((func(float32) bool)(nil))] = func(rvars []string, fct interface{}) (rfunc.Formula, error) {
		return rfunc.NewFuncF32ToBool(rvars, fct.(func(float32) bool)), nil
	}

	// (float64) -> bool
	funcs[reflect.TypeOf((func(float64) bool)(nil))] = func(rvars []string, fct interface{}) (rfunc.Formula, error) {
		return rfunc.NewFuncF64ToBool(rvars, fct.(func(float64) bool)), nil
	}

	// (int32) -> float64
	funcs[reflect.TypeOf((func(int32) float64)(nil))] = func(rvars []string, fct interface{}) (rfunc.Formula, error) {
		return rfunc.NewFuncI32ToF64(rvars, fct.(func(int32) float64)), nil
	}

	// (bool) -> bool
	funcs[reflect.TypeOf((func(bool) bool)(nil))] = newFuncBoolToBool

	// (int32) -> bool
	funcs[reflect.TypeOf((func(int32) bool)(nil))] = newFuncI32ToBool

	// (bool) -> float64
	funcs[reflect.TypeOf((func(bool) float64)(nil))] = newFuncBoolToF64

	// ([]float32, float64) -> float64[]
	funcs[reflect.TypeOf((func([]float32, float64) []float64)(nil))] = newFuncF32sF64ToF64s
	
	// ([]float32, float32) -> float64[]
	funcs[reflect.TypeOf((func([]float32, float32) []float64)(nil))] = newFuncF32sF32ToF64s

	// (float32, float32) -> float64
	funcs[reflect.TypeOf((func(float32, float32) float64)(nil))] = newFuncF32F32ToF64

	// ([]float32) -> []float64
	funcs[reflect.TypeOf((func([]float32) []float64)(nil))] = newFuncF32sToF64s

	// ([]float32, []float32, []float32) -> float64
	funcs[reflect.TypeOf((func([]float32, []float32, []float32) float64)(nil))] = newFuncF32sF32sF32sToF64

	// ([]int32, []float32) -> []float64
	funcs[reflect.TypeOf((func([]int32, []float32) []float64)(nil))] = newFuncI32sF32sToF64s

	// (bool, float32) -> bool
	funcs[reflect.TypeOf((func(bool, float32) bool)(nil))] = newFuncBoolF32ToBool

	// (bool, int32) -> bool
	funcs[reflect.TypeOf((func(bool, int32) bool)(nil))] = newFuncBoolI32ToBool

	// (int32, int32) -> float64
	funcs[reflect.TypeOf((func(int32, int32) float64)(nil))] = newFuncI32I32ToF64

	// (int32, float32) -> float64
	funcs[reflect.TypeOf((func(int32, float32) float64)(nil))] = newFuncI32F32ToF64

	// (int32, float64) -> float64
	funcs[reflect.TypeOf((func(int32, float64) float64)(nil))] = newFuncI32F64ToF64

	// (float32, float32, float32, float32) -> float64
	funcs[reflect.TypeOf((func(float32, float32, float32, float32) float64)(nil))] = newFuncF32F32F32F32ToF64

	// ([]float32) -> float64
	funcs[reflect.TypeOf((func([]float32) float64)(nil))] = newFuncF32sToF64

}
