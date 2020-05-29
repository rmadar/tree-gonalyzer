package ana

import (
	"fmt"
	"log"
	"strings"

	"image/color"

	"gonum.org/v1/plot/vg"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"

	"github.com/rmadar/hplot-style/style"
)

// Default color
var colorNil = color.NRGBA{R: 0, G: 0, B: 0, A: 0}

// Sample type
type sampleType int

const (
	data sampleType = iota
	bkg
	sig
)

// Sample contains all the information defining a single histogram
// of the final plot. A sample is made of (potentially) several
// components, or sub-sample. Each component can have different
// file/tree names, as well as additional cuts and weights,
// on top of the global ones. Concretly, the global cut is combined
// with the component cut with a AND, while weights are multiplied.
type Sample struct {

	// General settings
	Name       string             // Sample name.
	Type       string             // Sample type: 'data', 'bkg' or 'sig'.
	LegLabel   string             // Label used in the legend.
	Components []*SampleComponent // List of components included in the histogram.

	// Gobal weight and cut to be applied to
	// all components.
	CutFunc    TreeFunc
	WeightFunc TreeFunc

	// Cosmetic settings
	DataStyle         bool        // Enable data-like style (default: Type == 'data').
	LineColor         color.NRGBA // Line color of the histogram (default: blue).
	LineWidth         vg.Length   // Line width of the histogram (default: 1.5).
	LineDashes        []vg.Length // Line dashes format (default: continous).
	FillColor         color.NRGBA // Fill color of the histogram (default: none).
	CircleMarkers     bool        // Enable circle markers (default: false).
	CircleSize        vg.Length   // Circle size (default: 0).
	CircleColor       color.NRGBA // Circle color (default: transparent).
	Band              bool        // Enable error band display.
	YErrBars          bool        // Display error bars (default: false || DataStyle).
	YErrBarsLineWidth vg.Length   // Width of error bars line.
	YErrBarsCapWidth  vg.Length   // Width of horizontal bars of the error bars.

	// Internal
	sType sampleType
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

// NewSample creates a new empty sample, ie without any components,
// with the default options. Components can be then added using
// s.AddComponent(...) function.
func NewSample(sname, stype, sleg string, opts ...SampleOptions) *Sample {

	// Check/set sample type
	var sType sampleType
	switch strings.ToLower(stype) {
	case "data":
		sType = data
	case "background", "bkg", "bg":
		sType = bkg
	case "signal", "sig", "sg":
		sType = sig
	default:
		err := "Sample type not %v supported [%v]"
		log.Fatal(fmt.Sprint(err, stype, sname))
	}

	// Empty basic sample
	s := &Sample{
		Name:       sname,
		Type:       stype,
		LegLabel:   sleg,
		Components: []*SampleComponent{},
		sType:      sType,
	}

	// Configuration with defaults values for all optional fields
	cfg := newConfig(
		WithLineColor(color.NRGBA{R: 20, G: 20, B: 180, A: 255}),
		WithDataStyle(s.IsData()),
	)

	// Update the configuration looping over functional options
	for _, opt := range opts {
		opt(cfg)
	}

	// Set setting with the updated configuration
	s.WeightFunc = cfg.WeightFunc
	s.CutFunc = cfg.CutFunc
	s.LineColor = cfg.LineColor
	s.LineWidth = cfg.LineWidth
	s.LineDashes = cfg.LineDashes
	s.FillColor = cfg.FillColor
	s.CircleMarkers = cfg.CircleMarkers
	s.CircleSize = cfg.CircleSize
	s.CircleColor = cfg.CircleColor
	s.Band = cfg.Band
	s.YErrBars = cfg.YErrBars || cfg.DataStyle
	s.YErrBarsLineWidth = cfg.YErrBarsLineWidth
	s.YErrBarsCapWidth = cfg.YErrBarsCapWidth
	s.DataStyle = cfg.DataStyle

	return s
}

// CreateSample creates a non-empty sample having the default settings,
// with only one component. This function is a friendly API to
// ease single-component samples declaration. For multi-component samples,
// one can either add components on top with s.AddComponent(...), or start
// from an empty sample using NewSample(...) followed by few s.AddComponent(...).
func CreateSample(sname, stype, sleg, fname, tname string, opts ...SampleOptions) *Sample {

	// New empty sample
	s := NewSample(sname, stype, sleg, opts...)

	// Configuration
	cfg := newConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	// Add it to component the sample
	s.AddComponent(fname, tname)

	return s
}

// AddComponent adds a new sample component to the sample.
func (s *Sample) AddComponent(fname, tname string, opts ...SampleOptions) {

	// Manage default settings and passed options
	// FIXME(rmadar): most of SampleOption doesn't change a component.
	//                consider adding a protection against the ones which
	//                doesn't change the behaviour of the component?
	cfg := newConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	// Create a component
	c := &SampleComponent{
		FileName:   fname,
		TreeName:   tname,
		WeightFunc: cfg.WeightFunc,
		CutFunc:    cfg.CutFunc,
	}

	// Append it to the pointer-receiver sample
	s.Components = append(s.Components, c)
}

// CreateHisto returns a hplot.H1D with the sample style.
func (s Sample) CreateHisto(hdata *hbook.H1D, opts ...hplot.Options) *hplot.H1D {

	// Append sample-defined options
	opts = append(opts, hplot.WithYErrBars(s.YErrBars))
	opts = append(opts, hplot.WithBand(s.Band))

	// Create the plotable histo from histogrammed data
	h := hplot.NewH1D(hdata, opts...)

	// Line width
	h.LineStyle.Width = s.LineWidth

	// Line color
	if s.LineColor != colorNil {
		h.LineStyle.Color = s.LineColor
	}

	// Line dashes
	if len(s.LineDashes) > 0 {
		h.LineStyle.Dashes = s.LineDashes
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
	return s.sType == data
}

// IsData returns true it the sample Type is 'background'.
func (s *Sample) IsBkg() bool {
	return s.sType == bkg
}

// IsData returns true it the sample Type is 'signal'.
func (s *Sample) IsSig() bool {
	return s.sType == sig
}
