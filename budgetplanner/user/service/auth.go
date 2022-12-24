package service

import (
	"github.com/shaileshhb/budget-planner-go/budgetplanner/errors"
	userModal "github.com/shaileshhb/budget-planner-go/budgetplanner/models/user"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/repository"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/security"
	"gorm.io/gorm"
)

// Authentication service provides methods to update, delete, add, get method for Authentication.
type Authentication struct {
	db   *gorm.DB
	repo repository.Repository
	auth *security.Authentication
}

// create new difficulty serviec
func NewAuthentication(db *gorm.DB, repo repository.Repository, auth *security.Authentication) *Authentication {
	return &Authentication{
		db:   db,
		repo: repo,
		auth: auth,
	}
}

// Register will register new user in the system.
func (ser *Authentication) Register(user *userModal.User) error {

	err := ser.validateUser(user)
	if err != nil {
		return err
	}

	password, err := security.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(password)

	uow := repository.NewUnitOfWork(ser.db)
	defer uow.RollBack()

	err = ser.repo.Add(uow, user)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// validateUser will check if it is unique user.
func (ser *Authentication) validateUser(user *userModal.User) error {
	exist, err := repository.DoesRecordExist(ser.db, userModal.User{}, repository.Filter("users.`id` != ?"+
		" AND users.`email` = ? AND users.`username` = ?", user.ID, user.Email, user.Username))
	if err != nil {
		return err
	}
	if exist {
		return errors.NewValidationError("Email or Username already exist. Try loggin in.")
	}
	return nil
}
