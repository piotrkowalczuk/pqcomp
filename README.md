# pqcomp [![GoDoc](https://godoc.org/github.com/piotrkowalczuk/pqcomp?status.svg)](http://godoc.org/github.com/piotrkowalczuk/pqcomp) [![Build Status](https://travis-ci.org/piotrkowalczuk/pqcomp.svg)](https://travis-ci.org/piotrkowalczuk/pqcomp)&nbsp;[![codecov.io](https://codecov.io/github/piotrkowalczuk/pqcomp/coverage.svg?branch=master)](https://codecov.io/github/piotrkowalczuk/pqcomp?branch=master)

Dead simple query builder that support null types from [sql](https://golang.org/pkg/database/sql/) package, but also provide interface [Appearer](https://godoc.org/github.com/piotrkowalczuk/pqcomp#Appearer). 
Detailed information and [examples](https://godoc.org/github.com/piotrkowalczuk/pqcomp#example-package) can be found in the [documentation](https://godoc.org/github.com/piotrkowalczuk/pqcomp#Composer). 



## Benchmarks

```
BenchmarkComposer_AddExpr-4      	 5000000	       470 ns/op	     100 B/op	       1 allocs/op
BenchmarkComposer_Placeholder-4  	10000000	       130 ns/op	       4 B/op	       2 allocs/op
BenchmarkComposer_Args-4         	300000000	         4.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkComposer_New-4          	  200000	     15311 ns/op	    5120 B/op	       4 allocs/op
```