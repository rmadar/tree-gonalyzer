package ana_test

import (
	"image/color"
	"math"
	"testing"

	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/vg"

	"github.com/rmadar/tree-gonalyzer/ana"
)

func TestSimpleUseCase(t *testing.T) {
	cmpimg.CheckPlot(Example_aSimpleUseCase, t,
		"Plots_simpleUseCase/Mttbar.png",
		"Plots_simpleUseCase/DphiLL.png",
	)
}

func TestShapeComparison(t *testing.T) {
	cmpimg.CheckPlot(Example_shapeComparison, t,
		"Plots_shapeComparison/TopPt.png",
		"Plots_shapeComparison/DphiLL.png",
	)
}

func TestSystVariations(t *testing.T) {
	cmpimg.CheckPlot(Example_systematicVariations, t,
		"Plots_systVariations/Mttbar.png",
		"Plots_systVariations/DphiLL.png",
	)
}

func TestShapeDistortion(t *testing.T) {
	cmpimg.CheckPlot(Example_shapeDistortion, t,
		"Plots_shapeDistortion/Mttbar.png",
		"Plots_shapeDistortion/DphiLL.png",
	)
}

func TestSliceVariables(t *testing.T) {
	cmpimg.CheckPlot(Example_withSliceVariables, t,
		"Plots_withSliceVariables/hitTimes.png",
	)
}

func TestWithSignals(t *testing.T) {
	cmpimg.CheckPlot(Example_withSignals, t,
		"Plots_withSignals/Mttbar.png",
		"Plots_withSignals/DphiLL.png",
	)
}

func TestWithStackedSignals(t *testing.T) {
	cmpimg.CheckPlot(Example_withStackedSignals, t,
		"Plots_withStackedSignals/Mttbar.png",
		"Plots_withStackedSignals/DphiLL.png",
	)
}

func TestWithTreeDumping(t *testing.T) {
	cmpimg.CheckPlot(Example_withTreeDumping, t,
		"Plots_withTreeDumping/LowM/Mttbar.png",
		"Plots_withTreeDumping/LowM/DphiLL.png",
	)
}

// Creation of the default analysis maker type with
// single-component samples.
func Example_aSimpleUseCase() {
	// Define samples
	samples := []*ana.Sample{
		ana.CreateSample("data", "data", `Data`, fData, tName),
		ana.CreateSample("bkg1", "bkg", `Proc 1`, fBkg1, tName, ana.WithWeight(w1)),
		ana.CreateSample("bkg2", "bkg", `Proc 2`, fBkg2, tName, ana.WithWeight(w2)),
		ana.CreateSample("bkg3", "bkg", `Proc 3`, fBkg1, tName, ana.WithWeight(w2)),
	}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("Mttbar", ana.TreeVarF32("ttbar_m"), 25, 350, 1000,
			ana.WithAxisLabels("M(t,t) [GeV]", "Events Yields"),
		),
		ana.NewVariable("DphiLL", ana.TreeVarF64("truth_dphi_ll"), 10, 0, math.Pi,
			ana.WithAxisLabels("dPhi(l,l)", "Events Yields"),
			ana.WithLegLeft(true)),
	}

	// Create analyzer object
	analyzer := ana.New(samples, variables,
		ana.WithSavePath("testdata/Plots_simpleUseCase"),
	)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

func Example_withSignals() {
	// Define samples
	samples := []*ana.Sample{
		ana.CreateSample("data", "data", `Data`, fData, tName),
		ana.CreateSample("bkg1", "bkg", `Proc 1`, fBkg1, tName, ana.WithWeight(w1)),
		ana.CreateSample("bkg2", "bkg", `Proc 2`, fBkg2, tName, ana.WithWeight(w2)),
		ana.CreateSample("bkg3", "bkg", `Proc 3`, fBkg1, tName, ana.WithWeight(w2)),
		ana.CreateSample("sig1", "sig", `Sig 1`, fBkg2, tName, ana.WithWeight(wSigM(500, 0.04)),
			ana.WithLineColor(darkRed),
			ana.WithLineDashes([]vg.Length{2, 3, 2}),
			ana.WithLineWidth(3),
		),
		ana.CreateSample("sig2", "sig", `Sig 2`, fBkg2, tName, ana.WithWeight(wSigM(650, 0.02)),
			ana.WithLineColor(darkGreen),
			ana.WithLineDashes([]vg.Length{1, 1, 1}),
			ana.WithLineWidth(3),
		),
		ana.CreateSample("sig3", "sig", `Sig 2`, fBkg2, tName, ana.WithWeight(wSigM(800, 0.01)),
			ana.WithLineColor(darkBlue),
			ana.WithLineDashes([]vg.Length{3, 3, 3}),
			ana.WithLineWidth(3),
		),
	}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("Mttbar", ana.TreeVarF32("ttbar_m"), 100, 350, 1000,
			ana.WithAxisLabels("M(t,t) [GeV]", "Events Yields"),
		),
		ana.NewVariable("DphiLL", ana.TreeVarF64("truth_dphi_ll"), 10, 0, math.Pi,
			ana.WithAxisLabels("dPhi(l,l)", "Events Yields"),
			ana.WithLegLeft(true)),
	}

	// Create analyzer object
	analyzer := ana.New(samples, variables,
		ana.WithSavePath("testdata/Plots_withSignals"),
	)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

