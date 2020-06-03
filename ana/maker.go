// Package allowing to wrap all needed element of a TTree plotting analysis
package ana

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"sync"
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
	Samples   []*Sample    // List of samples on which to run.
	Variables []*Variable  // List of variables to plot.
	KinemCuts []*Selection // List of cuts to apply (default: no cut).
	Nevts     int64        // Maximum number of events per components.
	SampleMT   bool        // Enable concurency accross samples
	
	// Ouputs
	SavePath     string // Path to which plot will be saved (default: 'outputs').
	SaveFormat   string // Plot file extension: 'png' (default), 'pdf' or 'png'.
	CompileLatex bool   // On-the-fly latex compilation (default: true).
	DumpTree     bool   // Dump a TTree in a file for each sample (default: false).
	PlotHisto    bool   // Enable histogram plotting (default: true).

	// Plots
	AutoStyle    bool        // Enable automatic styling (default: true).
	PlotTitle    string      // General plot title (default: 'TTree GOnalyzer').
	RatioPlot    bool        // Enable ratio plot (default: true).
	HistoStack   bool        // Enable histogram stacking (default: true).
	SignalStack  bool        // Enable signal stack (default: false).
	HistoNorm    bool        // Normalize distributions to unit area (default: false).
	TotalBand    bool        // Enable total error band in stack mode (default: true).
	ErrBandColor color.NRGBA // Color for the uncertainty band (default: gray).

	// Histograms for {samples x selections x variables}
	HbookHistos [][][]*hbook.H1D
	HplotHistos [][][]*hplot.H1D

	// Internal: tree dumping
	nVars       int         // number of variables
	nEvtsSample []int64     // number of events per sample

	// Internal management
	cutIdx      map[string]int // Linking cut name and cut index
	samIdx      map[string]int // Linking sample name and sample index
	varIdx      map[string]int // Linking variable name and variable index
	histoFilled bool           // true if histograms are filled.
	nEvents     int64          // Number of processed events
	timeLoop    time.Duration  // Processing time for filling histograms (event loop over samples x cuts x histos)
	timePlot    time.Duration  // Processing time for plotting histogram
}

type dumper struct { 
	Var   []float64   // Storing the F64 variables values to dump the TTree.
	Vars  [][]float64 // Storing the F64s variable values to dump the TTree.
	VarsN []int32     // Storing the number of object in the F64s to dump the TTree.
}

// New creates a default analysis maker from a list of sample
// and a list of variables.
func New(s []*Sample, v []*Variable, opts ...Options) Maker {

	// Create the object
	a := Maker{
		Samples:      s,
		Variables:    v,
		Nevts:        -1,
		PlotHisto:    true,
		SampleMT:     true,
		AutoStyle:    true,
		SavePath:     "outputs",
		SaveFormat:   "png",
		PlotTitle:    "TTree GOnalyzer",
		CompileLatex: true,
		HistoStack:   true,
		RatioPlot:    true,
		TotalBand:    true,
		ErrBandColor: color.NRGBA{A: 100},
		KinemCuts:    []*Selection{EmptySelection()},
		nVars:        len(v),
	}

	// Configuration with default values for all optional fields
	cfg := newConfig()
	
	// Update the configuration looping over functional options
	for _, opt := range opts {
		opt(cfg)
	}

	// Set fields with updaded configuration
	if cfg.KinemCuts.usr {
		a.KinemCuts = cfg.KinemCuts.val
	}
	if cfg.Nevts.usr {
		a.Nevts = cfg.Nevts.val
	}
	if cfg.SampleMT.usr {
		a.SampleMT = cfg.SampleMT.val
	}
	if cfg.SavePath.usr {
		a.SavePath = cfg.SavePath.val
	}
	if cfg.SaveFormat.usr {
		a.SaveFormat = cfg.SaveFormat.val
	}
	if cfg.DumpTree.usr {
		a.DumpTree = cfg.DumpTree.val
	}
	if cfg.PlotHisto.usr {
		a.PlotHisto = cfg.PlotHisto.val
	}
	if cfg.AutoStyle.usr {
		a.AutoStyle = cfg.AutoStyle.val
	}
	if cfg.PlotTitle.usr {
		a.PlotTitle = cfg.PlotTitle.val
	}
	if cfg.CompileLatex.usr {
		a.CompileLatex = cfg.CompileLatex.val
	}
	if cfg.RatioPlot.usr {
		a.RatioPlot = cfg.RatioPlot.val
	}
	if cfg.HistoStack.usr {
		a.HistoStack = cfg.HistoStack.val
	}
	if cfg.SignalStack.usr {
		a.SignalStack = cfg.SignalStack.val
	}
	if cfg.HistoNorm.usr {
		a.HistoNorm = cfg.HistoNorm.val
	}
	if cfg.TotalBand.usr {
		a.TotalBand = cfg.TotalBand.val
	}
	if cfg.ErrBandColor.usr {
		a.ErrBandColor = cfg.ErrBandColor.val
	}

	// Get mappings between slice indices and object names
	a.samIdx = getIdxMap(a.Samples, &Sample{})
	a.varIdx = getIdxMap(a.Variables, &Variable{})
	a.cutIdx = getIdxMap(a.KinemCuts, &Selection{})

	// Managing event number with concurrency
	a.nEvtsSample = make([]int64, len(a.Samples))
	
	// Build hbook and hplot H1D containers
	a.initHistoContainers()

	// Build the slice of values to store
	// FIX-ME(rmadar): this is not so clean to assess slice of not
	//                 by doing a loop over variables for the first
	//                 component of the first sample to fill v.isSlice.
	a.assessVariableTypes()

	return a
}

