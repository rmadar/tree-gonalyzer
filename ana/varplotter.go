package ana

import (
	"image/color"
	"log"
	"os"
	"sync"
	"time"
	//"fmt"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/hplot/htex"

	"github.com/rmadar/hplot-style/style"
)

// PlotVariables loops over all filled histograms and produce one plot
// for each variable and selection, including all sample histograms.
func (ana *Maker) PlotVariables() error {

	if !ana.PlotHisto {
		return nil
	}

	// Start timing
	start := time.Now()

	// Set histogram styles
	if ana.AutoStyle {
		ana.setAutoStyle()
	}

	// Return an error if HbookHistos is empty
	if !ana.histoFilled {
		err := "There is no histograms. Please make sure that"
		err += "'FillHistos()' is called before 'PlotVariables()'"
		log.Fatalf(err)
	}

	// Compute all normalizations beforehand
	ana.normHists, ana.normTotal = ana.Normalizations()

	// Handle on-the-fly LaTeX compilation
	var latex htex.Handler = htex.NoopHandler{}
	if ana.CompileLatex {
		latex = htex.NewGoHandler(-1, "pdflatex")
	}

	// Loop over variables and cuts
	var wg sync.WaitGroup
	wg.Add(len(ana.Variables) * len(ana.KinemCuts))
	for iv := range ana.Variables {
		for ic := range ana.KinemCuts {
			go ana.concurrentPlotVar(iv, ic, latex, &wg)
		}
	}
	wg.Wait()

	// Handle latex compilation
	if latex, ok := latex.(*htex.GoHandler); ok {
		if err := latex.Wait(); err != nil {
			log.Fatalf("could not compiler latex document(s): %+v", err)
		}
	}

	// End timing
	ana.timePlot = time.Since(start)

	return nil
}

func (ana *Maker) concurrentPlotVar(iVar, iCut int, latex htex.Handler, wg *sync.WaitGroup) {

	// Handle concurrency
	defer wg.Done()

	// Fill the histo
	ana.plotVar(iVar, iCut, latex)
}

func (ana *Maker) plotVar(iVar, iCut int, latex htex.Handler) {

	// Current variable
	v := ana.Variables[iVar]

	var (
		drw       hplot.Drawer
		plt       = hplot.New()
		figWidth  = 6 * vg.Inch
		figHeight = 4.5 * vg.Inch
	)

	// Post-normalization hbook histograms
	bhistos := ana.getNormHbookHistos(iCut, iVar)

	// hplot histograms
	phistos := ana.getHplotH1D(bhistos, v.LogY)

	// Stack signal/bkg histograms
	phBkgs := hplotHistoFromIdx(phistos, ana.idxBkgs)
	phSigs := hplotHistoFromIdx(phistos, ana.idxSigs)
	stack := ana.stackHistograms(phBkgs, phSigs, v.LogY)

	// Data
	phData := hplotHistoFromIdx(phistos, ana.idxData)

	// Add histograms to the legend
	for i, s := range ana.Samples {
		plt.Legend.Add(s.LegLabel, phistos[i])
	}

	// Add total error band to the legend
	if ana.HistoStack && ana.TotalBand {
		hBand := hplot.NewH1D(hbook.NewH1D(1, 0, 1), hplot.WithBand(true))
		hBand.Band = stack.Band
		hBand.Band.FillColor = ana.TotalBandColor
		hBand.LineStyle.Width = 0
		plt.Legend.Add("Uncer.", hBand)
	}

	// Add histogram and stacks to the plot
	if stack != nil {
		plt.Add(stack)
	}
	if !ana.SignalStack {
		for _, hs := range phSigs {
			plt.Add(hs)
		}
	}
	if len(phData) > 0 {
		plt.Add(phData[0])
	}

	// Apply common and user-defined style for this variable
	plt.Title.Text = ana.PlotTitle
	style.ApplyToPlot(plt)
	v.setPlotStyle(plt)
	if v.LogY {
		plt.Y.Scale = plot.LogScale{}
		plt.Y.Tick.Marker = plot.LogTicks{}
	}
	drw = plt

	// Addition of the ratio plot
	if ana.RatioPlot {

		// Create a ratio plot and style it using plt
		rp := hplot.NewRatioPlot()
		style.ApplyToRatioPlot(rp, plt)

		// Update the drawer and figure size
		figWidth, figHeight = 6*vg.Inch, 4.5*vg.Inch
		drw = rp

		// Compute and add ratios to the plot
		ana.addRatioToPlot(rp, bhistos, phistos)

		// Adjust ratio plot scale
		if v.RatioYmin != v.RatioYmax {
			rp.Bottom.Y.Min = v.RatioYmin
			rp.Bottom.Y.Max = v.RatioYmax
		}
	}

	// Create the figure
	f := hplot.Figure(drw)
	style.ApplyToFigure(f)
	f.Latex = latex

	// Save the figure
	path := ana.SavePath + "/" + ana.KinemCuts[iCut].Name
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
	outputname := path + "/" + v.SaveName + "." + ana.SaveFormat
	if err := hplot.Save(f, figWidth, figHeight, outputname); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}

