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
		plt             hplot.Drawer
		p               = hplot.New()
		figWidth        = 6 * vg.Inch
		figHeight       = 4.5 * vg.Inch
		bhData          = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
		bhBkgTot        = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
		bhSigTot        = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
		bhBkgs_postnorm []*hbook.H1D
		phBkgs          []*hplot.H1D
		bhSigs_postnorm []*hbook.H1D
		phSigs          []*hplot.H1D
		phData          *hplot.H1D
	)

	// Normalize histograms
	if ana.HistoNorm {
		ana.normalizeHistos(iCut, iVar)
	}

	// FIX-ME [rmadar]: create few functions to get: 
	//  -> bhBkgs, bhSigs, bhData = ana.getNormalizedH1D()
	//  -> bhBkgTot, bhSigTot = XXX() altough not sure bhSigTot is needed

	// Function to get hplot histo
	// and slice to be stacked later.
	
	// Sample Loop
	for iSample, sample := range ana.Samples {

		// Get the current histogram
		h := ana.HbookHistos[iSample][iCut][iVar]
		
		// Keep data appart
		if sample.IsData() {
			bhData = h
		}
			
		// Get plottable histogram and add it to the legend
		ana.HplotHistos[iSample][iCut][iVar] = sample.CreateHisto(h, hplot.WithLogY(v.LogY))
		p.Legend.Add(sample.LegLabel, ana.HplotHistos[iSample][iCut][iVar])

		// Keep track of different histo given their type
		switch sample.sType {

		// Keep data appart from backgrounds, style it.
		case data:
			phData = ana.HplotHistos[iSample][iCut][iVar]
			if sample.DataStyle {
				style.ApplyToDataHist(phData)
				if sample.CircleSize > 0 {
					phData.GlyphStyle.Radius = sample.CircleSize
				}
				if sample.YErrBarsLineWidth > 0 {
					phData.YErrs.LineStyle.Width = sample.YErrBarsLineWidth
				}
				if sample.YErrBarsCapWidth > 0 {
					phData.YErrs.CapWidth = sample.YErrBarsCapWidth
				}
			}

		// Sum-up normalized bkg and store all bkgs in a slice for the stack
		case bkg:
			phBkgs = append(phBkgs, ana.HplotHistos[iSample][iCut][iVar])
			bhBkgs_postnorm = append(bhBkgs_postnorm, h)
			bhBkgTot = hbook.AddH1D(h, bhBkgTot)

		//
		case sig:
			phSigs = append(phSigs, ana.HplotHistos[iSample][iCut][iVar])
			bhSigs_postnorm = append(bhSigs_postnorm, h)
			bhSigTot = hbook.AddH1D(h, bhSigTot)
		}
	}

	// Manage background stack plotting
	if len(phBkgs)+len(phSigs) > 0 {

		// Put all backgrounds in the stack
		phStack := make([]*hplot.H1D, len(phBkgs))
		copy(phStack, phBkgs)

		// Reverse the order so that legend and plot order matches
		for i, j := 0, len(phStack)-1; i < j; i, j = i+1, j-1 {
			phStack[i], phStack[j] = phStack[j], phStack[i]
		}

		// Add signals if asked (after the order reversering to have
		// signals on top of the bkg).
		if ana.SignalStack {
			for _, hs := range phSigs {
				phStack = append(phStack, hs)
			}
		}

		// Stacking the background histo
		stack := hplot.NewHStack(phStack, hplot.WithBand(ana.TotalBand), hplot.WithLogY(v.LogY))
		if ana.HistoStack && ana.TotalBand {
			stack.Band.FillColor = ana.TotalBandColor
			hBand := hplot.NewH1D(hbook.NewH1D(1, 0, 1), hplot.WithBand(true))
			hBand.Band = stack.Band
			hBand.LineStyle.Width = 0
			p.Legend.Add("Uncer.", hBand)
		} else {
			stack.Stack = hplot.HStackOff
		}

		// Add the stack to the plot
		p.Add(stack)
	}

	// Adding hplot.H1D data to the plot, set the drawer to the current plot
	if bhData.Entries() > 0 {
		p.Add(phData)
	}

	// Add individual signals (if not stacked) after the data
	if !ana.SignalStack {
		for _, hs := range phSigs {
			p.Add(hs)
		}
	}

	// Apply common and user-defined style for this variable
	// FIX-ME (rmadar): the v.setPlotStyle(v) command doesn't update
	//                  y-axis scale if it is put before the samples
	//                  loop and I am not sure why.
	// Add plot title
	p.Title.Text = ana.PlotTitle
	style.ApplyToPlot(p)
	v.setPlotStyle(p)

	// Manage log scale after settings
	if v.LogY {
		p.Y.Scale = plot.LogScale{}
		p.Y.Tick.Marker = plot.LogTicks{}
	}

	plt = p

	// Addition of the ratio plot
	if ana.RatioPlot {

		// Create a ratio plot, init top and bottom plots with current plot p
		rp := hplot.NewRatioPlot()
		style.ApplyToBottomPlot(rp.Bottom)
		rp.Bottom.X = p.X
		rp.Top = p
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
		plt = rp

		// Compute and store the ratio (type hbook.S2D)
		switch {
		case ana.HistoStack:

			if bhData.Entries() > 0 {
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

		default:
			// FIX-ME (rmadar): Ratio wrt data (or first bkg if data is empty)
			//                    -> to be specied as an option?
			for ib, h := range bhBkgs_postnorm {

				href := bhData
				if bhData.Entries() == 0 {
					href = bhBkgs_postnorm[0]
				}

				hbs2d_ratio, err := hbook.DivideH1D(h, href, hbook.DivIgnoreNaNs())
				if err != nil {
					log.Fatal("cannot divide histo for the ratio plot")
				}

				hps2d_ratio := hplot.NewS2D(hbs2d_ratio,
					hplot.WithBand(phBkgs[ib].Band != nil),
					hplot.WithStepsKind(hplot.HiSteps),
				)
				style.CopyStyleH1DtoS2D(hps2d_ratio, phBkgs[ib])
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
	f := hplot.Figure(plt)
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
