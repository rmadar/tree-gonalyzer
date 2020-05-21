// Package allowing to wrap all needed element of a TTree plotting analysis
package ana

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/hplot/htex"

	"github.com/rmadar/hplot-style/style"
)

// Analysis maker type is the main object of the ana package.
// It contains all samples and variables on which
// to loop to produce and plot histograms. Several options
// can be specified to normalize and/or stack histograms,
// add a ratio plot, or a list of kinematic selections.
type Maker struct {

	// Inputs
	Samples      []*Sample    // List of samples on which to run.
	Variables    []*Variable  // List of variables to plot.
	KinemCuts    []*Selection // List of cuts to apply (default: no cut).

	// Figures
	SavePath     string // Path to which plot will be saved (default: 'plots').
	SaveFormat   string // Plot file extension: 'tex' (default), 'pdf' or 'png'.
	CompileLatex bool   // On-the-fly latex compilation (default: true).

	// Plots
	AutoStyle    bool        // Enable automatic styling (default: false).
	PlotTitle    string      // General plot title (default: 'TTree GOnalyzer').
	RatioPlot    bool        // Enable ratio plot (default: true).
	HistoStack   bool        // Enable histogram stacking (default: true).
	HistoNorm    bool        // Normalize distributions to unit area (default: false).
	ErrBandColor color.NRGBA // Color for the uncertainty band (default: gray).

	// Histograms for {variables x samples x selection}
	HbookHistos [][][]*hbook.H1D 
	HplotHistos [][][]*hplot.H1D 

	// Fields for benchmarking (TEMP)
	WithVarsTreeFormula bool
	NoTreeFormula       bool
	NoFuncCall          bool

	// Internal fields
	cutIdx      map[string]int // Linking cut name and cut index
	samIdx      map[string]int // Linking sample name and sample index
	varIdx      map[string]int // Linking variable name and variable index
	histoFilled bool           // true if histograms are filled.
	nEvents     int64          // Number of processed events
	timeLoop    time.Duration  // Processing time for filling histograms (event loop over samples x cuts x histos)
	timePlot    time.Duration  // Processing time for plotting histogram

}

// Creating a new object
func New(s []*Sample, v []*Variable, opts ...Options) Maker {

	// Create the object
	a := Maker{
		Samples:   s,
		Variables: v,
	}

	// Configuration with default values for all optional fields
	cfg := newConfig(
		WithSavePath("plots"),
		WithSaveFormat("tex"),
		WithPlotTitle(`TTree GOnalyzer`),
		WithCompileLatex(true),
		WithHistoStack(true),
		WithHistoNorm(false),
		WithRatioPlot(true),
		WithErrBandColor(color.NRGBA{A: 100}),
		WithKinemCuts([]*Selection{NewSelection()}),
	)

	// Update the configuration looping over functional options
	for _, opt := range opts {
		opt(cfg)
	}

	// Set fields with updaded configuration
	a.KinemCuts = cfg.KinemCuts
	a.SavePath = cfg.SavePath
	a.SaveFormat = cfg.SaveFormat
	a.AutoStyle = cfg.AutoStyle
	a.PlotTitle = cfg.PlotTitle
	a.CompileLatex = cfg.CompileLatex
	a.RatioPlot = cfg.RatioPlot
	a.HistoStack = cfg.HistoStack
	a.HistoNorm = cfg.HistoNorm
	a.ErrBandColor = cfg.ErrBandColor

	// Get mappings between slice indices and object names
	a.samIdx = getIdxMap(a.Samples)
	a.varIdx = getIdxMap(a.Variables)
	a.cutIdx = getIdxMap(a.KinemCuts)

	// Build hbook and hplot H1D containers
	a.initHistoContainers()

	return a
}

// Helper function creating the mapping between name and objects
func getIdxMap(obj interface{}) map[string]int {
	// obj = []*Variables, []*Sample, []*Selections
	return make(map[string]int, 10)
}

