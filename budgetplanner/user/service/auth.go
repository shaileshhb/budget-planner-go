package service

import (
	"github.com/shaileshhb/budget-planner-go/budgetplanner/errors"
	userModal "github.com/shaileshhb/budget-planner-go/budgetplanner/models/user"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/repository"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/security"
	"gorm.io/gorm"
)

type AuthenticationService interface {
	Register(user *userModal.User) error
}

// AuthenticationService service provides methods to update, delete, add, get method for AuthenticationService.
type authenticationService struct {
	db   *gorm.DB
	repo repository.Repository
	auth *security.Authentication
}

// NewAuthenticationService create new AuthenticationService
func NewAuthenticationService(db *gorm.DB, repo repository.Repository, auth *security.Authentication) AuthenticationService {
	return &authenticationService{
		db:   db,
		repo: repo,
		auth: auth,
	}
}

// Register will register new user in the system.
func (ser *authenticationService) Register(user *userModal.User) error {

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
func (ser *authenticationService) validateUser(user *userModal.User) error {
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
