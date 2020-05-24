// Showing how gonalyzer/ana package works
package main

import (
	"flag"
	"math"

	"github.com/rmadar/tree-gonalyzer/ana"
)

// Run the analyzer
func main() {

	// Options passed by command lines.
	var pFormat = flag.String("f", "tex", "Select output plot format")
	var doLatex = flag.Bool("l", false, "On-the-fly LaTeX compilation to produce figures")
	var doRatio = flag.Bool("r", false, "Enable ratio plot")
	var doStack = flag.Bool("s", false, "Enable histogram stacking")
	var doNorm = flag.Bool("n", false, "Enable histogram normalization")
	flag.Parse()

	// Samples
	samples := []*ana.Sample{
		ana.CreateSample("data", "data", `Data`, file1, tname),
		ana.CreateSample("bkg1", "bkg", `Proc 1`, file2, tname, ana.WithWeight(w1)),
		ana.CreateSample("bkg2", "bkg", `Proc 2`, file3, tname, ana.WithWeight(w1)),
		ana.CreateSample("bkg4", "bkg", `Proc 3`, file1, tname, ana.WithWeight(w2)),
	}

	// Variables
	variables := []*ana.Variable{
		ana.NewVariable("truth_dphi_ll", "truth_dphi_ll", new(float64), 15, 0, math.Pi),
		ana.NewVariable("m_tt", "ttbar_m", new(float32), 25, 300, 1000),
	}

	// Selections
	selections := []*ana.Selection{sel0, sel1, sel2}

	// Create analyzer object with some selections, enabeling automatic style
	analyzer := ana.New(samples, variables, ana.WithKinemCuts(selections), ana.WithAutoStyle(true))

	// Command line options
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
		
	// Weights 
	w1 = ana.NewTreeFuncValF64(0.5)
	w2 = ana.NewTreeFuncValF64(0.5)

	// Selections
	sel0 = ana.NewSelection()
	sel1 = &ana.Selection{
		Name: "m_gt_800",
		TreeFunc: ana.TreeFunc{
			VarsName: []string{"ttbar_m"},
			Fct:      func(m float32) bool { return m > 800 },
		},
	}
	sel2 = &ana.Selection{
		Name: "dphi_lg_1",
		TreeFunc: ana.TreeFunc{
			VarsName: []string{"truth_dphi_ll"},
			Fct:      func(dphi float64) bool { return dphi < 1.0 },
		},
	}
)
