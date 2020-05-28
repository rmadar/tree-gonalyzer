package ana

import (
	"image/color"

	"gonum.org/v1/plot/vg"
)

// Options encodes the various settings to pass to
// an analysis, ie a Maker, as arguments of ana.New().
type Options func(cfg *config)

// VariableOptions encodes various settings to pass
// to a Variable type, as arguments of ana.NewVariable().
type VariableOptions func(cfg *config)

// SampleOptions encodes the various settings to pass
// to a Sample type, as arguments of ana.NewSample().
type SampleOptions func(cfg *config)

// config contains all the possible options and their values.
type config struct {

	// Maker options
	KinemCuts    []*Selection // List of cuts.
	Nevts        int64        // Maximum of processed events.
	SavePath     string       // Path to which plot will be saved.
	SaveFormat   string       // Extension of saved figure 'tex', 'pdf', 'png'.
	CompileLatex bool         // Enable on-the-fly latex compilation of plots.
	AutoStyle    bool         // Enable auto style of the histogram.
	PlotTitle    string       // General plot title.
	RatioPlot    bool         // Enable ratio plot.
	HistoStack   bool         // Disable histogram stacking (e.g. compare various processes).
	HistoNorm    bool         // Normalize distributions to unit area (when stacked, the total is normalized).
	TotalBand    bool         // Enable total error band for stacks.
	ErrBandColor color.NRGBA  // Color for the uncertainty band.

	// Sample options
	WeightFunc        TreeFunc    // Weight applied to the sample
	CutFunc           TreeFunc    // Cut applied to the sample
	LineColor         color.NRGBA // Line color of the sample histogram
	LineWidth         vg.Length   // Line width of the sample histogram
	LineDashes        []vg.Length // Line dashes format
	FillColor         color.NRGBA // Fill color of the sample histogram
	CircleMarkers     bool        // Use of circled marker
	CircleSize        vg.Length   // Size of the markers
	CircleColor       color.NRGBA // Color of the markers
	Band              bool        // Enable error band display
	YErrBars          bool        // Use of y error bars
	YErrBarsLineWidth vg.Length   // Line width of the y error bar
	YErrBarsCapWidth  vg.Length   // Width of the y error bar caps
	DataStyle         bool        // Use default data style histogram

	// Variable options
	SaveName                 string   // Filename for a variable plot
	TreeVar                  TreeFunc // TreeFunc for a computed variable (ie not a single branch)
	XLabel, YLabel           string   // Axis labels
	XTickFormat, YTickFormat string   // Ticks formating
	RangeXmin, RangeXmax     float64  // X-axis ranges
	RangeYmin, RangeYmax     float64  // Y-axis ranges
	RatioYmin, RatioYmax     float64  // Ratio Y-axis range
	LegPosTop, LegPosLeft    bool     // Legend position
}

