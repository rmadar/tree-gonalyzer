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
	KinemCuts struct {
		val []*Selection // List of cuts.
		usr bool
	}
	NevtsMax struct {
		val int64 // Maximum of processed events.
		usr bool
	}
	Lumi struct {
		val float64 // Integrated luminosity.
		usr bool
	}
	SampleMT struct {
		val bool // Enable concurency over samples.
		usr bool
	}
	SavePath struct {
		val string // Path to which plot will be saved.
		usr bool
	}
	SaveFormat struct {
		val string // Extension of saved figure 'tex', 'pdf', 'png'.
		usr bool
	}
	CompileLatex struct {
		val bool // Enable on-the-fly latex compilation of plots.
		usr bool
	}
	DumpTree struct {
		val bool // Enable Tree dumping
		usr bool
	}
	PlotHisto struct {
		val bool // Enable histograms plotting
		usr bool
	}
	AutoStyle struct {
		val bool // Enable auto style of histograms.
		usr bool
	}
	PlotTitle struct {
		val string // General plot title.
		usr bool
	}
	RatioPlot struct {
		val bool // Enable ratio plot.
		usr bool
	}
	HistoStack struct {
		val bool // Disable histogram stacking (e.g. compare various processes).
		usr bool
	}
	HistoNorm struct {
		val bool // Normalize distributions to unit area (when stacked, the total is normalized).
		usr bool
	}
	SignalStack struct {
		val bool // Enable signal histogram stacking.
		usr bool
	}
	TotalBand struct {
		val bool // Enable total error band for stacks.
		usr bool
	}
	ErrBandColor struct {
		val color.NRGBA // Color for the uncertainty band.
		usr bool
	}

	// Sample options
	WeightFunc struct {
		val TreeFunc // Weight applied to the sample/component
		usr bool
	}
	CutFunc struct {
		val TreeFunc // Cut applied to the sample/component
		usr bool
	}
	JointTrees struct {
		val []input // slice of file/tree name of joints trees
		usr bool
	}
	Xsec struct {
		val float64 // cross-section of this sample/component
		usr bool
	}
	Ngen struct {
		val float64 // Number of (weighted) generated events
		usr bool
	}
	LineColor struct { // Line color of the sample histogram
		val color.NRGBA
		usr bool
	}
	LineWidth struct {
		val vg.Length // Line width of the sample histogram
		usr bool
	}
	LineDashes struct {
		val []vg.Length // Line dashes format
		usr bool
	}
	FillColor struct {
		val color.NRGBA // Fill color of the sample histogram
		usr bool
	}
	CircleMarkers struct {
		val bool // Use of circled marker
		usr bool
	}
	CircleSize struct {
		val vg.Length // Size of the markers
		usr bool
	}
	CircleColor struct {
		val color.NRGBA // Color of the markers
		usr bool
	}
	Band struct {
		val bool // Enable error band display
		usr bool
	}
	YErrBars struct {
		val bool // Use of y error bars
		usr bool
	}
	YErrBarsLineWidth struct {
		val vg.Length // Line width of the y error bar
		usr bool
	}
	YErrBarsCapWidth struct {
		val vg.Length // Width of the y error bar caps
		usr bool
	}
	DataStyle struct {
		val bool // Use default data style histogram
		usr bool
	}

	// Variable options
	SaveName struct {
		val string // Filename for a variable plot
		usr bool
	}
	TreeVar struct {
		val TreeFunc // TreeFunc for a computed variable (ie not a single branch)
		usr bool
	}
	XLabel struct {
		val string
		usr bool
	}
	YLabel struct {
		val string // Axis labels
		usr bool
	}
	XTickFormat struct {
		val string
		usr bool
	}
	YTickFormat struct {
		val string // Ticks formating
		usr bool
	}
	RangeXmin struct {
		val float64
		usr bool
	}
	RangeXmax struct {
		val float64 // X-axis ranges
		usr bool
	}
	RangeYmin struct {
		val float64 // X-axis ranges
		usr bool
	}
	RangeYmax struct {
		val float64 // X-axis ranges
		usr bool
	}
	RatioYmax struct {
		val float64 // Y-axis ranges
		usr bool
	}
	RatioYmin struct {
		val float64 // Y-axis ranges
		usr bool
	}
	LegPosTop struct {
		val bool
		usr bool
	}
	LegPosLeft struct {
		val bool // Legend position
		usr bool
	}
}

