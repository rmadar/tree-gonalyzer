package ana_test

import (
	"github.com/rmadar/tree-gonalyzer/ana"
)

func ExampleSample_default() {
	// Data sample
	sData := ana.NewSample("DATA", "data", `pp data`, "myfile.root", "mytree")

	// Background sample, say QCD pp->ttbar production
	sBkg := ana.NewSample("ttbar", "bkg", `QCD prod.`, "myfile.root", "mytree")

	// Signal sample, say pp->H->ttbar production
	sSig := ana.NewSample("Htt", "sig", `Higgs prod.`, "myfile.root", "mytree")

	// New analysis
	ana.New([]*ana.Sample{sData, sBkg, sSig}, []*ana.Variable{})
}

func ExampleSample_withWeight() {
	// Weight computed from several branches
	w := ana.TreeFunc{
		VarsName: []string{"w1", "w2", "w3"},
		Fct: func(w1, w2, w3 float64) float64 {
			return w1 * w2 * w3
		},
	}

	// Sample with computed weight
	ana.NewSample("proc", "bkg", `leg`, "myfile.root", "mytree",
		ana.WithWeight(w),
	)

	// Sample with single branch weight
	ana.NewSample("proc", "bkg", `leg`, "myfile.root", "mytree",
		ana.WithWeight(ana.NewTreeFuncVarF64("evtWght")),
	)
}

func ExampleSample_withCut() {
	// Selection criteria computed from several branches
	sel := ana.TreeFunc{
		VarsName: []string{"pt", "eta", "m"},
		Fct: func(pt, eta, m float64) bool {
			return (pt > 150 && eta > 0) || m < 125
		},
	}

	// Sample with computed boolean
	ana.NewSample("proc", "bkg", `leg`, "myfile.root", "mytree",
		ana.WithCut(sel),
	)

	// Sample with single branch boolean
	ana.NewSample("proc", "bkg", `leg`, "myfile.root", "mytree",
		ana.WithWeight(ana.NewTreeFuncVarBool("passCriteria")),
	)
}

func ExampleSample_severalComponents() {
	// ttbar background starting with one component, say ttbar->dilepton
	ttbarIncl := ana.NewSample("ttbar", "bkg", `Inclusive`, "dilep.root", "mytree")

	// Adding l+jets decay, weighted by BR(ttbar->l+jets)
	ttbarIncl.AddComponent("ljets.root", "mytree", ana.WithWeight(ana.NewTreeFuncValF64(0.3)))
	
	// Adding full hadronic decay, applying a cut to make sure of the decay
	ttbarIncl.AddComponent("fullhad.root", "mytree", ana.WithCut(ana.NewTreeFuncVarBool("isHadronic")))
}
