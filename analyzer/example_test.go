package analyzer_test

import (
	"log"
	"fmt"
	"math"
	"image/color"
	
	"gonum.org/v1/plot/vg"
	
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
	
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/hplot/htex"
	
	"github.com/rmadar/hplot-style/style"

	"tree-gonalyzer/sample"
	"tree-gonalyzer/variable"
)

// All samples
var all_samples = []sample.Spl {
	sample.Spl{
		FileName: "../../../../data/outputs/MC16a.410472.PhPy8EG.TruthOnly.root",
		TreeName: "truth",
		LegLabel: `$t\bar{t}+$jets PhP8`,
		LineWidth: 0,
		LineColor: color.NRGBA{R:  50, G:  20, B: 150, A: 20},
		FillColor: color.NRGBA{R:  50, G:  20, B: 150, A: 20},
		CircleMarkers: false,
		WithYErrBars: false,
	},
		
	sample.Spl{
		FileName: "../../../../data/outputs/ttbar_0j_parton_MG_MadSpin.root",
		TreeName: "truth",
		LegLabel: `$t\bar{t}$ MG w/ spin`,
		LineColor: color.NRGBA{R:  180, G:  30, B: 50, A: 255},
		LineWidth: 2,
		CircleMarkers: false,
		CircleSize: 3,
		WithYErrBars: false,
	},
	
	sample.Spl{
		FileName: "../../../../data/outputs/ttbar_0j_parton_MG_MadSpinCorrOff.root",
		TreeName: "truth",
		LegLabel: `$t\bar{t}$ MG w/o spin`,
		LineColor: color.NRGBA{R:  50, G:  50, B: 180, A: 255},
		LineWidth: 2,
		CircleMarkers: false,
		CircleSize: 1.5,
		WithYErrBars: false,
	},

	sample.Spl{
		FileName: "../../../../data/outputs/ttbar_0j_parton_MG_fullME.root",
		TreeName: "truth",
		LegLabel: `$t\bar{t}$ MG full ME`,
		LineColor: color.NRGBA{R:  50, G:  180, B: 180, A: 255},
		LineWidth: 2,
		CircleMarkers: false,
		CircleSize: 1.5,
		WithYErrBars: false,
	},		
}


