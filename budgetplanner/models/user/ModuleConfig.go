package user

import (
	"sync"

	"github.com/shaileshhb/budget-planner-go/budgetplanner/log"
	"gorm.io/gorm"
)

// ModuleConfig use for Automigrant Tables.
type ModuleConfig struct {
	db *gorm.DB
}

// NewUserModuleConfig Return New Module Config.
func NewUserModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		db: db,
	}
}

// TableMigration Update Table Structure with Latest Version.
func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {
	var models []interface{} = []interface{}{
		&User{},
	}

	for _, model := range models {
		err := config.db.Debug().AutoMigrate(model)
		if err != nil {
			log.GetLogger().Errorf("Auto Migration ==> %s", err.Error())
		}
	}

	log.GetLogger().Info("User Module Configured.")
}