// Helper function to setup the automatic style.
func (ana *Maker) setAutoStyle() {

	ic := 0
	for _, s := range ana.Samples {

		// Color
		r, g, b, a := plotutil.Color(ic).RGBA()
		c := color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}

		switch s.sType {
		case data:
			s.DataStyle = true
		case bkg:
			// Fill for stacked histo, lines otherwise
			if ana.HistoStack {
				s.FillColor = c
				s.LineWidth = 0.
			} else {
				s.FillColor = color.NRGBA{}
				s.LineColor = c
				s.LineWidth = 2.
			}
			ic += 1
		case sig:
			// Fill for stacked histo, lines otherwise
			if ana.SignalStack {
				s.FillColor = c
				s.LineWidth = 0.
			} else {
				s.FillColor = color.NRGBA{}
				s.LineColor = c
				s.LineWidth = 2.
			}
			ic += 1
		}

		// Apply user-defined setting on top of default ones.
		s.applyConfig()
	}
}

// Helper function computing the normalisation of
// of all samples for a given cut
func (ana *Maker) Normalizations() ([][]float64, []float64) {

	// Initialization
	nTot := make([]float64, len(ana.KinemCuts))
	norms := make([][]float64, len(ana.KinemCuts))
	for i := range norms {
		norms[i] = make([]float64, len(ana.Samples))
	}

	// If no normalization is needed, compute nothing.
	if !ana.HistoNorm {
		for ic := range ana.KinemCuts {
			nTot[ic] = 1.0
			for is := range ana.Samples {
				norms[ic][is] = 1.0
			}
		}
		return norms, nTot
	}

	// Otherwise, loop over cuts and samples.
	for ic, _ := range ana.KinemCuts {
		for is, s := range ana.Samples {

			// Individual normalization including under/over-flows
			n := ana.HbookHistos[is][ic][0].Integral()
			norms[ic][is] = n

			// Cumulate backgrounds for the total
			if s.IsBkg() {
				nTot[ic] += n
			}

			// Cumulate signals for the total, it stacked
			if s.IsSig() && ana.SignalStack {
				nTot[ic] += n
			}
		}
	}

	return norms, nTot
}

// Helper function to normalize and return hbook histograms
// of a given cut  and variable, for bkgs, sigs and data.
func (ana *Maker) getNormHbookHistos(iCut, iVar int) []*hbook.H1D {

	// Get normalization
	nHistos, nTot := ana.normHists[iCut], ana.normTotal[iCut]

	// Prepare histo maps
	bhistos := make([]*hbook.H1D, len(ana.Samples))

	// Loop over sample
	for i, s := range ana.Samples {

		// Get a clone of the histo
		h := ana.HbookHistos[i][iCut][iVar].Clone()

		// Normalize
		if ana.HistoNorm {
			switch s.sType {
			case data:
				h.Scale(1 / nHistos[i])
			case bkg, sig:
				if ana.HistoStack {
					h.Scale(1. / nTot)
				} else {
					h.Scale(1. / nHistos[i])
				}
			}
		}

		// Store
		bhistos[i] = h
	}

	return bhistos
}

// Helper function to get hplot histograms from hbook histograms.
// It returns two maps [string]hplot.H1D (bkgs), [string]hplot.H1D (sigs)
// and one histogram hplot.H1D (data).
func (ana *Maker) getHplotH1D(hs []*hbook.H1D, LogY bool) []*hplot.H1D {

	// Prepare the output
	phistos := make([]*hplot.H1D, len(ana.Samples))

	// Loop over backgrounds
	for _, ib := range ana.idxBkgs {
		s := ana.Samples[ib]
		phistos[ib] = s.CreateHisto(hs[ib], hplot.WithLogY(LogY))
	}

	// Loop over signals
	for _, is := range ana.idxSigs {
		s := ana.Samples[is]
		phistos[is] = s.CreateHisto(hs[is], hplot.WithLogY(LogY))
	}

	// Loop over signals
	for _, id := range ana.idxData {
		s := ana.Samples[id]
		phistos[id] = s.CreateHisto(hs[id], hplot.WithLogY(LogY))
		if s.DataStyle {
			style.ApplyToDataHist(phistos[id])
			if s.CircleSize > 0 {
				phistos[id].GlyphStyle.Radius = s.CircleSize
			}
			if s.YErrBarsLineWidth > 0 {
				phistos[id].YErrs.LineStyle.Width = s.YErrBarsLineWidth
			}
			if s.YErrBarsCapWidth > 0 {
				phistos[id].YErrs.CapWidth = s.YErrBarsCapWidth
			}
		}
	}

	return phistos
}

