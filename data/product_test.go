package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "test",
		Price: 1.00,
		SKU:   "ABC123",
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
