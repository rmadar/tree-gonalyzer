package ana_test

import (
	"math"
	"image/color"
	
	"github.com/rmadar/tree-gonalyzer/ana"
)

var (
	// Some files and trees names
	fData = "../testdata/file1.root"
	fBkg1 = "../testdata/file2.root"
	fBkg2 = "../testdata/file3.root"
	tName = "truth"

	// Some weights
	w1 = ana.NewTreeFuncValF64(1.0)
	w2 = ana.NewTreeFuncValF64(0.5)
	w3 = ana.TreeFunc{
		VarsName: []string{"t_pt"},
		Fct:      func(pt float32) float64 { return 1.0 + float64(pt)/50. },
	}

	// Some colors
	noColor = color.NRGBA{}
	shadowBlue = color.NRGBA{R: 50, G: 20, B: 150, A: 20}
	darkRed = color.NRGBA{R: 180, G: 30, B: 50, A: 200}
)

// Creation of the default analysis maker type with
// single-component samples.
func ExampleMaker_aSimpleUseCase() {
	// Define samples
	samples := []*ana.Sample{
		ana.CreateSample("data", "data", `Data`, fData, tName),
		ana.CreateSample("bkg1", "bkg", `Proc 1`, fBkg1, tName, ana.WithWeight(w1)),
		ana.CreateSample("bkg2", "bkg", `Proc 2`, fBkg2, tName, ana.WithWeight(w2)),
		ana.CreateSample("bkg3", "bkg", `Proc 3`, fBkg1, tName, ana.WithWeight(w2)),
	}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("Mttbar", "ttbar_m", new(float32), 25, 350, 1000,
			ana.WithAxisLabels("M(t,t) [GeV]", "Events Yields"),
		),
		ana.NewVariable("DphiLL", "truth_dphi_ll", new(float64), 10, 0, math.Pi,
			ana.WithAxisLabels("dPhi(l,l)", "Events Yields"),
			ana.WithLegLeft(true)),
	}

	// Create analyzer object
	analyzer := ana.New(samples, variables,
		ana.WithAutoStyle(true),
		ana.WithSaveFormat("png"),
		ana.WithSavePath("testdata/Plots_simpleUseCase"),
	)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

// Creation of the default analysis maker type with
// multi-component samples. The files, trees and
// variables are dummy, they are here just for the example.
func ExampleMaker_multiComponentSamples() {
	// Weights and cuts
	w := ana.NewTreeFuncVarF32("weight")
	isQQ := ana.NewTreeFuncVarBool("init_qq")

	// Data sample.
	data := ana.NewSample("data", "data", `Data 18-20`)
	data.AddComponent(fData, tName)
	data.AddComponent(fBkg1, tName)

	// Background A sample including three components.
	bkgA := ana.NewSample("BkgTotA", "bkg", `Total Bkg A`, ana.WithWeight(w))
	bkgA.AddComponent(fBkg1, tName)
	bkgA.AddComponent(fBkg2, tName)
	bkgA.AddComponent(fBkg1, tName, ana.WithCut(isQQ))

	// Background B sample including two components.
	bkgB := ana.NewSample("BkgTotB", "bkg", `Total Bkg B`, ana.WithWeight(w))
	bkgB.AddComponent(fBkg1, tName)
	bkgB.AddComponent(fBkg2, tName)

	// Put samples together.
	samples := []*ana.Sample{data, bkgA, bkgB}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("Mttbar", "ttbar_m", new(float32), 25, 350, 1000),
		ana.NewVariable("DphiLL", "truth_dphi_ll", new(float64), 10, 0, math.Pi),
	}

	// Create analyzer object with normalized histograms.
	analyzer := ana.New(samples, variables,
		ana.WithAutoStyle(true),
		ana.WithHistoNorm(true),
		ana.WithSaveFormat("png"),
		ana.WithSavePath("testdata/Plots_multiComponents"),
	)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

func ExampleMaker_shapeComparison() {
	// Define samples
	samples := []*ana.Sample{
		ana.CreateSample("proc1", "bkg", `Proc 1`, fBkg1, tName,
			ana.WithFillColor(shadowBlue),
		),
		ana.CreateSample("proc2", "bkg", `Proc 2`, fBkg2, tName,
			ana.WithWeight(w3),
			ana.WithLineColor(darkRed),
			ana.WithLineWidth(2),
			ana.WithBand(true),
		),
	}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("Mttbar", "ttbar_m", new(float32), 25, 350, 1000),
		ana.NewVariable("DphiLL", "truth_dphi_ll", new(float64), 10, 0, math.Pi, ana.WithLegLeft(true)),
	}

	// Create analyzer object
	analyzer := ana.New(samples, variables,
		ana.WithHistoStack(false),
		ana.WithHistoNorm(true),
		ana.WithSaveFormat("png"),
		ana.WithSavePath("testdata/Plots_shapeComparison"),
	)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}

}

func ExampleMaker_systematicVariations() {

}

func ExampleMaker_shapeDistortion() {

}

func ExampleMaker_withKinemCuts() {

}
