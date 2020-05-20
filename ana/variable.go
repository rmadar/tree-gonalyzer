// Package managing the variables to plot
package ana

import (
	"fmt"

	"go-hep.org/x/hep/hplot"
)

type Variable struct {
	Name        string
	SaveName    string
	TreeName    string
	Value       interface{}
	TreeVar     TreeFunc
	Nbins       int
	Xmin, Xmax  float64
	XLabel      string
	YLabel      string
	XTickFormat string
	YTickFormat string
	RangeXmin   float64
	RangeXmax   float64
	RangeYmin   float64
	RangeYmax   float64
	LegPosTop   bool
	LegPosLeft  bool
}

// Create a new type variable
func NewVariable(name, tname string, value interface{},
	nbins int, xmin, xmax float64, opts ...Options) *Variable {

	// Create the object
	v := &Variable{
		Name:     name,
		TreeName: tname,
		Value:    value,
		Nbins:    nbins,
		Xmin:     xmin,
		Xmax:     xmax,
	}

	// Configuration with default values for all optional fields
	cfg := newConfig(
		WithSaveName(v.Name),
		WithAxisLabels(`Variable`, `Events`),
		WithLegPosition(true, false),
	)

	// Update the configuration looping over functional options
	for _, opt := range opts {
		opt(cfg)
	}

	// Set fields with updaded configuration
	v.SaveName = cfg.SaveName
	v.TreeVar = cfg.TreeVar
	v.XLabel = cfg.XLabel
	v.YLabel = cfg.YLabel
	v.XTickFormat = cfg.XTickFormat
	v.YTickFormat = cfg.YTickFormat
	v.RangeXmin = cfg.RangeXmin
	v.RangeXmax = cfg.RangeXmax
	v.RangeYmin = cfg.RangeYmin
	v.RangeYmax = cfg.RangeYmax
	v.LegPosTop = cfg.LegPosTop
	v.LegPosLeft = cfg.LegPosLeft

	return v
}

// Create a new variable from a TreeFunc object
func NewVariableFromTreeFunc(name string, f TreeFunc, nbins int,
	xmin, xmax float64, opts ...Options) *Variable {
	v := NewVariable(name, "", nil, nbins, xmin, xmax, opts...)
	v.TreeVar = f
	return v
}

// Get a value according to it's type
func (v Variable) GetValue() float64 {
	switch v := v.Value.(type) {
	case *float64:
		return *v
	case *float32:
		return float64(*v)
	case *bool:
		return map[bool]float64{true: 1, false: 0}[*v]
	default:
		panic(fmt.Errorf("invalid variable value-type %T", v))
	}
}

// Set user-specified style on the plot
func (v Variable) SetPlotStyle(p *hplot.Plot) {

	// Plot labels
	if v.XLabel != "" {
		p.X.Label.Text = v.XLabel
	}
	if v.YLabel != "" {
		p.Y.Label.Text = v.YLabel
	}

	// Axis ranges
	if v.RangeXmin != v.RangeXmax {
		p.X.Min = v.RangeXmin
		p.X.Max = v.RangeXmax
	}
	if v.RangeYmin != v.RangeYmax {
		p.Y.Min = v.RangeYmin
		p.Y.Max = v.RangeYmax
	}

	// Axis ticks tuning
	if v.XTickFormat != "" {
		p.X.Tick.Marker = hplot.Ticks{N: 10, Format: v.XTickFormat}
	}
	if v.YTickFormat != "" {
		p.Y.Tick.Marker = hplot.Ticks{N: 10, Format: v.YTickFormat}
	}

	// Legend position: basic
	p.Legend.Top = v.LegPosTop
	p.Legend.Left = v.LegPosLeft

	// Legend position: fine tuning
	if p.Legend.Top {
		p.Legend.YOffs = -5
	} else {
		p.Legend.YOffs = 5
	}
	if p.Legend.Left {
		p.Legend.XOffs = 5
	} else {
		p.Legend.XOffs = -5
	}

}