// All variable.variables
var (
	var_dphi= &variable.Var{
		OutputName: "truth_dphi_ll.tex",
		TreeName: "truth_dphi_ll",
		Value: new(float64),
		Type: "float64",
		Nbins: 25,
		Xmin: 0, 
		Xmax: math.Pi,
		PlotTitle: `\textbf{ATLAS} Simulation -- $\sqrt{s}=13\,$TeV`,
		XLabel: `$\Delta\phi_{\ell\ell}$`,
		YLabel: `PDF($\Delta\phi_{\ell\ell}$)`,
		LegPosTop: true,
		LegPosLeft: true,
		RangeYmax: 0.08,
	}
	
	var_Ckk = &variable.Var{
		OutputName: "truth_Ckk.tex",
		TreeName: "truth_Ckk",
		Value: new(float64),
		Type: "float64",
		Nbins: 50,
		Xmin: -1, 
		Xmax:  1,
		PlotTitle: `\textbf{ATLAS} Simulation -- $\sqrt{s}=13\,$TeV`,
		XLabel: `$\cos\theta^{+}_{k} \: \cos\theta^{-}_{k}$`,
		YLabel: `PDF($\cos\theta^{+}_{k} \cos\theta^{-}_{k}$)`,
		RangeXmin:  -1.5,
		RangeXmax:  1,
		LegPosTop: true,
		LegPosLeft: true,
	}

	var_Crr = &variable.Var{
		OutputName: "truth_Crr.tex",
		TreeName: "truth_Crr",
		Value: new(float64),
		Type: "float64",
		Nbins: 50,
		Xmin: -1, 
		Xmax:  1,
		PlotTitle: `\textbf{ATLAS} Simulation -- $\sqrt{s}=13\,$TeV`,
		XLabel: `$\cos\theta^{+}_{r} \: \cos\theta^{-}_{r}$`,
		YLabel: `PDF($\cos\theta^{+}_{r} \cos\theta^{-}_{r}$)`,
		RangeXmin:  -1.5,
		RangeXmax:  1,
		LegPosTop: true,
		LegPosLeft: true,
	}
	
	var_Cnn = &variable.Var{
		OutputName: "truth_Cnn.tex",
		TreeName: "truth_Cnn",
		Value: new(float64),
		Type: "float64",
		Nbins: 50,
		Xmin: -1, 
		Xmax:  1,
		PlotTitle: `\textbf{ATLAS} Simulation -- $\sqrt{s}=13\,$TeV`,
		XLabel: `$\cos\theta^{+}_{n} \: \cos\theta^{-}_{n}$`,
		YLabel: `PDF($\cos\theta^{+}_{n} \cos\theta^{-}_{n}$)`,
		RangeXmin:  -1.5,
		RangeXmax:  1,
		LegPosTop: true,
		LegPosLeft: true,
	}

	var_pt_lep = &variable.Var{
		OutputName: "pt_lep.tex",
		TreeName: "l_pt",
		Value: new(float32),
		Type: "float32",
		Nbins: 50,
		Xmin: 0, 
		Xmax: 500,
		PlotTitle: `\textbf{ATLAS} Simulation -- $\sqrt{s}=13\,$TeV`,
		XLabel: `$p^{\ell}_{T}$ [GeV]`,
		YLabel: `PDF($p^{\ell}_{T}$)`,
		LegPosTop: true,
		LegPosLeft: false,

	}
	
	var_eta_lep = &variable.Var{
		OutputName: "eta_lep.tex",
		TreeName: "l_eta",
		Value: new(float32),
		Type: "float32",
		Nbins: 50,
		Xmin: -5, 
		Xmax: 5,
		PlotTitle: `\textbf{ATLAS} Simulation -- $\sqrt{s}=13\,$TeV`,
		XLabel: `$\eta^{\ell}$`,
		YLabel: `PDF($\eta^{\ell}$)`,
		LegPosTop: true,
		LegPosLeft: false,
		RangeXmax: 8,

	}

	var_pt_b = &variable.Var{
		OutputName: "pt_b.tex",
		TreeName: "b_pt",
		Value: new(float32),
		Type: "float32",
		Nbins: 50,
		Xmin: 0, 
		Xmax: 500,
		PlotTitle: `\textbf{ATLAS} Simulation -- $\sqrt{s}=13\,$TeV`,
		XLabel: `$p^{b}_{T}$ [GeV]`,
		YLabel: `PDF($p^{b}_{T}$)`,
		LegPosTop: true,
		LegPosLeft: false,
	}

	var_eta_b = &variable.Var{
		OutputName: "eta_b.tex",
		TreeName: "b_eta",
		Value: new(float32),
		Type: "float32",
		Nbins: 50,
		Xmin: -5, 
		Xmax: 5,
		PlotTitle: `\textbf{ATLAS} Simulation -- $\sqrt{s}=13\,$TeV`,
		XLabel: `$\eta^{b}$`,
		YLabel: `PDF($\eta^{b}$)`,
		LegPosTop: true,
		LegPosLeft: false,
		RangeXmax:  8,
	}

	var_pt_vsum = &variable.Var{
		OutputName: "pt_vsum.tex",
		TreeName: "vsum_pt",
		Value: new(float32),
		Type: "float32",
		Nbins: 50,
		Xmin: 0, 
		Xmax: 500,
		PlotTitle: `\textbf{ATLAS} Simulation -- $\sqrt{s}=13\,$TeV`,
		XLabel: `Truth $E^{\mathrm{miss}}_{T} \; \equiv \; |\vec{p}^{\,\nu}_T + \vec{p}^{\,\bar{\nu}}_T|$`,
		YLabel: `PDF($E^{\mathrm{miss}}_{T}$)`,
		LegPosTop: true,
		LegPosLeft: false,
	}
	
	var_pt_t = &variable.Var{
		OutputName: "pt_t.tex",
		TreeName: "t_pt",
		Value: new(float32),
		Type: "float32",
		Nbins: 50,
		Xmin: 0, 
		Xmax: 500,
		PlotTitle: `\textbf{ATLAS} Simulation -- $\sqrt{s}=13\,$TeV`,
		XLabel: `$p^{t}_{T}$ [GeV]`,
		YLabel: `PDF($p^{t}_{T}$)`,
		LegPosTop: true,
		LegPosLeft: false,
	}

	var_eta_t = &variable.Var{
		OutputName: "eta_t.tex",
		TreeName: "t_eta",
		Value: new(float32),
		Type: "float32",
		Nbins: 50,
		Xmin: -5, 
		Xmax: 5,
		PlotTitle: `\textbf{ATLAS} Simulation -- $\sqrt{s}=13\,$TeV`,
		XLabel: `$\eta^{t}$`,
		YLabel: `PDF($\eta^{t}$)`,
		LegPosTop: true,
		LegPosLeft: false,
		RangeXmax:  8,
	}

	var_m_tt = &variable.Var{
		OutputName: "m_tt.tex",
		TreeName: "ttbar_m",
		Value: new(float32),
		Type: "float32",
		Nbins: 50,
		Xmin: 300, 
		Xmax: 1000,
		PlotTitle: `\textbf{ATLAS} Simulation -- $\sqrt{s}=13\,$TeV`,
		XLabel: `$m_{t\bar{t}}$ [GeV]`,
		YLabel: `PDF($m_{t\bar{t}}$)`,
		LegPosTop: true,
		LegPosLeft: false,
		RangeXmin: 300,
	}

	var_pt_tt = &variable.Var{
		OutputName: "pt_tt.tex",
		TreeName: "ttbar_pt",
		Value: new(float32),
		Type: "float32",
		Nbins: 50,
		Xmin: 0, 
		Xmax: 150,
		PlotTitle: `\textbf{ATLAS} Simulation -- $\sqrt{s}=13\,$TeV`,
		XLabel: `$p^{t\bar{t}}_T$ [GeV]`,
		YLabel: `PDF($p^{t\bar{t}}_T$)`,
		LegPosTop: true,
		LegPosLeft: false,
	} 
)

