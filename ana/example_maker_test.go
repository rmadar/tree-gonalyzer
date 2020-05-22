package ana_test

import (
	"github.com/rmadar/tree-gonalyzer/ana"
)

// Creation of the default analysis maker type
func ExampleMaker_simpleCase() {
	// Define samples
	samples := []*ana.Sample{
		ana.NewSample("data", "data", `Data`, "data.root", "mytree"),
		ana.NewSample("bkg1", "bkg", `Proc 1`, "proc1.root", "mytree"),
		ana.NewSample("bkg2", "bkg", `Proc 2`, "proc2.root", "mytree"),
		ana.NewSample("bkg3", "bkg", `Proc 3`, "proc3.root", "mytree"),
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

// Creation of the default analysis maker type
func ExampleMaker_complexCase() {
	// Define useful TreeFunc for weights and cuts
	w := ana.NewTreeFuncVarF64("evtWeight")
	isProc4 := ana.NewTreeFuncVarBool("IsProc4")

	// Define data sample
	data := ana.NewSample("data", "data", `Data 18-20`, "data2018.root", "mytree")
	data.AddComponent("data2019.root", "mytree")
	data.AddComponent("data2020.root", "mytree")

	// Define a single sample for the total background
	bkg := ana.NewSample("bkgTot", "bkg", `Total Bkg`, "proc1.root", "mytree",
		ana.WithWeight(w))
	bkg.AddComponent("proc2.root", "mytree")
	bkg.AddComponent("proc3.root", "mytree")
	bkg.AddComponent("proc4.root", "mytree",
		ana.WithCut(isProc4))

	// Define a single sample for the total signal
	sig := ana.NewSample("sigTot", "sig", `Total signal`, "sig1.root", "mytree",
		ana.WithWeight(w))
	sig.AddComponent("sig2.root", "mytree")
	sig.AddComponent("sig3.root", "mytree")

	// Put samples together
	samples := []*ana.Sample{data, bkg, sig}
	
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
