# Performance study

|   setup   | variable | cut & weights | 60 vars | 40 vars | 20 vars | 1 vars |
|:---------:|:--------:|:-------------:|:-------:|:-------:|:-------:|:------:|
|   **1**   |  direct  |     direct    |         |         |         |        |
|   **2**   |  formula |     direct    |         |         |         |        |
|   **3**   |  direct  |    formula    |         |         |         |        |
|   **4**   |  formula |    formula    |         |         |         |        |

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