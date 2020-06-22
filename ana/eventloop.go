package ana

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"

	"go-hep.org/x/hep/hbook"
)

// RunEventLoops runs one event loop per sample to fill
// histograms for each variables and selections. If DumpTree
// is true, a tree is also dumped with all variables and one
// branch per selection.
func (ana *Maker) RunEventLoops() error {

	// Start timing
	start := time.Now()

	// Initialize hbook H1D as N[samples] 2D-slices.
	ana.hbookHistos = make([][][]*hbook.H1D, len(ana.Samples))
	
	// Loop over the samples
	if ana.SampleMT {
		var wg sync.WaitGroup
		wg.Add(len(ana.Samples))
		for i := range ana.Samples {
			go ana.concurrentSampleEventLoop(i, &wg)
		}
		wg.Wait()
	} else {
		for i := range ana.Samples {
			ana.sampleEventLoop(i)
		}
	}

	for _, n := range ana.nEvtsSample {
		ana.nEvents += n

	}
	// Histograms are now filled.
	ana.histoFilled = true

	// End timing.
	ana.timeLoop = time.Since(start)

	return nil
}

func (ana *Maker) concurrentSampleEventLoop(sampleIdx int, wg *sync.WaitGroup) {

	// Handle concurrency
	defer wg.Done()

	// Fill the histo
	ana.sampleEventLoop(sampleIdx)
}

