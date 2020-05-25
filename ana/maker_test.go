package ana_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestSimpleUseCase(t *testing.T) {
	cmpimg.CheckPlot(ExampleMaker_aSimpleUseCase, t,
		"Plots_simpleUseCase/Mttbar.png",
		"Plots_simpleUseCase/DphiLL.png",
	)
}

func TestMultiComponentSamples(t *testing.T) {
	cmpimg.CheckPlot(ExampleMaker_multiComponentSamples, t,
		"Plots_multiComponents/Mttbar.png",
		"Plots_multiComponents/DphiLL.png",
	)
}