// newConfig returns a config type with a set of passed options.
func newConfig() *config {
	cfg := new(config)
	return cfg
}

// WithKinemCuts sets the list of kinematic cuts to run on.
func WithKinemCuts(c []*Selection) Options {
	return func(cfg *config) {
		cfg.KinemCuts.val = c
		cfg.KinemCuts.usr = true
	}
}

// WithNevtsMax sets the maximum processed event for
// each sample component.
func WithNevtsMax(n int64) Options {
	return func(cfg *config) {
		cfg.NevtsMax.val = n
		cfg.NevtsMax.usr = true

	}
}

// WithLumi sets the integrated luminosity [1/fb]
// The full normalisation factor is (xsec*lumi)/ngen. 'ngen' and
// 'xsec' are given for each sample/component, via ana.CreateSample()
// or s.AddComponent(). While the 'lumi' is given to
// via ana.New(). By default, lumi = 1/fb.
func WithLumi(l float64) Options {
	return func(cfg *config) {
		cfg.Lumi.val = l
		cfg.Lumi.usr = true
	}
}

// WithSampleMT enables a concurent sample processing.
func WithSampleMT(b bool) Options {
	return func(cfg *config) {
		cfg.SampleMT.val = b
		cfg.SampleMT.usr = true
	}
}

// WithSavePath sets the path to save plots.
func WithSavePath(p string) Options {
	return func(cfg *config) {
		cfg.SavePath.val = p
		cfg.SavePath.usr = true
	}
}

// WithSaveFormat sets the format for the plots.
func WithSaveFormat(f string) Options {
	return func(cfg *config) {
		cfg.SaveFormat.val = f
		cfg.SaveFormat.usr = true
	}
}

// WithCompileLatex enables automatic latex compilation.
func WithCompileLatex(b bool) Options {
	return func(cfg *config) {
		cfg.CompileLatex.val = b
		cfg.CompileLatex.usr = true
	}
}

// WithDumpTree enables tree dumping (one per sample).
func WithDumpTree(b bool) Options {
	return func(cfg *config) {
		cfg.DumpTree.val = b
		cfg.DumpTree.usr = true
	}
}

// WithPlotHisto enables histogram plotting. It can be
// set to false to only dump trees.
func WithPlotHisto(b bool) Options {
	return func(cfg *config) {
		cfg.PlotHisto.val = b
		cfg.PlotHisto.usr = true
	}
}

// WithAutoStyle enables automatic styling of the histograms.
func WithAutoStyle(b bool) Options {
	return func(cfg *config) {
		cfg.AutoStyle.val = b
		cfg.AutoStyle.usr = true
	}
}

// WithPlotTitle sets the general plot title.
func WithPlotTitle(t string) Options {
	return func(cfg *config) {
		cfg.PlotTitle.val = t
		cfg.PlotTitle.usr = true
	}
}

// WithRatioPlot enables the ratio plot panel.
func WithRatioPlot(b bool) Options {
	return func(cfg *config) {
		cfg.RatioPlot.val = b
		cfg.RatioPlot.usr = true
	}
}

// WithHistoStack enables histogram stacking for
// bkg-typed samples.
func WithHistoStack(b bool) Options {
	return func(cfg *config) {
		cfg.HistoStack.val = b
		cfg.HistoStack.usr = true
	}
}

// WithSignalStack enables histogram stacking for
// bkg-typed samples.
func WithSignalStack(b bool) Options {
	return func(cfg *config) {
		cfg.SignalStack.val = b
		cfg.SignalStack.usr = true
	}
}

// WithHistoStack enables histogram stacking for
// bkg-typed samples.
func WithTotalBand(b bool) Options {
	return func(cfg *config) {
		cfg.TotalBand.val = b
		cfg.TotalBand.usr = true
	}
}

