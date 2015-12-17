package pqcomp_test

import (
	"testing"

	"github.com/piotrkowalczuk/pqcomp"
)

var (
	benchmarkKey, benchmarkExpr, benchmarkPlaceholder string
	benchmarkAll                                      []interface{}
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
	b.ReportAllocs()
	comp := make([]*pqcomp.Composer, 0, b.N)

	key, expression, arg := "column", pqcomp.E, []byte("value")
	for n := 0; n < b.N; n++ {
		comp = append(comp, pqcomp.New(args, pexpr, cexpr...))
	}
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
	b.ReportAllocs()
	exprs := []struct {
		column, ope string
		arg         interface{}
	}{
		{"u.first_name", pqcomp.E, "johnsnow@gmail.com"},
		{"u.last_name", pqcomp.E, "Snow"},
		{"u.is_superuser", pqcomp.E, false},
		{"u.is_staff", pqcomp.E, false},
		{"u.details", pqcomp.E, "some information"},
		{"u.phone", pqcomp.E, "+48123123123"},
		{"u.address", pqcomp.E, "Strasse 1"},
		{"u.city", pqcomp.E, "Berlin"},
		{"u.zipcode", pqcomp.E, "123456"},
	}
	args := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		comp := pqcomp.New(10, 0, len(exprs), len(exprs))
		where := comp.Compose()
		update := comp.Compose()
		b.StartTimer()

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
		benchmarkAll = comp.Args()
	}
}
