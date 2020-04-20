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


type Sample struct {
	Name string
	Type string
	FileName string
	TreeName string
	Weight string
	Cut string
	LegLabel string
	LineColor color.NRGBA
	LineWidth vg.Length
	FillColor color.NRGBA
	CircleMarkers bool
	CircleSize vg.Length
	CircleColor color.NRGBA
	WithYErrBars bool
	YErrBarsLineWidth vg.Length
	YErrBarsCapWidth vg.Length
}

// Return a hplot.H1D with the proper style
func (s Sample) CreateHisto(hdata *hbook.H1D) *hplot.H1D {

	// Create the plotable histo from histogrammed data
	h := hplot.NewH1D(hdata, hplot.WithYErrBars(s.WithYErrBars))

	// Line and fill cosmetics
	h.LineStyle.Width = s.LineWidth
	h.LineStyle.Color = s.LineColor
	h.FillColor = s.FillColor

	// Markers
	if s.CircleMarkers {
		style.SetCircleMarkersToHist(h)
		if &s.CircleColor != nil {
			h.GlyphStyle.Color = s.CircleColor
		} else {
			h.GlyphStyle.Color = s.LineColor
		}
		if &s.CircleSize != nil {
			h.GlyphStyle.Radius = s.CircleSize
		}
	}

	// Error bars
	if s.WithYErrBars {
		if &s.CircleColor != nil {
			h.YErrs.LineStyle.Color = s.CircleColor
		} else {
			h.YErrs.LineStyle.Color = s.LineColor
		}
		
		if &s.YErrBarsLineWidth != nil {
			h.YErrs.LineStyle.Width = s.YErrBarsLineWidth
		}

		if &s.YErrBarsCapWidth != nil {
			h.YErrs.CapWidth = s.YErrBarsCapWidth
		}
	}
	return h
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

