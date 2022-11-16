package records

import (
	"reflect"
	"testing"

	"github.com/michaelrknutson/loanprogo/models"
)

type TestRecord struct {
	record models.Record
}

func TestCheckBalance(t *testing.T) {
	r := &TestRecord{}
	r2 := &TestRecord{}
	r.record = models.Record{
		UserBalance: 275.00,
	}
	r2.record = models.Record{
		UserBalance: 275.00,
	}

	r2.record.ID = 1

	expectr := 300.00
	expectr2 := 275.00
	ansr := checkBalance(r.record)
	ansr2 := checkBalance(r2.record)

	if expectr != ansr {
		t.Log("expected 300")
		t.FailNow()
	}

	if expectr2 != ansr2 {
		t.Log("expected 300")
		t.FailNow()
	}
}
func TestCheckOperation(t *testing.T) {
	r := &TestRecord{}
	r.record = models.Record{
		UserBalance: 275.00,
		Operation: models.Operation{
			Cost: 1.25,
			Type: "Addition",
		},
	}

	expectr := "3"

	res, _ := checkOperation(r.record.Operation, r.record.UserBalance, 1, 2)

	if expectr != res {
		t.Log("expected 3")
		t.FailNow()
	}

}

func TestGenerateStrings(t *testing.T) {
	response := generateStrings()

	if reflect.TypeOf(response).Kind() != reflect.String {
		t.Log("no string generated")
		t.FailNow()
	}
}

func TestPerformOperation(t *testing.T) {
	want := "300"

	add := performOperation("Addition", 200, 100)
	sub := performOperation("Subtraction", 301, 1)
	mult := performOperation("Multiplication", 150, 2)
	div := performOperation("Division", 600, 2)
	sqrt := performOperation("Square_root", 4, 0)

	sqans := "2"

	if add != want {
		t.Log("add wanted 300")
		t.FailNow()
	}

	if sub != want {
		t.Log("sub wanted 300")
		t.FailNow()
	}

	if mult != want {
		t.Log("mult wanted 300")
		t.FailNow()
	}

	if div != want {
		t.Log("div wanted 300")
		t.FailNow()
	}

	if sqrt != sqans {
		t.Log("sqrt wanted 2")
		t.FailNow()
	}
}