// WithHistoNorm enables histogram normalization to unity.
func WithHistoNorm(b bool) Options {
	return func(cfg *config) {
		cfg.HistoNorm.val = b
		cfg.HistoNorm.usr = true
	}
}

// WithErrBandColor sets the color for the error band of
// total histogram (and ratio).
func WithErrBandColor(c color.NRGBA) Options {
	return func(cfg *config) {
		cfg.ErrBandColor.val = c
		cfg.ErrBandColor.usr = true
	}
}

// WithWeight sets the weight to be used for this sample,
// as defined by the TreeFunc f, which must return a float64.
// Maker.FillHisto() will panic otherwise.
func WithWeight(f TreeFunc) SampleOptions {
	return func(cfg *config) {
		cfg.WeightFunc.val = f
		cfg.WeightFunc.usr = true
	}
}

// WithCut sets the cut to be applied to the sample, as
// defined by the TreeFunc f, which must return a bool.
// Maker.FillHisto() will panic otherwise.
func WithCut(f TreeFunc) SampleOptions {
	return func(cfg *config) {
		cfg.CutFunc.val = f
		cfg.CutFunc.usr = true
	}
}

// WithJointTree adds a tree to be joint for this sample/component.
// Branches of joint tree are added to the list of available variables.
// Several joint trees can be added using  WithJointTree() option several
// times.
func WithJointTree(fname, tname string) SampleOptions {
	return func(cfg *config) {
		cfg.JointTrees.val = append(cfg.JointTrees.val, input{FileName: fname, TreeName: tname})
		cfg.JointTrees.usr = true
	}
}

// WithXsect sets the cross-section [pb] to the sample/component.
// The full normalisation factor is (xsec*lumi)/ngen. 'ngen' and
// 'xsec' are given by sample/component while 'lumi' is given to
// via ana.New(). By default, xsec = 1 pb. This option cannot be passed to
// a sample (or component) of type "data".
func WithXsec(s float64) SampleOptions {
	return func(cfg *config) {
		cfg.Xsec.val = s
		cfg.Xsec.usr = true
	}
}

// WithNgen sets the total number of generated events for the
// sample/component. The full normalisation factor is (xsec*lumi)/ngen.
// 'ngen' and 'xsec' are given by sample/component while 'lumi' is given to
// via ana.New(). By default, Ngen = 1. This option cannot be passed to
// a sample (or component) of type "data".
func WithNgen(n float64) SampleOptions {
	return func(cfg *config) {
		cfg.Ngen.val = n
		cfg.Ngen.usr = true
	}
}

// WithLineColor sets the line color of the histogram.
func WithLineColor(c color.NRGBA) SampleOptions {
	return func(cfg *config) {
		cfg.LineColor.val = c
		cfg.LineColor.usr = true
	}
}

// WithLineWidth sets line width of the sample histogram.
func WithLineWidth(w vg.Length) SampleOptions {
	return func(cfg *config) {
		cfg.LineWidth.val = w
		cfg.LineWidth.usr = true
	}

}

// WithLineWidth sets line width of the sample histogram.
func WithLineDashes(s []vg.Length) SampleOptions {
	return func(cfg *config) {
		cfg.LineDashes.val = s
		cfg.LineDashes.usr = true
	}

}

// WithFillColor sets the color with which the histo will be filled.
func WithFillColor(c color.NRGBA) SampleOptions {
	return func(cfg *config) {
		cfg.FillColor.val = c
		cfg.FillColor.usr = true
	}
}

// WithCircleMarkers enables the use of circle markers (as for data histogram).
func WithCircleMarkers(b bool) SampleOptions {
	return func(cfg *config) {
		cfg.CircleMarkers.val = b
		cfg.CircleMarkers.usr = true
	}
}

// WithCircleSize sets the size of circle markers.
func WithCircleSize(s vg.Length) SampleOptions {
	return func(cfg *config) {
		cfg.CircleSize.val = s
		cfg.CircleSize.usr = true
	}
}

