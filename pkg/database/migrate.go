package database

import (
	"github.com/aragornz325/piloto-api/internal/profile/model"
	"github.com/aragornz325/piloto-api/internal/user/model"
	"github.com/aragornz325/piloto-api/pkg/logger"
)

func ExecuteMigrations() {
	// Migrate the schema
	logger.Log.Info("Migrating database...")
	if err := DB.AutoMigrate(
		&userModel.User{},
		&profileModel.Profile{},
	); err != nil {
		panic("failed to migrate database: " + err.Error())
	}
	logger.Log.Info("Database migrated successfully")
}