func Example_withStackedSignals() {
	// Define samples
	samples := []*ana.Sample{
		ana.CreateSample("data", "data", `Data`, fData, tName,
			ana.WithCircleSize(2.8),
			ana.WithYErrBarsCapWidth(0.01),
		),
		ana.CreateSample("bkg1", "bkg", `Proc 1`, fBkg1, tName, ana.WithWeight(w1)),
		ana.CreateSample("bkg2", "bkg", `Proc 2`, fBkg2, tName, ana.WithWeight(w2)),
		ana.CreateSample("bkg3", "bkg", `Proc 3`, fBkg1, tName, ana.WithWeight(w2)),
		ana.CreateSample("sig1", "sig", `Sig 1`, fBkg2, tName, ana.WithWeight(wSigM(500, 0.04))),
		ana.CreateSample("sig2", "sig", `Sig 2`, fBkg2, tName, ana.WithWeight(wSigM(650, 0.02))),
		ana.CreateSample("sig3", "sig", `Sig 2`, fBkg2, tName, ana.WithWeight(wSigM(800, 0.01))),
	}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("Mttbar", ana.TreeVarF32("ttbar_m"), 100, 350, 1000,
			ana.WithAxisLabels("M(t,t) [GeV]", "Events Yields"),
		),
		ana.NewVariable("DphiLL", ana.TreeVarF64("truth_dphi_ll"), 10, 0, math.Pi,
			ana.WithAxisLabels("dPhi(l,l)", "Events Yields"),
			ana.WithLegLeft(true)),
	}

	// Create analyzer object
	analyzer := ana.New(samples, variables,
		ana.WithSignalStack(true),
		ana.WithSavePath("testdata/Plots_withStackedSignals"),
	)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

// Creation of the default analysis maker type with
// multi-component samples. The files, trees and
// variables are dummy, they are here just for the example.
func Example_multiComponentSamples() {
	// Weights and cuts
	w := ana.TreeVarF32("weight")
	isQQ := ana.TreeVarBool("init_qq")

	// Data sample.
	data := ana.NewSample("data", "data", `Data 18-20`)
	data.AddComponent(fData, tName)
	data.AddComponent(fBkg1, tName)

	// Background A sample including three components.
	bkgA := ana.NewSample("BkgTotA", "bkg", `Total Bkg A`, ana.WithWeight(w))
	bkgA.AddComponent(fBkg1, tName)
	bkgA.AddComponent(fBkg2, tName)
	bkgA.AddComponent(fBkg1, tName, ana.WithCut(isQQ))

	// Background B sample including two components.
	bkgB := ana.NewSample("BkgTotB", "bkg", `Total Bkg B`, ana.WithWeight(w))
	bkgB.AddComponent(fBkg1, tName)
	bkgB.AddComponent(fBkg2, tName)

	// Put samples together.
	samples := []*ana.Sample{data, bkgA, bkgB}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("Mttbar", ana.TreeVarF32("ttbar_m"), 25, 350, 1000),
		ana.NewVariable("DphiLL", ana.TreeVarF64("truth_dphi_ll"), 10, 0, math.Pi),
	}

	// Create analyzer object with normalized histograms.
	analyzer := ana.New(samples, variables,
		ana.WithHistoNorm(true),
		ana.WithSavePath("testdata/Plots_multiComponents"),
	)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

