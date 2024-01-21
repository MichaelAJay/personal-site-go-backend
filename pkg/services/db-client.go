package services

import (
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var once sync.Once

func DbClient(dsn string) *gorm.DB {
	once.Do(func() {
		log.Printf("DSN: %s\n", dsn)

		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true,
		})

		if err != nil {
			panic("failed to connect database")
		}

		// Migrate the schema
		// db.AutoMigrate()

		log.Println("Database connection established")
	})
	return db
}
