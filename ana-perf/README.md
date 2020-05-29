# Performance study

![benchmarking](perf.png)

For 2M events and 60 variables, a comparison with similar ROOT-based code
(using `t->Draw()`) gives the following numbers:
 + `ROOT  -> 6.2 ms/kEvts`
 + `GOHEP -> 5.1 ms/kEvts`
 
Testing on only one variable to avoid event-loop repetition
in case of `t->Draw()` (even if I think it's quite optimized):
 + `ROOT  -> 0.39 ms/kEvts`
 + `GOHEP -> 0.27 ms/kEvts`