func Example_shapeComparison() {
	// Define samples
	samples := []*ana.Sample{
		ana.CreateSample("data", "data", `Data`, fBkg1, tName),
		ana.CreateSample("proc1", "bkg", `Simulation A`, fBkg1, tName,
			ana.WithWeight(w3),
			ana.WithLineColor(darkBlue),
			ana.WithLineWidth(2),
			ana.WithBand(true),
		),
		ana.CreateSample("proc2", "bkg", `Simulation B`, fBkg2, tName,
			ana.WithWeight(w4),
			ana.WithLineColor(darkRed),
			ana.WithLineWidth(2),
			ana.WithBand(true),
		),
	}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("TopPt", ana.TreeVarF32("t_pt"), 10, 0, 500),
		ana.NewVariable("DphiLL", ana.TreeVarF64("truth_dphi_ll"), 10, 0, math.Pi,
			ana.WithLegLeft(true)),
	}

	// Create analyzer object
	analyzer := ana.New(samples, variables,
		ana.WithNevts(500),
		ana.WithHistoStack(false),
		ana.WithHistoNorm(true),
		ana.WithRatioPlot(false),
		ana.WithSavePath("testdata/Plots_shapeComparison"),
	)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}

}

func Example_systematicVariations() {
	// Samples
	samples := []*ana.Sample{
		ana.CreateSample("nom", "bkg", `Nominal`, fBkg1, tName,
			ana.WithLineColor(softBlack),
			ana.WithLineWidth(2.0),
			ana.WithBand(true),
		),
		ana.CreateSample("up", "bkg", `Up`, fBkg1, tName,
			ana.WithWeight(w3),
			ana.WithLineColor(darkRed),
			ana.WithLineWidth(1.5),
			ana.WithLineDashes([]vg.Length{3, 2}),
		),
		ana.CreateSample("down", "bkg", `Down`, fBkg1, tName,
			ana.WithWeight(w4),
			ana.WithLineColor(darkBlue),
			ana.WithLineWidth(1.5),
			ana.WithLineDashes([]vg.Length{3, 2}),
		),
	}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("Mttbar", ana.TreeVarF32("ttbar_m"), 25, 350, 1500,
			ana.WithRatioYRange(0.7, 1.3)),
		ana.NewVariable("DphiLL", ana.TreeVarF64("truth_dphi_ll"), 10, 0, math.Pi,
			ana.WithRatioYRange(0.7, 1.3),
			ana.WithYRange(0, 0.2),
			ana.WithLegLeft(true),
		),
	}

	// Create analyzer object
	analyzer := ana.New(samples, variables,
		ana.WithRatioPlot(true),
		ana.WithHistoStack(false),
		ana.WithHistoNorm(true),
		ana.WithSavePath("testdata/Plots_systVariations"),
	)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

func Example_shapeDistortion() {
	// Selections as TreeFunc's
	ptTopGT50 := ana.TreeFunc{
		VarsName: []string{"t_pt"},
		Fct:      func(pt float32) bool { return pt > 50. },
	}
	ptTopGT100 := ana.TreeFunc{
		VarsName: []string{"t_pt"},
		Fct:      func(pt float32) bool { return pt > 100. },
	}
	ptTopGT200 := ana.TreeFunc{
		VarsName: []string{"t_pt"},
		Fct:      func(pt float32) bool { return pt > 200. },
	}

	// Samples
	samples := []*ana.Sample{
		ana.CreateSample("noCut", "bkg", `No cut`, fBkg1, tName,
			ana.WithFillColor(shadowBlue),
		),
		ana.CreateSample("cut1", "bkg", `pT>50`, fBkg1, tName,
			ana.WithCut(ptTopGT50),
			ana.WithLineColor(darkRed),
			ana.WithLineWidth(2),
		),
		ana.CreateSample("cut2", "bkg", `pT>100`, fBkg1, tName,
			ana.WithCut(ptTopGT100),
			ana.WithLineColor(darkBlue),
			ana.WithLineWidth(2),
		),
		ana.CreateSample("cut3", "bkg", `pT>200`, fBkg1, tName,
			ana.WithCut(ptTopGT200),
			ana.WithLineColor(darkGreen),
			ana.WithLineWidth(2),
		),
	}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("Mttbar", ana.TreeVarF32("ttbar_m"), 25, 350, 1500,
			ana.WithAxisLabels("M(t,t) [GeV]", "PDF"),
		),
		ana.NewVariable("DphiLL", ana.TreeVarF64("truth_dphi_ll"), 10, 0, math.Pi,
			ana.WithLegLeft(true),
			ana.WithAxisLabels("dPhi(l,l)", "PDF"),
			ana.WithYRange(0, 0.3),
		),
	}

	// Create analyzer object
	analyzer := ana.New(samples, variables,
		ana.WithAutoStyle(false),
		ana.WithHistoStack(false),
		ana.WithRatioPlot(false),
		ana.WithHistoNorm(true),
		ana.WithSavePath("testdata/Plots_shapeDistortion"),
	)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

