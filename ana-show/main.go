// Example to run gonalyzer package
package main

import (
	"flag"
	"math"

	"image/color"

	"github.com/rmadar/tree-gonalyzer/ana"
)

// Run the analyzer
func main() {

	// Options passed by command lines.
	var doLatex = flag.Bool("latex", false, "On-the-fly LaTeX compilation of produced figure")
	var doRatio = flag.Bool("r", false, "Enable ratio plot")
	flag.Parse()
	
	// Samples
	samples :=  []*ana.Sample{&splData, &splBkg1, &splBkg2, &splBkg3}

	// Variables
	variables := []*ana.Variable{
		ana.NewVariable("truth_dphi_ll", "truth_dphi_ll", new(float64), 15, 0, math.Pi),
		ana.NewVariable("m_tt", "ttbar_m", new(float32), 25, 300, 1000),
	}

	// Selections
	selections := []*ana.Selection{ana.NewSelection(), sel1, sel2}
	
	// Create analyzer object with options
	analyzer := ana.New(samples, variables,
		ana.WithKinemCuts(selections),
		ana.WithPlotTitle(`{\tt TTree} {\bf GO}nalyzer -- Demo`),
		ana.WithSavePath("plots"),
		ana.WithCompileLatex(*doLatex),
		ana.WithRatioPlot(*doRatio),
		ana.WithHistoNorm(true),
		ana.WithHistoStack(true),
	)
	
	// Run the analyzer and produce all plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}


var (
	splData = ana.NewSample("data", "data", `Pseudo-data`, "../testdata/ttbar_MadSpinOff.root", "truth",
		ana.WithDataStyle(true),
	)	
	
	splBkg1 = ana.NewSample("bkg1", "bkg", `Process 1 (gg)`, "../testdata/ttbar_MadSpinOn_1.root", "truth",
		ana.WithWeight(ana.NewTreeFuncValF64(0.5)),
		ana.WithCut(ana.NewTreeFuncVarBool("init_gg")),
		ana.WithFillColor(color.NRGBA{R: 0, G: 102, B: 255, A: 230}),
		ana.WithLineWidth(0),
	)
	
	splBkg2 = ana.NewSample("bkg2", "bkg", `Process 2 (qq)`, "../testdata/ttbar_MadSpinOn_1.root", "truth",
		ana.WithWeight(ana.NewTreeFuncValF64(2)),
		ana.WithCut(ana.NewTreeFuncVarBool("init_qq")),
		ana.WithFillColor(color.NRGBA{R: 20, G: 20, B: 170, A: 230}),
		ana.WithLineWidth(0),
	)
	
	splBkg3 = ana.NewSample("bkg3", "bkg", `Process 3 (qg)`, "../testdata/ttbar_MadSpinOn_1.root", "truth",
		ana.WithWeight(ana.NewTreeFuncValF64(0.5)),
		ana.WithFillColor(color.NRGBA{R: 255, G: 102, B: 0, A: 200}),
		ana.WithLineWidth(0),
	)
	
	
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
