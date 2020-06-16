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
