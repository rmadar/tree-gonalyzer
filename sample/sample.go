// Manage included samples
package sample

import (
	"image/color"

	"gonum.org/v1/plot/vg"	

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"

	"github.com/rmadar/hplot-style/style"
)


type Spl struct {
	FileName string
	TreeName string
	LegLabel string
	LineColor color.NRGBA
	LineWidth vg.Length
	FillColor color.NRGBA
	CircleMarkers bool
	CircleSize vg.Length
	WithYErrBars bool
}

// Return a hplot.H1D with the proper style
func (s Spl) CreateHisto(hdata *hbook.H1D) *hplot.H1D {

	// Create the plotable histo from histogrammed data
	h := hplot.NewH1D(hdata, hplot.WithYErrBars(s.WithYErrBars))

	// Line and fill cosmetics
	h.LineStyle.Width = s.LineWidth
	h.LineStyle.Color = s.LineColor
	h.FillColor = s.FillColor

	// Markers
	if s.CircleMarkers {
		style.SetCircleMarkersTo(h)
		h.GlyphStyle.Color = s.LineColor
		if s.CircleSize > 0 {
			h.GlyphStyle.Radius = s.CircleSize
		}
	}

	// Error bars
	if s.WithYErrBars {
		h.YErrs.LineStyle.Color = s.LineColor
		if s.LineWidth>0 {
			h.YErrs.LineStyle.Width = s.LineWidth
		}
	}
	return h
}



