module github.com/rmadar/tree-gonalyzer

go 1.14

require (
	github.com/rmadar/hplot-style v0.0.0-20200420135036-6c906c218e02
	go-hep.org/x/hep v0.26.1-0.20200421090732-5a84c0078ea0
	gonum.org/v1/gonum v0.7.1-0.20200330111830-e98ce15ff236
	gonum.org/v1/plot v0.7.1-0.20200414075901-f4e1939a9e7a
)

// replace github.com/rmadar/hplot-style => /home/rmadar/cernbox/goDev/hplot-style
replace go-hep.org/x/hep => /home/rmadar/cernbox/goDev/gohep_dev/hep