// WithCircleColor sets the color of circle markers.
func WithCircleColor(c color.NRGBA) SampleOptions {
	return func(cfg *config) {
		cfg.CircleColor.val = c
		cfg.CircleColor.usr = true
	}
}

// WithBand enables y error band.
func WithBand(b bool) SampleOptions {
	return func(cfg *config) {
		cfg.Band.val = b
		cfg.Band.usr = true
	}
}

// WithYErrBars enables y error bars.
func WithYErrBars(b bool) SampleOptions {
	return func(cfg *config) {
		cfg.YErrBars.val = b
		cfg.YErrBars.usr = true
	}
}

// WithYErrBarsLineWidth sets the width of the error bars line
func WithYErrBarsLineWidth(w vg.Length) SampleOptions {
	return func(cfg *config) {
		cfg.YErrBarsLineWidth.val = w
		cfg.YErrBarsLineWidth.usr = true
	}
}

// WithYErrBarsCapsWidth sets the width of the y error bars caps.
func WithYErrBarsCapWidth(w vg.Length) SampleOptions {
	return func(cfg *config) {
		cfg.YErrBarsCapWidth.val = w
		cfg.YErrBarsCapWidth.usr = true
	}
}

// WithDataStyle enables the default data histogram style.
func WithDataStyle(b bool) SampleOptions {
	return func(cfg *config) {
		cfg.DataStyle.val = b
		cfg.DataStyle.usr = true
	}
}

// WithSaveName sets the file name of the plot.
func WithSaveName(n string) VariableOptions {
	return func(cfg *config) {
		cfg.SaveName.val = n
		cfg.SaveName.usr = true
	}
}

// WithTreeVar sets a TreeFunc object for an on-the-fly
// computed variable.
func WithTreeVar(f TreeFunc) VariableOptions {
	return func(cfg *config) {
		cfg.TreeVar.val = f
		cfg.TreeVar.usr = true
	}
}

// WithAxisLabels sets the x- and y-axis label.
func WithAxisLabels(xlab, ylab string) VariableOptions {
	return func(cfg *config) {
		cfg.XLabel.val = xlab
		cfg.XLabel.usr = true
		cfg.YLabel.val = ylab
		cfg.YLabel.usr = true
	}
}

// WithAxisLabels sets the x- and y-axis labels.
func WithTickFormats(xticks, yticks string) VariableOptions {
	return func(cfg *config) {
		cfg.XTickFormat.val = xticks
		cfg.XTickFormat.usr = true
		cfg.YTickFormat.val = yticks
		cfg.YTickFormat.usr = true
	}
}

// WithXRange sets the x-axis min and max.
func WithXRange(min, max float64) VariableOptions {
	return func(cfg *config) {
		cfg.RangeXmin.val = min
		cfg.RangeXmin.usr = true
		cfg.RangeXmax.val = max
		cfg.RangeXmax.usr = true
	}
}

// WithYRange sets the y-axis min and max.
func WithYRange(min, max float64) VariableOptions {
	return func(cfg *config) {
		cfg.RangeYmin.val = min
		cfg.RangeYmin.usr = true
		cfg.RangeYmax.val = max
		cfg.RangeYmax.usr = true
	}
}

// WithRatioYRange sets the y-axis min and max for the ratio plot.
func WithRatioYRange(min, max float64) VariableOptions {
	return func(cfg *config) {
		cfg.RatioYmin.val = min
		cfg.RatioYmin.usr = true
		cfg.RatioYmax.val = max
		cfg.RatioYmax.usr = true
	}
}

// WithLegLeft sets the legend left/right position on the plot.
func WithLegLeft(left bool) VariableOptions {
	return func(cfg *config) {
		cfg.LegPosLeft.val = left
		cfg.LegPosLeft.usr = true
	}
}

// WithLegTop sets the legend top/bottom position on the plot.
func WithLegTop(top bool) VariableOptions {
	return func(cfg *config) {
		cfg.LegPosTop.val = top
		cfg.LegPosTop.usr = true
	}
}
