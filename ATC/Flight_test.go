package main

import (
	"testing"
)

func TestFlightValidParams(t *testing.T) {
	f := Flight{
		Name: "Sample Flight",
	}
	err := f.Validate()
	if err != nil {
		t.Errorf("unexpected error with valid flight object:\n%w", err)
	}
}

func TestFlightMissingName(t *testing.T) {
	f := Flight{}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected error with missing name parameter")
	}
}

func TestFlightEmptyName(t *testing.T) {
	f := Flight{
		Name: "",
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected error with missing name parameter")
	}
}
