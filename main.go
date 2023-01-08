package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/shaileshhb/budget-planner-go/budgetplanner"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/config"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/db"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/log"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/repository"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/security"
	"github.com/shaileshhb/budget-planner-go/module"
)

var production = "false"

func main() {
	isAppInProduction := false
	if production == "true" {
		isAppInProduction = true
	}

	// creates new instance of Logger
	log := log.GetLogger()

	file, err := os.OpenFile("./logs/out.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	}
	defer file.Close()

	// creates new instance of Config
	envconfig := config.NewConfig(isAppInProduction)

	// Create New Instace of DB
	db := db.NewDBConnection(log, envconfig)
	if db == nil {
		log.Fatal("Db connection failed.")
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal(err)
		}
		sqlDB.Close()
		// db.Close()
		log.Info("Db closed")
	}()

	midlware := security.NewAuthentication(db, envconfig)
	var wg sync.WaitGroup

	var repository = repository.NewGormRepository()

	app := budgetplanner.NewApp("Money wisely", db, log, envconfig,
		&wg, midlware, isAppInProduction, repository)

	module.CreateRouterInstance(app, repository)
	module.Configure(app)

	err = app.Start()
	if err != nil {
		log.Fatal(err)
		stopApp(app)
	}

	// Stop Server On System Call or Interrupt.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	stopApp(app)
}

func stopApp(app *budgetplanner.App) {
	app.Stop()
	app.WG.Wait()
	fmt.Println("After wait")
	os.Exit(0)
}
