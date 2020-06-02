# TTree GOnalyzer

[![Documentation](https://godoc.org/github.com/rmadar/tree-gonalyzer?status.svg)](https://godoc.org/github.com/rmadar/tree-gonalyzer)

This is a tool written in go to produce publication-quality plots from ROOT TTrees in an flexible and easy way.
This tool is built on top of [go-hep.org](https://go-hep.org). The main supported features are:
 - histograming over many samples, selections and variables
 - computation of new (high-level) variables
 - tree dumping
 - concurent sample processings.

## In a nutshell

```go
// Define samples
samples := []*ana.Sample{
	ana.CreateSample("data", "data", `Data`, "data.root", "mytree"),
	ana.CreateSample("bkg1", "bkg", `Proc 1`, "proc1.root", "mytree"),
	ana.CreateSample("bkg2", "bkg", `Proc 2`, "proc2.root", "mytree"),
	ana.CreateSample("bkg3", "bkg", `Proc 3`, "proc3.root", "mytree"),
}

// Define variables
variables := []*ana.Variable{
	ana.NewVariable("plot1", ana.TreeVarBool("branchBool"), 2, 0, 2),
	ana.NewVariable("plot2", ana.TreeVarF32("branchF32"), 25, 0, 1000),
	ana.NewVariable("plot3", ana.TreeVarF64("branchF64"), 50, 0, 1000),
}

// Create analyzer object with some options
analyzer := ana.New(samples, variables, ana.WithHistoNorm(true))

// Produce plots
analyzer.Run()

```

## Gallery

<table>
  <tr>
    <td><p align="center"><img src="ana/testdata/Plots_simpleUseCase/Mttbar_golden.png">
	Data/Background <a href="https://godoc.org/github.com/rmadar/tree-gonalyzer/ana#example-package--ASimpleUseCase">[code]</a></p>
    </td>
    <td><p align="center"><img src="ana/testdata/Plots_withSignals/Mttbar_golden.png">
	Unstacked signals <a href="https://godoc.org/github.com/rmadar/tree-gonalyzer/ana#example-package--WithSignals">[code]</a></p>
    </td>
    <td><p align="center"><img src="ana/testdata/Plots_withStackedSignals/Mttbar_golden.png">
	Stacked signals <a href="https://godoc.org/github.com/rmadar/tree-gonalyzer/ana#example-package--WithStackedSignals">[code]</a></p>
    </td>
  </tr>
	
  <tr>
    <td><p align="center"><img src="ana/testdata/Plots_shapeDistortion/DphiLL_golden.png">
	Shape distortion <a href="https://godoc.org/github.com/rmadar/tree-gonalyzer/ana#example-package--ShapeDistortion">[code]</a></p>
    </td>
    <td><p align="center"><img src="ana/testdata/Plots_shapeComparison/TopPt_golden.png">
	Shape comparison <a href="https://godoc.org/github.com/rmadar/tree-gonalyzer/ana#example-package--ShapeComparison">[code]</a></p>
    </td>
    <td><p align="center"><img src="ana/testdata/Plots_systVariations/DphiLL_golden.png">
	Systematic variation <a href="https://godoc.org/github.com/rmadar/tree-gonalyzer/ana#example-package--SystematicVariations">[code]</a></p>
    </td>
  </tr>


 </table>

## Performances

![benchmarking](ana-perf/perf.png)

For 2M events and 60 variables, a comparison with similar ROOT-based code
(using `t->Draw()`) gives:
 + `ROOT  -> 6.2 ms/kEvts`
 + `GOHEP -> 2.0 ms/kEvts`
 
Testing on only one variable to avoid event-loop repetition
in case of `t->Draw()` (even if it's probably not like doing N times the loop):
 + `ROOT  -> 0.39 ms/kEvts`
 + `GOHEP -> 0.11 ms/kEvts`
