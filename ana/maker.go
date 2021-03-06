// Package allowing to wrap all needed element of a TTree plotting analysis
package ana

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"go-hep.org/x/hep/hbook"
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
	NevtsMax  int64        // Maximum event number per components (default: -1),
	Lumi      float64      // Integrated luminosity en 1/fb (default: 1/pb).
	SampleMT  bool         // Enable concurency accross samples (default: true).

	// Ouputs
	SavePath     string // Path to which plot will be saved (default: 'outputs').
	SaveFormat   string // Plot file extension: 'png' (default), 'pdf' or 'png'.
	CompileLatex bool   // On-the-fly latex compilation (default: true).
	DumpTree     bool   // Dump a TTree in a file for each sample (default: false).
	PlotHisto    bool   // Enable histogram plotting (default: true).

	// Plots
	AutoStyle      bool        // Enable automatic styling (default: true).
	PlotTitle      string      // General plot title (default: 'TTree GOnalyzer').
	HistoStack     bool        // Enable histogram stacking (default: true).
	SignalStack    bool        // Enable signal stack (default: false).
	HistoNorm      bool        // Normalize distributions to unit area (default: false).
	TotalBand      bool        // Enable total error band in stack mode (default: true).
	TotalBandColor color.NRGBA // Color for the uncertainty band (default: gray).

	// Enable ratio plot (default: true).
	// If stack is on, the ratio is defined as data over total bkg.
	// If stack is off, ratios are defined as sample[i] / data when
	// a data sample is defined, sample[i] / sample[0] otherwise.
	RatioPlot bool

	// Histograms for {samples x selections x variables}
	hbookHistos [][][]*hbook.H1D

	// tree dumping
	nVars       int     // number of variables
	nEvtsSample []int64 // number of events per sample

	// Normalisation of each sample for each cut:
	// {selections x samples}
	normHists [][]float64

	// Normalisation of total background (and signal if
	// stacked) for each cut
	normTotal []float64

	idxData     []int         // Indices of data samples in []*Sample slice
	idxBkgs     []int         // Indices of bkg samples in []*Sample slice
	idxSigs     []int         // Indices of sig samples in []*Sample slice
	histoFilled bool          // true if histograms are filled.
	nEvents     int64         // Number of processed events
	timeLoop    time.Duration // Processing time for filling histograms (event loop over samples x cuts x histos)
	timePlot    time.Duration // Processing time for plotting histogram
}

// New creates a default analysis maker from a list of sample
// and a list of variables.
func New(s []*Sample, v []*Variable, opts ...Options) Maker {

	// Create the object
	a := Maker{
		Samples:        s,
		Variables:      v,
		NevtsMax:       -1,
		Lumi:           1e-3,
		PlotHisto:      true,
		SampleMT:       true,
		AutoStyle:      true,
		SavePath:       "outputs",
		SaveFormat:     "png",
		PlotTitle:      "TTree GOnalyzer",
		CompileLatex:   true,
		HistoStack:     true,
		RatioPlot:      true,
		TotalBand:      true,
		TotalBandColor: color.NRGBA{A: 100},
		KinemCuts:      []*Selection{EmptySelection()},
		nVars:          len(v),
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
	if cfg.NevtsMax.usr {
		a.NevtsMax = cfg.NevtsMax.val
	}
	if cfg.Lumi.usr {
		a.Lumi = cfg.Lumi.val
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
	if cfg.TotalBandColor.usr {
		a.TotalBandColor = cfg.TotalBandColor.val
	}

	// Get ordered lists of background and signal names
	a.idxData, a.idxBkgs, a.idxSigs = a.getSampleProc()

	// Managing event number with concurrency
	a.nEvtsSample = make([]int64, len(a.Samples))

	// Build the slice of values to store
	// FIX-ME(rmadar): this is not so clean to assess slice or not
	//                 by doing a loop over variables for the first
	//                 component of the first sample to fill v.isSlice.
	a.assessVariableTypes()

	return a
}

// Helper function to get ordered list of background
// and signal names
func (ana *Maker) getSampleProc() ([]int, []int, []int) {
	iData, iBkgs, iSigs := []int{}, []int{}, []int{}
	for i, s := range ana.Samples {
		switch s.sType {
		case data:
			iData = append(iData, i)
		case bkg:
			iBkgs = append(iBkgs, i)
		case sig:
			iSigs = append(iSigs, i)
		}
	}

	if len(iData) > 1 {
		log.Fatalf("More than one data sample")
	}

	return iData, iBkgs, iSigs
}

// PrintReport prints some general information about the number
// of processed samples, events and produced histograms.
func (ana Maker) PrintReport() {

	// Event and histo info
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

// PrintSlowTreeFuncs prints the list of TreeFunc which relies on
// a generic groot/rfunc formula, ie based on 'reflect' calls. These
// function are ~ 5 times slower than the one defined using this example
// https://godoc.org/go-hep.org/x/hep/groot/rtree#example-Reader--WithFormulaFromUser
func (ana *Maker) PrintSlowTreeFuncs() {

	appendSlow := func(fs *[]TreeFunc, f TreeFunc) {
		if f.IsSlow() {
			*fs = append(*fs, f)
		}
	}

	// Store all slow function
	slowFs := &[]TreeFunc{}

	// Variables
	for _, v := range ana.Variables {
		appendSlow(slowFs, v.TreeFunc)
	}

	// Kinematic cuts.
	for _, c := range ana.KinemCuts {
		appendSlow(slowFs, c.TreeFunc)
	}

	// Samples and component cuts & weights
	for _, s := range ana.Samples {
		appendSlow(slowFs, s.CutFunc)
		appendSlow(slowFs, s.WeightFunc)
		for _, c := range s.components {
			appendSlow(slowFs, c.CutFunc)
			appendSlow(slowFs, c.WeightFunc)
		}
	}

	// Print only if there is at least one slow func
	if len(*slowFs) > 0 {
		fmt.Println(" List of slow TreeFuncs:")
		for _, f := range *slowFs {
			fmt.Printf("    - %T --> args = %v \n", f.Fct, f.VarsName)
		}
		fmt.Println("")
	}

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

	// Print slow functions
	ana.PrintSlowTreeFuncs()

	// Return
	return nil
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