// Run the event loop to fill all histo across samples / variables / cuts (and later: systematics)
func (ana *Maker) MakeHistos() error {

	// Start timing
	start := time.Now()

	// Loop over the samples
	for is, s := range ana.Samples {

		// Anonymous function to avoid memory-leaks due to 'defer'
		func(j int) error {

			// Get the file and tree
			f, t := getTreeFromFile(s.FileName, s.TreeName)
			defer f.Close()

			// Prepare variables to explicitely load
			var rvars []rtree.ReadVar
			if !ana.WithVarsTreeFormula {
				for _, v := range ana.Variables {
					rvars = append(rvars, rtree.ReadVar{Name: v.TreeName, Value: v.Value})
				}
			}

			// Get the tree reader
			r, err := rtree.NewReader(t, rvars)
			if err != nil {
				log.Fatal("could not create tree reader: %w", err)
			}
			defer r.Close()

			varFormula := make([]func() float64, len(ana.Variables))
			if ana.WithVarsTreeFormula {
				for i, v := range ana.Variables {
					varFormula[i] = v.TreeVar.GetFuncF64(r)
				}
			}

			// Prepare the weight
			getWeight := func() float64 { return float64(1.0) }
			if s.WeightFunc.Fct != nil && !ana.NoTreeFormula {
				getWeight = s.WeightFunc.GetFuncF64(r)
			}

			// Prepare the sample cut
			passSampleCut := func() bool { return true }
			if s.CutFunc.Fct != nil && !ana.NoTreeFormula {
				passSampleCut = s.CutFunc.GetFuncBool(r)
			}

			// Prepare the cut string for kinematics
			passKinemCut := make([]func() bool, len(ana.KinemCuts))
			for ic, cut := range ana.KinemCuts {
				idx := ic
				if !ana.NoTreeFormula {
					passKinemCut[idx] = cut.TreeFunc.GetFuncBool(r)
				} else {
					_, passKinemCut[idx] = cut, func() bool { return true }
				}
			}

			// Read the tree (event loop)
			err = r.Read(func(ctx rtree.RCtx) error {

				// No call to function at all
				if ana.NoFuncCall {
					for ic := range ana.KinemCuts {
						for iv, v := range ana.Variables {
							ana.HbookHistos[iv][ic][is].Fill(v.GetValue(), 1.0)
						}
					}

				} else {

					// Sample-level cut
					if !passSampleCut() {
						return nil
					}

					// Get the event weight
					w := getWeight()

					// Loop over selection and variables
					for ic := range ana.KinemCuts {

						if !passKinemCut[ic]() {
							continue
						}

						for iv, v := range ana.Variables {
							val := 0.0
							if ana.WithVarsTreeFormula {
								val = varFormula[iv]()
							} else {
								val = v.GetValue()
							}
							ana.HbookHistos[iv][ic][is].Fill(val, w)
						}
					}
				}

				return nil
			})

			// Keep track of the number of processed events
			ana.nEvents += t.Entries()

			return nil
		}(is)
	}

	ana.histoFilled = true

	// End timing
	ana.timeLoop = time.Now().Sub(start)

	return nil
}

