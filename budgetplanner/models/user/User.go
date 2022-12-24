package user

import (
	"strings"

	"github.com/shaileshhb/budget-planner-go/budgetplanner/errors"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/models/general"
)

// User will store all the information required of a user.
type User struct {
	general.Base
	Name         string  `json:"name" gorm:"type:varchar(100)"`
	Username     string  `json:"username" gorm:"type:varchar(200)" sql:"index"`
	Email        string  `json:"email" gorm:"type:varchar(255)" sql:"index"`
	Password     string  `json:"password" gorm:"type:varchar(255)" sql:"index"`
	DateOfBirth  *string `json:"dateOfBirth" gorm:"type:varchar(10)"`
	Gender       *string `json:"gender" gorm:"type:varchar(20)"`
	Contact      *string `json:"contact" gorm:"type:varchar(15)"`
	ProfileImage *string `json:"profileImage" gorm:"type:varchar(255)"`
	IsVerified   bool    `json:"isVerified" gorm:"type:tinyint;default:0"`
}

// TableName will specify table name for user struct.
func (*User) TableName() string {
	return "users"
}

// User contains fields for DTO specifically.
type UserDTO struct {
	general.BaseDTO
	Name         string `json:"name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"-"`
	DateOfBirth  string `json:"dateOfBirth"`
	Gender       string `json:"gender"`
	Contact      string `json:"contact"`
	ProfileImage string `json:"profileImage"`
	IsVerified   bool   `json:"isVerified"`
}

// TableName will specify table name for user struct.
func (*UserDTO) TableName() string {
	return "users"
}

func (u *User) Validate() error {
	if len(strings.TrimSpace(u.Name)) == 0 {
		return errors.NewValidationError("name must be specified")
	}

	u.Name = strings.TrimSpace(u.Name)

	if len(strings.TrimSpace(u.Username)) == 0 {
		return errors.NewValidationError("username must be specified")
	}

	u.Username = strings.TrimSpace(u.Username)

	if len(strings.TrimSpace(u.Email)) == 0 {
		return errors.NewValidationError("email must be specified")
	}

	u.Email = strings.TrimSpace(u.Email)

	if len(strings.TrimSpace(u.Password)) == 0 {
		return errors.NewValidationError("password must be specified")
	}

	u.Password = strings.TrimSpace(u.Password)
	if u.Contact != nil {
		*u.Contact = strings.TrimSpace(*u.Contact)
	}

	return nil
}
