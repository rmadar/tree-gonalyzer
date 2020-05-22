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

// Sample contains all the information defining a single histogram
// of the final plot.
type Sample struct {

	// General settings
	Name       string             // Sample name.
	Type       string             // Sample type: 'data', 'bkg' or 'sig'.
	LegLabel   string             // Label used in the legend.
	Components []*SampleComponent // List of components included in the histogram.

	// Cosmetic settings
	DataStyle         bool        // Enable data-like style (default: Type == 'data').
	LineColor         color.NRGBA // Line color of the histogram (default: blue).
	LineWidth         vg.Length   // Line width of the histogram (default: 1.5).
	FillColor         color.NRGBA // Fill color of the histogram (default: none).
	CircleMarkers     bool        // Enable circle markers (default: false).
	CircleSize        vg.Length   // Circle size (default: 0).
	CircleColor       color.NRGBA // Circle color (default: transparent).
	YErrBars          bool        // Display y-error bars (default: false).
	YErrBarsLineWidth vg.Length   // Width of y-error bars.
	YErrBarsCapWidth  vg.Length   // Width of horizontal bars of the y-error bars.
}

// SampleComponent contains the needed information
// to fill the final histogram for a given component
// (or sub-sample). This includes a file, a tree and
// possibly a specific weight and cut.
type SampleComponent struct {
	FileName   string
	TreeName   string
	WeightFunc TreeFunc
	CutFunc    TreeFunc
}

// NewSample creates a sample with one sub-sample based
// on the default settings. This function is to be used
// for single-component samples.
func NewSample(sname, stype, sleg, fname, tname string, opts ...SampleOptions) *Sample {

	// New empty sample
	s := NewEmptySample(sname, stype, sleg, opts...)

	// Configuration
	cfg := newConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	
	// Create a component
	c := &SampleComponent{
		FileName:   fname,
		TreeName:   tname,
		WeightFunc: cfg.Weight,
		CutFunc:    cfg.Cut,
	}

	// Append it to the pointer-receiver sample
	s.Components = append(s.Components, c)
	
	return s
}

// NewEmptySample creates a new sample without any components.
// This function is to be favoured in case of several sub-samples,
// than can be added using s.AddComponent().
func NewEmptySample(sname, stype, sleg string, opts ...SampleOptions) *Sample {

	// Empty basic sample
	s := &Sample{
		Name:     sname,
		Type:     stype,
		LegLabel: sleg,
		Components: []*SampleComponent{},
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

	// Set cosmetic setting with the updated configuration
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

// AddComponent adds a new sample component to the sample.
// In case the sample is created using NewSample() function, default weights
// and cuts are the same as defined as the one passed to NewSample()
// function. They will be overwritten if new weights/cuts are passed to AddComponent().
//
// In order to avoid confusion, it is discouraged to use
// this function with a sample created via s := NewSample(). Instead,
// the use of the function NewEmptySample() is encouraged, where weight
// and cut of each component is explicitely given via s.AddComponent().
// 
// It should be possible to finding a better.
func (s *Sample) AddComponent(fname, tname string, opts ...SampleOptions) {

	// Manage default settings and passed options
	// FIXME(rmadar): most of SampleOption doesn't change a component.
	//                consider adding a protection against the ones whic
	//                doesn't change the behaviour of the component?
	cfg := newConfig(
		WithWeight(s.Components[0].WeightFunc),
		WithCut(s.Components[0].CutFunc),
	)
	for _, opt := range opts {
		opt(cfg)
	}

	// Create a component
	c := &SampleComponent{
		FileName:   fname,
		TreeName:   tname,
		WeightFunc: cfg.Weight,
		CutFunc:    cfg.Cut,
	}

	// Append it to the pointer-receiver sample
	s.Components = append(s.Components, c)
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
