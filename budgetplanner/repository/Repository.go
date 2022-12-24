package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines all methods to be present in repository.
type Repository interface {
	Get(uow *UnitOfWork, id uuid.UUID, out interface{}, queryProcessor ...QueryProcessor) error
	GetAll(uow *UnitOfWork, out interface{}, queryProcessor ...QueryProcessor) error
	GetRecord(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error
	GetAllInOrder(uow *UnitOfWork, out, orderBy interface{}, queryProcessor ...QueryProcessor) error

	GetCount(uow *UnitOfWork, out, count *int64, queryProcessors ...QueryProcessor) error
	GetCountUnscoped(uow *UnitOfWork, out, count *int64, queryProcessors ...QueryProcessor) error

	// Other CRUD operations.
	Add(uow *UnitOfWork, out interface{}) error
	Updates(uow *UnitOfWork, out interface{}) error
	UpdateWithMap(uow *UnitOfWork, model interface{}, value map[string]interface{}, queryProcessors ...QueryProcessor) error
	// BatchUpdate(uow *UnitOfWork, value, condition, out interface{}) error

	Save(uow *UnitOfWork, value interface{}) error
	Delete(uow *UnitOfWork, out interface{}, where ...interface{}) error

	RemoveAssociations(uow *UnitOfWork, out interface{}, associationName string, associations ...interface{}) error
	ReplaceAssociations(uow *UnitOfWork, out interface{}, associationName string, associations ...interface{}) error

	// Exec(uow *UnitOfWork, sql string, values ...interface{}) error

	Scan(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error
	// SubQuery(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) (*gorm.SqlExpr, error)
}

// GormRepository will implement repository interface.
type GormRepository struct{}

// NewGormRepository returns new instance of GormRepository.
func NewGormRepository() *GormRepository {
	return &GormRepository{}
}

// Get returns record from table by ID.
func (repository *GormRepository) Get(uow *UnitOfWork, id uuid.UUID, out interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().First(out, "id = ?", id).Error
}

// GetAll returns all records from the table.
func (repository *GormRepository) GetAll(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().Find(out).Error
}

// GetRecord returns a specific record from table with the given filter.
func (repository *GormRepository) GetRecord(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().First(out).Error
}

// GetAllInOrder returns all records from table in specified order.
func (repository *GormRepository) GetAllInOrder(uow *UnitOfWork, out, orderBy interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB

	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().Order(orderBy).Find(out).Error
}

// GetCount gives number of records in database.
func (repository *GormRepository) GetCount(uow *UnitOfWork, out, count *int64, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().Model(out).Count(count).Error
}

// GetCountUnscoped gives number of records in database.
func (repository *GormRepository) GetCountUnscoped(uow *UnitOfWork, out, count *int64, queryProcessors ...QueryProcessor) error {
	db := uow.DB.Unscoped()
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().Model(out).Count(count).Error
}

// Add adds record to table.
func (repository *GormRepository) Add(uow *UnitOfWork, out interface{}) error {
	return uow.DB.Create(out).Error
}

// Update updates the record in table.
func (repository *GormRepository) Updates(uow *UnitOfWork, out interface{}) error {
	return uow.DB.Model(out).Updates(out).Error
}

// UpdateWithMap updates the record in table using map.
//
//	UpdateWithMap(uow, user{id="101"},map[string]interface{}{"name":"Ramesh"}
//
// It will filter by ID only if value has a primary key.
//
//	Query: UPDATE users WHERE `id`="101" SET `name`="Ramesh";
func (repository *GormRepository) UpdateWithMap(uow *UnitOfWork, model interface{}, value map[string]interface{},
	queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, value, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().Model(model).Updates(value).Error
}

// Save updates the record in table. If value doesn't have primary key, new record will be inserted.
func (repository *GormRepository) Save(uow *UnitOfWork, value interface{}) error {
	return uow.DB.Save(value).Error
}

// Delete deletes a record from table.
func (repository *GormRepository) Delete(uow *UnitOfWork, out interface{}, where ...interface{}) error {
	return uow.DB.Delete(out, where...).Error
}

// ReplaceAssociations replaces associations from the given entity.
func (repository *GormRepository) ReplaceAssociations(uow *UnitOfWork, out interface{}, associationName string, associations ...interface{}) error {
	return uow.DB.Model(out).Association(associationName).Replace(associations...)
}

// RemoveAssociations removes associations from the given entity.
func (repository *GormRepository) RemoveAssociations(uow *UnitOfWork, out interface{}, associationName string, associations ...interface{}) error {
	return uow.DB.Model(out).Association(associationName).Delete(associations...)
}

// Scan will fill the out interface with data(fields) based on the given QP conditions.
func (repository *GormRepository) Scan(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Scan(out).Error
}

// executeQueryProcessors executes all queryProcessor func.
func executeQueryProcessors(db *gorm.DB, out interface{}, queryProcessors ...QueryProcessor) (*gorm.DB, error) {
	var err error
	for _, query := range queryProcessors {
		if query != nil {
			db, err = query(db, out)
			if err != nil {
				return db, err
			}
		}
	}
	return db, nil
}
