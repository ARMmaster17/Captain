package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	db2 "github.com/ARMmaster17/Captain/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRESTAirspaceRAll(t *testing.T) {
	HelperDeleteDBIfExistsForAPI()
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	db := HelperAPIInitDB()
	require.NotNil(t, db)
	req, err := http.NewRequest("GET", "/airspaces", nil)
	require.NoError(t, err)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusOK, response)

	var m []map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &m)
	require.NoError(t, err)
	assert.Equal(t, len(m), 1)
	firstEntry := m[0]
	assert.Equal(t, firstEntry["ID"], float64(1))
	assert.Equal(t, firstEntry["HumanName"], "Default Airspace")
	assert.Equal(t, firstEntry["NetName"], "default")
}

func TestRESTAirspaceR(t *testing.T) {
	HelperDeleteDBIfExistsForAPI()
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	db := HelperAPIInitDB()
	req, err := http.NewRequest("GET", "/airspace/1", nil)
	require.NoError(t, err)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusOK, response)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &m))
	assert.Equal(t, m["ID"], float64(1))
	assert.Equal(t, m["HumanName"], "Default Airspace")
	assert.Equal(t, m["NetName"], "default")
}

func TestRESTAirspaceC(t *testing.T) {
	HelperDeleteDBIfExistsForAPI()
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	db := HelperAPIInitDB()
	payload := map[string]string{}
	payload["HumanName"] = "test2a"
	payload["NetName"] = "test2b"
	data, err := json.Marshal(payload)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/airspace", bytes.NewReader(data))
	require.NoError(t, err)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusCreated, response)
	var fmtresponse map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &fmtresponse))
	assert.Equal(t, float64(2), fmtresponse["ID"])
	assert.Equal(t, "test2a", fmtresponse["HumanName"])
	assert.Equal(t, "test2b", fmtresponse["NetName"])
}

func TestRESTAirspaceU(t *testing.T) {
	HelperDeleteDBIfExistsForAPI()
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	db := HelperAPIInitDB()
	payload := map[string]string{}
	payload["HumanName"] = "test6a"
	payload["NetName"] = "test6b"
	data, err := json.Marshal(payload)
	require.NoError(t, err)
	req, err := http.NewRequest("PUT", "/airspace/1", bytes.NewReader(data))
	require.NoError(t, err)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusCreated, response)
	var fmtresponse map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &fmtresponse))
	assert.Equal(t, "test6a", fmtresponse["HumanName"])
	assert.Equal(t, "test6b", fmtresponse["NetName"])
}

func TestRESTFlightRAll(t *testing.T) {
	HelperDeleteDBIfExistsForAPI()
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	db := HelperAPIInitDB()
	id := HelperAPICreateSampleFlight(t, db)
	req, _ := http.NewRequest("GET", "/flights", nil)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusOK, response)
	var m []map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &m))
	assert.Equal(t, len(m), 1)
	firstEntry := m[0]
	assert.Equal(t, firstEntry["ID"], float64(id))
	assert.Equal(t, firstEntry["Name"], "test")
}

func TestRESTFlightR(t *testing.T) {
	HelperDeleteDBIfExistsForAPI()
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	db := HelperAPIInitDB()
	id := HelperAPICreateSampleFlight(t, db)
	req, err := http.NewRequest("GET", fmt.Sprintf("/flight/%d", id), nil)
	require.NoError(t, err)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusOK, response)
	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &m))
	assert.Equal(t, m["Name"], "test")
}

func TestRESTFlightC(t *testing.T) {
	HelperDeleteDBIfExistsForAPI()
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	db := HelperAPIInitDB()
	payload := map[string]interface{}{
		"AirspaceID": 1,
		"Name": "test3a",
	}
	data, err := json.Marshal(payload)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/flight", bytes.NewReader(data))
	require.NoError(t, err)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusCreated, response)
	var fmtresponse map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &fmtresponse))
	assert.Equal(t, float64(1), fmtresponse["ID"])
	assert.Equal(t, float64(1), fmtresponse["AirspaceID"])
	assert.Equal(t, "test3a", fmtresponse["Name"])
}

func TestRESTFlightU(t *testing.T) {
	HelperDeleteDBIfExistsForAPI()
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	db := HelperAPIInitDB()
	flightID := HelperAPICreateSampleFlight(t, db)
	payload := map[string]interface{}{
		"Name": "test5a",
	}
	data, err := json.Marshal(payload)
	require.NoError(t, err)
	req, err := http.NewRequest("PUT", fmt.Sprintf("/flight/%d", flightID), bytes.NewReader(data))
	require.NoError(t, err)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusCreated, response)
	var fmtresponse map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &fmtresponse))
	assert.Equal(t, "test5a", fmtresponse["Name"])
}

