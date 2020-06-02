// Package featuring a realistic example of analysis using ana package.
package main

import (
	"flag"

	"github.com/rmadar/tree-gonalyzer/ana"
)

// Run the analyzer
func main() {

	// Options passed by command lines.
	var (
		nMax    = flag.Int64("nevts", -1, "Maximum number of processed event per sample component.")
		pFormat = flag.String("f", "tex", "Select output plot format")
		doLatex = flag.Bool("l", false, "On-the-fly LaTeX compilation to produce figures")
		doRatio = flag.Bool("r", false, "Enable ratio plot")
		doStack = flag.Bool("s", false, "Enable histogram stacking")
		doNorm  = flag.Bool("n", false, "Enable histogram normalization")
	)
	flag.Parse()

	// Data
	sData := ana.CreateSample("data", "data", `Data`, file1, tname)

	// Background 1
	sBkg1 := ana.NewSample("bkg1", "bkg", `Proc1+Proc2`, ana.WithWeight(wptPos))
	sBkg1.AddComponent(file2, tname)
	sBkg1.AddComponent(file3, tname)

	// Background 2
	sBkg2 := ana.NewSample("bkg2", "bkg", `Proc3+Proc4`, ana.WithWeight(wptNeg))
	sBkg2.AddComponent(file2, tname, ana.WithWeight(w1))
	sBkg2.AddComponent(file3, tname, ana.WithWeight(w1))

	// Background 3
	sBkg3 := ana.NewSample("bkg3", "bkg", `Proc5+Proc6`, ana.WithWeight(w2))
	sBkg3.AddComponent(file2, tname, ana.WithCut(ana.TreeCutBool("init_qq")))
	sBkg3.AddComponent(file3, tname)

	// Putting samples together
	samples := []*ana.Sample{sData, sBkg1, sBkg2, sBkg3}

	// Variables
	variables := []*ana.Variable{
		ana.NewVariable("dPhi", ana.TreeVarF64("truth_dphi_ll"), 15, 0, 3.14),
		ana.NewVariable("m_tt", ana.TreeVarF32("ttbar_m"), 25, 300, 1000),
		ana.NewVariable("isGG", ana.TreeVarBool("init_gg"), 2, 0, 1),
		ana.NewVariable("combVar", combFunc, 100, 0, 100),
	}

	// Selections
	selections := []*ana.Selection{sel0, sel1, sel2}

	// Create analyzer object with some selections, enabeling automatic style
	analyzer := ana.New(samples, variables,
		ana.WithKinemCuts(selections),
		ana.WithDumpTree(true),
		ana.WithSampleMT(true),
	)

	// Command line options
	analyzer.Nevts = *nMax
	analyzer.SaveFormat = *pFormat
	analyzer.CompileLatex = *doLatex
	analyzer.RatioPlot = *doRatio
	analyzer.HistoNorm = *doNorm
	analyzer.HistoStack = *doStack

	// Run the analyzer and produce all plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

// Few definition
var (
	// Files and tree names
	file1 = "../testdata/file1.root"
	file2 = "../testdata/file2.root"
	file3 = "../testdata/file3.root"
	tname = "truth"

	// TreeFunc: variables, weights and cuts
	combFunc = ana.TreeFunc{
		VarsName: []string{"t_pt", "truth_dphi_ll"},
		//Fct: func(pt float32, dphi float64) float64 {
		//	return float64(pt) / 10.
		Fct: func(pt float32, dphi float64) []float64 {
			return []float64{dphi, float64(pt) / 10.}
		},
	}
	w1     = ana.TreeValF64(0.5)
	w2     = ana.TreeValF64(0.25)
	wptPos = ana.TreeFunc{
		VarsName: []string{"t_pt"},
		Fct:      func(pt float32) float64 { return 1.0 + float64(pt)/200. },
	}
	wptNeg = ana.TreeFunc{
		VarsName: []string{"t_pt"},
		Fct:      func(pt float32) float64 { return 1.0 - float64(pt)/200. },
	}
	mtGT800 = ana.TreeFunc{
		VarsName: []string{"ttbar_m"},
		Fct:      func(m float32) bool { return m > 800 },
	}
	dphiLT1 = ana.TreeFunc{
		VarsName: []string{"truth_dphi_ll"},
		Fct:      func(dphi float64) bool { return dphi < 1.0 },
	}

	// Selections
	sel0 = ana.NewSelection()
	sel1 = &ana.Selection{
		Name:     "m_gt_800",
		TreeFunc: mtGT800,
	}
	sel2 = &ana.Selection{
		Name:     "dphi_lg_1",
		TreeFunc: dphiLT1,
	}
)
