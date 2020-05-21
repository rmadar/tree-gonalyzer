package ana_test

func ExampleSample_default() {
	// Data sample
	sData := ana.NewSample("DATA", "data", `pp data`, "myfile.root", "mytree")

	// Background sample, say QCD pp->ttbar production
	sBkg := ana.NewSample("ttbar", "bkg", `QCD prod.`, "myfile.root", "mytree")

	// Signal sample, say pp->H->ttbar production
	sSig := ana.NewSample("Htt", "sig", `Higgs prod.`, "myfile.root", "mytree")
}


func ExampleSample_withWeight() {
	// Weight computed from several branches 
	w := ana.TreeFunc{
		VarsName: []string{"w1", "w2", "w3"},
		Fct: func (w1, w2, w3 float64) float64 {
			return w1*w2*w3
		},
	}

	// Computed weight
	s2 := ana.NewSample("proc", "bkg", `leg`, "myfile.root", "mytree",
		ana.WithWeight(w),
	)
	
	// Single branch weight
	s2 := ana.NewSample("proc", "bkg", `leg`, "myfile.root", "mytree",
		ana.WithWeight(ana.NewTreeFruncF64("evt_weight")),
	)
}

func ExampleSample_withCut() {
	
}

func ExampleSample_withSubSamples() {
	// To BE IMPLEMENTED FIRST
}
