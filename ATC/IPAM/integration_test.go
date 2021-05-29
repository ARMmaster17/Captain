package IPAM

import (
	"github.com/ARMmaster17/Captain/ATC/DB"
	"testing"
)

func TestIPAM_ReserveAddress(t *testing.T) {
	dbt, err := DB.ConnectToDB()
	if err != nil {
		t.Errorf("unable to connect to db:\n%w", err)
		return
	}
	
}
