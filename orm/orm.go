package orm

import (
	"errors"
	"math"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //to load dialect
)

var instance *gorm.DB

// NoLimit is literal of unlimited search
const NoLimit = math.MaxInt32

// ErrorRecordNotFound is an error when record not found
var ErrorRecordNotFound = gorm.ErrRecordNotFound

// Init opens database
func Init(databasePath string) (err error) {
	opened, err := gorm.Open("sqlite3", databasePath)
	if err != nil {
		return
	}
	instance = opened
	instance.DB().SetMaxIdleConns(5)
	return
}

// GetDB returns opend database of gorm
func GetDB() *gorm.DB {
	return instance
}

// Migrate create tables of models
func Migrate(models ...interface{}) error {
	db := GetDB()
	if db == nil {
		return errors.New("Database is not initilazed. Must call orm.Init() before calling this method")
	}
	return db.AutoMigrate(models...).Error
}
