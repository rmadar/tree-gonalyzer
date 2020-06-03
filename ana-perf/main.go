// Package allowing to benchmark performances of the ana package.
package main

import (
	"fmt"
	"log"
	"math"

	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"

	"github.com/rmadar/hplot-style/style"
	"github.com/rmadar/tree-gonalyzer/ana"
)

// Run all the tests
func main() {

	// Scan the number of variables
	nVars := []float64{1, 5, 14, 20, 25, 30, 35, 40, 45, 50, 55, 60}

	// Scan number of events
	nEvts1, nEvts2, nEvts3, nEvts4 := 10, 20, 50, 100

	// Containers
	t1 := make([]float64, len(nVars))
	t2 := make([]float64, len(nVars))
	t3 := make([]float64, len(nVars))
	t4 := make([]float64, len(nVars))

	// Run all test
	for i, n := range nVars {
		//if n != 60 {
		//	continue
		//}
		fmt.Println("Running for nVars =", n)
		t1[i] = runTest(nEvts1, int(n))
		fmt.Printf("  -> 0.5M evts: %.2f ms/kEvts\n", t1[i])
		t2[i] = runTest(nEvts2, int(n))
		fmt.Printf("  -> 1.0M evts: %.2f ms/kEvts\n", t2[i])
		t3[i] = runTest(nEvts3, int(n))
		fmt.Printf("  -> 2.0M evts: %.2f ms/kEvts\n", t3[i])
		t4[i] = runTest(nEvts4, int(n))
		fmt.Printf("  -> 2.5M evts: %.2f ms/kEvts\n", t4[i])
	}

	// Plot benchmarks
	plotBenchmarks(t1, t2, t3, t4, nVars)

}

