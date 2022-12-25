package service

import (
	"github.com/google/uuid"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/errors"
	userModal "github.com/shaileshhb/budget-planner-go/budgetplanner/models/user"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/repository"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/security"
	"gorm.io/gorm"
)

// AuthenticationService consist of all methods AuthenticationService should implement.
type AuthenticationService interface {
	Register(user *userModal.User) error
	Login(login *userModal.Login, auth *userModal.Authentication) error
	UpdateUser(user *userModal.User) error
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

// Login will verify user details and login into the system
func (ser *authenticationService) Login(login *userModal.Login, auth *userModal.Authentication) error {

	uow := repository.NewUnitOfWork(ser.db)
	defer uow.RollBack()

	user := userModal.User{}
	err := ser.repo.GetRecord(uow, &user, repository.Filter("users.`username` = ?", login.Username))
	if err != nil {
		return err
	}

	err = security.ComparePassword(user.Password, login.Password)
	if err != nil {
		return errors.NewValidationError("Invalid username or password.")
	}

	auth.UserID = user.ID
	auth.Name = user.Name
	auth.Email = user.Email

	err = ser.auth.GenerateLoginToken(auth)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// UpdateUser will update user details.
func (ser *authenticationService) UpdateUser(user *userModal.User) error {

	err := ser.validateUserID(user.ID)
	if err != nil {
		return err
	}

	err = ser.validateUser(user)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.db)
	defer uow.RollBack()

	tempUser := userModal.User{}

	err = ser.repo.GetRecord(uow, &tempUser, repository.Filter("users.`id` = ?", user.ID),
		repository.Select("`created_by`"))
	if err != nil {
		return err
	}

	user.CreatedBy = tempUser.CreatedBy

	err = ser.repo.Save(uow, &user)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// validateUserID will verify if userID exist or not.
func (ser *authenticationService) validateUserID(userID uuid.UUID) error {

	exist, err := repository.DoesRecordExist(ser.db, userModal.User{},
		repository.Filter("users.`id` = ?", userID))
	if err != nil {
		return err
	}
	if !exist {
		return errors.NewValidationError("User not found")
	}
	return nil
}

// validateUer will check if it is unique user.
func (ser *authenticationService) validateUser(user *userModal.User) error {
	exist, err := repository.DoesRecordExist(ser.db, userModal.User{},
		repository.Filter("users.`id` != ? AND users.`email` = ?"+
			" AND users.`username` = ?", user.ID, user.Email, user.Username))
	if err != nil {
		return err
	}
	if exist {
		return errors.NewValidationError("Email or Username already exist. Try loggin in.")
	}
	return nil
}
