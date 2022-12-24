package db

import (
	"fmt"
	"time"

	"github.com/shaileshhb/budget-planner-go/budgetplanner/config"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDBConnection Return DB Instace
func NewDBConnection(log log.Logger, conf config.ConfReader) *gorm.DB {

	// dsn -> data source name
	dsn := getConnectionString(conf)

	log.Infof("HERE IS THE OPEN URL: %s:*****@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true\n",
		conf.GetString(config.DBUser),
		conf.GetString(config.DBHost),
		conf.GetString(config.DBPort),
		conf.GetString(config.DBName))

	db, err := gorm.Open(mysql.New(
		mysql.Config{
			DSN: dsn,
		}), &gorm.Config{})
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	// sqlDB is the underlying mysql DB. It is needed to specify connection restrictions.
	sqlDB, err := db.DB()
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	sqlDB.SetMaxIdleConns(90)
	sqlDB.SetMaxOpenConns(400)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	// db.LogMode(true)
	// utf8_general_ci is the default collate for utf8 and it is okay to not specify it.
	// ci means case insensitive.
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci")

	// gorm logger interface needs to be implemented by the 3rd party logger for decorted output.
	// db.SetLogger(log)
	// blocks update without a where clause.
	// db.AllowGlobalUpdate = false

	return db
}

func getConnectionString(conf config.ConfReader) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		conf.GetString(config.DBUser),
		conf.GetString(config.DBPass),
		conf.GetString(config.DBHost),
		conf.GetString(config.DBPort),
		conf.GetString(config.DBName))
}