// Helper function creating the mapping between name and objects
func getIdxMap(objs interface{}, objType interface{}) map[string]int {
	res := make(map[string]int)
	switch objType.(type) {
	case *Variable:
		for i, obj := range objs.([]*Variable) {
			res[obj.Name] = i
		}
	case *Sample:
		for i, obj := range objs.([]*Sample) {
			res[obj.Name] = i
		}
	case *Selection:
		for i, obj := range objs.([]*Selection) {
			res[obj.Name] = i
		}
	default:
		panic(fmt.Errorf("invalid variable value-type %T", objType))
	}
	return res
}

// RunEventLoops runs one event loop per sample to fill
// histograms for each variables and selections. If DumpTree
// is true, a tree is also dumped with all variables and one
// branch per selection.
func (ana *Maker) RunEventLoops() error {

	// Start timing
	start := time.Now()

	// Loop over the samples
	if ana.SampleMT {
		var wg sync.WaitGroup
		for i := range ana.Samples {
			wg.Add(1)
			go ana.concurrentSampleEventLoop(i, &wg)
		}
		wg.Wait()
	} else {
		for i := range ana.Samples {
			ana.sampleEventLoop(i)
		}
	}

	for _, n := range ana.nEvtsSample {
		ana.nEvents += n

	}
	// Histograms are now filled.
	ana.histoFilled = true

	// End timing.
	ana.timeLoop = time.Since(start)

	return nil
}

func (ana *Maker) concurrentSampleEventLoop(sampleIdx int, wg *sync.WaitGroup) {

	// Handle concurrency
	defer wg.Done()

	// Fill the histo
	ana.sampleEventLoop(sampleIdx)
}

