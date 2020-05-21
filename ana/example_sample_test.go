package ana_test

func ExampleNewSample() {
	// Data sample
	sData := ana.NewSample("DATA", "data", `pp data`, "myfile.root", "mytree")

	// Background sample, say QCD pp->ttbar production
	sBkg := ana.NewSample("ttbar", "bkg", `QCD prod.`, "myfile.root", "mytree")

	// Signal sample, say pp->H->ttbar production
	sSig := ana.NewSample("Htt", "sig", `Higgs prod.`, "myfile.root", "mytree")
}


func ExampleSample_withWeight() {

	// Define a 'computed' weight 
	wComputed := ana.TreeFunc{
		VarsName: []string{"w1", "w2", "w3"},
		Fct: func (w1, w2, w3 float64) float64 {return w1*w2*w3},
	}
	
	sComputed := ana.NewSample(
		"Htt", "sig", `Higgs prod.`,
		"myfile.root", "mytree",
		ana.WithWeight(w),
	)
}

func ExampleSample_withCut() {
	
}

func ExampleSample_withSubSamples() {
	// To BE IMPLEMENTED FIRST
}
