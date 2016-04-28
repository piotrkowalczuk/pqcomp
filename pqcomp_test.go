package pqcomp_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/piotrkowalczuk/pqcomp"
)

func TestNew(t *testing.T) {
	comp := pqcomp.New(1, 1, 0)
	if comp == nil {
		t.Errorf("composer should not be nil")
	}
}

func TestComposer_AddArg(t *testing.T) {
	max := 100

	parent := pqcomp.New(max, max, 3)
	where := parent.Compose()
	update := parent.Compose()

	for i := 0; i < max; i++ {
		switch i % 3 {
		case 0:
			parent.AddArg(i)
		case 1:
			where.AddArg(i)
		case 2:
			update.AddArg(i)
		}
	}

	if len(parent.Args()) != max {
		t.Errorf("unexpected number of arguments, expect %d but got %d", max, len(parent.Args()))
	}
}

func TestComposer_AddExpr(t *testing.T) {
	max := 99
	expected := max / 3

	parent := pqcomp.New(0, 0)
	where := parent.Compose()
	update := parent.Compose()

	for i := 0; i < max; i++ {
		switch i % 3 {
		case 0:
			parent.AddExpr(fmt.Sprintf("column_parent_%d", i), pqcomp.Equal, i)
		case 1:
			where.AddExpr(fmt.Sprintf("column_where_%d", i), pqcomp.Equal, i)
		case 2:
			update.AddExpr(fmt.Sprintf("column_update_%d", i), pqcomp.Equal, i)
		}
	}

	if parent.Len() != expected {
		t.Errorf("parent expression length mismatch, expected %d but got %d", expected, parent.Len())
	}

	if where.Len() != expected {
		t.Errorf("where expression length mismatch, expected %d but got %d", expected, where.Len())
	}

	if update.Len() != expected {
		t.Errorf("where expression length mismatch, expected %d but got %d", expected, update.Len())
	}
}

func TestComposer_AddExpr_types(t *testing.T) {
	cases := map[string]interface{}{
		"string":          "text",
		"int":             1,
		"int8":            int8(2),
		"int16":           int16(3),
		"int32":           int32(4),
		"int64":           int64(5),
		"uint":            uint(6),
		"uint8":           uint8(7),
		"uint16":          uint16(8),
		"uint32":          uint32(9),
		"uint64":          uint64(10),
		"float32":         float32(11.11),
		"float64":         float64(12.12),
		"sql.NullString":  sql.NullString{String: "null-text", Valid: true},
		"sql.NullInt64":   sql.NullInt64{Int64: 13, Valid: true},
		"sql.NullFloat64": sql.NullFloat64{Float64: 14.14, Valid: true},
		"sql.NullBool":    sql.NullBool{Bool: true, Valid: true},
		"[]byte":          []byte("abc"),
		"byte":            []byte("a")[0],
		"rune":            rune("r"[0]),
	}

CasesLoop:
	for k, v := range cases {
		comp := pqcomp.New(0, 0)
		comp.AddExpr(k, pqcomp.Equal, v)
		args := comp.Args()

		if len(args) != 1 {
			t.Errorf("wrong args slice length, expected %d but got %d", 0, len(args))
			continue CasesLoop
		}
		if !reflect.DeepEqual(v, args[0]) {
			t.Errorf("wrong argument, expected %d but got %d", v, args[0])
			continue CasesLoop
		}
		t.Logf("%s: pass for type: %T and value: %v", k, v, v)
	}
}

func TestComposer_AddArg_slices(t *testing.T) {
	comp := pqcomp.New(0, 0)
	comp.AddExpr("[]string", pqcomp.IN, []string{"1", "2", "3"})
	comp.AddExpr("[]int64", pqcomp.IN, []int64{1, 2, 3})

	if len(comp.Args()) != 6 {
		t.Fatalf("wrong number of arguments, expected %d but got %d", 6, len(comp.Args()))
	}
}

