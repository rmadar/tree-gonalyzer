package ana

import (
	"image/color"
	"log"
	"os"
	"time"
	"sync"
	
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/hplot/htex"

	"github.com/rmadar/hplot-style/style"
)



func (ana *Maker) oldPlotVariables() error {

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

	// Preparing the final figure
	var plt hplot.Drawer
	figWidth, figHeight := 6*vg.Inch, 4.5*vg.Inch

	// Handle on-the-fly LaTeX compilation
	var latex htex.Handler = htex.NoopHandler{}
	if ana.CompileLatex {
		latex = htex.NewGoHandler(-1, "pdflatex")
	}

	// Loop over variables
	for _, iVar := range ana.varIdx {

		// Current variable
		v := ana.Variables[iVar]

		// Loop over selections
		for _, iCut := range ana.cutIdx {

			var (
				p               = hplot.New()
				bhData          = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
				bhBkgTot        = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
				bhSigTot        = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
				norm_histos     = make([]float64, 0, len(ana.Samples))
				normTot         = 0.0
				bhBkgs_postnorm []*hbook.H1D
				phBkgs          []*hplot.H1D
				bhSigs_postnorm []*hbook.H1D
				phSigs          []*hplot.H1D
				phData          *hplot.H1D
			)

			// Add plot title
			p.Title.Text = ana.PlotTitle

			// First sample loop: compute normalisation, sum bkg bh, keep data bh
			for iSample := range ana.Samples {

				// Get the current histogram
				h := ana.HbookHistos[iSample][iCut][iVar]

				// Compute the integral of the current histo
				n := h.Integral()

				// Properly store individual normalization
				norm_histos = append(norm_histos, n)

				// For background only
				if ana.Samples[iSample].IsBkg() {
					normTot += n
				}

				// Add signals if stacked
				if ana.SignalStack {
					normTot += n
				}

				// Keep data apart
				if ana.Samples[iSample].IsData() {
					bhData = h
				}
			}

			// Second sample loop: normalize bh, prepare background stack
			for iSample := range ana.Samples {

				// Get the current histogram
				h := ana.HbookHistos[iSample][iCut][iVar]

				// Deal with normalization
				if ana.HistoNorm {
					switch ana.Samples[iSample].sType {
					case data:
						h.Scale(1 / norm_histos[iSample])
					case bkg, sig:
						if ana.HistoStack {
							h.Scale(1. / normTot)
						} else {
							h.Scale(1. / norm_histos[iSample])
						}
					}
				}

				// Get plottable histogram and add it to the legend
				ana.HplotHistos[iSample][iCut][iVar] = ana.Samples[iSample].CreateHisto(h, hplot.WithLogY(v.LogY))
				p.Legend.Add(ana.Samples[iSample].LegLabel, ana.HplotHistos[iSample][iCut][iVar])

				// Keep track of different histo given their type
				switch ana.Samples[iSample].sType {

				// Keep data appart from backgrounds, style it.
				case data:
					phData = ana.HplotHistos[iSample][iCut][iVar]
					if ana.Samples[iSample].DataStyle {
						style.ApplyToDataHist(phData)
						if ana.Samples[iSample].CircleSize > 0 {
							phData.GlyphStyle.Radius = ana.Samples[iSample].CircleSize
						}
						if ana.Samples[iSample].YErrBarsLineWidth > 0 {
							phData.YErrs.LineStyle.Width = ana.Samples[iSample].YErrBarsLineWidth
						}
						if ana.Samples[iSample].YErrBarsCapWidth > 0 {
							phData.YErrs.CapWidth = ana.Samples[iSample].YErrBarsCapWidth
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
	//var wg sync.WaitGroup
	//wg.Add(len(ana.Variables)*len(ana.KinemCuts))
	for iv := range ana.Variables {
		for ic := range ana.KinemCuts {
			//go ana.concurrentPlotVar(iv, ic, latex, &wg)
			ana.plotVar(iv, ic, latex)
		}
	}
	//wg.Wait()

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

	// Compute total histogram
	bhBkgs := make([]*hbook.H1D, len(ana.idxBkgs))
	for i, ib := range ana.idxBkgs {
		bhBkgs[i] = bhistos[ib]
	}
	bhBkgTot := histTot(bhBkgs)

	// hplot histograms
	phistos := ana.getHplotH1D(bhistos, v.LogY)

	// Stack signal/bkg histograms
	hType := func(src []*hplot.H1D, indices []int) []*hplot.H1D {
		dst := make([]*hplot.H1D, len(indices))
		for i, idx := range indices {
			dst[i] = src[idx]
		}
		return dst
	}

	phBkgs := hType(phistos, ana.idxBkgs)
	phSigs := hType(phistos, ana.idxSigs)
	stack := ana.stackHistograms(phBkgs, phSigs, v.LogY)

	// Data
	phData := hType(phistos, ana.idxData)
	
	// Add legends
	for i, s := range ana.Samples {
		plt.Legend.Add(s.LegLabel, phistos[i])
	}

	// Add histogram and stacks to the plot
	if stack != nil {
		plt.Add(stack)
	}
	if len(phData) > 0 {
		plt.Add(phData[0])
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
func (ana *Maker) getHplotH1D(hs[]*hbook.H1D, LogY bool) []*hplot.H1D {	

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

// Helper function returning the summed histogram.
func histTot(hs []*hbook.H1D) *hbook.H1D {
	hTot := hs[0]
	for _, h := range hs[1:] {
		hTot = hbook.AddH1D(hTot, h)
	}
	return hTot
}
