package database

import (
	"test/features/users/data"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(data.User{})
}