# TTree GOnalyzer

[![Documentation](https://godoc.org/github.com/rmadar/tree-gonalyzer?status.svg)](https://godoc.org/github.com/rmadar/tree-gonalyzer)

This is a tool written in go to produce publication-quality plots from ROOT TTrees in an flexible and easy way.
This tool is built on top of [go-hep.org](https://go-hep.org).

## In a nutshell

```go
// Define samples
samples := []*ana.Sample{
	ana.CreateSample("data", "data", `Data 2020`, "data.root", "mytree"),
	ana.CreateSample("bkg1", "bkg", `Proc 1`, "proc1.root", "mytree"),
	ana.CreateSample("bkg2", "bkg", `Proc 2`, "proc2.root", "mytree"),
	ana.CreateSample("bkg3", "bkg", `Proc 3`, "proc3.root", "mytree"),
}

// Define variables
variables := []*ana.Variable{
	ana.NewVariable("plot1", "branch1", new(float32), 25, 0, 1000),
	ana.NewVariable("plot2", "branch2", new(float64), 50, 0, 1000),
}

// Create analyzer object with some options
analyzer := ana.New(samples, variables,
	      ana.WithAutoStyle(true),
	      ana.WithHistoNorm(true),
)

// Produce plots
analyzer.Run()

```

## Gallery

<table>
  <tr>
    <td><p align="center"><img src="ana/testdata/Plots_simpleUseCase/Mttbar_golden.png">
	Data/Background <a href="https://github.com/rmadar/tree-gonalyzer/blob/master/ana/example_maker_test.go#L33" target="_blank">[code]</a></p>
    </td>
    <td><p align="center"><img src="ana/testdata/Plots_shapeComparison/DphiLL_golden.png">
	Shape comparison <a href="https://github.com/rmadar/tree-gonalyzer/blob/master/ana/example_maker_test.go#L112">[code]</a></p>
    </td>
  </tr>


  <tr>
    <td><p align="center"><img src="ana/testdata/Plots_XXX/Mttbar_golden.png">
	Shape distortion <a href="https://github.com/rmadar/tree-gonalyzer/blob/master/ana/XXX">[code]</a></p>
    </td>
    <td><p align="center"><img src="ana/testdata/Plots_XXX/DphiLL_golden.png">
	Systematic variation <a href="https://github.com/rmadar/tree-gonalyzer/blob/master/ana/XXX">[code]</a></p>
    </td>
  </tr>


 </table>

