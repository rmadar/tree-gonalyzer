// Package allowing to wrap all needed element of a TTree plotting analysis
package ana

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/plot/vg"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/hplot/htex"

	"github.com/rmadar/hplot-style/style"
)

// Analyzer type
type Maker struct {
	Samples      []Sample         // Sample on which to run
	SamplesGroup string           // Specify how to group samples together
	Variables    []*Variable      // List of variables to plot
	Cuts         []Selection      // List of cuts
	SaveFormat   string           // Extension of saved figure 'tex', 'pdf', 'png'
	RatioPlot    bool             // Enable ratio plot
	CompileLatex bool             // Enable on-the-fly latex compilation of plots
	HistosData   [][][]*hbook.H1D // Currently 3D histo container, later: n-dim [var, sample, cut, syst]
	HistosPlot   [][][]*hplot.H1D // Currently 3D histo container, later: n-dim [var, sample, cut, syst]
	DontStack    bool             // Disable histogram stacking (e.g. compare various processes)
	Normalize    bool             // Normalize distributions to unit area (when stacked, the total is nomarlized)

	WithTreeFormula bool // TEMP for benchmarking

	nEvents  int64         // Number of processed events
	timeLoop time.Duration // Processing time for filling histograms (event loop over samples x cuts x histos)
	timePlot time.Duration // Processing time for plotting histogram

}

// Initialize histograms container shape
func (ana *Maker) initHistosData() {

	if len(ana.Cuts) == 0 {
		ana.Cuts = append(ana.Cuts, Selection{Name: "No-cut", TreeName: "true"})
	}

	ana.HistosData = make([][][]*hbook.H1D, len(ana.Variables))
	for iv := range ana.HistosData {
		ana.HistosData[iv] = make([][]*hbook.H1D, len(ana.Cuts))
		for isel := range ana.Cuts {
			ana.HistosData[iv][isel] = make([]*hbook.H1D, len(ana.Samples))
			v := ana.Variables[iv]
			for isample := range ana.HistosData[iv][isel] {
				ana.HistosData[iv][isel][isample] = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
			}
		}
	}
}

// Run the event loop to fill all histo across samples / variables / cuts (and later: systematics)
func (ana *Maker) MakeHistos() error {

	// Start timing
	start := time.Now()

	// Build hbook histograms container
	ana.initHistosData()

	// Loop over the samples
	for is, s := range ana.Samples {

		// Anonymous function to avoid memory-leaks due to 'defer'
		func(j int) error {

			// Get the file and tree
			f, t := getTreeFromFile(s.FileName, s.TreeName)
			defer f.Close()

			tree := rtree.Chain(t) /*, t, t, t,
			t, t, t, t, t, t, t, t,
			t, t, t, t, t, t, t, t,
			t, t, t, t, t, t, t, t,
			t, t, t, t, t, t, t, t,
			t, t, t, t, t, t, t, t,
			t, t, t, t, t, t, t, t)*/

			var rvars []rtree.ReadVar
			if !ana.WithTreeFormula {
				for _, v := range ana.Variables {
					rvars = append(rvars, rtree.ReadVar{Name: v.TreeName, Value: v.Value})
				}
			}

			// Get the tree reader
			r, err := rtree.NewReader(tree, rvars)
			if err != nil {
				return fmt.Errorf("could not create tree reader: %w", err)
			}
			defer r.Close()

			var_formula := make([]rtree.Formula, len(ana.Variables))
			if ana.WithTreeFormula {
				var errForm error
				for i, v := range ana.Variables {
					var_formula[i], errForm = r.Formula("float64("+v.TreeName+")", nil)
					if errForm != nil {
						log.Fatalf("could not create formula: %+v", errForm)
					}
				}
			}

			// Prepare the weight
			var wform rtree.Formula
			getWeight := func() float64 { return float64(1.0) }
			if s.Weight != "" {
				wform, err = r.Formula("float64("+s.Weight+")", nil)
				if err != nil {
					log.Fatalf("could not create sample weight formula: %+v", err)
				}
				getWeight = func() float64 { return wform.Eval().(float64) }
			}

			// Prepare the sample cut
			var cutSampleform rtree.Formula
			passSampleCut := func() bool { return true }
			if s.Cut != "" {
				cutSampleform, err = r.Formula("bool("+s.Cut+")", nil)
				if err != nil {
					log.Fatalf("could not create sample cut formula: %+v", err)
				}
				passSampleCut = func() bool { return cutSampleform.Eval().(bool) }
			}

			// Prepare the cut string for kinematics
			cutKinem := make([]rtree.Formula, len(ana.Cuts))
			for ic, cut := range ana.Cuts {
				cutKinem[ic], err = r.Formula("bool("+cut.TreeName+")", nil)
				if err != nil {
					log.Fatalf("could not create kinem cut formula: %+v", err)
				}
			}

			// FIXME(rmadar): this doesn't work and keep only the last rformula
			passKinemCut := make([]func() bool, len(ana.Cuts))
			for ic, cut := range ana.Cuts {
				passKinemCut[ic] = func() bool { return true }
				if cut.TreeName != "true" {
					passKinemCut[ic] = func() bool { return cutKinem[ic].Eval().(bool) }
				}
			}

			// Read the tree (event loop)
			err = r.Read(func(ctx rtree.RCtx) error {

				// Sample-level cut
				if !passSampleCut() {
					return nil
				}

				// Get the event weight
				w := getWeight()

				// Loop over selection and variables
				for ic := range ana.Cuts {

					// FIXME(rmadar): this doesn't work and keep only the last rformula
					// This would allow to avoid calling rformula when this is the default
					// ie no cut, ie TreeName = "true".
					// if !passKinemCut[isel]() { continue }

					// Instead calling .Eval() in every cases
					if !cutKinem[ic].Eval().(bool) {
						continue
					}

					for iv, v := range ana.Variables {
						val := v.GetValue()
						if ana.WithTreeFormula {
							val = var_formula[iv].Eval().(float64)
						}
						ana.HistosData[iv][ic][is].Fill(val, w)
					}
				}

				return nil
			})

			// Keep track of the number of processed events
			ana.nEvents += tree.Entries()

			return nil
		}(is)
	}

	// End timing
	ana.timeLoop = time.Now().Sub(start)

	return nil
}