// Plotting all histograms
func (ana *Maker) PlotHistos() error {

	// Start timing
	start := time.Now()

	// Set histogram styles
	if ana.AutoStyle {
		ana.setAutoStyle()
	}

	// Return an error if HbookHistos is empty
	if !ana.histoFilled {
		log.Fatalf("There is no histograms. Please make sure that 'MakeHistos()' is called before 'PlotHistos()'")
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
	for iv, h_sel_samples := range ana.HbookHistos {

		// Current variable
		v := ana.Variables[iv]

		// Loop over selections
		for isel, hsamples := range h_sel_samples {

			var (
				p               = hplot.New()
				bhBkgTot        = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
				bhData          = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
				norm_histos     = make([]float64, 0, len(hsamples))
				norm_bkgtot     = 0.0
				bhBkgs_postnorm []*hbook.H1D
				phBkgs          []*hplot.H1D
				phData          *hplot.H1D
			)

			// Add plot title
			p.Title.Text = ana.PlotTitle

			// First sample loop: compute normalisation, sum bkg bh, keep data bh
			for is, h := range hsamples {

				// Compute the integral of the current histo
				n := h.Integral()

				// Properly store individual normalization
				norm_histos = append(norm_histos, n)

				// For background only
				if ana.Samples[is].IsBkg() {
					norm_bkgtot += n
				}

				// Keep data apart
				if ana.Samples[is].IsData() {
					bhData = h
				}

			}

			// Second sample loop: normalize bh, prepare background stack
			for is, h := range hsamples {

				// Deal with normalization
				// FIX-ME (rmadar): use 'switch' and enum
				if ana.HistoNorm {
					if ana.Samples[is].IsData() {
						h.Scale(1 / norm_histos[is])
					}
					if ana.Samples[is].IsBkg() {
						if ana.HistoStack {
							h.Scale(1. / norm_bkgtot)
						} else {
							h.Scale(1. / norm_histos[is])
						}
					}
				}

				// Get plottable histogram and add it to the legend
				withBand := false
				if !ana.Samples[is].IsData() {
					withBand = true
				}
				hplt := ana.Samples[is].CreateHisto(h, hplot.WithBand(withBand))
				p.Legend.Add(ana.Samples[is].LegLabel, hplt)
				ana.HplotHistos[iv][isel][is] = hplt

				// Keep data appart from backgrounds
				if ana.Samples[is].IsData() {
					phData = ana.HplotHistos[iv][isel][is]
					if ana.Samples[is].DataStyle {
						style.ApplyToDataHist(phData)
					}

				}

				// Sum-up normalized bkg and store all bkgs in a slice for the stack
				if ana.Samples[is].IsBkg() {
					bhBkgs_postnorm = append(bhBkgs_postnorm, h)
					bhBkgTot = hbook.AddH1D(h, bhBkgTot)
					phBkgs = append(phBkgs, ana.HplotHistos[iv][isel][is])
				}
			}

			// Manage background stack plotting
			if len(phBkgs) > 0 {

				// Reverse the order so that legend and plot order matches
				for i, j := 0, len(phBkgs)-1; i < j; i, j = i+1, j-1 {
					phBkgs[i], phBkgs[j] = phBkgs[j], phBkgs[i]
				}

				// Stacking the background histo
				stack := hplot.NewHStack(phBkgs, hplot.WithBand(true))
				stack.Band.FillColor = ana.ErrBandColor
				if ana.HistoStack {
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

			// Apply common and user-defined style for this variable
			// FIX-ME (rmadar): the v.SetPlotStyle(v) command doesn't update
			//                  y-axis scale if it is put before the samples
			//                  loop and I am not sure why.
			style.ApplyToPlot(p)
			v.SetPlotStyle(p)
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
					// Data to MC
					hbs2d_ratio, err := hbook.DivideH1D(bhData, bhBkgTot, hbook.DivIgnoreNaNs())
					if err != nil {
						log.Fatal("cannot divide histo for the ratio plot")
					}
					hps2d_ratio := hplot.NewS2D(hbs2d_ratio, hplot.WithYErrBars(true),
						hplot.WithStepsKind(hplot.HiSteps),
					)
					// TODO: copy style from bhData when a function hplot-style will be ready.
					//if ana.Samples[is].DataStyle {
					style.ApplyToDataS2D(hps2d_ratio)
					//}

					// MC to MC
					hbs2d_ratio1, err := hbook.DivideH1D(bhBkgTot, bhBkgTot, hbook.DivIgnoreNaNs())
					if err != nil {
						log.Fatal("cannot divide histo for the ratio plot")
					}
					hps2d_ratio1 := hplot.NewS2D(hbs2d_ratio1, hplot.WithBand(true),
						hplot.WithStepsKind(hplot.HiSteps),
					)
					//if ana.Samples[is].DataStyle {
					//style.ApplyToDataS2D(hps2d_ratio1)
					//}
					hps2d_ratio1.GlyphStyle.Radius = 0
					hps2d_ratio1.LineStyle.Width = 0.0
					hps2d_ratio1.Band.FillColor = ana.ErrBandColor
					rp.Bottom.Add(hps2d_ratio1)
					rp.Bottom.Add(hps2d_ratio)

				default:
					// [FIX-ME 0 (rmadar)] Ratio wrt data (or 1 bkg if data is empty) -> to be specied as an option?
					// [FIX-ME 1 (rmadar)] loop is over bhBkgs_postnorm while 'ana.Samples[is]' runs also over data.
					for is, h := range bhBkgs_postnorm {

						href := bhData
						if bhData.Entries() == 0 {
							href = bhBkgs_postnorm[0]
						}

						hbs2d_ratio, err := hbook.DivideH1D(h, href, hbook.DivIgnoreNaNs())
						if err != nil {
							log.Fatal("cannot divide histo for the ratio plot")
						}
						hps2d_ratio := hplot.NewS2D(hbs2d_ratio, hplot.WithBand(true),
							hplot.WithStepsKind(hplot.HiSteps),
						)
						hps2d_ratio.GlyphStyle.Radius = 0
						hps2d_ratio.LineStyle.Color = ana.Samples[is].LineColor
						ana.Samples[is].SetBandStyle(hps2d_ratio.Band)
						rp.Bottom.Add(hps2d_ratio)
					}
				}

				// Adjust ratio plot scale
				rp.Bottom.Y.Min = 0.7
				rp.Bottom.Y.Max = 1.3
			}

			// Save the figure
			f := hplot.Figure(plt)
			style.ApplyToFigure(f)
			f.Latex = latex

			path := ana.SavePath + "/" + ana.KinemCuts[isel].Name
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
	ana.timePlot = time.Now().Sub(start)

	return nil
}

// Print processing report
func (ana Maker) PrintReport() {

	// Event, histo info
	nvars, nsamples, ncuts := len(ana.Variables), len(ana.Samples), len(ana.KinemCuts)
	nhist := nvars * nsamples
	if ncuts > 0 {
		nhist *= ncuts
	}
	nkevt := float64(ana.nEvents) / 1e3

	// Timing info
	dtLoop := float64(ana.timeLoop) / float64(time.Millisecond)
	dtPlot := float64(ana.timePlot) / float64(time.Millisecond)

	// Formating
	str_template := "\n Processing report:\n"
	str_template += "    - %v histograms filled over %.0f kEvts (%v samples, %v variables, %v selections)\n"
	str_template += "    - running time: %.1f ms/kEvt (%s for %.0f kEvts)\n"
	str_template += "    - time fraction: %.0f%% (event loop), %.0f%% (plotting)\n\n"

	fmt.Printf(str_template,
		nhist, nkevt, nsamples, nvars, ncuts,
		(dtLoop+dtPlot)/nkevt, fmtDuration(ana.timeLoop+ana.timePlot), nkevt,
		dtLoop/(dtLoop+dtPlot)*100., dtPlot/(dtLoop+dtPlot)*100.,
	)
}

// Run the analysis in one function
func (ana *Maker) Run() error {

	// Create histograms via an event loop
	err := ana.MakeHistos()
	if err != nil {
		return err
	}

	// Plot them on the same canvas
	err = ana.PlotHistos()
	if err != nil {
		return err
	}

	// Print processing report
	ana.PrintReport()

	// Return
	return nil
}

// Initialize histogram containers
func (ana *Maker) initHistoContainers() {

	// Initialize hbook H1D
	ana.HbookHistos = make([][][]*hbook.H1D, len(ana.Variables))
	for iv := range ana.HbookHistos {
		ana.HbookHistos[iv] = make([][]*hbook.H1D, len(ana.KinemCuts))
		for isel := range ana.KinemCuts {
			ana.HbookHistos[iv][isel] = make([]*hbook.H1D, len(ana.Samples))
			v := ana.Variables[iv]
			for isample := range ana.HbookHistos[iv][isel] {
				ana.HbookHistos[iv][isel][isample] = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
			}
		}
	}

	// Inititialize hplot H1D
	ana.HplotHistos = make([][][]*hplot.H1D, len(ana.Variables))
	for iv := range ana.HplotHistos {
		ana.HplotHistos[iv] = make([][]*hplot.H1D, len(ana.KinemCuts))
		for isel := range ana.KinemCuts {
			ana.HplotHistos[iv][isel] = make([]*hplot.H1D, len(ana.Samples))
		}
	}

}

func (ana *Maker) setAutoStyle() {

	for i, s := range ana.Samples {

		// Color
		r, g, b, a := plotutil.Color(i).RGBA()
		c := color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}

		// Apply data style automatically
		if s.IsData() {
			s.DataStyle = true
		}

		// Fill for stacked histo, lines otherwise
		if ana.HistoStack {
			s.FillColor = c
			s.LineWidth = 0.
		} else {
			s.FillColor = color.NRGBA{}
			s.LineColor = c
			s.LineWidth = 2.
		}
	}
}

// Helper to get a tree from a file
func getTreeFromFile(filename, treename string) (*groot.File, rtree.Tree) {

	// Get the file
	f, err := groot.Open(filename)
	if err != nil {
		log.Fatal("could not open ROOT file: %w", err)
	}

	// Get the tree
	obj, err := f.Get(treename)
	if err != nil {
		log.Fatal("could not retrieve object: %w", err)
	}

	return f, obj.(rtree.Tree)
}

// Helper duration formating: return a string 'hh:mm:ss' for a time.Duration object
func fmtDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
