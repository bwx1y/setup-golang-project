package config

import (
	"setup-preoject/app/model/entity"

	"gorm.io/gorm"
)

func MigrationDatabase(db *gorm.DB) {
	db.AutoMigrate(
		&entity.User{},
	)
}
