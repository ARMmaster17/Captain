package main

import (
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRESTAirspaceCRUD(t *testing.T) {
	db := HelperAPIInitMemoryDB()
	req, _ := http.NewRequest("GET", "/airspaces", nil)
	response := HelperAPIExecuteRequest(req, db)
	HelperAPICheckResponseCode(t, http.StatusOK, response)
}

func HelperAPIInitMemoryDB() *gorm.DB {
	db, _ := ConnectToDB()
	_ = initAirspaces(db)
	return db
}

func HelperAPIExecuteRequest(req *http.Request, db *gorm.DB) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	apiServer := APIServer{
		db: db,
	}
	_ = apiServer.Start()
	apiServer.router.ServeHTTP(rr, req)
	return rr
}

func HelperAPICheckResponseCode(t *testing.T, expected int, actual *httptest.ResponseRecorder) {
	if expected != actual.Code {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual.Code)
	}
}
