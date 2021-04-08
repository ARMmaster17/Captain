package main

import (
	"testing"
)

func TestGetFQDN(t *testing.T) {
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
		t.Errorf("Expected FQDN tf1.testing.com, got %s", testFQDN)
	}
}
