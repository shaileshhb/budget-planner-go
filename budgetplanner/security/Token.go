package security

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/config"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/errors"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/log"
	userModal "github.com/shaileshhb/budget-planner-go/budgetplanner/models/user"
)

// GenerateToken take userID, email, tablename as Role  Return Token
func (auth *Authentication) generateToken(claims jwt.Claims) (string, error) {

	// NewWithClaims returns token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// access token string based on token
	tokenString, err := token.SignedString([]byte(auth.Config.GetString(config.JWTKey)))
	if err != nil {
		log.GetLogger().Error(err.Error())
		return "", errors.NewHTTPError("unable to generate token", http.StatusInternalServerError)
	}
	return tokenString, nil
}

// GenerateLoginToken will create new login token
func (auth *Authentication) GenerateLoginToken(a *userModal.Authentication) error {

	// Create a claims map
	// claims based on which token should be created
	// claims := jwt.MapClaims{
	// 	"userID":  userID,
	// 	"name":    name,
	// 	"emailID": email,
	// 	"exp":     time.Now().Add(time.Hour * 20).Unix(),
	// }

	registeredClaims := jwt.RegisteredClaims{
		Issuer:  "budget-planner",
		Subject: "login",
		Audience: jwt.ClaimStrings{
			a.Name, a.Email,
		},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 20)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        a.UserID.String(),
	}

	var err error
	a.Token, err = auth.generateToken(registeredClaims)

	return err
}
