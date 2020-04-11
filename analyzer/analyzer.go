// Package allowing to wrap all needed element of a TTree plotting analysis
package analyzer

import (
	"log"
	"fmt"
	
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

// Analyzer type
type Obj struct {
	Samples      []sample.Obj     // sample on which to run
	SamplesGroup string           // specify how to group samples together
	Variables    []*variable.Obj  // variables to plot
	Selections   []string         // implement a type selection ?
	HistosData   [][]*hbook.H1D   // Currently 2D histo container, later: n-dim [var, sample, cut, syst]
	HistosPlot   [][]*hplot.H1D   // Currently 2D histo container, later: n-dim [var, sample, cut, syst]
}


// Initialize histograms container shape
func (ana *Obj) initHistosData(){
	ana.HistosData = make([][]*hbook.H1D, len(ana.Variables))
	for iv := range ana.HistosData {
		ana.HistosData[iv] = make([]*hbook.H1D, len(ana.Samples))
		v := ana.Variables[iv]
		for is := range ana.HistosData[iv] {
			ana.HistosData[iv][is] = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
		}
	}
}

// Data histogram accessor
//func (ana Obj) getHistoData(ivar, ispl int) *hbook.1HD {
//	return ana.HistosData[ivar][ispl]
//}

// Run the event loop to fill all histo across samples / variables (and later: cut / systematics)
func (ana *Obj) MakeHistos() error {

	// Build hbook histograms container
	ana.initHistosData()
	
	// Loop over the samples
	for is, s := range ana.Samples {
		
		// Anonymous function to avoid memory-leaks due to 'defer'
		func(j int) error { 
		
			// Get the file and tree
			f, tree := getTreeFromFile(s.FileName, s.TreeName)
			defer f.Close()
			
			// Prepare the variables to read
			rvars := []rtree.ReadVar{}
			for _, v := range ana.Variables {
				rvars = append(rvars, rtree.ReadVar{Name: v.TreeName, Value: v.Value})
			}
			
			// Get the tree reader
			r, err := rtree.NewReader(tree, rvars)
			if err != nil {
				return fmt.Errorf("could not create tree reader: %w", err)
			}
			defer r.Close()
			
			// Read the tree
			err = r.Read(func(ctx rtree.RCtx) error {

				// Later, add a loop on cuts here
				for iv, v := range ana.Variables {
					ana.HistosData[iv][is].Fill(v.GetValue(), 1)
				}

				return nil
			})
			
			return nil	

		}(is)
	}	

	return nil
}


// Plotting all histograms
func (ana *Obj) PlotHistos() error {

	ana.HistosPlot = make([][]*hplot.H1D, len(ana.Variables))
        for iv := range ana.HistosPlot {
		ana.HistosPlot[iv] = make([]*hplot.H1D, len(ana.Samples))
	}
	
	// Loop over variables and get histo for all samples
	for iv, hsamples := range ana.HistosData { 

		// Manipulate the current variable
		v := ana.Variables[iv]

		// Create a new styled plot
		p := hplot.New()
		p.Latex = htex.DefaultHandler
		style.ApplyToPlot(p)
		v.SetPlotStyle(p)

		// Additionnal legend
		p.Legend.Padding = 0.1 * vg.Inch
		p.Legend.ThumbnailWidth = 25
		p.Legend.TextStyle.Font.Size = 12
		
		// Loop over samples and turn hook.H1D into styled plottable histo
		for is, h := range hsamples {
			s := ana.Samples[is]
			h.Scale(1.0/h.Integral())
			ana.HistosPlot[iv][is] = s.CreateHisto(h)
			p.Legend.Add(s.LegLabel, ana.HistosPlot[iv][is])
			if is>0 { // Write data (assumed to be the first sample here) at last
				p.Add(ana.HistosPlot[iv][is])
			}
		}
		// Write data (assumed to be the first sample here) at last
		p.Add(ana.HistosPlot[iv][0])
		
		// Save the plot
		if err := p.Save(5.5*vg.Inch, 4*vg.Inch, "results/"+v.SaveName); err != nil {
			log.Fatalf("error saving plot: %v\n", err)
		}
	}
	
	return nil
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