// Plotting all histograms
func (ana *Maker) PlotHistos() error {

	// Start timing
	start := time.Now()

	// Plot format
	format := "tex"
	if ana.SaveFormat != "" {
		format = ana.SaveFormat
	}

	// Handle on-the-fly LaTeX compilation
	var latex htex.Handler = htex.NoopHandler{}
	if ana.CompileLatex {
		latex = htex.NewGoHandler(-1, "pdflatex")
	}

	// Inititialize histograms
	ana.HistosPlot = make([][][]*hplot.H1D, len(ana.Variables))
	for iv := range ana.HistosPlot {
		ana.HistosPlot[iv] = make([][]*hplot.H1D, len(ana.Cuts))
		for isel := range ana.Cuts {
			ana.HistosPlot[iv][isel] = make([]*hplot.H1D, len(ana.Samples))
		}
	}

	// Return an error if all normalization are 0
	if len(ana.HistosData) == 0 {
		log.Fatalf("There is no histograms. Please make sure that 'MakeHistos()' is called before 'PlotHistos()'")
	}

	// Loop over variables and get histo for all samples
	for iv, h_sel_samples := range ana.HistosData {

		// Manipulate the current variable
		v := ana.Variables[iv]

		// Loop over selections
		for isel, hsamples := range h_sel_samples {

			// Create a new styled plot and figure
			p := hplot.New()
			style.ApplyToPlot(p)

			// Propagate the user-defined style of the variable to the plot
			v.SetPlotStyle(p)

			// Additionnal legend style
			p.Legend.Padding = 0.1 * vg.Inch
			p.Legend.ThumbnailWidth = 25
			p.Legend.TextStyle.Font.Size = 12

			// Prepare histogram (possible) stacking via []*hplot.H1D
			var (
				hbkgs []*hplot.H1D
				hdata *hplot.H1D
				norms []float64
			)

			// Keep track of the normalization for every sample
			for _, h := range hsamples {
				norms = append(norms, h.Integral())
			}
			Nbkg := floats.Sum(norms[1:]) // still assume that data is i==1

			// Loop over samples
			hTotData := hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
			hDataData := hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
			for is, h := range hsamples {

				// Deal with normalize option for data
				if is == 0 && ana.Normalize {
					h.Scale(1 / norms[is])
				}

				// Deal with normalize option of non-data
				if is > 0 && ana.Normalize {
					if ana.DontStack {
						h.Scale(1 / norms[is])
					} else {
						h.Scale(1. / Nbkg)
					}
				}

				// Get the sum of all histos
				if !ana.Samples[is].IsData() {
					hTotData = hbook.AddH1D(h, hTotData)
				}

				// Get plottable histogram
				ana.HistosPlot[iv][isel][is] = ana.Samples[is].CreateHisto(h)

				// Prepare the legend
				p.Legend.Add(ana.Samples[is].LegLabel, ana.HistosPlot[iv][isel][is])

				// Keep data appart from backgrounds (FIX-ME: assumed to be the first sample for now)
				if ana.Samples[is].IsData() {
					hdata = ana.HistosPlot[iv][isel][is]
					hDataData = h
				}
				if ana.Samples[is].IsBkg() {
					hbkgs = append(hbkgs, ana.HistosPlot[iv][isel][is])
				}
			}

			// Manage stack
			if len(hbkgs) > 1 {
				for i, j := 0, len(hbkgs)-1; i < j; i, j = i+1, j-1 {
					hbkgs[i], hbkgs[j] = hbkgs[j], hbkgs[i]
				}
				stack := hplot.NewHStack(hbkgs)
				if ana.DontStack {
					stack.Stack = hplot.HStackOff
				}

				// Add the histgrams (possibly stack) and data
				p.Add(stack)
			}
			p.Add(hdata)

			var plt hplot.Drawer
			plt = p
			figWidth, figHeight := 6*vg.Inch, 4.5*vg.Inch

			// Add ratio plot
			if ana.RatioPlot {

				// Update figure sizes
				figWidth, figHeight = 6*vg.Inch, 4.5*vg.Inch

				// Build up the ratio histo
				s2d_ratio_data, err := hbook.DivideH1D(hDataData, hTotData, hbook.DivIgnoreNaNs())
				if err != nil {
					log.Fatal("cannot divide histo for the ratio plot")
				}
				s2d_ratio := hplot.NewS2D(s2d_ratio_data, hplot.WithYErrBars(true))
				style.ApplyToDataS2D(s2d_ratio)

				// Create ratio plot type
				rp := hplot.NewRatioPlot()

				// Deal with bottom pannel
				style.ApplyToBottomPlot(rp.Bottom)
				rp.Bottom.Add(s2d_ratio)
				rp.Bottom.X = p.X
				//rp.Bottom.Y.Min = 0.95
				//rp.Bottom.Y.Max = 1.05

				// Deal with Top pannel
				rp.Top = p
				rp.Top.HideX()
				rp.Top.X.Label.Text = ""

				// Update the drawer
				plt = rp
			}

			// Save the figure
			f := hplot.Figure(plt)
			style.ApplyToFigure(f)
			f.Latex = latex

			path := "results/" + ana.Cuts[isel].Name
			if _, err := os.Stat(path); os.IsNotExist(err) {
				os.MkdirAll(path, 0755)
			}
			outputname := path + "/" + v.SaveName + "." + format
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

// Helper to get a tree from a file
func getTreeFromFile(filename, treename string) (*groot.File, rtree.Tree) {

	// Get the file
	f, err := groot.Open(filename)
	if err != nil {
		err := fmt.Sprintf("could not open ROOT file %q: %w", filename, err)
		panic(err)
	}

	// Get the tree
	obj, err := f.Get(treename)
	if err != nil {
		err := fmt.Sprintf("could not retrieve object: %w", err)
		panic(err)
	}

	return f, obj.(rtree.Tree)
}

// Print processing report
func (ana Maker) PrintReport() {

	// Event, histo info
	nvars, nsamples, ncuts := len(ana.Variables), len(ana.Samples), len(ana.Cuts)
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

// Helper duration formating: return a string 'hh:mm:ss' for a time.Duration object
func fmtDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// Helper to add two histograms - uncertainty propagation is missing
func addHistos(h1, h2 *hbook.H1D, m float64) *hbook.H1D {
	hres := hbook.NewH1D(h1.Len(), h1.XMin(), h1.XMax())
	for i := 0; i < hres.Len(); i++ {
		v1 := h1.Value(i)
		v2 := h2.Value(i)
		x1, _ := h1.XY(i)
		hres.Fill(x1, v1+m*v2)
	}
	return hres
}

// Helper to divide two histograms - uncertainty propagation is missing
func divideHistos(hnum, hden *hbook.H1D) *hbook.H1D {
	hres := hbook.NewH1D(hnum.Len(), hnum.XMin(), hnum.XMax())
	for i := 0; i < hres.Len(); i++ {
		vnum := hnum.Value(i)
		vden := hden.Value(i)
		x, _ := hnum.XY(i)
		ratio := vnum / vden
		if math.IsNaN(ratio) || math.IsInf(ratio, 0) {
			hres.Fill(x, 0)
		} else {
			hres.Fill(x, ratio)
		}
	}
	return hres
}
