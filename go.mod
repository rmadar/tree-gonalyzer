module github.com/rmadar/tree-gonalyzer

go 1.14

require (
	github.com/rmadar/hplot-style v0.0.0-20200411142456-95ac5bcd1735
	go-hep.org/x/hep v0.26.1-0.20200416105033-399c28b79b69
	gonum.org/v1/gonum v0.7.1-0.20200330111830-e98ce15ff236
	gonum.org/v1/plot v0.7.1-0.20200414075901-f4e1939a9e7a
)

replace (
     //github.com/rmadar/hplot-style => /home/rmadar/cernbox/goDev/hplot-style
     go-hep.org/x/hep => /home/rmadar/cernbox/goDev/test_gohep/hep
)
