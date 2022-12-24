package repository

import "gorm.io/gorm"

// UnitOfWork represent connection
type UnitOfWork struct {
	DB        *gorm.DB
	Committed bool
}

// NewUnitOfWork creates new instance of UnitOfWork.
func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	commit := false

	return &UnitOfWork{
		DB:        db.Begin(),
		Committed: commit,
	}
}

// Commit use to commit after a successful transaction.
func (uow *UnitOfWork) Commit() {
	if !uow.Committed {
		uow.Committed = true
		uow.DB.Commit()
	}
}

// RollBack is used to rollback a transaction on failure.
func (uow *UnitOfWork) RollBack() {
	// This condition can be used if Rollback() is defered as soon as UOW is created.
	// So we only rollback if it's not committed.
	if !uow.Committed {
		uow.DB.Rollback()
	}
}
