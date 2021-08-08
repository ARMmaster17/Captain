package models

import (
	"github.com/ARMmaster17/Captain/Shared/framework"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

type CRUDTestObject struct {
	gorm.Model
}

func (C *CRUDTestObject) Create(db *gorm.DB) error {
	return db.Create(C).Error
}

func (C *CRUDTestObject) GetByID(db *gorm.DB, id int) error {
	return db.First(&C, id).Error
}

func (C *CRUDTestObject) Update(db *gorm.DB) error {
	return db.Save(&C).Error
}

func (C *CRUDTestObject) Delete(db *gorm.DB) error {
	return db.Delete(&C).Error
}

func Test_CRUDCreate(t *testing.T) {
	f, err := framework.NewFramework("test")
	require.NoError(t, err)
	mockDB := HelperInitMockDB(t, f)
	cto := CRUDTestObject{}
	cto.CreatedAt = time.Now()
	cto.UpdatedAt = time.Now()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("INSERT INTO `crud_test_objects` (`created_at`,`updated_at`,`deleted_at`) VALUES (?,?,?)")).
		WithArgs(cto.CreatedAt, cto.UpdatedAt, cto.DeletedAt).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mockDB.ExpectCommit()
	mockDB.ExpectRollback()
	err = cto.Create(f.DB)
	assert.NoError(t, err)
}

func Test_CRUDGetByID(t *testing.T) {
	f, err := framework.NewFramework("test")
	require.NoError(t, err)
	mockDB := HelperInitMockDB(t, f)
	mockDB.NewRows([]string{"id"}).AddRow("0")
	cto := CRUDTestObject{}
	cto.CreatedAt = time.Now()
	cto.UpdatedAt = time.Now()
	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `crud_test_objects` WHERE `crud_test_objects`.`id` = ? AND `crud_test_objects`.`deleted_at` IS NULL ORDER BY `crud_test_objects`.`id` LIMIT 1")).
		WithArgs(0).WillReturnRows(mockDB.NewRows([]string{"id"}).AddRow("0"))
	mockDB.ExpectCommit()
	mockDB.ExpectRollback()
	err = cto.GetByID(f.DB, 0)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), cto.ID)
}

func Test_CRUDUpdate(t *testing.T) {
	f, err := framework.NewFramework("test")
	require.NoError(t, err)
	mockDB := HelperInitMockDB(t, f)
	cto := CRUDTestObject{}
	cto.CreatedAt = time.Now()
	cto.UpdatedAt = time.Now()
	cto.ID = 1
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("UPDATE `crud_test_objects` SET `created_at`=?,`updated_at`=?,`deleted_at`=? WHERE `id` = ?")).
		WithArgs(cto.CreatedAt, sqlmock.AnyArg(), cto.DeletedAt, cto.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mockDB.ExpectCommit()
	mockDB.ExpectRollback()
	err = cto.Update(f.DB)
	require.NoError(t, err)
}

func Test_CRUDDelete(t *testing.T) {
	f, err := framework.NewFramework("test")
	require.NoError(t, err)
	mockDB := HelperInitMockDB(t, f)
	cto := CRUDTestObject{}
	cto.CreatedAt = time.Now()
	cto.UpdatedAt = time.Now()
	cto.ID = 1
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta("UPDATE `crud_test_objects` SET `deleted_at`=? WHERE `crud_test_objects`.`id` = ? AND `crud_test_objects`.`deleted_at` IS NULL")).
		WithArgs(sqlmock.AnyArg(), cto.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mockDB.ExpectCommit()
	mockDB.ExpectRollback()
	err = cto.Delete(f.DB)
	require.NoError(t, err)
}

// HelperInitMockDB injects a mocking framework into the framework that mocks all database transactions.
func HelperInitMockDB(t *testing.T, f *framework.Framework) sqlmock.Sqlmock {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)
	f.DB = gormDB
	return mock
}
