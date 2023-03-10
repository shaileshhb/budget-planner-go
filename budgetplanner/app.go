package budgetplanner

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/config"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/log"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/repository"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/security"
	"gorm.io/gorm"
)

// Controller is implemented by the controllers.
type Controller interface {
	RegisterRoutes(router *gin.RouterGroup)
}

// ModuleConfig needs to be implemented by every module.
type ModuleConfig interface {
	TableMigration(wg *sync.WaitGroup)
}

// App Struct For Start the tsam service.
type App struct {
	sync.Mutex
	Name           string
	Engine         *gin.Engine
	RouterGroup    *gin.RouterGroup
	DB             *gorm.DB
	Log            log.Logger
	Config         config.ConfReader
	Server         *http.Server
	WG             *sync.WaitGroup
	Auth           *security.Authentication
	Repository     repository.Repository
	IsInProduction bool
	// EventPool      *event.Pool
}

// NewApp returns app.
func NewApp(name string, db *gorm.DB, log log.Logger, conf config.ConfReader, wg *sync.WaitGroup,
	auth *security.Authentication, isProd bool, repo repository.Repository) *App {
	return &App{
		Name:           name,
		DB:             db,
		Log:            log,
		Config:         conf,
		WG:             wg,
		Auth:           auth,
		IsInProduction: isProd,
		Repository:     repo,
		// EventPool:      pool,
	}
}

// InitializeRouter Register the route.
// # new router
func (app *App) InitializeRouter() {
	app.Log.Info(app.Name + " App Route initializing")

	app.Engine = gin.New()
	app.initializeServer()
}

// initializeServer will initialize server with the given config.
func (app *App) initializeServer() {
	app.Engine.Use(gin.Recovery())
	app.Engine.Use(gin.Logger())

	app.Engine.Use(cors.New(cors.Config{
		AllowOriginFunc: app.checkOrigin,
		AllowMethods: []string{
			http.MethodPost, http.MethodPut, http.MethodGet, http.MethodDelete, http.MethodOptions,
		},
		AllowHeaders: []string{
			"Content-Type", "X-Total-Count", "Authorization",
		},
	}))
}

func (app *App) checkOrigin(origin string) bool {
	// origin will be the actual origin from which the request is made.
	if !app.IsInProduction {
		return true
	}

	switch origin {
	case "https://shaileshhb.github.io/budget-planner-frontend":
		return true
	}

	return false
}

// RegisterControllerRoutes will register the specified routes in controllers.
func (app *App) RegisterControllerRoutes(controllers []Controller) {
	app.Lock()
	defer app.Unlock()

	// controllers registering routes.
	for _, controller := range controllers {
		// need to check if gin can register routes using go routine
		controller.RegisterRoutes(app.Engine.Group("/api/v1/budget-planner"))
	}
}

// MigrateTables will do a table table migration for all modules.
func (app *App) MigrateTables(configs []ModuleConfig) {
	app.WG.Add(len(configs))
	for _, config := range configs {
		config.TableMigration(app.WG)
		app.WG.Done()
	}
	app.WG.Wait()
	app.Log.Info("End of Migration")
}

func (app *App) getPort() string {
	return app.Config.GetString(config.PORT)
}

// Start will start the app.
func (app *App) Start() error {

	app.Log.Info("Server Time: ", time.Now())
	app.Log.Info("Server Running on port: ", app.getPort())

	// if err := app.Server.ListenAndServe(); err != nil {
	// 	app.Log.Error("Listen and serve error: ", err)
	// 	return err
	// }

	app.Engine.Run(fmt.Sprintf(":%s", app.getPort()))
	return nil
}

// Stop stops the app.
func (app *App) Stop() {
	// Stopping scheduler.
	context, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	// Closing db

	sqlDB, err := app.DB.DB()
	if err != nil {
		app.Log.Fatal("Fail to close db...")
		return
	}

	sqlDB.Close()
	app.Log.Info("Db closed")

	// Stopping Server.
	err = app.Server.Shutdown(context)
	if err != nil {
		app.Log.Fatal("Fail to Stop Server...")
		return
	}
	app.Log.Info("Server shutdown gracefully.")
}