func Example_withKinemCuts() {

}

func Example_withTreeDumping() {
	// Weights and cuts
	w := ana.TreeVarF32("weight")
	isQQ := ana.TreeCutBool("init_qq")

	// Data sample.
	data := ana.NewSample("data", "data", `Data 18-20`)
	data.AddComponent(fData, tName)
	data.AddComponent(fBkg1, tName)

	// Background A sample including three components.
	bkgA := ana.NewSample("BkgTotA", "bkg", `Total Bkg A`, ana.WithWeight(w))
	bkgA.AddComponent(fBkg1, tName)
	bkgA.AddComponent(fBkg2, tName)
	bkgA.AddComponent(fBkg1, tName, ana.WithCut(isQQ))

	// Background B sample including two components.
	bkgB := ana.NewSample("BkgTotB", "bkg", `Total Bkg B`, ana.WithWeight(w))
	bkgB.AddComponent(fBkg1, tName)
	bkgB.AddComponent(fBkg2, tName)

	// Put samples together.
	samples := []*ana.Sample{data, bkgA, bkgB}

	// Define variables
	variables := []*ana.Variable{
		ana.NewVariable("Mttbar", ana.TreeVarF32("ttbar_m"), 25, 350, 1000),
		ana.NewVariable("DphiLL", ana.TreeVarF64("truth_dphi_ll"), 10, 0, math.Pi),
	}

	// Define some selections
	selections := []*ana.Selection{
		ana.NewSelection("LowM", cutMlt500),
		ana.NewSelection("HighM", cutMgt500),
	}

	// Create analyzer object with normalized histograms.
	analyzer := ana.New(samples, variables,
		ana.WithKinemCuts(selections),
		ana.WithDumpTree(true),
		ana.WithSavePath("testdata/Plots_withTreeDumping"),
	)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

func Example_withSliceVariables() {
	// File and tree names
	fName, tName := "../testdata/fileSlices.root", "modules"

	// Samples
	samples := []*ana.Sample{
		ana.CreateSample("HGTD", "bkg", `w/o calib.`, fName, tName),
	}

	// Variables
	variables := []*ana.Variable{
		ana.NewVariable("hitTimes", ana.TreeVarF32s("hits_time_mc"), 100, 10, 15),
	}

	// Analyzer
	analyzer := ana.New(samples, variables,
		ana.WithHistoStack(false),
		ana.WithRatioPlot(false),
		ana.WithSavePath("testdata/Plots_withSliceVariables"),
	)

	// Run the analyzer to produce all the plots
	if err := analyzer.Run(); err != nil {
		panic(err)
	}
}

var (
	// Some files and trees names
	fData = "../testdata/file1.root"
	fBkg1 = "../testdata/file2.root"
	fBkg2 = "../testdata/file3.root"
	tName = "truth"

	// Some weights and cuts TreeFunc's
	w1 = ana.TreeValF64(1.0)
	w2 = ana.TreeValF64(0.5)
	w3 = ana.TreeFunc{
		VarsName: []string{"t_pt"},
		Fct:      func(pt float32) float64 { return 1.0 + float64(pt)/50. },
	}
	w4 = ana.TreeFunc{
		VarsName: []string{"t_pt"},
		Fct:      func(pt float32) float64 { return 1.0 - float64(pt)/250. },
	}
	wSigM = func(mass, relWidth float32) ana.TreeFunc {
		return ana.TreeFunc{
			VarsName: []string{"ttbar_m"},
			Fct: func(m float32) float64 {
				dM2 := float64((m - mass) * (m - mass))
				sigma2 := float64(mass * relWidth * mass * relWidth)
				return 50 * 1 / float64(mass*relWidth) * math.Exp(-dM2/sigma2)
			},
		}
	}
	cutMlt500 = ana.TreeFunc{
		VarsName: []string{"ttbar_m"},
		Fct:      func(m float32) bool { return m < 500 },
	}
	cutMgt500 = ana.TreeFunc{
		VarsName: []string{"ttbar_m"},
		Fct:      func(m float32) bool { return m >= 500 },
	}

	// Some colors
	noColor    = color.NRGBA{}
	softBlack  = color.NRGBA{R: 50, G: 30, B: 50, A: 200}
	shadowBlue = color.NRGBA{R: 50, G: 20, B: 150, A: 20}
	darkRed    = color.NRGBA{R: 180, G: 30, B: 50, A: 200}
	darkGreen  = color.NRGBA{G: 180, R: 30, B: 50, A: 200}
	darkBlue   = color.NRGBA{B: 180, G: 30, R: 50, A: 200}
)
