package ana

import (
	"image/color"
	"log"
	"os"
	"time"
	
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
	for _, iv := range ana.varIdx {
		for _, ic := range ana.cutIdx {
			ana.plotVar(iv, ic, latex)
		}
	}

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

func (ana *Maker) plotVar(iVar, iCut int, latex htex.Handler) {

	// Current variable
	v := ana.Variables[iVar]

	var (
		drw       hplot.Drawer
		plt       = hplot.New()
		figWidth  = 6 * vg.Inch
		figHeight = 4.5 * vg.Inch
	)

	// Normalize histograms
	if ana.HistoNorm {
		ana.normalizeHistos(iCut, iVar)
	}

	// Post-normalization hbook histograms
	bhBkgs, bhSigs, bhData := ana.getHbookH1D(iCut, iVar)

	// Compute total histogram
	bhBkgTot := histTot(bhBkgs)

	// hplot histograms
	phBkgs, phSigs, phData := ana.getHplotH1D(bhBkgs, bhSigs, bhData, v.LogY)

	// Stack histograms  plotting
	hSlices := func(hm map[string]*hplot.H1D, names []string) (hs []*hplot.H1D) {
		hs = make([]*hplot.H1D, len(hm))
		for i, n := range names {
			hs[i] = hm[n]
		}
		return
	}
	phBkgsSlice := hSlices(phBkgs, ana.bkgNames)
	phSigsSlice := hSlices(phSigs, ana.sigNames)
	stack := ana.stackHistograms(phBkgsSlice, phSigsSlice, phData, v.LogY)

	// Add legends
	for _, s := range ana.Samples {
		switch s.sType {
		case data:
			plt.Legend.Add(s.LegLabel, phData)
		case bkg:
			plt.Legend.Add(s.LegLabel, phBkgs[s.Name])
		case sig:
			plt.Legend.Add(s.LegLabel, phSigs[s.Name])
		}
	}

	// Add histogram and stacks to the plot
	if stack != nil {
		plt.Add(stack)
	}
	if bhData != nil {
		plt.Add(phData)
	}
	if !ana.SignalStack {
		for _, hs := range phSigs {
			plt.Add(hs)
		}
	}
	if ana.HistoStack && ana.TotalBand {
		hBand := hplot.NewH1D(hbook.NewH1D(1, 0, 1), hplot.WithBand(true))
		hBand.Band = stack.Band
		hBand.Band.FillColor = ana.TotalBandColor
		hBand.LineStyle.Width = 0
		plt.Legend.Add("Uncer.", hBand)
	}

	// Apply common and user-defined style for this variable
	plt.Title.Text = ana.PlotTitle
	style.ApplyToPlot(plt)
	v.setPlotStyle(plt)

	// Manage log scale after settings
	if v.LogY {
		plt.Y.Scale = plot.LogScale{}
		plt.Y.Tick.Marker = plot.LogTicks{}
	}


	// -----------------------
	// TO RE-ORGANIZE FROM HERE
	// ------------------------
	
	drw = plt

	// Addition of the ratio plot
	if ana.RatioPlot {

		// Create a ratio plot, init top and bottom plots with current plot p
		rp := hplot.NewRatioPlot()
		style.ApplyToBottomPlot(rp.Bottom)
		rp.Bottom.X = plt.X
		rp.Top = plt
		rp.Top.X.Tick.Label.Font.Size = 0
		rp.Top.X.Tick.Label.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 0}
		rp.Top.X.Tick.LineStyle.Width = 0.5
		rp.Top.X.Tick.LineStyle.Color = color.NRGBA{R: 120, G: 120, B: 120, A: 255}
		rp.Top.X.Tick.Length = 5
		rp.Top.X.LineStyle.Width = 0.8
		rp.Top.X.LineStyle.Color = color.NRGBA{R: 120, G: 120, B: 120, A: 255}
		rp.Top.X.Label.Text = ""

		// Update the drawer and figure size
		figWidth, figHeight = 6*vg.Inch, 4.5*vg.Inch
		drw = rp

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

			if bhData != nil {
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
			for name, h := range bhBkgs {

				href := bhData
				if bhData == nil {	
					href = bhBkgs[ana.bkgNames[0]]
				}
								
				hbs2d_ratio, err := hbook.DivideH1D(h, href, hbook.DivIgnoreNaNs())
				if err != nil {
					log.Fatal("cannot divide histo for the ratio plot")
				}

				hps2d_ratio := hplot.NewS2D(hbs2d_ratio,
					hplot.WithBand(phBkgs[name].Band != nil),
					hplot.WithStepsKind(hplot.HiSteps),
				)
				style.CopyStyleH1DtoS2D(hps2d_ratio, phBkgs[name])
				rp.Bottom.Add(hps2d_ratio)
			}
		}

		// Adjust ratio plot scale
		if v.RatioYmin != v.RatioYmax {
			rp.Bottom.Y.Min = v.RatioYmin
			rp.Bottom.Y.Max = v.RatioYmax
		}
	}

	// Save the figure
	f := hplot.Figure(drw)
	style.ApplyToFigure(f)
	f.Latex = latex

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

// Helper function to normalize histograms of a given cut
// and variable.
func (ana *Maker) normalizeHistos(iCut, iVar int) {
	nHistos, nTot := ana.normHists[iCut], ana.normTotal[iCut]
	for iSample, sample := range ana.Samples {

		// Get the current histogram
		h := ana.HbookHistos[iSample][iCut][iVar]

		// Normalize depending on type / stack
		switch sample.sType {
		case data:
			h.Scale(1 / nHistos[iSample])
		case bkg, sig:
			if ana.HistoStack {
				h.Scale(1. / nTot)
			} else {
				h.Scale(1. / nHistos[iSample])
			}
		}
	}
}

// Helper function to get (hbook) histograms.
// It returns two maps [bkg.Name]hbook.H1D (bkgs), [sig.Name]hbook.H1D (sigs)
// and one histogram hbook.H1D (data).
func (ana *Maker) getHbookH1D(iCut, iVar int) (map[string]*hbook.H1D, map[string]*hbook.H1D, *hbook.H1D) {

	bhBkgs := make(map[string]*hbook.H1D)
	bhSigs := make(map[string]*hbook.H1D)
	var bhData *hbook.H1D

	for i, s := range ana.Samples {
		h := ana.HbookHistos[i][iCut][iVar]
		switch s.sType {
		case data:
			bhData = h
		case bkg:
			bhBkgs[s.Name] = h
		case sig:
			bhSigs[s.Name] = h
		}
	}

	return bhBkgs, bhSigs, bhData
}

// Helper function to get hplot histograms from hbook histograms.
// It returns two maps [string]hplot.H1D (bkgs), [string]hplot.H1D (sigs)
// and one histogram hplot.H1D (data).
func (ana *Maker) getHplotH1D(
	hBkgs map[string]*hbook.H1D,
	hSigs map[string]*hbook.H1D,
	hData *hbook.H1D, LogY bool) (map[string]*hplot.H1D, map[string]*hplot.H1D, *hplot.H1D) {

	phBkgs := make(map[string]*hplot.H1D)
	phSigs := make(map[string]*hplot.H1D)
	var phData *hplot.H1D

	// Loop over sample
	for _, s := range ana.Samples {

		switch s.sType {
		case data:
			phData = s.CreateHisto(hData, hplot.WithLogY(LogY))
			if s.DataStyle {
				style.ApplyToDataHist(phData)
				if s.CircleSize > 0 {
					phData.GlyphStyle.Radius = s.CircleSize
				}
				if s.YErrBarsLineWidth > 0 {
					phData.YErrs.LineStyle.Width = s.YErrBarsLineWidth
				}
				if s.YErrBarsCapWidth > 0 {
					phData.YErrs.CapWidth = s.YErrBarsCapWidth
				}
			}

		case bkg:
			phBkgs[s.Name] = s.CreateHisto(hBkgs[s.Name], hplot.WithLogY(LogY))

		case sig:
			phSigs[s.Name] = s.CreateHisto(hSigs[s.Name], hplot.WithLogY(LogY))

		}

	}

	return phBkgs, phSigs, phData
}

// Helper function to stack histograms
func (ana *Maker) stackHistograms(
	hBkgs []*hplot.H1D,
	hSigs []*hplot.H1D,
	hData *hplot.H1D,
	LogY bool) *hplot.HStack {

	if len(hBkgs)+len(hSigs) == 0 {
		return nil
	}

	// Put all backgrounds in the stack
	phStack := make([]*hplot.H1D, len(hBkgs))
	copy(phStack, hBkgs)

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

// Helper function returning the summed histogram.
func histTot(hs map[string]*hbook.H1D) *hbook.H1D {

	hSlice := make([]*hbook.H1D, 0, len(hs))
	for _, h := range hs {
		hSlice = append(hSlice, h)
	}

	hTot := hSlice[0]
	for _, h := range hSlice[1:] {
		hTot = hbook.AddH1D(hTot, h)
	}

	return hTot
}
