package framework

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_UseMemDBOnNoConfig(t *testing.T) {
	db, err := InitDB("")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.Equal(t, "sqlite", db.Name())
}

func Test_CreateMemDB(t *testing.T) {
	db, err := InitDB("TEST")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.Equal(t, "sqlite", db.Name())
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
