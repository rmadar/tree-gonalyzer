// Package allowing to wrap all needed element of a TTree plotting analysis
package ana

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

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

	// Inputs info
	Samples      []Sample         // Sample on which to run
	SamplesGroup string           // Specify how to group samples together
	Variables    []*Variable      // List of variables to plot
	Cuts         []Selection      // List of cuts

	// Figure related setup
	SaveFormat   string           // Extension of saved figure 'tex', 'pdf', 'png'
	CompileLatex bool             // Enable on-the-fly latex compilation of plots

	// Plot related setup
	RatioPlot    bool             // Enable ratio plot
	DontStack    bool             // Disable histogram stacking (e.g. compare various processes)
	Normalize    bool             // Normalize distributions to unit area (when stacked, the total is normalized)

	// Histograms
	HbookHistos  [][][]*hbook.H1D // Currently 3D histo container
	HplotHistos  [][][]*hplot.H1D // Currently 3D histo container

	// Temp
	WithTreeFormula bool   // TEMP for benchmarking

	cutIdx        map[string]int // Linking cut name and cut index
	sampleIdx     map[string]int // Linking sample name and sample index
	variableIdx   map[string]int // Linking variable name and variable index
	histoFilled bool           // true if histograms are filled.
	nEvents  int64             // Number of processed events
	timeLoop time.Duration     // Processing time for filling histograms (event loop over samples x cuts x histos)
	timePlot time.Duration     // Processing time for plotting histogram

}

// Creating a new object
// TO-DO: switch to all pointers []*Sample, []*Variables
func New(s []Sample, v []*Variable) Maker {
	return Maker{
		Samples: s,
		Variables: v,
	}
}

// Helper function creating the mapping between name and objects
func getNameIndices(obj interface{}) map[string]int {
	// obj = []variables, []samples, []cuts
	return make(map[string]int, 10)
}

