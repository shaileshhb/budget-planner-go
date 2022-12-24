package main

import (
	"github.com/shaileshhb/budget-planner-go/budgetplanner/config"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/db"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/log"
)

func main() {

	// creates new instance of Logger
	log := log.GetLogger()

	// creates new instance of Config
	envconfig := config.NewConfig(false)

	// Create New Instace of DB
	db := db.NewDBConnection(log, envconfig)
	if db == nil {
		log.Fatal("Db connection failed.")
	}
	defer func() {
		// sqlDB, err := db.DB()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// sqlDB.Close()
		// db.Close()
		log.Info("Db closed")
	}()

}
