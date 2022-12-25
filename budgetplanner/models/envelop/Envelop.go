package envelop

import (
	"strings"

	"github.com/google/uuid"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/errors"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/models/general"
	userModel "github.com/shaileshhb/budget-planner-go/budgetplanner/models/user"
)

// Envelop will consist of data related to user envelops.
type Envelop struct {
	general.Base
	Name   string         `json:"name" gorm:"type:varchar(100);not_null"`
	User   userModel.User `json:"-" gorm:"foreignKey:UserID"` // added to create foregin key. can't create using constraint
	UserID uuid.UUID      `json:"userID" gorm:"type:char(36);index:idx_user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Amount float64        `json:"amount" gorm:"type:decimal(10, 2);not_null"`
}

// TableName will specify table name for envelop struct.
func (*Envelop) TableName() string {
	return "envelops"
}

// Validate will verify compulsory fields of envelop.
func (e *Envelop) Validate() error {

	if len(strings.TrimSpace(e.Name)) == 0 {
		return errors.NewValidationError("envelop name must be specified")
	}

	e.Name = strings.TrimSpace(e.Name)

	if e.UserID == uuid.Nil {
		return errors.NewValidationError("user must be specified")
	}

	if e.Amount == 0 {
		return errors.NewValidationError("amount must be greater than 0")
	}

	return nil
}

// EnvelopDTO contains fields for DTO specifically.
type EnvelopDTO struct {
	general.BaseDTO
	Name   string    `json:"name"`
	UserID uuid.UUID `json:"userID"`
	Amount float64   `json:"amount"`
}

// TableName will specify table name for envelop struct.
func (*EnvelopDTO) TableName() string {
	return "envelops"
}
