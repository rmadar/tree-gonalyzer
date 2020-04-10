// Package managing the variables to plot
package variable

import (
	"fmt"
	
	"go-hep.org/x/hep/hplot"
)

type Var struct {
	Name string
	SaveName string
	TreeName string
	Value interface{}
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

// Get a value according to it's type
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

// Set user-specified style on the plot
func (v Var) SetPlotStyle(p *hplot.Plot) {

	// Plot labels
	if &v.PlotTitle != nil {
		p.Title.Text = v.PlotTitle
	}
	if &v.XLabel != nil {
		p.X.Label.Text = v.XLabel
	}
	if &v.YLabel != nil {
		p.Y.Label.Text = v.YLabel
	}

	// Axis ranges
	if &v.RangeXmin != nil {
		p.X.Min = v.RangeXmin
	}
	if &v.RangeXmax != nil {
		p.X.Max = v.RangeXmax
	}
	if &v.RangeYmin != nil {
		p.X.Min = v.RangeXmin
	}
	if &v.RangeYmax != nil {
		p.X.Max = v.RangeXmax
	}
	
	// Legend setup
	if &v.LegPosTop != nil {
		p.Legend.Top = v.LegPosTop
		if p.Legend.Top {
			p.Legend.YOffs = -5
		} else {
			p.Legend.YOffs =  5
		}
	}
	if &p.Legend.Left != nil {
		p.Legend.Left = v.LegPosLeft
		if p.Legend.Left {
			p.Legend.XOffs = 5
		} else {
			p.Legend.XOffs = -5
		}
	}
}
