package main

import (
	"testing"
)

func TestPlaneGetFQDN(t *testing.T) {
	p := Plane{
		Num: 1,
		Formation: Formation{
			Name:        "test formation",
			BaseName:    "tf",
			Domain:      "testing.com",
		},
	}
	testFQDN := p.getFQDN()
	if testFQDN != "tf1.testing.com" {
		t.Errorf("expected FQDN tf1.testing.com, got %s", testFQDN)
	}
}

func TestPlaneIsHealthy(t *testing.T) {
	p := Plane{}
	healthCheck, _ := p.isHealthy(nil)
	if !healthCheck {
		t.Errorf("expected plane to be healthy, got FALSE")
	}
}

func TestPlaneIsNotHealthy(t *testing.T) {
	t.Skip()
}

func TestPlaneNewInvalidNum(t *testing.T) {
	p, err := NewPlane(0)
	if p != nil {
		t.Errorf("expected plane to be nil with invalid Num parameter")
	}
	if err == nil {
		t.Errorf("expected error to be thrown with invalid Num parameter for plane")
	}
}

func TestPlaneNewValidNum(t *testing.T) {
	p, err := NewPlane(1)
	if p == nil {
		t.Errorf("expected plane to be created with valid Num parameter")
	}
	if err != nil {
		t.Errorf("expected no error to be thrown with valid Num parameter")
	}
}
