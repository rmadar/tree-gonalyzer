package ana_test

// Default samples for data, background and signal. 
func ExampleNewSample() {
	// Data sample
	sData := ana.NewSample("DATA", "data", `pp data`, "myfile.root", "mytree")

	// Background sample, say QCD pp->ttbar production
	sBkg := ana.NewSample("ttbar", "bkg", `QCD prod.`, "myfile.root", "mytree")

	// Signal sample, say  pp->H->ttbar production
	sSig := ana.NewSample("Htt", "sig", `Higgs prod.`, "myfile.root", "mytree")
}