func TestComposer_AddExpr_sql(t *testing.T) {
	comp := pqcomp.New(0, 0)
	comp.AddExpr("int64-valid", pqcomp.Equal, &sql.NullInt64{Int64: 1, Valid: true})
	comp.AddExpr("int64-invalid", pqcomp.Equal, &sql.NullInt64{Int64: 2, Valid: false})
	comp.AddExpr("string-valid", pqcomp.Equal, &sql.NullString{String: "3", Valid: true})
	comp.AddExpr("string-invalid", pqcomp.Equal, &sql.NullString{String: "4", Valid: false})
	comp.AddExpr("float64-valid", pqcomp.Equal, &sql.NullFloat64{Float64: 5, Valid: true})
	comp.AddExpr("float64-invalid", pqcomp.Equal, &sql.NullFloat64{Float64: 6, Valid: false})
	comp.AddExpr("bool-valid", pqcomp.Equal, &sql.NullBool{Bool: true, Valid: true})
	comp.AddExpr("bool-invalid", pqcomp.Equal, &sql.NullBool{Bool: true, Valid: false})

	if len(comp.Args()) != 4 {
		t.Fatalf("wrong number of arguments, expected %d but got %d", 4, len(comp.Args()))
	}
}

func TestComposer_AddExpr_time(t *testing.T) {
	var tt *time.Time
	now := time.Now()

	comp := pqcomp.New(0, 0)
	comp.AddExpr("time", pqcomp.Equal, now)
	comp.AddExpr("time-zero", pqcomp.Equal, time.Time{})
	comp.AddExpr("time-pointer", pqcomp.Equal, &now)
	comp.AddExpr("time-pointer-zero", pqcomp.Equal, &time.Time{})
	comp.AddExpr("time-pointer-zero", pqcomp.Equal, tt)

	if len(comp.Args()) != 2 {
		t.Fatalf("wrong number of arguments, expected %d but got %d", 2, len(comp.Args()))
	}

}

func TestComposer_AddExpr_nil(t *testing.T) {
	comp := pqcomp.New(0, 0)
	func(comp *pqcomp.Composer, s *sql.NullString, i *sql.NullInt64) {
		comp.AddExpr("v1", pqcomp.Equal, nil)
		comp.AddExpr("v2", pqcomp.Equal, s)
		comp.AddExpr("v3", pqcomp.Equal, i)
	}(comp, nil, nil)

	if comp.Len() != 0 {
		t.Errorf("length mismatch, expected 0 but got %d: %#v", comp.Len(), comp.Args())
	}
}

func TestComposer_Len(t *testing.T) {
	comp := pqcomp.New(1, 1, 2)

	compA := comp.Compose()
	compA.AddExpr("column", pqcomp.Equal, "value")
	compA.AddExpr("column", pqcomp.Equal, "value")

	compB := comp.Compose()
	compB.AddExpr("column", pqcomp.Equal, "value")
	compB.AddExpr("column", pqcomp.Equal, "value")
	compB.AddExpr("column", pqcomp.Equal, "value")
	compB.AddExpr("column", pqcomp.Equal, "value")

	if comp.Len() != 0 {
		t.Errorf("wrong parent composer length, got %d but expected %d", comp.Len(), 0)
	}

	if compA.Len() != 2 {
		t.Errorf("wrong composer A length, got %d but expected %d", compA.Len(), 2)
	}

	if compB.Len() != 4 {
		t.Errorf("wrong composer B length, got %d but expected %d", compB.Len(), 4)
	}
}

func TestComposer_PlaceHolder(t *testing.T) {
	_, compA, compB := prepareComposers(10, 20)

	j := 0
	for compA.Next() {
		expected := fmt.Sprintf("$%d", j+1)
		if compA.PlaceHolder() != expected {
			t.Errorf("wrong placeholder for composer A, got %s but expected %s", compA.PlaceHolder(), expected)
		}
		j++
	}
	for compB.Next() {
		expected := fmt.Sprintf("$%d", j+1)
		if compB.PlaceHolder() != expected {
			t.Errorf("wrong placeholder for composer B, got %s but expected %s", compB.PlaceHolder(), expected)
		}
		j++
	}
}

