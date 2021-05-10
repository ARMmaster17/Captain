package main

import (
	"testing"
)

func TestFormationValidParams(t *testing.T) {
	f := Formation{
		Name:        "SQL Cluster 1",
		CPU:         1,
		RAM:         256,
		Disk:		 8,
		BaseName:    "sql",
		Domain:      "example.com",
		TargetCount: 1,
	}
	err := f.Validate()
	if err != nil {
		t.Errorf("unexpected error in valid formation params: %w", err)
	}
}

func TestFormationInvalidName(t *testing.T) {
	f := Formation{
		Name:		 "",
		CPU:         1,
		RAM:         256,
		Disk:		 8,
		BaseName:    "sql",
		Domain:      "example.com",
		TargetCount: 1,
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected an error with invalid Name")
	}
}

func TestFormationMissingName(t *testing.T) {
	f := Formation{
		CPU:         1,
		RAM:         256,
		Disk:		 8,
		BaseName:    "sql",
		Domain:      "example.com",
		TargetCount: 1,
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected an error with invalid Name")
	}
}

func TestFormationInvalidCPU(t *testing.T) {
	f := Formation{
		Name:		"SQL Cluster 1",
		CPU:         0,
		RAM:         256,
		Disk:		 8,
		BaseName:    "sql",
		Domain:      "example.com",
		TargetCount: 1,
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected an error with invalid CPU")
	}
}

func TestFormationInvalidMinRAM(t *testing.T) {
	f := Formation{
		Name:		"SQL Cluster 1",
		CPU:         0,
		RAM:         0,
		Disk:		 8,
		BaseName:    "sql",
		Domain:      "example.com",
		TargetCount: 1,
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected an error with invalid RAM")
	}
}

func TestFormationInvalidMaxRAM(t *testing.T) {
	f := Formation{
		Name:		"SQL Cluster 1",
		CPU:         0,
		RAM:         307201,
		Disk:		 8,
		BaseName:    "sql",
		Domain:      "example.com",
		TargetCount: 1,
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected an error with invalid RAM")
	}
}

func TestFormationMissingBaseName(t *testing.T) {
	f := Formation{
		Name:		"SQL Cluster 1",
		CPU:         0,
		RAM:         307201,
		Disk:		 8,
		BaseName:    "",
		Domain:      "example.com",
		TargetCount: 1,
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected an error with invalid BaseName")
	}
}

func TestFormationBadBaseName(t *testing.T) {
	f := Formation{
		Name:		"SQL Cluster 1",
		CPU:         0,
		RAM:         307201,
		Disk:		 8,
		BaseName:    "!@#$%^&*()",
		Domain:      "example.com",
		TargetCount: 1,
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected an error with invalid BaseName")
	}
}

func TestFormationMissingDomain(t *testing.T) {
	f := Formation{
		Name:		"SQL Cluster 1",
		CPU:         0,
		RAM:         307201,
		Disk:		 8,
		BaseName:    "sql",
		TargetCount: 1,
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected an error with invalid domain")
	}
}

func TestFormationEmptyDomain(t *testing.T) {
	f := Formation{
		Name:		"SQL Cluster 1",
		CPU:         0,
		RAM:         307201,
		Disk:		 8,
		BaseName:    "sql",
		Domain:      "",
		TargetCount: 1,
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected an error with invalid domain")
	}
}

func TestFormationBadDomain(t *testing.T) {
	f := Formation{
		Name:		"SQL Cluster 1",
		CPU:         0,
		RAM:         307201,
		Disk:		 8,
		BaseName:    "sql",
		Domain:      "!@#$%^&*()",
		TargetCount: 1,
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected an error with invalid domain")
	}
}

func TestFormationMissingTargetCount(t *testing.T) {
	f := Formation{
		Name:		"SQL Cluster 1",
		CPU:         0,
		RAM:         307201,
		Disk:		 8,
		BaseName:    "sql",
		Domain:      "example.com",
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected an error with invalid TargetCount")
	}
}

func TestFormationInvalidTargetCount(t *testing.T) {
	f := Formation{
		Name:		"SQL Cluster 1",
		CPU:         0,
		RAM:         307201,
		Disk:		 8,
		BaseName:    "sql",
		Domain:      "example.com",
		TargetCount: -1,
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected an error with invalid TargetCount")
	}
}

func TestFormationMissingDisk(t *testing.T) {
	f := Formation{
		Name:		"SQL Cluster 1",
		CPU:         0,
		RAM:         307201,
		BaseName:    "sql",
		Domain:      "example.com",
		TargetCount: 1,
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected an error with missing disk parameter")
	}
}

func TestFormationBadDiskMin(t *testing.T) {
	f := Formation{
		Name:		"SQL Cluster 1",
		CPU:         0,
		RAM:         307201,
		Disk:		 0,
		BaseName:    "sql",
		Domain:      "example.com",
		TargetCount: 1,
	}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected an error with invalid disk parameter")
	}
}