// Run the event loop to fill all histo across samples / variables / cuts (and later: systematics)
func (ana *Maker) MakeHistos() error {

	// Start timing
	start := time.Now()

	// Build hbook histograms container
	ana.initHbookHistos()

	// Loop over the samples
	for is, s := range ana.Samples {

		// Anonymous function to avoid memory-leaks due to 'defer'
		func(j int) error {

			// Get the file and tree
			f, t := getTreeFromFile(s.FileName, s.TreeName)
			defer f.Close()
		
			var rvars []rtree.ReadVar
			if !ana.WithTreeFormula {
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
			if ana.WithTreeFormula {
				for i, v := range ana.Variables {
					varFormula[i] = v.TreeFunc.GetVarFunc(r)
				}
			}

			// Prepare the weight
			getWeight := func() float64 { return float64(1.0) }
			if s.Weight != "" {
				if ana.WithTreeFormula {
					getWeight = s.WeightFunc.GetVarFunc(r)
				}
			}
			
			// Prepare the sample cut
			passSampleCut := func() bool { return true }
			if s.Cut != "" {
				if ana.WithTreeFormula {
					passSampleCut = s.CutFunc.GetCutFunc(r)
				}
			}

			// Prepare the cut string for kinematics
			passKinemCut := make([]func() bool, len(ana.Cuts))
			for ic, cut := range ana.Cuts {
				passKinemCut[ic] = func() bool { return true }
				idx := ic
				if cut.TreeName != "true" {
					if ana.WithTreeFormula {
						passKinemCut[idx] = cut.TreeFunc.GetCutFunc(r)
					}
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
					
					if !passKinemCut[ic]() {
						continue
					}
					
					for iv, v := range ana.Variables {
						val := 0.0
						if ana.WithTreeFormula {
							val = varFormula[iv]()
						} else {
							val = v.GetValue()
						}
						ana.HbookHistos[iv][ic][is].Fill(val, w)
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

	// Return an error if HbookHistos is empty
	if !ana.histoFilled {
		log.Fatalf("There is no histograms. Please make sure that 'MakeHistos()' is called before 'PlotHistos()'")
	}

	// Preparing the final figure
	var plt hplot.Drawer
	figWidth, figHeight := 6*vg.Inch, 4.5*vg.Inch
	format := "tex"
	if ana.SaveFormat != "" {
		format = ana.SaveFormat
	}

	// Handle on-the-fly LaTeX compilation
	var latex htex.Handler = htex.NoopHandler{}
	if ana.CompileLatex {
		latex = htex.NewGoHandler(-1, "pdflatex")
	}

	// Inititialize all hplot.H1D histograms
	ana.HplotHistos = make([][][]*hplot.H1D, len(ana.Variables))
	for iv := range ana.HplotHistos {
		ana.HplotHistos[iv] = make([][]*hplot.H1D, len(ana.Cuts))
		for isel := range ana.Cuts {
			ana.HplotHistos[iv][isel] = make([]*hplot.H1D, len(ana.Samples))
		}
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

			// Apply common and user-defined style for this variable
			style.ApplyToPlot(p)
			v.SetPlotStyle(p)

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
				if ana.Normalize {
					if ana.Samples[is].IsData() {
						h.Scale(1 / norm_histos[is])
					}
					if ana.Samples[is].IsBkg() {
						if ana.DontStack {
							h.Scale(1. / norm_histos[is])
						} else {
							h.Scale(1. / norm_bkgtot)
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
				stack.Band.FillColor = color.NRGBA{R: 255, G: 255, B: 255, A: 150}
				if ana.DontStack {
					stack.Stack = hplot.HStackOff
				} else {
					hBand := hplot.NewH1D(hbook.NewH1D(1, 0, 1), hplot.WithBand(true))
					hBand.Band = stack.Band 
					p.Legend.Add("Uncer.", hBand)
				}

				// Add the stack to the plot
				p.Add(stack)
				
			}
			
			
			// Adding hplot.H1D data to the plot, set the drawer to the current plot
			if bhData.Entries() > 0 {
				p.Add(phData)
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
				case ana.DontStack:
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
				default:
					// Data to MC
					hbs2d_ratio, err := hbook.DivideH1D(bhData, bhBkgTot, hbook.DivIgnoreNaNs())
					if err != nil {
						log.Fatal("cannot divide histo for the ratio plot")
					}
					hps2d_ratio := hplot.NewS2D(hbs2d_ratio, hplot.WithYErrBars(true),
						hplot.WithStepsKind(hplot.HiSteps),
					)
					style.ApplyToDataS2D(hps2d_ratio)

					// MC to MC
					hbs2d_ratio1, err := hbook.DivideH1D(bhBkgTot, bhBkgTot, hbook.DivIgnoreNaNs())
					if err != nil {
						log.Fatal("cannot divide histo for the ratio plot")
					}
					hps2d_ratio1 := hplot.NewS2D(hbs2d_ratio1, hplot.WithBand(true),
						hplot.WithStepsKind(hplot.HiSteps),
					)
					style.ApplyToDataS2D(hps2d_ratio1)
					hps2d_ratio1.GlyphStyle.Radius = 0
					hps2d_ratio1.LineStyle.Width = 0.0
					hps2d_ratio1.LineStyle.Color = color.NRGBA{R: 140, G: 140, B: 140, A: 255}
					hps2d_ratio1.Band.FillColor = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
					rp.Bottom.Add(hps2d_ratio1)
					rp.Bottom.Add(hps2d_ratio)
				}

				// Adjust ratio plot scale
				rp.Bottom.Y.Min = 0.7
				rp.Bottom.Y.Max = 1.3
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

// Initialize histograms container shape
func (ana *Maker) initHbookHistos() {

	if len(ana.Cuts) == 0 {
		ana.Cuts = append(ana.Cuts, Selection{Name: "No-cut", TreeName: "true"})
	}

	ana.HbookHistos = make([][][]*hbook.H1D, len(ana.Variables))
	for iv := range ana.HbookHistos {
		ana.HbookHistos[iv] = make([][]*hbook.H1D, len(ana.Cuts))
		for isel := range ana.Cuts {
			ana.HbookHistos[iv][isel] = make([]*hbook.H1D, len(ana.Samples))
			v := ana.Variables[iv]
			for isample := range ana.HbookHistos[iv][isel] {
				ana.HbookHistos[iv][isel][isample] = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
			}
		}
	}
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

// Helper duration formating: return a string 'hh:mm:ss' for a time.Duration object
func fmtDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