func TestRESTFormationRAll(t *testing.T) {
	HelperDeleteDBIfExistsForAPI()
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	db := HelperAPIInitDB()
	flightID := HelperAPICreateSampleFlight(t, db)
	_ = HelperAPICreateSampleFormation(t, db, flightID)
	req, err := http.NewRequest("GET", fmt.Sprintf("/formations"), nil)
	require.NoError(t, err)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusOK, response)
	var m []map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &m))
	assert.Equal(t, len(m), 1)
	assert.Equal(t, m[0]["Name"], "test")
}

func TestRESTFormationR(t *testing.T) {
	HelperDeleteDBIfExistsForAPI()
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	db := HelperAPIInitDB()
	flightID := HelperAPICreateSampleFlight(t, db)
	formationID := HelperAPICreateSampleFormation(t, db, flightID)
	req, err := http.NewRequest("GET", fmt.Sprintf("/formation/%d", formationID), nil)
	require.NoError(t, err)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusOK, response)
	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &m))
	assert.Equal(t, m["Name"], "test")
}

func TestRESTFormationC(t *testing.T) {
	HelperDeleteDBIfExistsForAPI()
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	db := HelperAPIInitDB()
	flightID := HelperAPICreateSampleFlight(t, db)
	payload := map[string]interface{}{
		"FlightID": flightID,
		"Name": "test3a",
		"CPU": 1,
		"RAM": 128,
		"Disk": 8,
		"BaseName": "test",
		"Domain": "example.com",
		"TargetCount": 0,
	}
	data, err := json.Marshal(payload)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/formation", bytes.NewReader(data))
	require.NoError(t, err)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusCreated, response)
	var fmtresponse map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &fmtresponse))
	assert.Equal(t, float64(1), fmtresponse["ID"])
	assert.Equal(t, float64(flightID), fmtresponse["FlightID"])
	assert.Equal(t, "test3a", fmtresponse["Name"])
	assert.Equal(t, float64(1), fmtresponse["CPU"])
	assert.Equal(t, float64(128), fmtresponse["RAM"])
	assert.Equal(t, float64(8), fmtresponse["Disk"])
	assert.Equal(t, "test", fmtresponse["BaseName"])
	assert.Equal(t, "example.com", fmtresponse["Domain"])
	assert.Equal(t, float64(0), fmtresponse["TargetCount"])
}

func TestRESTFormationU(t *testing.T) {
	HelperDeleteDBIfExistsForAPI()
	require.NoError(t, helperSetupConfigFile("config_dummy_only.yaml"))
	db := HelperAPIInitDB()
	flightID := HelperAPICreateSampleFlight(t, db)
	formationID := HelperAPICreateSampleFormation(t, db, flightID)
	payload := map[string]interface{}{
		"TargetCount": 1,
	}
	data, err := json.Marshal(payload)
	require.NoError(t, err)
	req, err := http.NewRequest("PUT", fmt.Sprintf("/formation/%d", formationID), bytes.NewReader(data))
	require.NoError(t, err)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusOK, response)
	var fmtresponse map[string]interface{}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &fmtresponse))
	assert.Equal(t, float64(flightID), fmtresponse["FlightID"])
	assert.Equal(t, "test", fmtresponse["Name"])
	assert.Equal(t, float64(1), fmtresponse["CPU"])
	assert.Equal(t, float64(128), fmtresponse["RAM"])
	assert.Equal(t, float64(8), fmtresponse["Disk"])
	assert.Equal(t, "test", fmtresponse["BaseName"])
	assert.Equal(t, "example.com", fmtresponse["Domain"])
	assert.Equal(t, float64(1), fmtresponse["TargetCount"])
}

func HelperAPICreateSampleFlight(t *testing.T, db *gorm.DB) int {
	flight := Flight{
		Name: "test",
		AirspaceID: 1,
	}
	tx := db.Save(&flight)
	require.NoError(t, tx.Error)
	return int(flight.ID)
}

func HelperAPICreateSampleFormation(t *testing.T, db *gorm.DB, flightID int) int {
	formation := Formation{
		Name:        "test",
		CPU:         1,
		RAM:         128,
		Disk:        8,
		BaseName:    "test",
		Domain:      "example.com",
		TargetCount: 0,
		FlightID:    flightID,
	}
	tx := db.Save(&formation)
	require.NoError(t, tx.Error)
	return int(formation.ID)
}

func HelperAPIInitDB() *gorm.DB {
	db, _ := db2.ConnectToDB()
	_ = initAirspaces(db)
	return db
}

func HelperAPIExecuteRequest(req *http.Request, apiServer APIServer) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	apiServer.router.ServeHTTP(rr, req)
	return rr
}

func HelperAPIGetServerInstance(db *gorm.DB) APIServer {
	apiServer := APIServer{
		db: db,
	}
	_ = apiServer.Start()
	return apiServer
}

func HelperAPICheckResponseCode(t *testing.T, expected int, actual *httptest.ResponseRecorder) {
	if expected != actual.Code {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual.Code)
	}
}

func HelperDeleteDBIfExistsForAPI() {
	_ = os.Remove("./testing.db")
}