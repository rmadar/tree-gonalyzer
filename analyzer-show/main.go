// Example to run gonalyzer package
package main

import (
	"math"
	
	"image/color"
	
	"tree-gonalyzer/analyzer"
	"tree-gonalyzer/sample"
	"tree-gonalyzer/variable"

	"github.com/rmadar/hplot-style/style"
)

// Run the analyzer
func main(){

	// Create analyzer object
	ana := analyzer.Ana{
		
		Samples: []sample.Spl{
			spl_data,
			spl_bkg1,
			spl_bkg2,
			spl_alt,
		},
		
		Variables: []*variable.Var{
			var_dphi,
			var_Ckk,
			var_Crr,
			var_Cnn,
			var_pt_lep,
			var_eta_lep,
			var_pt_b,
			var_eta_b,
			var_pt_t,
			var_eta_t,
			var_pt_vsum,
			var_m_tt,
			var_pt_tt,
		},
	}
	
	// Create histograms via an event loop
	err := ana.MakeHistos()
	if err != nil {
		panic(err)
	}
	
	// Plot them on the same canvas
	err = ana.PlotHistos()
	if err != nil {
		panic(err)
	}

}

// Define all samples and variables of the analysis
var (
	// samples
	spl_data = sample.Spl{
		Name: "data",
		FileName: "../testdata/ttbar_ME.root",
		TreeName: "truth",
		LegLabel: `Pseudo-data`,
		CircleMarkers: true,
		CircleColor: style.SmoothBlack, 
		CircleSize: 2,
		WithYErrBars: true,
		YErrBarsLineWidth: 1,
		YErrBarsCapWidth: 5,
	}


	spl_bkg1 = sample.Spl{
		Name: "bkg1",
		FileName: "../testdata/ttbar_MadSpinOn_1.root",
		TreeName: "truth",
		LegLabel: `$t\bar{t}$ 1/2`,
		LineColor: color.NRGBA{R:  50, G:  20, B: 150, A: 20},
		FillColor: color.NRGBA{R:  50, G:  20, B: 150, A: 20},
		LineWidth: 0,
		CircleMarkers: false,
		CircleSize: 1.5,
		WithYErrBars: true,
	}

	spl_bkg2 = sample.Spl{
		Name: "bkg2",
		FileName: "../testdata/ttbar_MadSpinOn_2.root",
		TreeName: "truth",
		LegLabel: `$t\bar{t}$ 2/2 (nom)`,
		LineColor: color.NRGBA{R:  50, G:  20, B: 150, A: 20},
		FillColor: color.NRGBA{R:  50, G:  20, B: 150, A: 20},
		LineWidth: 0,
		CircleMarkers: false,
		CircleSize: 1.5,
		WithYErrBars: true,
	}		
	
	spl_alt = sample.Spl{
		Name: "spinoff",
		FileName: "../testdata/ttbar_MadSpinOff.root",
		TreeName: "truth",
		LegLabel: `$t\bar{t}$ no spin effects (alt.)`,
		LineColor: color.NRGBA{R:  50, G:  50, B: 180, A: 255},
		LineWidth: 2,
		CircleMarkers: false,
		CircleSize: 1.5,
		WithYErrBars: false,
	}

	// Variables
	var_dphi= &variable.Var{
		Name: "truth_dphi_ll",
		SaveName: "truth_dphi_ll.tex",
		TreeName: "truth_dphi_ll",
		Value: new(float64),
		Nbins: 15,
		Xmin: 0, 
		Xmax: math.Pi,
		PlotTitle: `{\tt TTree} {\bf GO}nalyzer -- $pp \to t\bar{t}$ @ $13\,$ TeV`,
		XLabel: `$\Delta\phi_{\ell\ell}$`,
		YLabel: `PDF($\Delta\phi_{\ell\ell}$)`,
		LegPosTop: true,
		LegPosLeft: true,
		RangeYmax: 0.08,
	}
	
	var_Ckk = &variable.Var{
		Name: "truth_Ckk",
		SaveName: "truth_Ckk.tex",
		TreeName: "truth_Ckk",
		Value: new(float64),
		Nbins: 25,
		Xmin: -1, 
		Xmax:  1,
		PlotTitle: `{\tt TTree} {\bf GO}nalyzer -- $pp \to t\bar{t}$ @ $13\,$ TeV`,
		XLabel: `$\cos\theta^{+}_{k} \: \cos\theta^{-}_{k}$`,
		YLabel: `PDF($\cos\theta^{+}_{k} \cos\theta^{-}_{k}$)`,
		RangeXmin:  -1.5,
		RangeXmax:  1,
		LegPosTop: true,
		LegPosLeft: true,
	}

	var_Crr = &variable.Var{
		Name: "truth_Crr",
		SaveName: "truth_Crr.tex",
		TreeName: "truth_Crr",
		Value: new(float64),
		Nbins: 25,
		Xmin: -1, 
		Xmax:  1,
		PlotTitle: `{\tt TTree} {\bf GO}nalyzer -- $pp \to t\bar{t}$ @ $13\,$ TeV`,
		XLabel: `$\cos\theta^{+}_{r} \: \cos\theta^{-}_{r}$`,
		YLabel: `PDF($\cos\theta^{+}_{r} \cos\theta^{-}_{r}$)`,
		RangeXmin:  -1.5,
		RangeXmax:  1,
		LegPosTop: true,
		LegPosLeft: true,
	}
	
	var_Cnn = &variable.Var{
		Name: "truth_Cnn",
		SaveName: "truth_Cnn.tex",
		TreeName: "truth_Cnn",
		Value: new(float64),
		Nbins: 25,
		Xmin: -1, 
		Xmax:  1,
		PlotTitle: `{\tt TTree} {\bf GO}nalyzer -- $pp \to t\bar{t}$ @ $13\,$ TeV`,
		XLabel: `$\cos\theta^{+}_{n} \: \cos\theta^{-}_{n}$`,
		YLabel: `PDF($\cos\theta^{+}_{n} \cos\theta^{-}_{n}$)`,
		RangeXmin:  -1.5,
		RangeXmax:  1,
		LegPosTop: true,
		LegPosLeft: true,
	}

	var_pt_lep = &variable.Var{
		Name: "pt_lep",
		SaveName: "pt_lep.tex",
		TreeName: "l_pt",
		Value: new(float32),
		Nbins: 25,
		Xmin: 0, 
		Xmax: 500,
		PlotTitle: `{\tt TTree} {\bf GO}nalyzer -- $pp \to t\bar{t}$ @ $13\,$ TeV`,
		XLabel: `$p^{\ell}_{T}$ [GeV]`,
		YLabel: `PDF($p^{\ell}_{T}$)`,
		LegPosTop: true,
		LegPosLeft: false,

	}
	
	var_eta_lep = &variable.Var{
		Name: "eta_lep",
		SaveName: "eta_lep.tex",
		TreeName: "l_eta",
		Value: new(float32),
		Nbins: 25,
		Xmin: -5, 
		Xmax: 5,
		PlotTitle: `{\tt TTree} {\bf GO}nalyzer -- $pp \to t\bar{t}$ @ $13\,$ TeV`,
		XLabel: `$\eta^{\ell}$`,
		YLabel: `PDF($\eta^{\ell}$)`,
		LegPosTop: true,
		LegPosLeft: false,
		RangeXmax: 8,

	}

	var_pt_b = &variable.Var{
		Name: "pt_b",
		SaveName: "pt_b.tex",
		TreeName: "b_pt",
		Value: new(float32),
		Nbins: 25,
		Xmin: 0, 
		Xmax: 500,
		PlotTitle: `{\tt TTree} {\bf GO}nalyzer -- $pp \to t\bar{t}$ @ $13\,$ TeV`,
		XLabel: `$p^{b}_{T}$ [GeV]`,
		YLabel: `PDF($p^{b}_{T}$)`,
		LegPosTop: true,
		LegPosLeft: false,
	}

	var_eta_b = &variable.Var{
		Name: "eta_b",
		SaveName: "eta_b.tex",
		TreeName: "b_eta",
		Value: new(float32),
		Nbins: 25,
		Xmin: -5, 
		Xmax: 5,
		PlotTitle: `{\tt TTree} {\bf GO}nalyzer -- $pp \to t\bar{t}$ @ $13\,$ TeV`,
		XLabel: `$\eta^{b}$`,
		YLabel: `PDF($\eta^{b}$)`,
		LegPosTop: true,
		LegPosLeft: false,
		RangeXmax:  8,
	}

	var_pt_vsum = &variable.Var{
		Name: "vsum_pt",
		SaveName: "pt_vsum.tex",
		TreeName: "vsum_pt",
		Value: new(float32),
		Nbins: 25,
		Xmin: 0, 
		Xmax: 250,
		PlotTitle: `{\tt TTree} {\bf GO}nalyzer -- $pp \to t\bar{t}$ @ $13\,$ TeV`,
		XLabel: `Truth $E^{\mathrm{miss}}_{T} \; \equiv \; |\vec{p}^{\,\nu}_T + \vec{p}^{\,\bar{\nu}}_T|$`,
		YLabel: `PDF($E^{\mathrm{miss}}_{T}$)`,
		LegPosTop: true,
		LegPosLeft: false,
	}
	
	var_pt_t = &variable.Var{
		Name: "t_pt",
		SaveName: "pt_t.tex",
		TreeName: "t_pt",
		Value: new(float32),
		Nbins: 100,
		Xmin: 0, 
		Xmax: 500,
		PlotTitle: `{\tt TTree} {\bf GO}nalyzer -- $pp \to t\bar{t}$ @ $13\,$ TeV`,
		XLabel: `$p^{t}_{T}$ [GeV]`,
		YLabel: `PDF($p^{t}_{T}$)`,
		LegPosTop: true,
		LegPosLeft: false,
	}

	var_eta_t = &variable.Var{
		Name: "eta_t",
		SaveName: "eta_t.tex",
		TreeName: "t_eta",
		Value: new(float32),
		Nbins: 25,
		Xmin: -5, 
		Xmax: 5,
		PlotTitle: `{\tt TTree} {\bf GO}nalyzer -- $pp \to t\bar{t}$ @ $13\,$ TeV`,
		XLabel: `$\eta^{t}$`,
		YLabel: `PDF($\eta^{t}$)`,
		LegPosTop: true,
		LegPosLeft: false,
		RangeXmax:  8,
	}

	var_m_tt = &variable.Var{
		Name: "m_tt",
		SaveName: "m_tt.tex",
		TreeName: "ttbar_m",
		Value: new(float32),
		Nbins: 25,
		Xmin: 300, 
		Xmax: 1000,
		PlotTitle: `{\tt TTree} {\bf GO}nalyzer -- $pp \to t\bar{t}$ @ $13\,$ TeV`,
		XLabel: `$m_{t\bar{t}}$ [GeV]`,
		YLabel: `PDF($m_{t\bar{t}}$)`,
		LegPosTop: true,
		LegPosLeft: false,
		RangeXmin: 300,
	}

	var_pt_tt = &variable.Var{
		Name: "pt_tt",
		SaveName: "pt_tt.tex",
		TreeName: "ttbar_pt",
		Value: new(float32),
		Nbins: 25,
		Xmin: 0, 
		Xmax: 150,
		PlotTitle: `{\tt TTree} {\bf GO}nalyzer -- $pp \to t\bar{t}$ @ $13\,$ TeV`,
		XLabel: `$p^{t\bar{t}}_T$ [GeV]`,
		YLabel: `PDF($p^{t\bar{t}}_T$)`,
		LegPosTop: true,
		LegPosLeft: false,
	} 
)
