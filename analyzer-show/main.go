// Example to run gonalyzer package
package main

import (
	"math"
	
	"image/color"
	
	"github.com/rmadar/hplot-style/style"
	
	"github.com/rmadar/tree-gonalyzer/analyzer"
	"github.com/rmadar/tree-gonalyzer/sample"
	"github.com/rmadar/tree-gonalyzer/variable"
)

// Run the analyzer
func main(){

	// Create analyzer object
	ana := analyzer.Obj{

		DontStack: false,
		
		Normalize: true,
		
		Samples: []sample.Obj{
			spl_data,
			spl_bkg1,
			spl_bkg2,
			spl_alt,
		},
		
		Variables: []*variable.Obj{
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
	spl_data = sample.Obj{
		Name: "data",
		FileName: "../testdata/ttbar_MadSpinOff.root",
		TreeName: "truth",
		LegLabel: `Pseudo-data`, 
		CircleMarkers: true,
		CircleColor: style.SmoothBlack, 
		CircleSize: 3,
		WithYErrBars: true,
		YErrBarsLineWidth: 2,
		YErrBarsCapWidth: 5,
	}


	spl_bkg1 = sample.Obj{
		Name: "bkg1",
		FileName: "../testdata/ttbar_MadSpinOn_1.root",
		TreeName: "truth",
		LegLabel: `$t\bar{t}$ contribution 1`,
		FillColor: color.NRGBA{R:  0, G: 102, B: 255, A: 230},
		LineWidth: 0,
		CircleMarkers: false,
		CircleSize: 1.5,
		WithYErrBars: true,
	}

	spl_bkg2 = sample.Obj{
		Name: "bkg2",
		FileName: "../testdata/ttbar_MadSpinOn_2.root",
		TreeName: "truth",
		LegLabel: `$t\bar{t}$ contribution 2`,
		FillColor: color.NRGBA{R: 255, G: 102, B: 0, A: 200},
		LineColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		LineWidth: 0,
		CircleMarkers: false,
		CircleSize: 1.5,
		WithYErrBars: true,
	}		

	spl_alt = sample.Obj{
		Name: "spinoff",
		FileName: "../testdata/ttbar_ME.root",
		TreeName: "truth",
		LegLabel: `$t\bar{t}$ alternative`,
		FillColor: color.NRGBA{R: 0, G:  204, B:  80, A: 200},
		LineColor: color.NRGBA{R: 255, G:  255, B: 255, A: 255},
		LineWidth: -1,
		CircleMarkers: false,
		CircleSize: 1.5,
		WithYErrBars: false,
	}
	
	// Variables
	var_dphi= &variable.Obj{
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
	
	var_Ckk = &variable.Obj{
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

	var_Crr = &variable.Obj{
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
	
	var_Cnn = &variable.Obj{
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

	var_pt_lep = &variable.Obj{
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
	
	var_eta_lep = &variable.Obj{
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

	var_pt_b = &variable.Obj{
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

	var_eta_b = &variable.Obj{
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

	var_pt_vsum = &variable.Obj{
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
		XTickFormat: "%2.0f",
		LegPosTop: true,
		LegPosLeft: false,
	}
	
	var_pt_t = &variable.Obj{
		Name: "t_pt",
		SaveName: "pt_t.pdf",
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

	var_eta_t = &variable.Obj{
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

	var_m_tt = &variable.Obj{
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

	var_pt_tt = &variable.Obj{
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
