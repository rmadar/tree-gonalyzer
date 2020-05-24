package ana_test

import (
	"github.com/rmadar/tree-gonalyzer/ana"
)

// Creation of the default analysis maker type with
// single-component samples. The files, trees and
// variables are dummy, they are here just for the example.
func ExampleMaker_defaultSingleComponent() {
	// Define samples
	samples := []*ana.Sample{
		ana.CreateSample("data", "data", `Data`, "data.root", "mytree"),
		ana.CreateSample("bkg1", "bkg", `Proc 1`, "proc1.root", "mytree"),
		ana.CreateSample("bkg2", "bkg", `Proc 2`, "proc2.root", "mytree"),
		ana.CreateSample("bkg3", "bkg", `Proc 3`, "proc3.root", "mytree"),
	}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("plot1", "branch1", new(float64), 15, 0, 10),
		ana.NewVariable("plot2", "branch2", new(float32), 25, 0, 10),
	}

	// Create analyzer object
	analyzer := ana.New(samples, variables)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

// Creation of the default analysis maker type with
// multi-component samples. The files, trees and
// variables are dummy, they are here just for the example.
func ExampleMaker_defaultMultiComponents() {
	// Weights and cuts
	w := ana.NewTreeFuncVarF64("evtWeight")
	isProc4 := ana.NewTreeFuncVarBool("IsProc4")

	// Data sample.
	data := ana.NewSample("data", "data", `Data 18-20`)
	data.AddComponent("data2018.root", "mytree")
	data.AddComponent("data2019.root", "mytree")
	data.AddComponent("data2020.root", "mytree")

	// Background sample including four components.
	bkg := ana.NewSample("bkgTot", "bkg", `Total Bkg`, ana.WithWeight(w))
	bkg.AddComponent("proc1.root", "mytree")
	bkg.AddComponent("proc2.root", "mytree")
	bkg.AddComponent("proc3.root", "mytree")
	bkg.AddComponent("proc4.root", "mytree", ana.WithCut(isProc4))

	// Signal sample including three components.
	sig := ana.NewSample("sigTot", "sig", `Total signal`, ana.WithWeight(w))
	sig.AddComponent("sig1.root", "mytree")
	sig.AddComponent("sig2.root", "mytree")
	sig.AddComponent("sig3.root", "mytree")
	
	// Put samples together.
	samples := []*ana.Sample{data, bkg, sig}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("plotName1", "branchName1", new(float64), 15, 0, 10),
		ana.NewVariable("plotName2", "branchName2", new(float32), 25, 0, 10),
	}

	// Create analyzer object
	analyzer := ana.New(samples, variables)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

// Creation of the Stack analysis maker type
func ExampleMaker_withStacked() {

}

// Creation of the normalized analysis maker type
func ExampleMaker_withNormalized() {

}

// Creation of the auto-styled analysis maker type
func ExampleMaker_withAutoStyle() {

}

// Creation of an analysis maker type including cuts
func ExampleMaker_kinematicCuts() {

}
