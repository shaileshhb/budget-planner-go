package envelop

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/errors"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/models/general"
	userModel "github.com/shaileshhb/budget-planner-go/budgetplanner/models/user"
)

// Transaction will contain all details related to user transactions.
type Transaction struct {
	general.Base
	User            userModel.User `json:"-" gorm:"foreignKey:UserID"` // added to create foregin key. can't create using constraint
	Envelop         Envelop        `json:"-" gorm:"foreignKey:EnvelopID"`
	UserID          uuid.UUID      `json:"userID" gorm:"type:char(36);index:idx_user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	EnvelopID       uuid.UUID      `json:"envelopID" gorm:"type:char(36);index:idx_envelop_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Payee           string         `json:"payee" gorm:"type:varchar(100);not_null"`
	Amount          float64        `json:"amount" gorm:"type:decimal(10,2);not_null"`
	Date            time.Time      `json:"date" gorm:"type:datetime;not_null"`
	TransactionType string         `json:"transactionType" gorm:"type:varchar(255)"`
	Description     *string        `json:"description" gorm:"type:varchar(1000)"`
}

// TableName will specify table name for transaction struct.
func (*Transaction) TableName() string {
	return "transactions"
}

// Validate will verfiy compuslory fields of transaction struct.
func (t *Transaction) Validate() error {

	if t.UserID == uuid.Nil {
		return errors.NewValidationError("user must be specified")
	}

	if t.EnvelopID == uuid.Nil {
		return errors.NewValidationError("envelop must be specified")
	}

	if len(strings.TrimSpace(t.Payee)) == 0 {
		return errors.NewValidationError("payee must be specified")
	}

	t.Payee = strings.TrimSpace(t.Payee)

	if t.Amount == 0 {
		return errors.NewValidationError("amount must be specified")
	}

	if len(t.TransactionType) == 0 {
		return errors.NewValidationError("transaction type must be specified")
	}

	if t.Date.IsZero() {
		return errors.NewValidationError("date must be specified")
	}

	return nil
}

// TransactionDTO contains fields for DTO specifically.
type TransactionDTO struct {
	general.BaseDTO
	Payee           string    `json:"payee"`
	Amount          float64   `json:"amount"`
	Date            time.Time `json:"date"`
	TransactionType string    `json:"transactionType"`
	Description     *string   `json:"description"`
}

// TableName will specify table name for transaction struct.
func (*TransactionDTO) TableName() string {
	return "transactions"
}
