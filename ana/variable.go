package ana

import (
	"go-hep.org/x/hep/hplot"
)

type Variable struct {
	Name       string      // Variable name.
	TreeFunc   TreeFunc    // Variable definition from branches & functions.
	Nbins      int         // Number of bins of final histograms.
	Xmin, Xmax float64     // Mininum and maximum values of the histogram.
	SaveName string   // Name of the plot to be saved default (default 'Name').
	XLabel, YLabel           string  // Axis labels (default: 'Variable', 'Events').
	XTickFormat, YTickFormat string  // Axis tick formatting (default: hplot default).
	RangeXmin, RangeXmax     float64 // X-axis range (default: hplot default).
	RangeYmin, RangeYmax     float64 // Y-axis range (default: hplot default).
	RatioYmin, RatioYmax     float64 // Ratio Y-axis range (default: hplot default).
	LegPosTop, LegPosLeft bool // Legend position (default: true, false)
}

// NewVariable creates a new variable value with
// default settings.
func NewVariable(name string, tFunc TreeFunc, nBins int, xMin, xMax float64, opts ...VariableOptions) *Variable {

	// Create the object
	v := &Variable{
		Name:     name,
		TreeFunc: tFunc,
		Nbins:    nBins,
		Xmin:     xMin,
		Xmax:     xMax,
	}

	// Configuration with default values for all optional fields
	cfg := newConfig(
		WithSaveName(v.Name),
		WithAxisLabels(`Variable`, `Events`),
		WithLegTop(true),
		WithLegLeft(false),
	)

	// Update the configuration looping over functional options
	for _, opt := range opts {
		opt(cfg)
	}

	// Set fields with updaded configuration
	v.SaveName = cfg.SaveName
	v.XLabel = cfg.XLabel
	v.YLabel = cfg.YLabel
	v.XTickFormat = cfg.XTickFormat
	v.YTickFormat = cfg.YTickFormat
	v.RangeXmin = cfg.RangeXmin
	v.RangeXmax = cfg.RangeXmax
	v.RangeYmin = cfg.RangeYmin
	v.RangeYmax = cfg.RangeYmax
	v.RatioYmin = cfg.RatioYmin
	v.RatioYmax = cfg.RatioYmax
	v.LegPosTop = cfg.LegPosTop
	v.LegPosLeft = cfg.LegPosLeft

	return v
}

// SetPlotStyle sets the user-specified style on
// the hplot.Plot value.
func (v Variable) setPlotStyle(p *hplot.Plot) {

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
