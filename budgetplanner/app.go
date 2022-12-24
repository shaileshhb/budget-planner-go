package budgetplanner

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/config"
	"gorm.io/gorm"
)

// Controller is implemented by the controllers.
type Controller interface {
	RegisterRoutes(router *mux.Router)
}

// ModuleConfig needs to be implemented by every module.
type ModuleConfig interface {
	TableMigration(wg *sync.WaitGroup)
}

// App Struct For Start the tsam service.
type App struct {
	sync.Mutex
	Name           string
	Router         *mux.Router
	DB             *gorm.DB
	Log            log.Logger
	Config         config.ConfReader
	Server         *http.Server
	WG             *sync.WaitGroup
	IsInProduction bool
	// EventPool      *event.Pool
	// Auth           *security.Authentication
	// Repository     repository.Repository
}
