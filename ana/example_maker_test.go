package ana_test

// Creation of the default analysis maker type
func ExampleMaker_default() {
	// Define Samples
	samples := []*ana.Sample{
		ana.NewSample("data", "data", `Data 2020`, "../testdata/ttbar_MadSpinOff.root", "truth"),
		ana.NewSample("bkg1", "bkg", `Proc 1`, "../testdata/ttbar_MadSpinOn_1.root", "truth"),
		ana.NewSample("bkg2", "bkg", `Proc 2`, "../testdata/ttbar_MadSpinOn_2.root", "truth"),
		ana.NewSample("bkg3", "bkg", `Proc 3`, "../testdata/ttbar_MadSpinOn_1.root", "truth"),
	}
	
	// Variables
	variables := []*ana.Variable{
		ana.NewVariable("truth_dphi_ll", "truth_dphi_ll", new(float64), 15, 0, math.Pi),
		ana.NewVariable("m_tt", "ttbar_m", new(float32), 25, 300, 1000),
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
