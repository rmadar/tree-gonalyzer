package ana

import (
	"image/color"

	"gonum.org/v1/plot/vg"
)

// Options encodes various options to pass to ana types (maker, sample, variable).
type Options func(cfg *config)

// config contains all the possible options than can be enable or not.
type config struct {

	// ana.Maker options
	KinemCuts    []*Selection // List of cuts
	SavePath     string       // Path to which plot will be saved
	SaveFormat   string       // Extension of saved figure 'tex', 'pdf', 'png'
	CompileLatex bool         // Enable on-the-fly latex compilation of plots
	RatioPlot    bool         // Enable ratio plot
	DontStack    bool         // Disable histogram stacking (e.g. compare various processes)
	Normalize    bool         // Normalize distributions to unit area (when stacked, the total is normalized)
	ErrBandColor color.NRGBA  // Color for the uncertainty band.

	// Sample options
	Weight            TreeFunc    // Weight applied to the sample
	Cut               TreeFunc    // Cut applied to the sample
	LineColor         color.NRGBA // Line color of the sample histogram
	LineWidth         vg.Length   // Line width of the sample histogram
	FillColor         color.NRGBA // Fill color of the sample histogram
	CircleMarkers     bool        // Use of circled marker
	CircleSize        vg.Length   // Size of the markers
	CircleColor       color.NRGBA // Color of the markers
	YErrBars          bool        // Use of y error bars
	YErrBarsLineWidth vg.Length   // Line width of the y error bar
	YErrBarsCapWidth  vg.Length   // Width of the y error bar caps
	DataStyle         bool        // Use default data style histogram
}

// newConfig returns a config type with a set of passed options.
func newConfig(opts ...Options) *config {
	cfg := new(config)
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// WithWeight sets the weight to be used, as defined by the TreeFunc f.
func WithWeight(f TreeFunc) Options {
	return func(cfg *config) {
		cfg.Weight = f
	}
}

// WithCut sets the cut to be applied, as defined by the TreeFunc f.
func WithCut(f TreeFunc) Options {
	return func(cfg *config) {
		cfg.Cut = f
	}
}

// WithLineColor sets line color of the sample histogram.
func WithLineColor(c color.NRGBA) Options {
	return func(cfg *config) {
		cfg.LineColor = c
	}
}

// WithLineWidth sets line width of the sample histogram.
func WithLineWidth(w vg.Length) Options {
	return func(cfg *config) {
		cfg.LineWidth = w
	}

}

// WithFillColor sets the color with which the histo will be filled.
func WithFillColor(c color.NRGBA) Options {
	return func(cfg *config) {
		cfg.FillColor = c
	}
}

// WithCircleMarkers enables the use of circle markers (as for data histogram).
func WithCircleMarkers(b bool) Options {
	return func(cfg *config) {
		cfg.CircleMarkers = b
	}
}

// WithCircleSize sets the size of circle markers.
func WithCircleSize(s vg.Length) Options {
	return func(cfg *config) {
		cfg.CircleSize = s
	}
}

// WithCircleColor sets the color of circle markers.
func WithCircleColor(c color.NRGBA) Options {
	return func(cfg *config) {
		cfg.CircleColor = c
	}
}

// WithYErrBars enables y error bars.
func WithYErrBars(b bool) Options {
	return func(cfg *config) {
		cfg.YErrBars = b
	}
}

// WithYErrBarsLineWidth sets the width of the error bars line
func WithYErrBarsLineWidth(w vg.Length) Options {
	return func(cfg *config) {
		cfg.YErrBarsLineWidth = w
	}
}

// WithYErrBarsCapsWidth sets the width of the y error bars caps.
func WithYErrBarsCapWidth(w vg.Length) Options {
	return func(cfg *config) {
		cfg.YErrBarsCapWidth = w
	}
}

// WithDataStyle enables the default data histogram style.
func WithDataStyle(b bool) Options {
	return func(cfg *config) {
		cfg.DataStyle = b
	}
}

// WithKinemCuts sets the list of kinematic cuts to run on.
func WithKinemCuts(c []*Selection) Options {
	return func(cfg *config) {
		cfg.KinemCuts = c
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

// CompileLatex enables automatic latex compilation.
func WithCompileLatex(b bool) Options {
	return func(cfg *config) {
		cfg.CompileLatex = b
	}
}

// WithRatioPlot enables the ratio plot panel.
func WithRatioPlot(b bool) Options {
	return func(cfg *config) {
		cfg.RatioPlot = b
	}
}

// WithHistoStack enables histogram stacking for bkg-typed samples.
func WithHistoStack(b bool) Options {
	return func(cfg *config) {
		cfg.DontStack = !b
	}
}

// WithHistoNorm enables histogram normalization to unity, to compare shapes.
func WithHistoNorm(b bool) Options {
	return func(cfg *config) {
		cfg.Normalize = b
	}
}

// WithHistoNorm enables histogram normalization to unity, to compare shapes.
func WithErrBandColor(c color.NRGBA) Options {
	return func(cfg *config) {
		cfg.ErrBandColor = c
	}
}
