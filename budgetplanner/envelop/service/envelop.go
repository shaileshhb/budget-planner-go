package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/errors"
	envelopModel "github.com/shaileshhb/budget-planner-go/budgetplanner/models/envelop"
	userModel "github.com/shaileshhb/budget-planner-go/budgetplanner/models/user"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/repository"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/security"
	"gorm.io/gorm"
)

// envelopService service provides methods to update, delete, add, get method for envelopService.
type EnvelopService interface {
	AddEnvelop(envelop *envelopModel.Envelop) error
	UpdateEnvelop(envelop *envelopModel.Envelop) error
	DeleteEnvelop(envelop *envelopModel.Envelop) error
	GetEnvelops(envelops *[]envelopModel.EnvelopDTO, userID uuid.UUID) error
}

// envelopService
type envelopService struct {
	db          *gorm.DB
	repo        repository.Repository
	auth        *security.Authentication
	MaxEnvelops int
}

// NewEnvelopService create new envelop service.
func NewEnvelopService(db *gorm.DB, repo repository.Repository, auth *security.Authentication) EnvelopService {
	return &envelopService{
		db:          db,
		repo:        repo,
		auth:        auth,
		MaxEnvelops: 20,
	}
}

// AddEnvelop will add new envelop for specified user.
func (ser *envelopService) AddEnvelop(envelop *envelopModel.Envelop) error {

	err := ser.validateUserID(envelop.UserID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.db)
	defer uow.RollBack()

	var totalCount int64

	err = ser.repo.GetCount(uow, envelopModel.Envelop{}, &totalCount,
		repository.Filter("envelops.`user_id` = ?", envelop.UserID))
	if err != nil {
		return err
	}

	if totalCount >= int64(ser.MaxEnvelops) {
		return errors.NewValidationError("Maximum envelops created")
	}

	err = ser.repo.Add(uow, envelop)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// UpdateEnvelop will update specified envelop.
func (ser *envelopService) UpdateEnvelop(envelop *envelopModel.Envelop) error {

	err := ser.validateUserID(envelop.UserID)
	if err != nil {
		return err
	}

	err = ser.validateEnvelopID(envelop.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.db)
	defer uow.RollBack()

	// using update because there is no nullable field in envelops table
	err = ser.repo.Updates(uow, envelop)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// DeleteEnvelop will delete specified envelop.
func (ser *envelopService) DeleteEnvelop(envelop *envelopModel.Envelop) error {

	err := ser.validateEnvelopID(envelop.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.db)
	defer uow.RollBack()

	// using update because there is no nullable field in envelops table
	err = ser.repo.UpdateWithMap(uow, envelopModel.Envelop{}, map[string]interface{}{
		"DeletedBy": envelop.DeletedBy,
		"DeletedAt": time.Now(),
	}, repository.Filter("envelops.`id` = ?", envelop.ID))
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// GetEnvelops will fetch all the envelops for specifed user.
func (ser *envelopService) GetEnvelops(envelops *[]envelopModel.EnvelopDTO, userID uuid.UUID) error {

	err := ser.validateUserID(userID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.db)
	defer uow.RollBack()

	err = ser.repo.GetAllInOrder(uow, envelops, "envelops.`name`", repository.Filter("envelops.`user_id` = ?", userID))
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// validateUserID will verify if userID exist or not.
func (ser *envelopService) validateUserID(userID uuid.UUID) error {

	exist, err := repository.DoesRecordExist(ser.db, userModel.User{},
		repository.Filter("users.`id` = ?", userID))
	if err != nil {
		return err
	}
	if !exist {
		return errors.NewValidationError("User not found")
	}
	return nil
}

// validateEnvelopID will verify if envelopID exist or not.
func (ser *envelopService) validateEnvelopID(envelopID uuid.UUID) error {

	exist, err := repository.DoesRecordExist(ser.db, envelopModel.Envelop{},
		repository.Filter("envelops.`id` = ?", envelopID))
	if err != nil {
		return err
	}
	if !exist {
		return errors.NewValidationError("Envelop not found")
	}
	return nil
}
