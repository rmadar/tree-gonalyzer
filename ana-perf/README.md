# Performance study

Use explicit loading for variables and `rtree.Formula` for weight and cuts:
```bash
go run ./main.go
```

Use `rtree.Formula` for variables and `rtree.Formula` for weights and cuts:
```bash
go run ./main.go -varFormula
```

Disable weights and cuts (and associated `rtree.Formula`:
```bash
go run ./main.go -noCutWeight
```