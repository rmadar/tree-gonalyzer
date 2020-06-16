package ana

import (
	"log"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

type dumper struct {
	Var   []float64   // Storing the F64 variables values to dump the TTree.
	Vars  [][]float64 // Storing the F64s variable values to dump the TTree.
	VarsN []int32     // Storing the number of object in the F64s to dump the TTree.
}

func (ana *Maker) newDumper() dumper {
	dVar := make([]float64, len(ana.Variables)+len(ana.KinemCuts))
	dVars := make([][]float64, len(ana.Variables))
	dVarsN := make([]int32, len(ana.Variables))
	for i, v := range ana.Variables {
		if v.isSlice {
			dVars[i] = []float64{}
			dVarsN[i] = 0
		} else {
			dVar[i] = 0
		}
	}
	for i := range ana.KinemCuts {
		dVar[ana.nVars+i] = 0
	}

	// Return the dumper
	return dumper{
		Var:   dVar,
		Vars:  dVars,
		VarsN: dVarsN,
	}
}

// Helper function to assess variables type, needed
// to instantiate a dumper before reading a tree.
func (ana *Maker) assessVariableTypes() {

	// Get the main tree
	fName := ana.Samples[0].components[0].FileName
	tName := ana.Samples[0].components[0].TreeName
	f, tMain := getTreeFromFile(fName, tName)
	defer f.Close()

	// Get associated trees
	trees := []rtree.Tree{tMain}
	for _, in := range ana.Samples[0].components[0].JointTrees {
		fJoin, tJoin := getTreeFromFile(in.FileName, in.TreeName)
		trees = append(trees, tJoin)
		defer fJoin.Close()
	}

	// Join them
	t, err := rtree.Join(trees...)
	if err != nil {
		log.Fatalf("could not join trees: %+v", err)
	}

	// Get reader associated to the final tree
	r, err := rtree.NewReader(t, rtree.NewReadVars(t))
	if err != nil {
		log.Fatal("could not create tree reader: %w", err)
	}
	defer r.Close()

	// Loop over variable to assess whether they are float64
	// or a slice of float64.
	for _, v := range ana.Variables {
		v.isSlice = false
		if _, ok := v.TreeFunc.GetFuncF64(r); !ok {
			v.isSlice = true
			if _, ok = v.TreeFunc.GetFuncF64s(r); !ok {
				err := "Type assertion failed [variable \"%v\"]:"
				err += " TreeFunc.Fct must return a float64 or a []float64."
				log.Fatalf(err, v.Name)
			}
		}
	}
}

// Helper function creating a file and tree to be dumped.
func (ana *Maker) getOutFileTree(fname, tname string, d dumper) (*groot.File, rtree.Writer) {

	// Create a new ROOT file
	f, err := groot.Create(fname)
	if err != nil {
		log.Fatalf("could not create ROOT file %v: %v", fname, err)
	}

	// Variables to save
	wvars := []rtree.WriteVar{}
	for i, v := range ana.Variables {
		if v.isSlice {
			wvars = append(wvars, rtree.WriteVar{
				Name:  v.Name + "N",
				Value: &d.VarsN[i]},
			)
			wvars = append(wvars, rtree.WriteVar{
				Name:  v.Name,
				Value: &d.Vars[i],
				Count: v.Name + "N"},
			)
		} else {
			wvars = append(wvars, rtree.WriteVar{
				Name:  v.Name,
				Value: &d.Var[i]},
			)
		}
	}
	for i, s := range ana.KinemCuts {
		wvars = append(wvars, rtree.WriteVar{
			Name:  "pass" + s.Name,
			Value: &d.Var[ana.nVars+i]},
		)
	}

	// Create a new TTree
	t, err := rtree.NewWriter(f, tname, wvars)
	if err != nil {
		log.Fatal("could not create tree writer: %w", err)
	}

	return f, t
}
