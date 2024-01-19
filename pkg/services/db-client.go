package services

import (
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var once sync.Once

func DbClient() *gorm.DB {
	once.Do(func() {
		dsn := "TODO"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true,
		})

		if err != nil {
			panic("failed to connect database")
		}

		// Migrate the schema
		db.AutoMigrate()
	})
	return db
}
