# Performance study

| N[vars] | variable | cut & weights | time (ms/kEvts) |
|:-------:|:--------:|:-------------:|:---------------:|
|    60   |          |               |                 |
|         |  direct  |     direct    |                 |
|         |  formula |     direct    |                 |
|         |  direct  |    formula    |                 |
|         |  formula |    formula    |                 |

1. Use explicit loading for variables and `rtree.Formula` for weight and cuts:
```
go run ./main.go
```

2. Use explicit loading for variables but disable `rtree.Formula` of weights and cuts (not applied):
```
go run ./main.go -noCutWeight
```

3. Use `rtree.Formula` for variables and disable `rtree.Formula` of weights and cuts:
```
go run ./main.go -varFormula -noCutWeight
```

4. Use `rtree.Formula` for variables and `rtree.Formula` for weight and cuts:
```
go run ./main.go -varFormula
```