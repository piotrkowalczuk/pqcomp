package pqcomp_test

import (
	"fmt"

	"github.com/piotrkowalczuk/pqcomp"
)

func Example() {
	var (
		uquery, wquery string
	)
	comp := pqcomp.New(1, 1, 1, 3)
	update := comp.Compose()
	where := comp.Compose()

	comp.AddArg(10)

	update.AddExpr("u.username", pqcomp.E, "johnsnow")
	update.AddExpr("u.first_name", pqcomp.E, "John")
	update.AddExpr("u.last_name", pqcomp.E, "Snow")

	where.AddExpr("u.id", pqcomp.E, 1)

	if update.Len() == 0 || where.Len() == 0 {
		return
	}

	for update.Next() {
		if update.First() {
			uquery += "SET "
		} else {
			uquery += ", "
		}
		uquery += fmt.Sprintf("%s %s %s", update.Key(), update.Oper(), update.PlaceHolder())
	}
	for where.Next() {
		if where.First() {
			wquery += "WHERE "
		} else {
			wquery += ", "
		}
		wquery += fmt.Sprintf("%s %s %s ", where.Key(), where.Oper(), where.PlaceHolder())
	}

	fmt.Println(where.Args())
	fmt.Println(update.Args())
	fmt.Println(comp.Args())
	fmt.Printf("UPDATE users AS u %s %sLIMIT $1 \n", uquery, wquery)

	// Output:
	// [1]
	// [johnsnow John Snow]
	// [10 johnsnow John Snow 1]
	// UPDATE users AS u SET u.username = $2, u.first_name = $3, u.last_name = $4 WHERE u.id = $5 LIMIT $1
}