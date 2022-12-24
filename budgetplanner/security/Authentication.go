package security

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/config"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/errors"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/log"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/web"
	"gorm.io/gorm"
)

// Authentication Provide Method AuthUser.
type Authentication struct {
	DB     *gorm.DB
	Config config.ConfReader
}

// NewAuthentication returns new instance of Authentication
func NewAuthentication(db *gorm.DB, config config.ConfReader) *Authentication {
	return &Authentication{
		DB:     db,
		Config: config,
	}
}

// Middleware will fetch jwt token from header and verify it.
func (auth *Authentication) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// tokenStr, err := request.HeaderExtractor{"token"}.ExtractToken(r)

		tokenStr := r.Header.Get("Authorization")
		if len(tokenStr) == 0 {
			log.GetLogger().Error("Token must be specified")
			web.RespondError(w, errors.NewHTTPError("Token must be specified", http.StatusUnauthorized))
			return
		}

		claims := jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(
			tokenStr, &claims,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
				}
				return []byte(auth.Config.GetString(config.JWTKey)), nil
			})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				log.GetLogger().Error(err.Error())
				web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusUnauthorized))
				return
			}

			log.GetLogger().Error(err.Error())
			web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
			return
		}

		if claims.VerifyExpiresAt(time.Now(), true) {
			log.GetLogger().Error("Session expired! Please login again")
			web.RespondError(w, errors.NewHTTPError("Session expired! Please login again", http.StatusUnauthorized))
			return
		}

		// prints all the claims
		// for key, val := range claims {
		// 	fmt.Printf("Key: %v, value: %v\n", key, val)
		// }

		// if token is valid then it will be redirected to the endpoint
		if token.Valid {
			if r.Method == "OPTION" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
			return
		}

		// returns error if token is not valid
		log.GetLogger().Error(err.Error())
		web.RespondError(w, errors.NewHTTPError("Session expired!!!!", http.StatusUnauthorized))
		return
	})
}