func (ana *Maker) sampleEventLoop(sampleIdx int) {
	
	// Current sample
	samp := ana.Samples[sampleIdx]

	// Initiate the structure of the histo container: h[iCut][iVar]
	h := make([][]*hbook.H1D, len(ana.KinemCuts))
	for iCut := range ana.KinemCuts {
		h[iCut] = make([]*hbook.H1D, len(ana.Variables))
		for iVar, v := range ana.Variables {
			if ana.PlotHisto {
				h[iCut][iVar] = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
			} else {
				h[iCut][iVar] = hbook.NewH1D(1, 0, 1)
			}
		}
	}

	// Output in case of TTree dumping
	path := ana.SavePath + "/ntuples/"
	if _, err := os.Stat(path); os.IsNotExist(err) && ana.DumpTree {
		os.MkdirAll(path, 0755)
	}
	var fOut *groot.File
	var tOut rtree.Writer
	dump := ana.newDumper()
	if ana.DumpTree {
		fOut, tOut = ana.getOutFileTree(path+samp.Name+".root", "GOtree", dump)
		defer fOut.Close()
		defer tOut.Close()
	}

	// Loop over the sample components
	for iComp, comp := range samp.components {

		// Anonymous function to avoid memory-leaks due to 'defer'
		func(j int) error {

			// Get the file and tree
			f, t := getTreeFromFile(comp.FileName, comp.TreeName)
			defer f.Close()

			// Get the tree reader
			r, err := rtree.NewReader(t, []rtree.ReadVar{}, rtree.WithRange(0, ana.Nevts))
			if err != nil {
				log.Fatal("could not create tree reader: %w", err)
			}
			defer r.Close()

			// Prepare variables
			ok := false
			getF64 := make([]func() float64, len(ana.Variables))
			getF64s := make([]func() []float64, len(ana.Variables))
			for iv, v := range ana.Variables {
				idx := iv
				if !v.isSlice {
					if getF64[idx], ok = v.TreeFunc.GetFuncF64(r); !ok {
						err := "Type assertion failed [variable \"%v\"]:"
						err += " this TreeFunc.Fct is supposed to return a float64."
						log.Fatalf(err, v.Name)
					}
				} else {
					if getF64s[idx], ok = v.TreeFunc.GetFuncF64s(r); !ok {
						err := "Type assertion failed [variable \"%v\"]:"
						err += " this TreeFunc.Fct is supposed to return a []float64."
						log.Fatalf(err, v.Name)
					}
				}
			}
			
			// Prepare the sample global weight
			getWeightSamp := func() float64 { return float64(1.0) }
			if samp.WeightFunc.Fct != nil {
				if getWeightSamp, ok = samp.WeightFunc.GetFuncF64(r); !ok {
					err := "Type assertion failed [weight of %v]:"
					err += " TreeFunc.Fct must return a float64."
					log.Fatalf(err, samp.Name)
				}
			}

			// Prepare the additional weight of the component
			getWeightComp := func() float64 { return float64(1.0) }
			if comp.WeightFunc.Fct != nil {
				if getWeightComp, ok = comp.WeightFunc.GetFuncF64(r); !ok {
					err := "Type assertion failed [weight of (%v, %v)]:"
					err += " TreeFunc.Fct must return a float64."
					log.Fatalf(err, samp.Name, comp.FileName)
				}
			}

			// Prepare the sample global cut
			passCutSamp := func() bool { return true }
			if samp.CutFunc.Fct != nil {
				if passCutSamp, ok = samp.CutFunc.GetFuncBool(r); !ok {
					err := "Type assertion failed [Cut of %v]:"
					err += " TreeFunc.Fct must return a bool.\n"
					err += "\t -> Make sure to use NewCutBool(), not NewVarBool()."
					log.Fatalf(err, samp.Name)
				}
			}

			// Prepare the component additional cut
			passCutComp := func() bool { return true }
			if comp.CutFunc.Fct != nil {
				if passCutComp, ok = comp.CutFunc.GetFuncBool(r); !ok {
					err := "Type assertion failed [Cut of (%v, %v)]:"
					err += " TreeFunc.Fct must return a bool.\n"
					err += "\t -> Make sure to use NewCutBool(), not NewVarBool()."
					log.Fatalf(err, samp.Name, comp.FileName)
				}
			}

			// Prepare the cut string for kinematics
			passKinemCut := make([]func() bool, len(ana.KinemCuts))
			for ic, cut := range ana.KinemCuts {
				idx := ic
				if passKinemCut[idx], ok = cut.TreeFunc.GetFuncBool(r); !ok {
					err := "Type assertion failed [selection \"%v\"]:"
					err += " TreeFunc.Fct must return a bool.\n"
					err += "\t -> Make sure to use NewCutBool(), not NewVarBool()."
					log.Fatal(fmt.Sprintf(err, cut.Name))
				}
			}

			// Read the tree (event loop)
			err = r.Read(func(ctx rtree.RCtx) error {

				// Sample-level and component-level cut
				if !(passCutSamp() && passCutComp()) {
					return nil
				}

				// Get the event weight
				w := getWeightSamp() * getWeightComp()

				// Loop over selection and variables
				for ic := range ana.KinemCuts {

					// Look at the next selection if the event is not selected.
					if !passKinemCut[ic]() {
						dump.Var[ana.nVars+ic] = 0.0
						continue
					} else {
						dump.Var[ana.nVars+ic] = 1.0
					}

					// Otherwise, loop over variables.
					for iv, v := range ana.Variables {

						// Fill histo (and fill tree) with full slices...
						if v.isSlice {
							xs := getF64s[iv]()
							for _, x := range xs {
								h[ic][iv].Fill(x, w)
							}
							if ana.DumpTree {
								dump.Vars[iv] = xs
								dump.VarsN[iv] = int32(len(xs))
							}

						} else {
							// ... or the single variable value.
							x := getF64[iv]()
							h[ic][iv].Fill(x, w)
							if ana.DumpTree {
								dump.Var[iv] = x
							}
						}
					}

				}
				
				if ana.DumpTree {
					_, err = tOut.Write()
					if err != nil {
						log.Fatalf("could not write event in a tree: %+v", err)
					}
				}

				return nil
			})

			// Error check of rtree.Reader
			if err != nil {
				log.Fatalf("could not read tree: %+v", err)
			}
			
			// Keep track of the number of processed events.
			switch ana.Nevts {
			case -1:
				ana.nEvtsSample[sampleIdx] += t.Entries()
			default:
				ana.nEvtsSample[sampleIdx] += ana.Nevts
			}
			
			return nil
		}(iComp)
	}

	ana.HbookHistos[sampleIdx] = h
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
				ana.HplotHistos[iSample][iCut][iVar] = ana.Samples[iSample].CreateHisto(h)
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
				stack := hplot.NewHStack(phStack, hplot.WithBand(ana.TotalBand))
				if ana.HistoStack && ana.TotalBand {
					stack.Band.FillColor = ana.ErrBandColor
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
					hps2d_ratioMC.Band.FillColor = ana.ErrBandColor
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

// PrintReport prints some general information about the number
// of processed samples, events and produced histograms.
func (ana Maker) PrintReport() {

	// Event, histo info
	nfiles := 0
	for _, s := range ana.Samples {
		for _ = range s.components {
			nfiles++
		}
	}
	nvars, ncuts := len(ana.Variables), len(ana.KinemCuts)
	nhist := nvars * nfiles
	if ncuts > 0 {
		nhist *= ncuts
	}

	// Time computation
	nkevt := float64(ana.nEvents) / 1e3
	dtLoop := float64(ana.timeLoop) / float64(time.Millisecond)
	dtPlot := float64(ana.timePlot) / float64(time.Millisecond)

	// Formating
	str_template := "\n Processing report:\n"
	str_template += "    - %v histograms filled over %.0f kEvts (%v files, %v variables, %v selections)\n"
	str_template += "    - running time: %.1f ms/kEvt (%s for %.0f kEvts)\n"
	str_template += "    - time fraction: %.0f%% (event loop), %.0f%% (plotting)\n\n"

	fmt.Printf(str_template,
		nhist, nkevt, nfiles, nvars, ncuts,
		(dtLoop+dtPlot)/nkevt, fmtDuration(ana.timeLoop+ana.timePlot), nkevt,
		dtLoop/(dtLoop+dtPlot)*100., dtPlot/(dtLoop+dtPlot)*100.,
	)
}

// RunTimePerKEvts returns the running time in millisecond per kEvents.
func (ana *Maker) RunTimePerKEvts() float64 {
	nkevt := float64(ana.nEvents) / 1e3
	dtLoop := float64(ana.timeLoop) / float64(time.Millisecond)
	dtPlot := float64(ana.timePlot) / float64(time.Millisecond)
	return (dtLoop + dtPlot) / nkevt
}

// Run performs the three steps in one function: fill histos, plot histos
// and print report.
func (ana *Maker) Run() error {

	// Create histograms via event loops
	err := ana.RunEventLoops()
	if err != nil {
		return err
	}

	// Plot each variable x selection, overlaying samples.
	err = ana.PlotVariables()
	if err != nil {
		return err
	}

	// Print processing report
	ana.PrintReport()

	// Return
	return nil
}

// Helper function to initialize histogram containers
func (ana *Maker) initHistoContainers() {

	// Initialize hbook H1D as N[samples] 2D-slices.
	// Cut x variable initialization is done in fillSampleHistos().
	ana.HbookHistos = make([][][]*hbook.H1D, len(ana.Samples))

	// Inititialize hplot H1D
	ana.HplotHistos = make([][][]*hplot.H1D, len(ana.Samples))
	for iSamp := range ana.Samples {
		ana.HplotHistos[iSamp] = make([][]*hplot.H1D, len(ana.KinemCuts))
		for iCut := range ana.KinemCuts {
			ana.HplotHistos[iSamp][iCut] = make([]*hplot.H1D, len(ana.Variables))
		}
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

func (ana *Maker) assessVariableTypes() {

	fName := ana.Samples[0].components[0].FileName
	tName := ana.Samples[0].components[0].TreeName
	f, t := getTreeFromFile(fName, tName)
	defer f.Close()
	r, err := rtree.NewReader(t, rtree.NewReadVars(t))
	if err != nil {
		log.Fatal("could not create tree reader: %w", err)
	}
	defer r.Close()
	for _, v := range ana.Variables {
		v.isSlice = false
		if _, ok := v.TreeFunc.GetFuncF64(r); !ok {
			v.isSlice = true
			if _, ok = v.TreeFunc.GetFuncF64s(r); !ok {
				err := "Type assertion failed [variable \"%v\"]:"
				err += " TreeFunc.Fct must return a float64 or a []float64."
				log.Fatalf(err, v.Name)
			}
		}
	}
}

func (ana *Maker) newDumper() dumper {
	dVar := make([]float64, len(ana.Variables)+len(ana.KinemCuts))
	dVars := make([][]float64, len(ana.Variables))
	dVarsN := make([]int32, len(ana.Variables))
	for i, v := range ana.Variables {
		if v.isSlice {
			dVars[i] = []float64{}
			dVarsN[i] = 0
		} else {
			dVar[i] = 0
		}
	}
	for i := range ana.KinemCuts {
		dVar[ana.nVars+i] = 0
	}

	// Return the dumper
	return dumper{
		Var: dVar,
		Vars: dVars,
		VarsN: dVarsN,
	}
}

func (ana *Maker) getOutFileTree(fname, tname string, d dumper) (*groot.File, rtree.Writer) {

	// Create a new ROOT file
	f, err := groot.Create(fname)
	if err != nil {
		log.Fatalf("could not create ROOT file %v: %v", fname, err)
	}

	// Variables to save
	wvars := []rtree.WriteVar{}
	for i, v := range ana.Variables {
		if v.isSlice {
			wvars = append(wvars, rtree.WriteVar{
				Name:  v.Name + "N",
				Value: &d.VarsN[i]},
			)
			wvars = append(wvars, rtree.WriteVar{
				Name:  v.Name,
				Value: &d.Vars[i],
				Count: v.Name + "N"},
			)
		} else {
			wvars = append(wvars, rtree.WriteVar{
				Name:  v.Name,
				Value: &d.Var[i]},
			)
		}
	}
	for i, s := range ana.KinemCuts {
		wvars = append(wvars, rtree.WriteVar{
			Name:  "pass" + s.Name,
			Value: &d.Var[ana.nVars+i]},
		)
	}

	// Create a new TTree
	t, err := rtree.NewWriter(f, tname, wvars)
	if err != nil {
		log.Fatal("could not create tree writer: %w", err)
	}

	return f, t
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
