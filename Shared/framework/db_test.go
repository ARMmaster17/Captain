package framework

import (
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
)

func Test_UseMemDBOnNoConfig(t *testing.T) {
	if os.Getenv("CAPTAIN_TEST_USEREALDB") != "TRUE" && runtime.GOOS == "windows" {
		t.Skip()
	}
	db, err := InitDB("")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.Equal(t, "sqlite", db.Name())
}

// This test asserts that the testing database is not initialized so that we can inject mocks.
func Test_CreateTestDB(t *testing.T) {
	db, err := InitDB("TEST")
	assert.NoError(t, err)
	assert.Nil(t, db)
}

func Test_ConnectToPostgres(t *testing.T) {
	if os.Getenv("CAPTAIN_TEST_USEREALDB") != "TRUE" {
		t.Skip()
	}
	db, err := InitDB("host=localhost user=captain password=captain dbname=captain port=9920 sslmode=disable")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.Equal(t, "postgres", db.Name())
}
