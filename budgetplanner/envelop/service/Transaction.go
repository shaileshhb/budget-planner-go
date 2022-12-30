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

// TransactionService service provides methods to update, delete, add, get method for TransactionService.
type TransactionService interface {
	AddTransaction(transaction *envelopModel.Transaction) error
	UpdateTransaction(transaction *envelopModel.Transaction) error
	DeleteTransaction(transaction *envelopModel.Transaction) error
}

// transactionService
type transactionService struct {
	db   *gorm.DB
	repo repository.Repository
	auth *security.Authentication
}

// NewTransactionService create new envelop service.
func NewTransactionService(db *gorm.DB, repo repository.Repository, auth *security.Authentication) TransactionService {
	return &transactionService{
		db:   db,
		repo: repo,
		auth: auth,
	}
}

// AddTransaction will add new transaction for user in specified envelop.
func (ser *transactionService) AddTransaction(transaction *envelopModel.Transaction) error {

	err := ser.validateUserID(transaction.UserID)
	if err != nil {
		return err
	}

	err = ser.validateEnvelopID(transaction.UserID, transaction.EnvelopID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.db)
	defer uow.RollBack()

	err = ser.repo.Add(uow, transaction)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// UpdateTransaction will update specified transaction of user.
func (ser *transactionService) UpdateTransaction(transaction *envelopModel.Transaction) error {

	err := ser.validateUserID(transaction.UserID)
	if err != nil {
		return err
	}

	err = ser.validateEnvelopID(transaction.UserID, transaction.EnvelopID)
	if err != nil {
		return err
	}

	err = ser.validateTransactionID(transaction.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.db)
	defer uow.RollBack()

	tempTransaction := envelopModel.Transaction{}

	err = ser.repo.GetRecord(uow, &tempTransaction, repository.Filter("`id` = ?", transaction.ID),
		repository.Select("`created_by`"))
	if err != nil {
		return err
	}

	transaction.CreatedBy = tempTransaction.CreatedBy

	err = ser.repo.Save(uow, transaction)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// DeleteTransaction will delete specified transaction of user.
func (ser *transactionService) DeleteTransaction(transaction *envelopModel.Transaction) error {

	err := ser.validateTransactionID(transaction.ID)
	if err != nil {
		return err
	}

	err = ser.validateUserID(transaction.UserID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(ser.db)
	defer uow.RollBack()

	err = ser.repo.UpdateWithMap(uow, transaction, map[string]interface{}{
		"DeletedAt": time.Now(),
		"DeletedBy": transaction.UserID,
	}, repository.Filter("`id` = ?", transaction.ID))
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// validateUserID will verify if userID exist or not.
func (ser *transactionService) validateUserID(userID uuid.UUID) error {

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
func (ser *transactionService) validateEnvelopID(userID, envelopID uuid.UUID) error {

	exist, err := repository.DoesRecordExist(ser.db, envelopModel.Envelop{},
		repository.Filter("envelops.`id` = ? AND envelops.`user_id` = ?", envelopID, userID))
	if err != nil {
		return err
	}
	if !exist {
		return errors.NewValidationError("Envelop not found")
	}
	return nil
}

// validateTransactionID will verify if transaction exist or not.
func (ser *transactionService) validateTransactionID(transactionID uuid.UUID) error {

	exist, err := repository.DoesRecordExist(ser.db, envelopModel.Transaction{},
		repository.Filter("transactions.`id` = ?", transactionID))
	if err != nil {
		return err
	}

	if !exist {
		return errors.NewValidationError("Transaction not found")
	}
	return nil
}
