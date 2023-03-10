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
	Username     string  `json:"username" gorm:"type:varchar(200);unique;index:idx_username"`
	Email        string  `json:"email" gorm:"type:varchar(255);unique;index:idx_email"`
	Password     string  `json:"password" gorm:"type:varchar(255);index:idx_password"`
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

// UserDTO contains fields for DTO specifically.
type UserDTO struct {
	general.BaseDTO
	Name         string  `json:"name"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	Password     string  `json:"-"`
	DateOfBirth  *string `json:"dateOfBirth"`
	Gender       *string `json:"gender"`
	Contact      *string `json:"contact"`
	ProfileImage *string `json:"profileImage"`
	IsVerified   bool    `json:"isVerified"`
}

// TableName will specify table name for user struct.
func (*UserDTO) TableName() string {
	return "users"
}

// ValidateRegistration will verify compulsory fields of user.
func (u *User) ValidateRegistration() error {
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

// ValidateUser will verify compulsory fields of user.
func (u *User) ValidateUser() error {
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

	if u.Contact != nil {
		*u.Contact = strings.TrimSpace(*u.Contact)
	}

	return nil
}

// Login contains details required for login.
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate will verify compulsory fields of login.
func (l *Login) Validate() error {
	if len(strings.TrimSpace(l.Username)) == 0 {
		return errors.NewValidationError("username must be specified")
	}

	if len(strings.TrimSpace(l.Password)) == 0 {
		return errors.NewValidationError("password must be specified")
	}

	l.Username = strings.TrimSpace(l.Username)
	l.Password = strings.TrimSpace(l.Password)

	return nil
}
