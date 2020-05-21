package ana

import (
	"strings"

	"image/color"

	"gonum.org/v1/plot/vg"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"

	"github.com/rmadar/hplot-style/style"
)

// Default color
var colorNil = color.NRGBA{R: 0, G: 0, B: 0, A: 0}

// Sample is the structure representing
type Sample struct {
	Name              string      // Sample name.
	Type              string      // Sample type: 'data', 'bkg' or 'sig'.
	LegLabel          string      // Label used in the legend.
	FileName          string      // Path the file of the sample.
	TreeName          string      // Name of the tree.
	WeightFunc        TreeFunc    // Weight to be applied to the sample (default: 1.0).
	CutFunc           TreeFunc    // Cut to be applied to the sample (default: no cut).
	DataStyle         bool        // Enable data-like style (default: Type == 'data')
	LineColor         color.NRGBA // Line color of the histogram (default: blue).
	LineWidth         vg.Length   // Line width of the histogram (default: 1.5).
	FillColor         color.NRGBA // Fill color of the histogram (default: none).
	CircleMarkers     bool        // Enable the use of circle markers (default: false).
	CircleSize        vg.Length   // Circle size (default: 0).
	CircleColor       color.NRGBA // Circle color (default: transparent).
	YErrBars          bool        // Display y-error bars (default: false).
	YErrBarsLineWidth vg.Length   // Width of y-error bars.
	YErrBarsCapWidth  vg.Length   // Width of horizontal bars of the y-error bars.
}

// NewSample creates a sample with the default settings.
func NewSample(sname, stype, sleg, fname, tname string, opts ...SampleOptions) *Sample {

	// Required fields
	s := &Sample{
		Name:     sname,
		Type:     stype,
		LegLabel: sleg,
		FileName: fname,
		TreeName: tname,
	}

	// Configuration with defaults values for all optional fields
	cfg := newConfig(
		WithFillColor(color.NRGBA{R: 20, G: 20, B: 180, A: 200}),
		WithDataStyle(s.IsData()),
	)

	// Update the configuration looping over functional options
	for _, opt := range opts {
		opt(cfg)
	}

	// Set all fields with the updated configuration
	s.WeightFunc = cfg.Weight
	s.CutFunc = cfg.Cut
	s.LineColor = cfg.LineColor
	s.LineWidth = cfg.LineWidth
	s.FillColor = cfg.FillColor
	s.CircleMarkers = cfg.CircleMarkers
	s.CircleSize = cfg.CircleSize
	s.CircleColor = cfg.CircleColor
	s.YErrBars = cfg.YErrBars
	s.YErrBarsLineWidth = cfg.YErrBarsLineWidth
	s.YErrBarsCapWidth = cfg.YErrBarsCapWidth
	s.DataStyle = cfg.DataStyle

	return s
}

// CreateHisto returns a hplot.H1D with the sample style.
func (s Sample) CreateHisto(hdata *hbook.H1D, opts ...hplot.Options) *hplot.H1D {

	// Append sample-defined options
	opts = append(opts, hplot.WithYErrBars(s.YErrBars))

	// Create the plotable histo from histogrammed data
	h := hplot.NewH1D(hdata, opts...)

	// Line width
	h.LineStyle.Width = s.LineWidth

	// Line color
	if s.LineColor != colorNil {
		h.LineStyle.Color = s.LineColor
	}

	// Fill color
	if s.FillColor != colorNil {
		h.FillColor = s.FillColor
	}

	// Markers
	if s.CircleMarkers {
		style.SetCircleMarkersToHist(h)
		if s.CircleColor != colorNil {
			h.GlyphStyle.Color = s.CircleColor
		} else {
			h.GlyphStyle.Color = s.LineColor
		}
		if s.CircleSize != 0.0 {
			h.GlyphStyle.Radius = s.CircleSize
		}
	}

	// Error bars
	if s.YErrBars {
		if s.CircleColor != colorNil {
			h.YErrs.LineStyle.Color = s.CircleColor
		} else {
			h.YErrs.LineStyle.Color = s.LineColor
		}

		if s.YErrBarsLineWidth != 0.0 {
			h.YErrs.LineStyle.Width = s.YErrBarsLineWidth
		}

		if s.YErrBarsCapWidth != 0.0 {
			h.YErrs.CapWidth = s.YErrBarsCapWidth
		}
	}

	// Band setup
	if h.Band != nil {
		s.setBandStyle(h.Band)
	}

	return h
}

// Helper function to set the error band style.
func (s Sample) setBandStyle(b *hplot.Band) {

	if s.FillColor != colorNil {
		b.FillColor = style.ChangeOpacity(s.FillColor, 150)
	}
	if s.LineColor != colorNil {
		b.FillColor = style.ChangeOpacity(s.LineColor, 150)
	}
}

// IsData returns true it the sample Type is 'data'.
func (s *Sample) IsData() bool {
	return strings.ToLower(s.Type) == "data"
}

// IsData returns true it the sample Type is 'background'.
func (s *Sample) IsBkg() bool {
	return strings.ToLower(s.Type) == "bkg" ||
		strings.ToLower(s.Type) == "bg" ||
		strings.ToLower(s.Type) == "background"
}

// IsData returns true it the sample Type is 'signal'.
func (s *Sample) IsSig() bool {
	return strings.ToLower(s.Type) == "sig" ||
		strings.ToLower(s.Type) == "signal"
}
