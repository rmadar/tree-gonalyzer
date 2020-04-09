// Package managing the variables to plot
package variable

import (
	"fmt"
)

type Var struct {
	OutputName string
	TreeName string
	Value interface{}
	Type string
	Nbins int
	Xmin, Xmax float64
	PlotTitle string
	XLabel string
	YLabel string
	RangeXmin float64
	RangeXmax float64
	RangeYmin float64
	RangeYmax float64
	LegPosTop bool
	LegPosLeft bool
}


func (v Var) GetValue() float64 {
        switch v := v.Value.(type) {
        case *float64:
                return *v
        case *float32:
                return float64(*v)
        default:
                panic(fmt.Errorf("invalid variable value-type %T", v))
        }
}

