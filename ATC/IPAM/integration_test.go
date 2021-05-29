package IPAM

import (
	"fmt"
	"github.com/ARMmaster17/Captain/ATC/DB"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"os"
	"sync"
	"testing"
)

func TestIPAM_ReserveAddress(t *testing.T) {
	ipam := helperInitTest(t)
	// Test reserving an IP.
	ip, err := ipam.GetNewAddress()
	assert.NoError(t, err)
	assert.NotNil(t, ip)
	assert.Equal(t, net.ParseIP("10.1.0.0"), ip)
	// Delete database.
	helperDeleteDBIfExists()
}

func TestIPAM_ReleaseAddress(t *testing.T) {
	ipam := helperInitTest(t)
	// Test reserving an IP.
	ip, err := ipam.GetNewAddress()
	require.NoError(t, err)
	require.NotNil(t, ip)
	require.Equal(t, net.ParseIP("10.1.0.0"), ip)
	err = ipam.ReleaseAddress(ip)
	assert.NoError(t, err)
	var addressCountInBlock int64
	result := ipam.db.Model(&ReservedAddress{}).Where("reserved_block_id = ?", 1).Count(&addressCountInBlock)
	require.NoError(t, result.Error)
	assert.Equal(t, int64(0), addressCountInBlock)
	// Delete database.
	helperDeleteDBIfExists()
}

func TestIPAM_ReserveTwoAddress(t *testing.T) {
	ipam := helperInitTest(t)
	// Test reserving an IP.
	ip, err := ipam.GetNewAddress()
	require.NoError(t, err)
	require.NotNil(t, ip)
	require.Equal(t, net.ParseIP("10.1.0.0"), ip)
	ip, err = ipam.GetNewAddress()
	assert.NoError(t, err)
	assert.NotNil(t, ip)
	assert.Equal(t, net.ParseIP("10.1.0.1"), ip)
	// Delete database.
	helperDeleteDBIfExists()
}

func TestIPAM_ReserveAddressWithRollover(t *testing.T) {
	ipam := helperInitTest(t)
	// Get to rollover range.
	for i := 0; i <= 255; i++ {
		ip, err := ipam.GetNewAddress()
		require.NoError(t, err)
		require.NotNil(t, ip)
		require.Equal(t, net.ParseIP(fmt.Sprintf("10.1.0.%d", i)), ip)
	}
	ip, err := ipam.GetNewAddress()
	assert.NoError(t, err)
	assert.NotNil(t, ip)
	assert.Equal(t, net.ParseIP("10.1.1.0"), ip)
	// Delete database.
	helperDeleteDBIfExists()
}

func helperInitTest(t *testing.T) IPAM {
	// Initialize DB connection.
	dbt, err := DB.ConnectToDB()
	if err != nil {
		require.NoError(t, err)
	}
	// Initialize IPAM service.
	ipam := IPAM{
		db:    dbt,
		mutex: &sync.Mutex{},
	}
	err = ipam.Initialize(dbt)
	require.NoError(t, err)
	err = ipam.syncDBBlocksWithConfig([]string{"10.1.0.0/16"})
	require.NoError(t, err)

	return ipam
}

func helperDeleteDBIfExists() {
	_ = os.Remove("./testing.db")
}