func (ana *Maker) sampleEventLoop(sampleIdx int) {

	// Current sample
	samp := ana.Samples[sampleIdx]

	// Initiate the structure of the histo container: h[iCut][iVar]
	h := make([][]*hbook.H1D, len(ana.KinemCuts))
	for iCut := range ana.KinemCuts {
		h[iCut] = make([]*hbook.H1D, len(ana.Variables))
		for iVar, v := range ana.Variables {
			if ana.PlotHisto {
				h[iCut][iVar] = hbook.NewH1D(v.Nbins, v.Xmin, v.Xmax)
			} else {
				h[iCut][iVar] = hbook.NewH1D(1, 0, 1)
			}
		}
	}

	// Output in case of TTree dumping
	path := ana.SavePath + "/ntuples/"
	if _, err := os.Stat(path); os.IsNotExist(err) && ana.DumpTree {
		os.MkdirAll(path, 0755)
	}
	var fOut *groot.File
	var tOut rtree.Writer
	dump := ana.newDumper()
	if ana.DumpTree {
		fOut, tOut = ana.getOutFileTree(path+samp.Name+".root", "GOtree", dump)
		defer fOut.Close()
		defer tOut.Close()
	}

	// Loop over the sample components
	for iComp, comp := range samp.components {

		// Anonymous function to avoid memory-leaks due to 'defer'
		func(j int) error {

			// Get the main file and tree
			f, tMain := getTreeFromFile(comp.FileName, comp.TreeName)
			defer f.Close()

			// Get the trees to be joint
			trees := []rtree.Tree{tMain}
			for _, in := range comp.JointTrees {
				fJoin, tJoin := getTreeFromFile(in.FileName, in.TreeName)
				trees = append(trees, tJoin)
				defer fJoin.Close()
			}
			t, err := rtree.Join(trees...)
			if err != nil {
				log.Fatalf("could not join trees: %+v", err)
			}

			// Get the tree reader
			r, err := rtree.NewReader(t, []rtree.ReadVar{}, rtree.WithRange(0, ana.NevtsMax))
			if err != nil {
				log.Fatal("could not create tree reader: %w", err)
			}
			defer r.Close()

			// Prepare variables
			ok := false
			getF64 := make([]func() float64, len(ana.Variables))
			getF64s := make([]func() []float64, len(ana.Variables))
			for iv, v := range ana.Variables {
				idx := iv
				if !v.isSlice {
					if getF64[idx], ok = v.TreeFunc.GetFuncF64(r); !ok {
						err := "Type assertion failed [variable \"%v\"]:"
						err += " this TreeFunc.Fct is supposed to return a float64."
						log.Fatalf(err, v.Name)
					}
				} else {
					if getF64s[idx], ok = v.TreeFunc.GetFuncF64s(r); !ok {
						err := "Type assertion failed [variable \"%v\"]:"
						err += " this TreeFunc.Fct is supposed to return a []float64."
						log.Fatalf(err, v.Name)
					}
				}
			}

			// Prepare the sample global weight
			getWeightSamp := func() float64 { return 1.0 }
			if samp.WeightFunc.Fct != nil {
				if getWeightSamp, ok = samp.WeightFunc.GetFuncF64(r); !ok {
					err := "Type assertion failed [weight of %v]:"
					err += " TreeFunc.Fct must return a float64."
					log.Fatalf(err, samp.Name)
				}
			}

			// Prepare the normalization weight of this component
			normWeight := ana.Lumi * 1000 * comp.Xsec / comp.Ngen
			if samp.sType == data {
				normWeight = 1.0
			}

			// Prepare the additional weight of the component
			getWeightComp := func() float64 { return 1.0 }
			if comp.WeightFunc.Fct != nil {
				if getWeightComp, ok = comp.WeightFunc.GetFuncF64(r); !ok {
					err := "Type assertion failed [weight of (%v, %v)]:"
					err += " TreeFunc.Fct must return a float64."
					log.Fatalf(err, samp.Name, comp.FileName)
				}
			}

			// Prepare the sample global cut
			passCutSamp := func() bool { return true }
			if samp.CutFunc.Fct != nil {
				if passCutSamp, ok = samp.CutFunc.GetFuncBool(r); !ok {
					err := "Type assertion failed [Cut of %v]:"
					err += " TreeFunc.Fct must return a bool.\n"
					err += "\t -> Make sure to use NewCutBool(), not NewVarBool()."
					log.Fatalf(err, samp.Name)
				}
			}

			// Prepare the component additional cut
			passCutComp := func() bool { return true }
			if comp.CutFunc.Fct != nil {
				if passCutComp, ok = comp.CutFunc.GetFuncBool(r); !ok {
					err := "Type assertion failed [Cut of (%v, %v)]:"
					err += " TreeFunc.Fct must return a bool.\n"
					err += "\t -> Make sure to use NewCutBool(), not NewVarBool()."
					log.Fatalf(err, samp.Name, comp.FileName)
				}
			}

			// Prepare the cut string for kinematics
			passKinemCut := make([]func() bool, len(ana.KinemCuts))
			for ic, cut := range ana.KinemCuts {
				idx := ic
				if passKinemCut[idx], ok = cut.TreeFunc.GetFuncBool(r); !ok {
					err := "Type assertion failed [selection \"%v\"]:"
					err += " TreeFunc.Fct must return a bool.\n"
					err += "\t -> Make sure to use NewCutBool(), not NewVarBool()."
					log.Fatal(fmt.Sprintf(err, cut.Name))
				}
			}

			// Read the tree (event loop)
			err = r.Read(func(ctx rtree.RCtx) error {

				// Sample-level and component-level cut
				if !(passCutSamp() && passCutComp()) {
					return nil
				}

				// Get the event weight
				w := getWeightSamp() * getWeightComp() * normWeight

				// Loop over selection and variables
				for ic := range ana.KinemCuts {

					// Look at the next selection if the event is not selected.
					if !passKinemCut[ic]() {
						dump.Var[ana.nVars+ic] = 0.0
						continue
					} else {
						dump.Var[ana.nVars+ic] = 1.0
					}

					// Otherwise, loop over variables.
					for iv, v := range ana.Variables {

						// Fill histo (and fill tree) with full slices...
						if v.isSlice {
							xs := getF64s[iv]()
							for _, x := range xs {
								h[ic][iv].Fill(x, w)
							}
							if ana.DumpTree {
								dump.Vars[iv] = xs
								dump.VarsN[iv] = int32(len(xs))
							}

						} else {
							// ... or the single variable value.
							x := getF64[iv]()
							h[ic][iv].Fill(x, w)
							if ana.DumpTree {
								dump.Var[iv] = x
							}
						}
					}
				}

				if ana.DumpTree {
					_, err = tOut.Write()
					if err != nil {
						log.Fatalf("could not write event in a tree: %+v", err)
					}
				}

				return nil
			})

			// Error check of rtree.Reader
			if err != nil {
				log.Fatalf("could not read tree: %+v", err)
			}

			// Keep track of the number of processed events.
			switch ana.NevtsMax {
			case -1:
				ana.nEvtsSample[sampleIdx] += t.Entries()
			default:
				ana.nEvtsSample[sampleIdx] += ana.NevtsMax
			}

			return nil
		}(iComp)
	}

	// Fill the histos for this sample
	ana.hbookHistos[sampleIdx] = h

	// Explicitely close file and tree
	if ana.DumpTree {
		if err := tOut.Close(); err != nil {
			log.Fatalf("could not close tree: %+v", err)
		}
		if err := fOut.Close(); err != nil {
			log.Fatalf("could not close root file: %+v", err)
		}
	}

}

// Helper to get a tree from a file
func getTreeFromFile(filename, treename string) (*groot.File, rtree.Tree) {

	// Get the file
	f, err := groot.Open(filename)
	if err != nil {
		log.Fatal("could not open ROOT file: %w", err)
	}

	// Get the tree
	obj, err := f.Get(treename)
	if err != nil {
		log.Fatal("could not retrieve object: %w", err)
	}

	return f, obj.(rtree.Tree)
}
