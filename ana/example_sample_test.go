package ana_test

// Default samples for data, background and signal. 
func ExampleNewSample() {

	// Data sample
	sData := ana.NewSample("DATA", "data", `pp data 2020`, "myfile.root", "mytree")

	// Background sample, say QCD pp->ttbar production
	sBkg := ana.NewSample("ttbar", "bkg", `$pp \to t\bar{t}`, "myfile.root", "mytree")

	// Signal sample, say  pp->H->ttbar production
	sSig := ana.NewSample("Htt", "sig", `$pp \to H \to t\bar{t}`, "myfile.root", "mytree")
}
