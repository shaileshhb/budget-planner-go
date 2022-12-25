package account

import (
	"strings"

	"github.com/google/uuid"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/errors"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/models/general"
	userModel "github.com/shaileshhb/budget-planner-go/budgetplanner/models/user"
)

// Account consist of all details regarding user accounts
type Account struct {
	general.Base
	Name   string         `json:"name" gorm:"type:varchar(100);not_null"`
	User   userModel.User `json:"-" gorm:"foreignKey:UserID"` // added to create foregin key. can't create using constraint
	UserID uuid.UUID      `json:"userID" gorm:"type:char(36);index:idx_user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Amount float64        `json:"amount" gorm:"type:decimal(10,2);not_null"`
}

// TableName will specify table name for envelop struct.
func (*Account) TableName() string {
	return "accounts"
}

// Validate will verify compulsory fields of envelop.
func (a *Account) Validate() error {

	if len(strings.TrimSpace(a.Name)) == 0 {
		return errors.NewValidationError("account name must be specified")
	}

	a.Name = strings.TrimSpace(a.Name)

	if a.UserID == uuid.Nil {
		return errors.NewValidationError("user must be specified")
	}

	if a.Amount == 0 {
		return errors.NewValidationError("amount must be greater than 0")
	}

	return nil
}

// AccountDTO contains fields for DTO specifically.
type AccountDTO struct {
	general.BaseDTO
	Name   string    `json:"name"`
	UserID uuid.UUID `json:"userID"`
	Amount float64   `json:"amount"`
}

// TableName will specify table name for envelop struct.
func (*AccountDTO) TableName() string {
	return "accounts"
}
