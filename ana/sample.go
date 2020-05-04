// Manage included samples
package ana

import (
	"strings"

	"image/color"

	"gonum.org/v1/plot/vg"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"

	"github.com/rmadar/hplot-style/style"
)

var colorNil = color.NRGBA{R: 0, G: 0, B: 0, A: 0}	

type Sample struct {
	Name              string
	Type              string
	FileName          string
	TreeName          string
	Weight            string
	WeightFunc        TreeFunc
	Cut               string
	CutFunc           TreeFunc
	LegLabel          string
	LineColor         color.NRGBA
	LineWidth         vg.Length
	FillColor         color.NRGBA
	CircleMarkers     bool
	CircleSize        vg.Length
	CircleColor       color.NRGBA
	WithYErrBars      bool
	YErrBarsLineWidth vg.Length
	YErrBarsCapWidth  vg.Length
}

// Return a hplot.H1D with the proper style
func (s Sample) CreateHisto(hdata *hbook.H1D, opts ...hplot.Options) *hplot.H1D {
	
	// Append sample-defined options
	opts = append(opts, hplot.WithYErrBars(s.WithYErrBars))
	
	// Create the plotable histo from histogrammed data
	h := hplot.NewH1D(hdata, opts...)

	// Line and fill cosmetics
	h.LineStyle.Width = s.LineWidth
	h.LineStyle.Color = s.LineColor
	h.FillColor = s.FillColor

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
	if s.WithYErrBars {
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
		s.SetBandStyle(h.Band)
	}
	
	return h
}

func (s Sample) SetBandStyle(b *hplot.Band) {

	if s.FillColor != colorNil {
		b.FillColor = s.FillColor
	}
	if s.LineColor != colorNil {
		b.FillColor = s.LineColor
	}
}

func (s *Sample) IsData() bool {
	return strings.ToLower(s.Type) == "data"
}

func (s *Sample) IsBkg() bool {
	return strings.ToLower(s.Type) == "bkg" ||
		strings.ToLower(s.Type) == "bg" ||
		strings.ToLower(s.Type) == "background"
}

func (s *Sample) IsSig() bool {
	return strings.ToLower(s.Type) == "sig" ||
		strings.ToLower(s.Type) == "signal"
}
