package repository

import "gorm.io/gorm"

// QueryProcessor allows to modify the query before it is executed
type QueryProcessor func(db *gorm.DB, out interface{}) (*gorm.DB, error)
