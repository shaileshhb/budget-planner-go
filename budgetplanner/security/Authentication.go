package security

import (
	"github.com/shaileshhb/budget-planner-go/budgetplanner/config"
	"gorm.io/gorm"
)

// Authentication Provide Method AuthUser.
type Authentication struct {
	DB     *gorm.DB
	Config config.ConfReader
}

// NewAuthentication returns new instance of Authentication
func NewAuthentication(db *gorm.DB, config config.ConfReader) *Authentication {
	return &Authentication{
		DB:     db,
		Config: config,
	}
}
