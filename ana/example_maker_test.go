package ana_test

import (
	"github.com/rmadar/tree-gonalyzer/ana"
)

// Creation of the default analysis maker type
func ExampleMaker_default() {
	// Define Samples
	samples := []*ana.Sample{
		ana.NewSample("data", "data", `Data`, "data.root", "mytree"),
		ana.NewSample("bkg1", "bkg", `Proc 1`, "proc1.root", "mytree"),
		ana.NewSample("bkg2", "bkg", `Proc 2`, "proc2.root", "mytree"),
		ana.NewSample("bkg3", "bkg", `Proc 3`, "proc3.root", "mytree"),
	}
	
	// Variables
	variables := []*ana.Variable{
		ana.NewVariable("plot1", "branch1", new(float64), 15, 0, 10),
		ana.NewVariable("plot2", "branch2", new(float32), 25, 0, 10),
	}

	// Create analyzer object
	analyzer := ana.New(samples, variables)
	
	// Run the analyzer and produce all plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

// Creation of the Stack analysis maker type
func ExampleMaker_stack() {

}

// Creation of the normalized analysis maker type
func ExampleMaker_norm() {

}

// Creation of the auto-styled analysis maker type
func ExampleMaker_autoStyle() {

}

// Creation of an analysis maker type including cuts
func ExampleMaker_kinematicCuts() {

}
