/**
Package pqcomp provides dead simple query builder that support null types from sql package, but also provide interface Appearer.


Benchmarks

To be comprehensive solution, query builder needs to be optimized. Some of the benchmark results:

	BenchmarkComposer_AddExpr-4      	 5000000	       470 ns/op	     100 B/op	       1 allocs/op
	BenchmarkComposer_Placeholder-4  	10000000	       130 ns/op	       4 B/op	       2 allocs/op
	BenchmarkComposer_Args-4         	300000000	         4.89 ns/op	       0 B/op	       0 allocs/op
	BenchmarkComposer_New-4          	  200000	     15311 ns/op	    5120 B/op	       4 allocs/op

*/
package pqcomp
