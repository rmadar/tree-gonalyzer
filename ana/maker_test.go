package ana_test

import (
	"testing"
	
	"gonum.org/v1/plot/cmpimg"
)

func TestSimpleUseCase(t *testing.T) {
	cmpimg.CheckPlot(ExampleMaker_aSimpleUseCase, t,
		"Plots_simpleUseCase/No-cut/Mttbar.png",
		"Plots_simpleUseCase/No-cut/DphiLL.png",
	)
}
