package goplaces

import (
	"errors"
	"testing"
)

func TestQuery(t *testing.T) {
	// Check query
	r, err := Query(Parameters{
		Query:     "Brickell Avenue, Miami, Florida",
		Countries: "us",
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	// Check length
	if len(r.Hits) == 0 {
		t.Error(errors.New("No results from query"))
		t.FailNow()
	}
}

func TestLabelFromAddress(t *testing.T) {
	addr := Address{
		Postcode: "33131",
		State:    "Florida",
		City:     "Miami",
		Street:   "Brickell Avenue",
	}
	lb := NewLabelFromAddress(addr)
	if lb != "Brickell Avenue, Miami, Florida, 33131" {
		t.Error(errors.New("Incorrect label"))
	}
}