func plotBenchmarks(s1, s2, s3, s4, n []float64) {

	// Plot
	p := hplot.New()
	style.ApplyToPlot(p)
	p.Title.Text = "GOnalyzer Benchmarks"
	p.X.Label.Text = "Number of variables"
	p.Y.Label.Text = "Running Time [ms / kEvts]"
	p.Y.Min = 0.0

	// Graph
	g1 := hplot.NewS2D(hbook.NewS2DFrom(n, s1))
	g2 := hplot.NewS2D(hbook.NewS2DFrom(n, s2))
	g3 := hplot.NewS2D(hbook.NewS2DFrom(n, s3))
	g4 := hplot.NewS2D(hbook.NewS2DFrom(n, s4))

	// Comsetics
	applyStyle(g1, 0)
	applyStyle(g2, 1)
	applyStyle(g3, 2)
	applyStyle(g4, 3)

	// Add graph to the legend
	p.Legend.Add(`0.5M Evts`, g1)
	p.Legend.Add(`1.0M Evts`, g2)
	p.Legend.Add(`2.5M Evts`, g3)
	p.Legend.Add(`5.0M Evts`, g4)
	p.Legend.Top = true
	p.Legend.Left = true
	p.Legend.XOffs = 12
	p.Legend.YOffs = -8

	// Add graph to the plot
	p.Add(g1)
	p.Add(g2)
	p.Add(g3)
	p.Add(g4)

	// Figure
	f := hplot.Figure(p)
	style.ApplyToFigure(f)
	if err := hplot.Save(f, 10*vg.Inch, 4*vg.Inch, "perf.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}

// Helper to set S2D style
func applyStyle(g *hplot.S2D, icolor int) {
	g.LineStyle.Width = 2
	g.GlyphStyle.Radius = 0
	g.LineStyle.Color = plotutil.Color(icolor)
}

// Run one test and returns the time in ms/kEvts
func runTest(n10kEvtsPerSample, nVariables int) float64 {

	// Data
	splData := ana.NewSample("data", "data", `Pseudo-data`)
	loadManyComponents(splData, n10kEvtsPerSample)

	// Background 1
	splBkg1 := ana.NewSample("bkg1", "bkg", `Proc 1`, ana.WithWeight(w1))
	loadManyComponents(splBkg1, n10kEvtsPerSample)

	// Background 2
	splBkg2 := ana.NewSample("bkg2", "bkg", `Proc 2`, ana.WithWeight(w2))
	loadManyComponents(splBkg2, n10kEvtsPerSample)

	// Background 3
	splBkg3 := ana.NewSample("bkg3", "bkg", `Proc 3`, ana.WithWeight(w1))
	loadManyComponents(splBkg3, n10kEvtsPerSample)

	// Background 4
	splBkg4 := ana.NewSample("bkg4", "bkg", `Proc 4`, ana.WithWeight(w2))
	loadManyComponents(splBkg4, n10kEvtsPerSample)

	// Group samples together
	samples := []*ana.Sample{splData, splBkg1, splBkg2, splBkg3, splBkg4}

	// Variables, organized in bunch of 17
	variables := []*ana.Variable{

		/*var_dphi, var_Crr, var_Cnn, var_Ckk, var_dphi,
		var_dphi, var_Crr, var_Cnn, var_Ckk, var_dphi,
		var_dphi, var_Crr, var_Cnn, var_Ckk, var_dphi,
		var_dphi, var_Crr, var_Cnn, var_Ckk, var_dphi,

		var_dphi, var_Crr, var_Cnn, var_Ckk, var_dphi,
		var_dphi, var_Crr, var_Cnn, var_Ckk, var_dphi,
		var_dphi, var_Crr, var_Cnn, var_Ckk, var_dphi,
		var_dphi, var_Crr, var_Cnn, var_Ckk, var_dphi,

		var_dphi, var_Crr, var_Cnn, var_Ckk, var_dphi,
		var_dphi, var_Crr, var_Cnn, var_Ckk, var_dphi,
		var_dphi, var_Crr, var_Cnn, var_Ckk, var_dphi,
		var_dphi, var_Crr, var_Cnn, var_Ckk, var_dphi,
		*/

		var_dphi, var_m_tt, var_eta_t, var_pt_lep, var_Ckk,
		var_Crr, var_Cnn, var_pt_lep, var_eta_lep, var_pt_b,
		var_eta_b, var_pt_vsum, var_pt_t, var_pt_tt, var_x1,
		var_x1x2, var_GG,

		var_dphi, var_m_tt, var_eta_t, var_pt_lep, var_Ckk,
		var_Crr, var_Cnn, var_pt_lep, var_eta_lep, var_pt_b,
		var_eta_b, var_pt_vsum, var_pt_t, var_pt_tt, var_x1,
		var_x1x2, var_GG,

		var_dphi, var_m_tt, var_eta_t, var_pt_lep, var_Ckk,
		var_Crr, var_Cnn, var_pt_lep, var_eta_lep, var_pt_b,
		var_eta_b, var_pt_vsum, var_pt_t, var_pt_tt, var_x1,
		var_x1x2, var_GG,

		var_dphi, var_m_tt, var_eta_t, var_pt_lep, var_Ckk,
		var_Crr, var_Cnn, var_pt_lep, var_eta_lep, var_pt_b,
		var_eta_b, var_pt_vsum, var_pt_t, var_pt_tt, var_x1,
		var_x1x2, var_GG,
	}

	// Protection for too high number of variables
	nVars := len(variables)
	if nVariables > -1 {
		nVars = nVariables
	}
	if nVars > len(variables) {
		panic(fmt.Errorf("Too much variables (max %v, got %v)", len(variables), nVars))
	}

	// Create analyzer object with options
	analyzer := ana.New(samples, variables[:nVars],
		//ana.WithDumpTree(true),
		//ana.WithPlotHisto(false),
		ana.WithSaveFormat("tex"),
		ana.WithCompileLatex(false),
		ana.WithHistoNorm(true),
		ana.WithHistoStack(true),
		ana.WithSampleMT(true),
	)

	// Run the analyzer and produce all plots
	if err := analyzer.RunEventLoops(); err != nil {
		log.Fatal("Cannot fill histos:", err)
	}
	if err := analyzer.PlotVariables(); err != nil {
		log.Fatal("Cannot plot histos:", err)
	}

	return analyzer.RunTimePerKEvts()
}

// Define all samples and variables of the analysis
var (
	// files/trenames
	file1 = "../testdata/file1.root"
	file2 = "../testdata/file2.root"
	file3 = "../testdata/file3.root"
	tname = "truth"

	// Some TreeFunc: weights and cuts
	w1   = ana.TreeValF64(0.5)
	w2   = ana.TreeValF64(2.0)
	isGG = ana.TreeVarBool("init_gg")
	isQQ = ana.TreeVarBool("init_qq")

	// Variables
	var_dphi = &ana.Variable{
		Name:       "truth_dphi_ll",
		SaveName:   "truth_dphi_ll",
		TreeFunc:   ana.TreeVarF64("truth_dphi_ll"),
		Nbins:      15,
		Xmin:       0,
		Xmax:       math.Pi,
		XLabel:     `$\Delta\phi_{\ell\ell}$`,
		YLabel:     `PDF($\Delta\phi_{\ell\ell}$)`,
		LegPosTop:  true,
		LegPosLeft: true,
		RangeYmax:  0.08,
	}

	var_Ckk = &ana.Variable{
		Name:       "truth_Ckk",
		SaveName:   "truth_Ckk",
		TreeFunc:   ana.TreeVarF64("truth_Ckk"),
		Nbins:      25,
		Xmin:       -1,
		Xmax:       1,
		XLabel:     `$\cos\theta^{+}_{k} \: \cos\theta^{-}_{k}$`,
		YLabel:     `PDF($\cos\theta^{+}_{k} \cos\theta^{-}_{k}$)`,
		RangeXmin:  -1.5,
		RangeXmax:  1,
		LegPosTop:  true,
		LegPosLeft: true,
	}

	var_Crr = &ana.Variable{
		Name:       "truth_Crr",
		SaveName:   "truth_Crr",
		TreeFunc:   ana.TreeVarF64("truth_Crr"),
		Nbins:      25,
		Xmin:       -1,
		Xmax:       1,
		XLabel:     `$\cos\theta^{+}_{r} \: \cos\theta^{-}_{r}$`,
		YLabel:     `PDF($\cos\theta^{+}_{r} \cos\theta^{-}_{r}$)`,
		RangeXmin:  -1.5,
		RangeXmax:  1,
		LegPosTop:  true,
		LegPosLeft: true,
	}

	var_Cnn = &ana.Variable{
		Name:       "truth_Cnn",
		SaveName:   "truth_Cnn",
		TreeFunc:   ana.TreeVarF64("truth_Cnn"),
		Nbins:      25,
		Xmin:       -1,
		Xmax:       1,
		XLabel:     `$\cos\theta^{+}_{n} \: \cos\theta^{-}_{n}$`,
		YLabel:     `PDF($\cos\theta^{+}_{n} \cos\theta^{-}_{n}$)`,
		RangeXmin:  -1.5,
		RangeXmax:  1,
		LegPosTop:  true,
		LegPosLeft: true,
	}

	var_pt_lep = &ana.Variable{
		Name:       "pt_lep",
		SaveName:   "pt_lep",
		TreeFunc:   ana.TreeVarF32("l_pt"),
		Nbins:      25,
		Xmin:       0,
		Xmax:       500,
		XLabel:     `$p^{\ell}_{T}$ [GeV]`,
		YLabel:     `PDF($p^{\ell}_{T}$)`,
		LegPosTop:  true,
		LegPosLeft: false,
	}

	var_eta_lep = &ana.Variable{
		Name:       "eta_lep",
		SaveName:   "eta_lep",
		TreeFunc:   ana.TreeVarF32("l_eta"),
		Nbins:      25,
		Xmin:       -5,
		Xmax:       5,
		XLabel:     `$\eta^{\ell}$`,
		YLabel:     `PDF($\eta^{\ell}$)`,
		LegPosTop:  true,
		LegPosLeft: false,
		RangeXmax:  8,
	}

	var_pt_b = &ana.Variable{
		Name:       "pt_b",
		SaveName:   "pt_b",
		TreeFunc:   ana.TreeVarF32("b_pt"),
		Nbins:      25,
		Xmin:       0,
		Xmax:       500,
		XLabel:     `$p^{b}_{T}$ [GeV]`,
		YLabel:     `PDF($p^{b}_{T}$)`,
		LegPosTop:  true,
		LegPosLeft: false,
	}

	var_eta_b = &ana.Variable{
		Name:       "eta_b",
		SaveName:   "eta_b",
		TreeFunc:   ana.TreeVarF32("b_eta"),
		Nbins:      25,
		Xmin:       -5,
		Xmax:       5,
		XLabel:     `$\eta^{b}$`,
		YLabel:     `PDF($\eta^{b}$)`,
		LegPosTop:  true,
		LegPosLeft: false,
		RangeXmax:  8,
	}

	var_pt_vsum = &ana.Variable{
		Name:        "vsum_pt",
		SaveName:    "pt_vsum",
		TreeFunc:    ana.TreeVarF32("vsum_pt"),
		Nbins:       25,
		Xmin:        0,
		Xmax:        250,
		XLabel:      `Truth $E^{\mathrm{miss}}_{T} \; \equiv \; |\vec{p}^{\,\nu}_T + \vec{p}^{\,\bar{\nu}}_T|$`,
		YLabel:      `PDF($E^{\mathrm{miss}}_{T}$)`,
		XTickFormat: "%2.0f",
		LegPosTop:   true,
		LegPosLeft:  false,
	}

	var_pt_t = &ana.Variable{
		Name:     "t_pt",
		SaveName: "pt_t",
		TreeFunc: ana.TreeVarF32("t_pt"),
		Nbins:    100,
		Xmin:     0,
		Xmax:     500,
		XLabel:   `$p^{t}_{T}$ [GeV]`,
		// YLabel: `Number of Entries`,
		YLabel:     `PDF($p^{t}_{T}$)`,
		LegPosTop:  true,
		LegPosLeft: false,
		//YTickFormat: "%2.0f",
	}

	var_eta_t = &ana.Variable{
		Name:       "eta_t",
		SaveName:   "eta_t",
		TreeFunc:   ana.TreeVarF32("t_eta"),
		Nbins:      25,
		Xmin:       -5,
		Xmax:       5,
		XLabel:     `$\eta^{t}$`,
		YLabel:     `PDF($\eta^{t}$)`,
		LegPosTop:  true,
		LegPosLeft: false,
		RangeXmax:  8,
		RangeYmax:  100,
	}

	var_m_tt = &ana.Variable{
		Name:       "m_tt",
		SaveName:   "m_tt",
		TreeFunc:   ana.TreeVarF32("ttbar_m"),
		Nbins:      25,
		Xmin:       300,
		Xmax:       1500,
		XLabel:     `$m_{t\bar{t}}$ [GeV]`,
		YLabel:     `PDF($m_{t\bar{t}}$)`,
		LegPosTop:  true,
		LegPosLeft: false,
	}

	var_pt_tt = &ana.Variable{
		Name:       "pt_tt",
		SaveName:   "pt_tt",
		TreeFunc:   ana.TreeVarF32("ttbar_pt"),
		Nbins:      25,
		Xmin:       0,
		Xmax:       150,
		XLabel:     `$p^{t\bar{t}}_T$ [GeV]`,
		YLabel:     `PDF($p^{t\bar{t}}_T$)`,
		LegPosTop:  true,
		LegPosLeft: false,
	}

	var_x1 = &ana.Variable{
		SaveName: "init_x1",
		TreeFunc: ana.TreeVarF32("init_x1"),
		Nbins:    25,
		Xmin:     0,
		Xmax:     1,
	}

	var_GG = &ana.Variable{
		SaveName: "init_gg",
		TreeFunc: ana.TreeVarBool("init_qq"),
		Nbins:    2,
		Xmin:     0,
		Xmax:     2,
	}

	var_x1x2 = &ana.Variable{
		SaveName: "x1x2",
		TreeFunc: ana.TreeFunc{
			VarsName: []string{"init_x1", "init_x2"},
			Fct: func(x1, x2 float32) float64 {
				return float64(x1 * x2)
			},
		},
		Nbins: 25,
		Xmin:  0,
		Xmax:  1,
	}

	sel1 = &ana.Selection{
		Name: "cut1",
		TreeFunc: ana.TreeFunc{
			VarsName: []string{"l_pt"},
			Fct:      func(pt float32) bool { return pt > -20 },
		},
	}

	sel2 = &ana.Selection{
		Name: "cut2",
		TreeFunc: ana.TreeFunc{
			VarsName: []string{"l_pt"},
			Fct:      func(pt float32) bool { return pt > -50 },
		},
	}

	sel3 = &ana.Selection{
		Name: "cut3",
		TreeFunc: ana.TreeFunc{
			VarsName: []string{"l_pt"},
			Fct:      func(pt float32) bool { return pt > -100 },
		},
	}

	sel4 = &ana.Selection{
		Name: "cut4",
		TreeFunc: ana.TreeFunc{
			VarsName: []string{"l_pt"},
			Fct:      func(pt float32) bool { return pt > -150 },
		},
	}
)

// Helper function to add many components (10k per components) to a sample
func loadManyComponents(s *ana.Sample, n10kEvts int) {
	for i := 0; i < n10kEvts; i++ {
		s.AddComponent(file2, tname)
	}
}
