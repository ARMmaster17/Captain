package main

import (
	"fmt"
	db2 "github.com/ARMmaster17/Captain/ATC/DB"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestE2EMonitoringLoop(t *testing.T) {
	helperDeleteDBIfExists()
	db, err := db2.ConnectToDB()
	require.NoError(t, err)
	require.NotNil(t, db)
	err = initAirspaces(db)
	assert.NoError(t, err)
	err = monitoringLoop(db)
	assert.NoError(t, err)
}

func TestPlaneScaleStateServer(t *testing.T) {
	helperDeleteDBIfExists()
	assert.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	db, err := db2.ConnectToDB()
	require.NoError(t, err)
	require.NotNil(t, db)
	err = initAirspaces(db)
	require.NoError(t, err)
	flight := Flight{
		Name: "test",
		AirspaceID: 1,
	}
	tx := db.Save(&flight)
	require.NoError(t, tx.Error)
	formation := Formation{
		Name:        "plane",
		CPU:         1,
		RAM:         512,
		Disk:        8,
		BaseName:    "plane",
		Domain:      "example.com",
		TargetCount: 1,
		FlightID: int(flight.ID),
	}
	tx = db.Save(&formation)
	require.NoError(t, tx.Error)
	err = monitoringLoop(db)
	assert.NoError(t, err)
	result := db.Where("formation_id = ?", formation.ID).Preload("Formation").Find(&formation.Planes)
	assert.NoError(t, result.Error)
	assert.Equal(t, 1, len(formation.Planes))
}

func TestPlaneScaleDownStateServer(t *testing.T) {
	helperDeleteDBIfExists()
	assert.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	db, err := db2.ConnectToDB()
	require.NoError(t, err)
	require.NotNil(t, db)
	err = initAirspaces(db)
	require.NoError(t, err)
	flight := Flight{
		Name: "test",
		AirspaceID: 1,
	}
	tx := db.Save(&flight)
	require.NoError(t, tx.Error)
	formation := Formation{
		Name:        "plane",
		CPU:         1,
		RAM:         512,
		Disk:        8,
		BaseName:    "plane",
		Domain:      "example.com",
		TargetCount: 1,
		FlightID: int(flight.ID),
	}
	tx = db.Save(&formation)
	require.NoError(t, tx.Error)
	err = monitoringLoop(db)
	require.NoError(t, err)
	result := db.Where("formation_id = ?", formation.ID).Preload("Formation").Find(&formation.Planes)
	require.NoError(t, result.Error)
	require.Equal(t, 1, len(formation.Planes))
	formation.TargetCount = 0
	tx = db.Save(&formation)
	require.NoError(t, tx.Error)
	err = monitoringLoop(db)
	assert.NoError(t, err)
	result = db.Where("formation_id = ?", formation.ID).Preload("Formation").Find(&formation.Planes)
	assert.NoError(t, result.Error)
	assert.Equal(t, 0, len(formation.Planes))
}

func helperSetupConfigFile(configFile string) error {
	viper.Reset()
	_ = os.Remove("/etc/captain/config.yaml")
	input, err := ioutil.ReadFile(fmt.Sprintf("./testing/%s", configFile))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("/etc/captain/config.yaml", input, 0644)
	if err != nil {
		return err
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/captain")
	return viper.ReadInConfig()
}

func helperDeleteDBIfExists() {
	_ = os.Remove("testing.db")
}
