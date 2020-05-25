package ana_test

import (
	"math"

	"github.com/rmadar/tree-gonalyzer/ana"
)

var (
	fData = "../testdata/file1.root"
	fBkg1 = "../testdata/file2.root"
	fBkg2 = "../testdata/file3.root"
	tName = "truth"
)

// Creation of the default analysis maker type with
// single-component samples.
func ExampleMaker_aSimpleUseCase() {
	// Define samples
	samples := []*ana.Sample{
		ana.CreateSample("data", "data", `Data`, fData, tName),
		ana.CreateSample("bkg1", "bkg", `Proc 1`, fBkg1, tName),
		ana.CreateSample("bkg2", "bkg", `Proc 2`, fBkg2, tName),
		ana.CreateSample("bkg3", "bkg", `Proc 3`, fBkg1, tName),
	}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("Mttbar", "ttbar_m", new(float32), 50, 0, 1000),
		ana.NewVariable("DphiLL", "truth_dphi_ll", new(float64), 25, 0, math.Pi),
	}

	// Create analyzer object
	analyzer := ana.New(samples, variables,
		ana.WithAutoStyle(true),
		ana.WithSaveFormat("png"),
		ana.WithSavePath("Plots_simpleUseCase"),
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
	w := ana.NewTreeFuncVarF64("weight")
	isQQ := ana.NewTreeFuncVarBool("init_qq")

	// Data sample.
	data := ana.NewSample("data", "data", `Data 18-20`)
	data.AddComponent(fData, tName)
	data.AddComponent(fBkg1, tName)

	// Background sample including four components.
	bkg := ana.NewSample("bkgTot", "bkg", `Total Bkg`, ana.WithWeight(w))
	bkg.AddComponent(fBkg1, tName)
	bkg.AddComponent(fBkg2, tName)
	bkg.AddComponent(fBkg1, tName, ana.WithCut(isQQ))

	// Signal sample including three components.
	sig := ana.NewSample("sigTot", "sig", `Total signal`, ana.WithWeight(w))
	sig.AddComponent(fBkg1, tName)
	sig.AddComponent(fBkg2, tName)

	// Put samples together.
	samples := []*ana.Sample{data, bkg, sig}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("plotName1", "l_pt", new(float32), 25, 0, 250),
		ana.NewVariable("plotName2", "v_pt", new(float32), 25, 0, 250),
	}

	// Create analyzer object
	analyzer := ana.New(samples, variables,
		ana.WithAutoStyle(true),
		ana.WithSaveFormat("png"),
		ana.WithSavePath("Plots_multiComponents"),		
	)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

func ExampleMaker_shapeComparison() {

}

func ExampleMaker_systematicVariations() {

}

func ExampleMaker_shapeDistortion() {

}

func ExampleMaker_withKinemCuts() {

}
