package ana

import (
	"go-hep.org/x/hep/hplot"
)

type Variable struct {
	Name                     string   // Variable name.
	TreeFunc                 TreeFunc // Variable definition from branches & functions.
	Nbins                    int      // Number of bins of final histograms.
	Xmin, Xmax               float64  // Mininum and maximum values of the histogram.
	SaveName                 string   // Name of the saved plot (default 'Name').
	XLabel, YLabel           string   // Axis labels (default: 'Variable', 'Events').
	XTickFormat, YTickFormat string   // Axis tick formatting (default: hplot default).
	RangeXmin, RangeXmax     float64  // X-axis range (default: hplot default).
	RangeYmin, RangeYmax     float64  // Y-axis range (default: hplot default).
	RatioYmin, RatioYmax     float64  // Ratio Y-axis range (default: hplot default).
	LegPosTop, LegPosLeft    bool     // Legend position (default: true, false)

	isSlice bool
}

// NewVariable creates a new variable value with
// default settings. The TreeFunc object should returns either
// a float64 or a []float64. Any other returned type will panic.
func NewVariable(name string, tFunc TreeFunc, nBins int, xMin, xMax float64, opts ...VariableOptions) *Variable {

	// Create the object
	v := &Variable{
		Name:      name,
		TreeFunc:  tFunc,
		Nbins:     nBins,
		Xmin:      xMin,
		Xmax:      xMax,
		SaveName:  name,
		XLabel:    `Variable`,
		YLabel:    `Events`,
		LegPosTop: true,
	}

	// Configuration with default values for all optional fields
	cfg := newConfig()

	// Update the configuration looping over functional options
	for _, opt := range opts {
		opt(cfg)
	}

	// Set fields with updaded configuration
	if cfg.SaveName.usr {
		v.SaveName = cfg.SaveName.val
	}
	if cfg.XLabel.usr {
		v.XLabel = cfg.XLabel.val
	}
	if cfg.YLabel.usr {
		v.YLabel = cfg.YLabel.val
	}
	if cfg.XTickFormat.usr {
		v.XTickFormat = cfg.XTickFormat.val
	}
	if cfg.YTickFormat.usr {
		v.YTickFormat = cfg.YTickFormat.val
	}
	if cfg.RangeXmin.usr {
		v.RangeXmin = cfg.RangeXmin.val
	}
	if cfg.RangeXmax.usr {
		v.RangeXmax = cfg.RangeXmax.val
	}
	if cfg.RangeYmin.usr {
		v.RangeYmin = cfg.RangeYmin.val
	}
	if cfg.RangeYmax.usr {
		v.RangeYmax = cfg.RangeYmax.val
	}
	if cfg.RatioYmin.usr {
		v.RatioYmin = cfg.RatioYmin.val
	}
	if cfg.RatioYmax.usr {
		v.RatioYmax = cfg.RatioYmax.val
	}
	if cfg.LegPosTop.usr {
		v.LegPosTop = cfg.LegPosTop.val
	}
	if cfg.LegPosLeft.usr {
		v.LegPosLeft = cfg.LegPosLeft.val
	}

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