func TestComposer_Key(t *testing.T) {
	lengthA, lengthB := 10, 20
	_, compA, compB := prepareComposers(lengthA, lengthB)

	j := 0
	for compA.Next() {
		expected := fmt.Sprintf("column_%d", j)
		if compA.Key() != expected {
			t.Errorf("wrong key for composer A, expected %s but got %s", expected, compA.Key())
		}
		j++
	}
	for compB.Next() {
		expected := fmt.Sprintf("column_%d", j-lengthA)
		if compB.Key() != expected {
			t.Errorf("wrong key for composer B, expected %s but got %s", expected, compB.Key())
		}
		j++
	}
}

func TestComposer_Expr(t *testing.T) {
	lengthA, lengthB := 10, 20
	_, compA, compB := prepareComposers(lengthA, lengthB)

	j := 0
	expected := pqcomp.E
	for compA.Next() {
		if compA.Oper() != expected {
			t.Errorf("wrong expression for composer A, expected %s but got %s", expected, compA.Oper())
		}
		j++
	}
	for compB.Next() {
		if compB.Oper() != expected {
			t.Errorf("wrong expression for composer B, expected %s but got %s", expected, compB.Oper())
		}
		j++
	}
}

func TestComposer_ExprOptional(t *testing.T) {
	var comp *pqcomp.Composer

	success := []interface{}{
		newAppearer("string"),
		&sql.NullString{String: "string", Valid: true},
		&sql.NullInt64{Int64: 1, Valid: true},
		&sql.NullFloat64{Float64: 1.0, Valid: true},
		&sql.NullBool{Bool: true, Valid: true},
		"something",
		1,
		1.0,
		false,
	}
	comp = pqcomp.New(0, len(success))

	for _, s := range success {
		comp.AddExpr("column", pqcomp.Equal, s)
	}

	if comp.Len() != len(success) {
		t.Errorf("unexpected ammount of expressions, expected %d but got %d", len(success), comp.Len())
	}

	failure := []interface{}{
		newAppearer(""),
		&sql.NullString{},
		&sql.NullInt64{},
		&sql.NullFloat64{},
		&sql.NullBool{},
	}
	comp = pqcomp.New(0, 0)

	for _, f := range failure {
		comp.AddExpr("column", pqcomp.Equal, f)
	}

	if comp.Len() != 0 {
		t.Errorf("unexpected ammount of expressions, expected %d but got %d", 0, comp.Len())
	}
}

func TestComposer_First(t *testing.T) {
	lengthA, lengthB := 10, 20
	_, compA, compB := prepareComposers(lengthA, lengthB)

	firstA, firstB := false, false
	for compA.Next() {
		if firstA && compA.First() {
			t.Errorf("first iteration already took place for composer A")
		}
		if compA.First() {
			firstA = true
		}
	}
	for compB.Next() {
		if firstB && compB.First() {
			t.Errorf("first iteration already took place for composer B")
		}
		if compB.First() {
			firstB = true
		}
	}
}

func prepareComposers(lengthA, lengthB int) (comp, compA, compB *pqcomp.Composer) {
	comp = pqcomp.New(0, 0, lengthA, lengthB)

	compA = comp.Compose()
	for i := 0; i < lengthA; i++ {
		compA.AddExpr(fmt.Sprintf("column_%d", i), pqcomp.Equal, "value")
	}

	compB = comp.Compose()
	for i := 0; i < lengthB; i++ {
		compB.AddExpr(fmt.Sprintf("column_%d", i), pqcomp.Equal, "value")
	}

	return
}

type appearer string

func newAppearer(s string) *appearer {
	a := appearer(s)
	return &a
}

// Appear implements Appearer interface.
func (a *appearer) Appear() bool {
	if a == nil {
		return false
	}
	return (*a) != ""
}
