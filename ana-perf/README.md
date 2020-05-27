# Performance study

![benchmarking](perf.png)


## How to run each setup?

**1.** Use explicit loading for variables and `rtree.Formula` for weight:
```go
t := runTest(n10kEvtsPerSample, nVars, false, false)
```

**2.** Use explicit loading for variables but disable `rtree.Formula` of weights (*i.e.* not applied):
```go
t := runTest(n10kEvtsPerSample, nVars, false, true)
```

**3.** Use `rtree.Formula` for variables and disable `rtree.Formula` of weights (*i.e.* not applied):
```go
t := runTest(n10kEvtsPerSample, nVars, true, true)
```

**4.** Use `rtree.Formula` for variables and `rtree.Formula` for weights:
```go
t := runTest(n10kEvtsPerSample, nVars, true, false)
```