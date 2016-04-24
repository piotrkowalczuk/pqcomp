package pqcomp_test

import (
	"testing"

	"github.com/piotrkowalczuk/pqcomp"
)

var (
	comp, where, update                               *pqcomp.Composer
	benchmarkKey, benchmarkExpr, benchmarkPlaceholder string
	benchmarkArgs                                     []interface{}
)

func BenchmarkComposer_few_tiny(b *testing.B)   { benchmarkComposer(b, 1, 1) }
func BenchmarkComposer_few_small(b *testing.B)  { benchmarkComposer(b, 5, 5) }
func BenchmarkComposer_few_medium(b *testing.B) { benchmarkComposer(b, 10, 10) }
func BenchmarkComposer_few_large(b *testing.B)  { benchmarkComposer(b, 100, 100) }
func BenchmarkComposer_few_huge(b *testing.B)   { benchmarkComposer(b, 1000, 1000) }

func BenchmarkComposer_medium_tiny(b *testing.B)   { benchmarkComposer(b, 1, 1, 1, 1) }
func BenchmarkComposer_medium_small(b *testing.B)  { benchmarkComposer(b, 5, 5, 5, 5) }
func BenchmarkComposer_medium_medium(b *testing.B) { benchmarkComposer(b, 10, 10, 10, 10) }
func BenchmarkComposer_medium_large(b *testing.B)  { benchmarkComposer(b, 100, 100, 100, 100) }
func BenchmarkComposer_medium_huge(b *testing.B)   { benchmarkComposer(b, 1000, 1000, 1000, 1000) }

func BenchmarkComposer_alot_tiny(b *testing.B)   { benchmarkComposer(b, 1, 1, 1, 1, 1, 1, 1, 1) }
func BenchmarkComposer_alot_small(b *testing.B)  { benchmarkComposer(b, 5, 5, 5, 5, 5, 5, 5, 5) }
func BenchmarkComposer_alot_medium(b *testing.B) { benchmarkComposer(b, 10, 10, 10, 10, 10, 10, 10, 10) }
func BenchmarkComposer_alot_large(b *testing.B) {
	benchmarkComposer(b, 100, 100, 100, 100, 100, 100, 100, 100)
}
func BenchmarkComposer_alot_huge(b *testing.B) {
	benchmarkComposer(b, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000)
}

func benchmarkComposer(b *testing.B, args, pexpr int, cexpr ...int) {
	comp := make([]*pqcomp.Composer, 0, b.N)
	key, expression, arg := "column", pqcomp.E, []byte("value")
	for n := 0; n < b.N; n++ {
		comp = append(comp, pqcomp.New(args, pexpr, cexpr...))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < args; i++ {
			comp[n].AddArg(arg)
		}

		for i := 0; i < pexpr; i++ {
			comp[n].AddExpr(key, expression, arg)
		}

		for _, expr := range cexpr {
			for i := 0; i < expr; i++ {
				comp[n].Compose().AddExpr(key, expression, arg)
			}
		}
	}
}

func BenchmarkComposer_real(b *testing.B) {
	type expression struct {
		column, ope string
		arg         interface{}
	}
	exprs := []expression{
		{"u.first_name", pqcomp.E, "johnsnow@gmail.com"},
		{"u.last_name", pqcomp.E, "Snow"},
		{"u.is_superuser", pqcomp.E, false},
		{"u.is_staff", pqcomp.E, false},
		{"u.details", pqcomp.E, "some information"},
		{"u.phone", pqcomp.E, "+48123123123"},
		{"u.address", pqcomp.E, "Kazimierza Wielkiego 1"},
		{"u.country", pqcomp.E, "Poland"},
		{"u.city", pqcomp.E, "Wrocław"},
		{"u.zipcode", pqcomp.E, "123456"},
	}
	args := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		comp = pqcomp.New(len(args), 0, len(exprs), len(exprs))
		where = comp.Compose()
		update = comp.Compose()

		for _, arg := range args {
			where.AddArg(arg)
		}

		for _, ex := range exprs {
			where.AddExpr(ex.column, ex.ope, ex.arg)
		}

		for _, ex := range exprs {
			update.AddExpr(ex.column, ex.ope, ex.arg)
		}

		for where.Next() {
			if where.First() {
			}

			benchmarkKey, benchmarkExpr, benchmarkPlaceholder = where.Key(), where.Oper(), where.PlaceHolder()
		}
		for update.Next() {
			if where.First() {
			}

			benchmarkKey, benchmarkExpr, benchmarkPlaceholder = where.Key(), where.Oper(), where.PlaceHolder()
		}
		benchmarkArgs = comp.Args()
	}
	b.Logf("root length: %d", comp.Len())
	b.Logf("where length: %d", where.Len())
	b.Logf("update length: %d", update.Len())
}

func BenchmarkComposer_AddExpr(b *testing.B) {
	column := "column"
	arg := "argument"
	comp := pqcomp.New(0, b.N)

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		comp.AddExpr(column, pqcomp.E, arg)
	}
}

func BenchmarkComposer_Placeholder(b *testing.B) {
	comp := pqcomp.New(0, b.N)

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		benchmarkPlaceholder = comp.PlaceHolder()
	}
}

func BenchmarkComposer_Args(b *testing.B) {
	comp := pqcomp.New(10, 10)

	comp.AddArg("arg1")
	comp.AddArg("arg2")
	comp.AddArg("arg3")
	comp.AddArg("arg4")
	comp.AddArg("arg5")
	comp.AddArg("arg6")
	comp.AddArg("arg7")
	comp.AddArg("arg8")
	comp.AddArg("arg9")
	comp.AddArg("arg10")
	comp.AddExpr("column1", pqcomp.E, "expr1")
	comp.AddExpr("column2", pqcomp.E, "expr2")
	comp.AddExpr("column3", pqcomp.E, "expr3")
	comp.AddExpr("column4", pqcomp.E, "expr4")
	comp.AddExpr("column5", pqcomp.E, "expr5")
	comp.AddExpr("column6", pqcomp.E, "expr6")
	comp.AddExpr("column7", pqcomp.E, "expr7")
	comp.AddExpr("column8", pqcomp.E, "expr8")
	comp.AddExpr("column9", pqcomp.E, "expr9")
	comp.AddExpr("column10", pqcomp.E, "expr10")

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		benchmarkArgs = comp.Args()
	}
}

func BenchmarkComposer_New(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		comp = pqcomp.New(100, 100)
	}
}
