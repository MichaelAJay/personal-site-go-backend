package db_client

import (
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB
var once sync.Once

func DbClient(dsn string) *gorm.DB {
	once.Do(func() {
		var err error
		Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true,
		})

		if err != nil {
			panic("failed to connect database")
		}
	})
	return Db
}