// Helper function to stack histograms
func (ana *Maker) stackHistograms(hBkgs, hSigs []*hplot.H1D, LogY bool) *hplot.HStack {

	if len(hBkgs)+len(hSigs) == 0 {
		return nil
	}

	// Put all backgrounds in the stack
	phStack := []*hplot.H1D{}
	for _, b := range hBkgs {
		phStack = append(phStack, b)
	}

	// Reverse the order so that legend and plot order matches
	for i, j := 0, len(phStack)-1; i < j; i, j = i+1, j-1 {
		phStack[i], phStack[j] = phStack[j], phStack[i]
	}

	// Add signals if asked (after the order reversering to have
	// signals on top of the bkg).
	if ana.SignalStack {
		for _, hs := range hSigs {
			phStack = append(phStack, hs)
		}
	}

	// Stacking the background histo
	stack := hplot.NewHStack(phStack, hplot.WithBand(ana.TotalBand), hplot.WithLogY(LogY))
	if !ana.HistoStack {
		stack.Stack = hplot.HStackOff
	}

	return stack
}

// Helper function computing the ratio and adding them to the plot.
func (ana *Maker) addRatioToPlot(rp *hplot.RatioPlot, bhistos []*hbook.H1D, phistos []*hplot.H1D) {

	// Get all histogram (hbook to compute ratio) and (hplot) for the style
	bhBkgs := hbookHistoFromIdx(bhistos, ana.idxBkgs)
	phBkgs := hplotHistoFromIdx(phistos, ana.idxBkgs)
	bhBkgTot := histTot(bhBkgs)

	// Compute and store the ratio (type hbook.S2D)
	switch {
	case ana.HistoStack:

		// MC to MC
		hbs2d_ratioMC, err := hbook.DivideH1D(bhBkgTot, bhBkgTot, hbook.DivIgnoreNaNs())
		if err != nil {
			log.Fatal("cannot divide histo for the ratio plot")
		}
		hps2d_ratioMC := hplot.NewS2D(hbs2d_ratioMC, hplot.WithBand(true),
			hplot.WithStepsKind(hplot.HiSteps),
		)
		hps2d_ratioMC.GlyphStyle.Radius = 0
		hps2d_ratioMC.LineStyle.Width = 0.0
		hps2d_ratioMC.Band.FillColor = ana.TotalBandColor
		rp.Bottom.Add(hps2d_ratioMC)

		if len(ana.idxData) > 0 {

			bhData := bhistos[ana.idxData[0]]
			phData := phistos[ana.idxData[0]]

			// Data to MC
			hbs2d_ratio, err := hbook.DivideH1D(bhData, bhBkgTot, hbook.DivIgnoreNaNs())
			if err != nil {
				log.Fatal("cannot divide histo for the ratio plot")
			}
			hps2d_ratio := hplot.NewS2D(hbs2d_ratio, hplot.WithYErrBars(true),
				hplot.WithStepsKind(hplot.HiSteps),
			)
			style.CopyStyleH1DtoS2D(hps2d_ratio, phData)
			rp.Bottom.Add(hps2d_ratio)
		}

	default:
		// FIX-ME (rmadar): Ratio wrt data (or first bkg if data is empty)
		//                    -> to be specied as an option?
		href := bhBkgs[ana.idxBkgs[0]]
		if len(ana.idxData) > 0 {
			href = bhistos[ana.idxData[0]]
		}
		for i, h := range bhBkgs {

			hbs2d_ratio, err := hbook.DivideH1D(h, href, hbook.DivIgnoreNaNs())
			if err != nil {
				log.Fatal("cannot divide histo for the ratio plot")
			}

			hps2d_ratio := hplot.NewS2D(hbs2d_ratio,
				hplot.WithBand(phBkgs[i].Band != nil),
				hplot.WithStepsKind(hplot.HiSteps),
			)
			style.CopyStyleH1DtoS2D(hps2d_ratio, phBkgs[i])
			rp.Bottom.Add(hps2d_ratio)
		}
	}
}

// Helper function returning a slice of hplot histo
// corresponding to a list of indices.
func hplotHistoFromIdx(src []*hplot.H1D, indices []int) []*hplot.H1D {
	dst := make([]*hplot.H1D, len(indices))
	for i, idx := range indices {
		dst[i] = src[idx]
	}
	return dst
}

// Helper function returning a slice of hbook histo
// corresponding to a list of indices.
func hbookHistoFromIdx(src []*hbook.H1D, indices []int) []*hbook.H1D {
	dst := make([]*hbook.H1D, len(indices))
	for i, idx := range indices {
		dst[i] = src[idx]
	}
	return dst
}

// Helper function returning the summed histogram.
func histTot(hs []*hbook.H1D) *hbook.H1D {
	hTot := hs[0]
	for _, h := range hs[1:] {
		hTot = hbook.AddH1D(hTot, h)
	}
	return hTot
}