var all_variables = []*variable.Var {
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
}

func main(){

	// Get histos
	histos_Nvars_Nsamples, err := makeHistos()
	if err != nil {
		panic(err)
	}

	// Plot them on the same canvas
	plotHistos(histos_Nvars_Nsamples)
}


func makeHistos() ([][]*hbook.H1D, error) {

	// Build a 2D-slices of histogram [nall_variables][nsamples]
	s2d_histos := make([][]*hbook.H1D, len(all_variables))
	for i := range s2d_histos{
		s2d_histos[i] = make([]*hbook.H1D, len(all_samples))
	}
	
	// Loop over the samples
	for is, s := range all_samples {
		
		// Create H1D objects for each variables, for each samples
		for iv, v := range all_variables {
			s2d_histos[iv][is] = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
		}

		func(j int) { // Anonymous function to avoid memory leak due to 'defer'
		
			// Get the file and tree
			f, tree := getTreeFromFile(s.FileName, s.TreeName)
			defer f.Close()
			
			// Prepare the variables to read
			rvars := []rtree.ReadVar{}
			for _, v := range all_variables {
				rvars = append(rvars, rtree.ReadVar{Name: v.TreeName, Value: v.Value})
			}
			
			// Get the tree reader
			r, err := rtree.NewReader(tree, rvars)
			if err != nil {
				log.Fatalf("could not create tree reader: %w", err)
			}
			defer r.Close()
			
			// Read the tree
			err = r.Read(func(ctx rtree.RCtx) error {
				
				for iv, v := range all_variables {
					s2d_histos[iv][is].Fill(v.GetValue(), 1)
				}
				
				return nil
			})
		}(is)
			
	}	

	return s2d_histos, nil
}



func plotHistos(histos [][]*hbook.H1D) {

	// Loop over variables and get histo for all samples
	for iv, hsamples := range histos { 
		
		// Create a new styled plot
		p := hplot.New()
		p.Latex = htex.DefaultHandler
		style.ApplyToPlot(p)

		// Plot labels
		p.Title.Text = all_variables[iv].PlotTitle
		p.X.Label.Text = all_variables[iv].XLabel
		p.Y.Label.Text = all_variables[iv].YLabel
		p.X.Min = all_variables[iv].RangeXmin
		p.X.Max = all_variables[iv].RangeXmax
		p.Y.Min = all_variables[iv].RangeYmin
		p.Y.Max = all_variables[iv].RangeYmax
		
		// Loop over samples and turn hook.H1D into styled plottable histo
		for is, h := range hsamples {
			h.Scale(1.0/h.Integral())
			hist := all_samples[is].CreateHisto(h)
			p.Legend.Add(all_samples[is].LegLabel, hist)
			p.Add(hist)
		}

		// Legend
		p.Legend.Top = all_variables[iv].LegPosTop
		p.Legend.Left = all_variables[iv].LegPosLeft
		p.Legend.YOffs = -5
		if p.Legend.Left {
			p.Legend.XOffs = 5
		} else {
			p.Legend.XOffs = -5
		}
		p.Legend.TextStyle.Font.Size = 12
		
		// Save the plot
		if err := p.Save(5.5*vg.Inch, 4*vg.Inch, "results/"+all_variables[iv].OutputName); err != nil {
			log.Fatalf("error saving plot: %v\n", err)
		}
	}
}


// Helper to get a tree from a file
func getTreeFromFile(filename, treename string) (*groot.File, rtree.Tree) {

	// Get the file
	f, err := groot.Open(filename)
	if err != nil {
		err := fmt.Sprintf("could not open ROOT file %q: %w", filename, err)
		panic(err)
	}
	
	// Get the tree
	obj, err := f.Get(treename)
	if err != nil {
		err := fmt.Sprintf("could not retrieve object: %w", err)
		panic(err)
	}		

	return f, obj.(rtree.Tree)
}
