// Package gonalyzer exposes types and functions to ease the analysis of ROOT trees.
//
//  // Define samples
//  samples := []*ana.Sample{
//    ana.CreateSample("data", "data", `Data`, "data.root", "mytree"),
//    ana.CreateSample("bkg1", "bkg", `Proc 1`, "proc1.root", "mytree"),
//    ana.CreateSample("bkg2", "bkg", `Proc 2`, "proc2.root", "mytree"),
//  }
//
//  // Define variables
//  variables := []*ana.Variable{
//    ana.NewVariable("plot1", ana.TreeVarBool("branchBool"), 2, 0, 2),
//    ana.NewVariable("plot2", ana.TreeVarF32("branchF32"), 25, 0, 1000),
//    ana.NewVariable("plot3", ana.TreeVarF64("branchF64"), 50, 0, 1000),
//  }
//
//  // Create analyzer object with some options
//  analyzer := ana.New(samples, variables)
//
//  // Produce plots
//  analyzer.Run()
//
// There is also the possibility of doing cutflows in a simple way. Few lines of code
// can produce this ASCII table:
//
//  | Cut name              | Raw Yields                 | Weighted Yields            |
//  |                       |                 Abs    Rel |                 Abs    Rel |
//  |-----------------------|----------------------------|----------------------------|
//  | Electron channel      |         5526   100%   100% |     28230.59   100%   100% |
//  | pT > 10 GeV           |         5281    96%    96% |     28065.97    99%    99% |
//  | Phi < 2.0 rad         |         4312    78%    82% |     22874.73    81%    82% |
//
package gonalyzer 
