package security

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/config"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/errors"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/log"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/web"
	"gorm.io/gorm"
)

// Authentication Provide Method AuthUser.
type Authentication struct {
	db                      *gorm.DB
	Config                  config.ConfReader
	authorizationTypeBearer string
}

// NewAuthentication returns new instance of Authentication
func NewAuthentication(db *gorm.DB, config config.ConfReader) *Authentication {
	return &Authentication{
		db:                      db,
		Config:                  config,
		authorizationTypeBearer: "bearer",
	}
}

// Middleware will fetch jwt token from header and verify it.
func (auth *Authentication) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authorizationHeader := ctx.GetHeader("Authorization")
		if len(authorizationHeader) == 0 {
			log.GetLogger().Error("authorization header not found")
			web.RespondErrorMessage(ctx, http.StatusUnauthorized, "authorization header not found")
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			log.GetLogger().Error("invalid authorization header provided")
			web.RespondErrorMessage(ctx, http.StatusUnauthorized, "invalid authorization header provided")
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != auth.authorizationTypeBearer {
			log.GetLogger().Error(fmt.Sprintf("unsupported authorization type %s", authorizationType))
			web.RespondErrorMessage(ctx, http.StatusUnauthorized, fmt.Sprintf("unsupported authorization type %s", authorizationType))
			return
		}

		accessToken := fields[1]
		claims := jwt.RegisteredClaims{}

		payload, err := jwt.ParseWithClaims(
			accessToken, &claims,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
				}
				return []byte(auth.Config.GetString(config.JWTKey)), nil
			})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				log.GetLogger().Error(err.Error())
				web.RespondErrorMessage(ctx, http.StatusUnauthorized, err.Error())
				return
			}

			log.GetLogger().Error(err)
			web.RespondErrorMessage(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		if claims.VerifyExpiresAt(time.Now(), true) {
			log.GetLogger().Error("Session expired! Please login again")
			web.RespondErrorMessage(ctx, http.StatusUnauthorized, "Session expired! Please login again")
			return
		}

		// prints all the claims
		// for key, val := range claims {
		// 	fmt.Printf("Key: %v, value: %v\n", key, val)
		// }

		// if token is valid then it will be redirected to the endpoint
		if payload.Valid {
			if ctx.Request.Method == "OPTION" {
				// w.WriteHeader(http.StatusOK)
				ctx.Writer.WriteHeader(http.StatusOK)
				return
			}

			ctx.Next()
			return
		}
	}
}

// Middleware will fetch jwt token from header and verify it.
// func (auth *Authentication) Middleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		// tokenStr, err := request.HeaderExtractor{"Authorization"}.ExtractToken(r)

// 		tokenStr := r.Header.Get("Authorization")
// 		if len(tokenStr) == 0 {
// 			log.GetLogger().Error("Token must be specified")
// 			web.RespondError(ctx, errors.NewHTTPError("Token must be specified", http.StatusUnauthorized))
// 			return
// 		}

// claims := jwt.RegisteredClaims{}

// token, err := jwt.ParseWithClaims(
// 	tokenStr, &claims,
// 	func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.NewHTTPError(errors.ErrorCodeInternalError, http.StatusInternalServerError)
// 		}
// 		return []byte(auth.Config.GetString(config.JWTKey)), nil
// 	})

// if err != nil {
// 	if err == jwt.ErrSignatureInvalid {
// 		log.GetLogger().Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusUnauthorized))
// 		return
// 	}

// 	log.GetLogger().Error(err.Error())
// 	web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
// 	return
// }

// if claims.VerifyExpiresAt(time.Now(), true) {
// 	log.GetLogger().Error("Session expired! Please login again")
// 	web.RespondError(w, errors.NewHTTPError("Session expired! Please login again", http.StatusUnauthorized))
// 	return
// }

// 		// prints all the claims
// 		// for key, val := range claims {
// 		// 	fmt.Printf("Key: %v, value: %v\n", key, val)
// 		// }

// 		// if token is valid then it will be redirected to the endpoint
// 		if token.Valid {
// 			if r.Method == "OPTION" {
// 				w.WriteHeader(http.StatusOK)
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		// returns error if token is not valid
// 		log.GetLogger().Error(err.Error())
// 		web.RespondError(w, errors.NewHTTPError("Session expired!!!!", http.StatusUnauthorized))
// 		return
// 	})
// }
