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

	GetCount(uow *UnitOfWork, out interface{}, count *int64, queryProcessors ...QueryProcessor) error
	GetCountUnscoped(uow *UnitOfWork, out interface{}, count *int64, queryProcessors ...QueryProcessor) error

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
func (repository *GormRepository) GetCount(uow *UnitOfWork, out interface{}, count *int64, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().Model(out).Count(count).Error
}

// GetCountUnscoped gives number of records in database.
func (repository *GormRepository) GetCountUnscoped(uow *UnitOfWork, out interface{}, count *int64, queryProcessors ...QueryProcessor) error {
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

// ******************************** All GormRepository methods above this line ********************************

// OrderBy specifies order when retrieving records from database, set reorder to `true` to overwrite defined conditions
//
//	Order("name DESC")
func OrderBy(value interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Order(value)
		return db, nil
	}
}

// Select specify fields that you want to retrieve from database when querying, by default, will select all fields;
// When creating/updating, specify fields that you want to save to database.
func Select(query interface{}, args ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Select(query, args...)
		return db, nil
	}
}

// Join specifies join conditions as query processors. (Use Find() or something similar to get results)
//
//	Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "tsam@example.org")
func Join(query string, args ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Joins(query, args...)
		return db, nil
	}
}

// Model specifies the model you would like to run db operations on
//
//	// update all users's name to `hello`
//	db.Model(&User{}).Update("name", "hello")
//	// if user's primary key is non-blank, will use it as condition, then will only update the user's name to `hello`
//	db.Model(&user).Update("name", "hello")
func Model(value interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Debug().Model(value)
		return db, nil
	}
}

// Filter will filter the results based on condition.
//
//	Filter("name= ?","Ramesh")
//
// Query : WHERE `name`= "Ramesh"
func Filter(condition string, args ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Debug().Where(condition, args...)
		return db, nil
	}
}

// CombineQueries will process slice of queryprocessors and return single queryprocessor.
func CombineQueries(queryProcessors []QueryProcessor) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		tempDB, err := executeQueryProcessors(db, out, queryProcessors...)
		return tempDB, err
	}
}

// DoesRecordExist returns true if the record exists.
//
//	If ID is to be checked then populate it in the model
func DoesRecordExist(db *gorm.DB, out interface{}, queryProcessors ...QueryProcessor) (bool, error) {
	var count int64 = 0
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return false, err
	}

	if err := db.Debug().Model(out).Count(&count).Error; err != nil {
		return false, err
	}
	if count <= 0 {
		return false, nil
	}
	return true, nil
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