// newConfig returns a config type with a set of passed options.
func newConfig(opts ...func(cfg *config)) *config {
	cfg := new(config)
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// WithKinemCuts sets the list of kinematic cuts to run on.
func WithKinemCuts(c []*Selection) Options {
	return func(cfg *config) {
		cfg.KinemCuts = c
	}
}

// WithNevts sets the maximum processed event for
// each sample component.
func WithNevts(n int64) Options {
	return func(cfg *config) {
		cfg.Nevts = n
	}
}

// WithSavePath sets the path to save plots.
func WithSavePath(p string) Options {
	return func(cfg *config) {
		cfg.SavePath = p
	}
}

// WithSaveFormat sets the format for the plots.
func WithSaveFormat(f string) Options {
	return func(cfg *config) {
		cfg.SaveFormat = f
	}
}

// WithCompileLatex enables automatic latex compilation.
func WithCompileLatex(b bool) Options {
	return func(cfg *config) {
		cfg.CompileLatex = b
	}
}

// WithAutoStyle enables automatic styling of the histograms.
func WithAutoStyle(b bool) Options {
	return func(cfg *config) {
		cfg.AutoStyle = b
	}
}

// WithPlotTitle sets the general plot title.
func WithPlotTitle(t string) Options {
	return func(cfg *config) {
		cfg.PlotTitle = t
	}
}

// WithRatioPlot enables the ratio plot panel.
func WithRatioPlot(b bool) Options {
	return func(cfg *config) {
		cfg.RatioPlot = b
	}
}

// WithHistoStack enables histogram stacking for
// bkg-typed samples.
func WithHistoStack(b bool) Options {
	return func(cfg *config) {
		cfg.HistoStack = b
	}
}

// WithHistoStack enables histogram stacking for
// bkg-typed samples.
func WithTotalBand(b bool) Options {
	return func(cfg *config) {
		cfg.TotalBand = b
	}
}

// WithHistoNorm enables histogram normalization to unity.
func WithHistoNorm(b bool) Options {
	return func(cfg *config) {
		cfg.HistoNorm = b
	}
}

// WithErrBandColor sets the color for the error band of
// total histogram (and ratio).
func WithErrBandColor(c color.NRGBA) Options {
	return func(cfg *config) {
		cfg.ErrBandColor = c
	}
}

// WithWeight sets the weight to be used for this sample,
// as defined by the TreeFunc f, which must return a float64.
// Maker.FillHisto() will panic otherwise.
func WithWeight(f TreeFunc) SampleOptions {
	return func(cfg *config) {
		cfg.WeightFunc = f
	}
}

// WithCut sets the cut to be applied to the sample, as
// defined by the TreeFunc f, which must return a bool.
// Maker.FillHisto() will panic otherwise.
func WithCut(f TreeFunc) SampleOptions {
	return func(cfg *config) {
		cfg.CutFunc = f
	}
}

// WithLineColor sets the line color of the histogram.
func WithLineColor(c color.NRGBA) SampleOptions {
	return func(cfg *config) {
		cfg.LineColor = c
	}
}

// WithLineWidth sets line width of the sample histogram.
func WithLineWidth(w vg.Length) SampleOptions {
	return func(cfg *config) {
		cfg.LineWidth = w
	}

}

// WithLineWidth sets line width of the sample histogram.
func WithLineDashes(s []vg.Length) SampleOptions {
	return func(cfg *config) {
		cfg.LineDashes = s
	}

}

// WithFillColor sets the color with which the histo will be filled.
func WithFillColor(c color.NRGBA) SampleOptions {
	return func(cfg *config) {
		cfg.FillColor = c
	}
}

// WithCircleMarkers enables the use of circle markers (as for data histogram).
func WithCircleMarkers(b bool) SampleOptions {
	return func(cfg *config) {
		cfg.CircleMarkers = b
	}
}

// WithCircleSize sets the size of circle markers.
func WithCircleSize(s vg.Length) SampleOptions {
	return func(cfg *config) {
		cfg.CircleSize = s
	}
}

// WithCircleColor sets the color of circle markers.
func WithCircleColor(c color.NRGBA) SampleOptions {
	return func(cfg *config) {
		cfg.CircleColor = c
	}
}

// WithBand enables y error band.
func WithBand(b bool) SampleOptions {
	return func(cfg *config) {
		cfg.Band = b
	}
}

// WithYErrBars enables y error bars.
func WithYErrBars(b bool) SampleOptions {
	return func(cfg *config) {
		cfg.YErrBars = b
	}
}

// WithYErrBarsLineWidth sets the width of the error bars line
func WithYErrBarsLineWidth(w vg.Length) SampleOptions {
	return func(cfg *config) {
		cfg.YErrBarsLineWidth = w
	}
}

// WithYErrBarsCapsWidth sets the width of the y error bars caps.
func WithYErrBarsCapWidth(w vg.Length) SampleOptions {
	return func(cfg *config) {
		cfg.YErrBarsCapWidth = w
	}
}

// WithDataStyle enables the default data histogram style.
func WithDataStyle(b bool) SampleOptions {
	return func(cfg *config) {
		cfg.DataStyle = b
	}
}

// WithSaveName sets the file name of the plot.
func WithSaveName(n string) VariableOptions {
	return func(cfg *config) {
		cfg.SaveName = n
	}
}

// WithTreeVar sets a TreeFunc object for an on-the-fly
// computed variable.
func WithTreeVar(f TreeFunc) VariableOptions {
	return func(cfg *config) {
		cfg.TreeVar = f
	}
}

// WithAxisLabels sets the x- and y-axis label.
func WithAxisLabels(xlab, ylab string) VariableOptions {
	return func(cfg *config) {
		cfg.XLabel = xlab
		cfg.YLabel = ylab
	}
}

// WithAxisLabels sets the x- and y-axis labels.
func WithTickFormats(xticks, yticks string) VariableOptions {
	return func(cfg *config) {
		cfg.XTickFormat = xticks
		cfg.YTickFormat = yticks
	}
}

// WithXRange sets the x-axis min and max.
func WithXRange(min, max float64) VariableOptions {
	return func(cfg *config) {
		cfg.RangeXmin = min
		cfg.RangeXmax = max
	}
}

// WithYRange sets the y-axis min and max.
func WithYRange(min, max float64) VariableOptions {
	return func(cfg *config) {
		cfg.RangeYmin = min
		cfg.RangeYmax = max
	}
}

// WithRatioYRange sets the y-axis min and max for the ratio plot.
func WithRatioYRange(min, max float64) VariableOptions {
	return func(cfg *config) {
		cfg.RatioYmin = min
		cfg.RatioYmax = max
	}
}

// WithLegLeft sets the legend left/right position on the plot.
func WithLegLeft(left bool) VariableOptions {
	return func(cfg *config) {
		cfg.LegPosLeft = left
	}
}

// WithLegTop sets the legend top/bottom position on the plot.
func WithLegTop(top bool) VariableOptions {
	return func(cfg *config) {
		cfg.LegPosTop = top
	}
}
