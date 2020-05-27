// Package allowing to benchmark performances of the ana package.
package main

import (
	"fmt"
	"flag"
	"math"

	"github.com/rmadar/tree-gonalyzer/ana"
)

// Run the analyzer
func main() {

	// Options passed by command lines.
	var (
		compileLatex = flag.Bool("latex", false, "On-the-fly LaTeX compilation of produced figure.")
		varFormula   = flag.Bool("varFormula", false, "Use TreeFormula for variables.")
		nVariables   = flag.Int("nVars", -1, "Number of variables to loop over.")
		noCutWeight  = flag.Bool("noCutWeight", false, "Disable cuts and weights, avoiding TreeFormula.")
	)
	flag.Parse()

	// Data
	splData := ana.NewSample("data", "data", `Pseudo-data`)
	loadManyComponents(splData)

	// Background 1
	splBkg1 := ana.NewSample("bkg1", "bkg", `Proc 1`, ana.WithWeight(w1), ana.WithCut(isGG))
	loadManyComponents(splBkg1)

	// Background 2
	splBkg2 := ana.NewSample("bkg2", "bkg", `Proc 2`, ana.WithWeight(w2), ana.WithCut(isQQ))
	loadManyComponents(splBkg2)

	// Background 3
	splBkg3 := ana.NewSample("bkg3", "bkg", `Proc 3`)
	loadManyComponents(splBkg3)

	// Background 4
	splBkg4 := ana.NewSample("bkg3", "bkg", `Proc 3`)
	loadManyComponents(splBkg4)

	// Group samples together
	samples := []*ana.Sample{splData, splBkg1, splBkg2, splBkg3, splBkg4}

	// Variables, organized in bunch of 15
	variables := []*ana.Variable{
		var_dphi, var_m_tt, var_eta_t, var_pt_lep, var_Ckk,
		var_Crr, var_Cnn, var_pt_lep, var_eta_lep, var_pt_b,
		var_eta_b, var_pt_vsum, var_pt_t, var_pt_tt, var_x1,

		var_dphi, var_m_tt, var_eta_t, var_pt_lep, var_Ckk,
		var_Crr, var_Cnn, var_pt_lep, var_eta_lep, var_pt_b,
		var_eta_b, var_pt_vsum, var_pt_t, var_pt_tt, var_x1,

		var_dphi, var_m_tt, var_eta_t, var_pt_lep, var_Ckk,
		var_Crr, var_Cnn, var_pt_lep, var_eta_lep, var_pt_b,
		var_eta_b, var_pt_vsum, var_pt_t, var_pt_tt, var_x1,

		var_dphi, var_m_tt, var_eta_t, var_pt_lep, var_Ckk,
		var_Crr, var_Cnn, var_pt_lep, var_eta_lep, var_pt_b,
		var_eta_b, var_pt_vsum, var_pt_t, var_pt_tt, var_x1,
	}

	// Protection for too high number of variables
	nVars := len(variables)
	if *nVariables > -1 {
		nVars = *nVariables
	}
	if nVars > len(variables) {
		panic(fmt.Errorf("Too much variables (max 60, got %v)", nVars))
	}
	
	// Create analyzer object with options
	analyzer := ana.New(samples, variables[:nVars],
		ana.WithAutoStyle(true),
		ana.WithPlotTitle(`GOnalyzer Performance Test`),
		ana.WithCompileLatex(*compileLatex),
		ana.WithHistoNorm(true),
		ana.WithHistoStack(true),
	)

	// Few handles for benchmarking
	analyzer.WithVarsTreeFormula = *varFormula
	analyzer.NoTreeFormula = *noCutWeight

	// Run the analyzer and produce all plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

// Define all samples and variables of the analysis
var (
	// files/trenames
	file1 = "../testdata/file1.root"
	file2 = "../testdata/file2.root"
	file3 = "../testdata/file3.root"
	tname = "truth"

	// Some TreeFunc: weights and cuts
	w1   = ana.NewTreeFuncValF64(0.5)
	w2   = ana.NewTreeFuncValF64(2.0)
	isGG = ana.NewTreeFuncVarBool("init_gg")
	isQQ = ana.NewTreeFuncVarBool("init_qq")

	// Variables
	var_dphi = &ana.Variable{
		Name:       "truth_dphi_ll",
		SaveName:   "truth_dphi_ll",
		TreeName:   "truth_dphi_ll",
		Value:      new(float64),
		TreeVar:    ana.NewTreeFuncVarF64("truth_dphi_ll"),
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
		TreeName:   "truth_Ckk",
		Value:      new(float64),
		TreeVar:    ana.NewTreeFuncVarF64("truth_Ckk"),
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
		TreeName:   "truth_Crr",
		Value:      new(float64),
		TreeVar:    ana.NewTreeFuncVarF64("truth_Crr"),
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
		TreeName:   "truth_Cnn",
		Value:      new(float64),
		TreeVar:    ana.NewTreeFuncVarF64("truth_Cnn"),
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
		TreeName:   "l_pt",
		Value:      new(float32),
		TreeVar:    ana.NewTreeFuncVarF32("l_pt"),
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
		TreeName:   "l_eta",
		Value:      new(float32),
		TreeVar:    ana.NewTreeFuncVarF32("l_eta"),
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
		TreeName:   "b_pt",
		Value:      new(float32),
		TreeVar:    ana.NewTreeFuncVarF32("b_pt"),
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
		TreeName:   "b_eta",
		Value:      new(float32),
		TreeVar:    ana.NewTreeFuncVarF32("b_eta"),
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
		TreeName:    "vsum_pt",
		Value:       new(float32),
		TreeVar:     ana.NewTreeFuncVarF32("vsum_pt"),
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
		TreeName: "t_pt",
		Value:    new(float32),
		TreeVar:  ana.NewTreeFuncVarF32("t_pt"),
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
		TreeName:   "t_eta",
		Value:      new(float32),
		TreeVar:    ana.NewTreeFuncVarF32("t_eta"),
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
		TreeName:   "ttbar_m",
		Value:      new(float32),
		TreeVar:    ana.NewTreeFuncVarF32("ttbar_m"),
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
		TreeName:   "ttbar_pt",
		Value:      new(float32),
		TreeVar:    ana.NewTreeFuncVarF32("ttbar_pt"),
		Nbins:      25,
		Xmin:       0,
		Xmax:       150,
		XLabel:     `$p^{t\bar{t}}_T$ [GeV]`,
		YLabel:     `PDF($p^{t\bar{t}}_T$)`,
		LegPosTop:  true,
		LegPosLeft: false,
	}

	var_x1 = &ana.Variable{
		TreeName: "init_x1",
		Value:    new(float32),
		SaveName: "init_x1",
		TreeVar:  ana.NewTreeFuncVarF32("init_x1"),
		Nbins:    25,
		Xmin:     0,
		Xmax:     1,
	}

	var_x1x2 = &ana.Variable{
		SaveName: "x1x2",
		TreeVar: ana.TreeFunc{
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
			Fct:      func(pt float32) bool { return pt > 20 },
		},
	}

	sel2 = &ana.Selection{
		Name: "cut2",
		TreeFunc: ana.TreeFunc{
			VarsName: []string{"l_pt"},
			Fct:      func(pt float32) bool { return pt > 50 },
		},
	}

	sel3 = &ana.Selection{
		Name: "cut3",
		TreeFunc: ana.TreeFunc{
			VarsName: []string{"l_pt"},
			Fct:      func(pt float32) bool { return pt > 100 },
		},
	}

	sel4 = &ana.Selection{
		Name: "cut4",
		TreeFunc: ana.TreeFunc{
			VarsName: []string{"l_pt"},
			Fct:      func(pt float32) bool { return pt > 150 },
		},
	}
)

// Helper function to add many components to a sample
func loadManyComponents(s *ana.Sample) {
	s.AddComponent(file1, tname)
	s.AddComponent(file2, tname)
	s.AddComponent(file3, tname)
	s.AddComponent(file1, tname)
	s.AddComponent(file2, tname)
	s.AddComponent(file3, tname)
	s.AddComponent(file1, tname)
	s.AddComponent(file2, tname)
	s.AddComponent(file3, tname)
	s.AddComponent(file1, tname)
	s.AddComponent(file2, tname)
	s.AddComponent(file3, tname)
	s.AddComponent(file1, tname)
	s.AddComponent(file2, tname)
	s.AddComponent(file3, tname)
	s.AddComponent(file1, tname)
	s.AddComponent(file2, tname)
	s.AddComponent(file3, tname)
}
