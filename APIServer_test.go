package main

import (
	"encoding/json"
	db2 "github.com/ARMmaster17/Captain/db"
	"github.com/go-playground/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRESTAirspaceRAll(t *testing.T) {
	db := HelperAPIInitMemoryDB()
	req, _ := http.NewRequest("GET", "/airspaces", nil)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusOK, response)

	var m []map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, len(m), 1)
	firstEntry := m[0]
	assert.Equal(t, firstEntry["ID"], float64(1))
	assert.Equal(t, firstEntry["HumanName"], "Default Airspace")
	assert.Equal(t, firstEntry["NetName"], "default")
}

func TestRESTAirspaceR(t *testing.T) {
	db := HelperAPIInitMemoryDB()
	req, _ := http.NewRequest("GET", "/airspace/1", nil)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusOK, response)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, m["ID"], float64(1))
	assert.Equal(t, m["HumanName"], "Default Airspace")
	assert.Equal(t, m["NetName"], "default")
}

func TestRESTFlightRAll(t *testing.T) {
	db := HelperAPIInitMemoryDB()

	flight := Flight{
		AirspaceID: 1,
		Name: "FirstFlight",
	}
	db.Create(&flight)

	req, _ := http.NewRequest("GET", "/flights", nil)
	api := HelperAPIGetServerInstance(db)
	response := HelperAPIExecuteRequest(req, api)
	HelperAPICheckResponseCode(t, http.StatusOK, response)

	var m []map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, len(m), 1)
	firstEntry := m[0]
	assert.Equal(t, firstEntry["ID"], float64(flight.ID))
	assert.Equal(t, firstEntry["Name"], "FirstFlight")

	db.Delete(&Flight{}, &flight.ID)
}

func TestRESTFlightR(t *testing.T) {
	t.Skip()
}

func HelperAPIInitMemoryDB() *gorm.DB {
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
