package captain

import (
	"testing"
)

func Test_NewPlane(t *testing.T) {
	plane_helper(t, "test", 1, 512, 8, false)
}

func plane_helper(t *testing.T, name string, cpu int, ram int, storage int, shouldError bool) {
	plane, err := NewPlane(name, cpu, ram, storage)
	if err != nil && !shouldError {
		t.Logf("expected err to be nil with input {%s,%d,%d,%d}, got %s", name, cpu, ram, storage, err)
		t.Fail()
	}
	if err != nil {
		return
	}
	if plane.Name != name {
		t.Logf("expected Name to be %s with input {%s,%d,%d,%d}, got %s", name, name, cpu, ram, storage, plane.Name)
		t.Fail()
	}
	if plane.CPU != cpu {
		t.Logf("expected CPU to be %d with input {%s,%d,%d,%d}, got %d", cpu, name, cpu, ram, storage, plane.CPU)
		t.Fail()
	}
	if plane.RAM != ram {
		t.Logf("expected RAM to be %d with input {%s,%d,%d,%d}, got %d", ram, name, cpu, ram, storage, plane.RAM)
		t.Fail()
	}
	if plane.Storage != storage {
		t.Logf("expected CPU to be %d with input {%s,%d,%d,%d}, got %d", cpu, name, cpu, ram, storage, plane.CPU)
		t.Fail()
	}
}