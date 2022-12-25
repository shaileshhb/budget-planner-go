package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/errors"
)

// RespondJSON Make response with json formate.
func RespondJSON(ctx *gin.Context, code int, payload interface{}) {
	// response, err := json.Marshal(payload)
	// if err != nil {
	// 	// w.WriteHeader(http.StatusInternalServerError)
	// 	// w.Write([]byte(error.Error()))
	// 	ctx.JSON(http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(code)
	// w.Write([]byte(response))
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(code, payload)
}

// RespondJSONWithXTotalCount Make response with json format and add X-Total-Count header.
func RespondJSONWithXTotalCount(ctx *gin.Context, code int, count int, payload interface{}) {
	// response, err := json.Marshal(payload)
	// if err != nil {
	// 	// w.WriteHeader(http.StatusInternalServerError)
	// 	// w.Write([]byte(error.Error()))
	// 	ctx.JSON(http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// w.Header().Set("Content-Type", "application/json")
	// SetNewHeader(w, "X-Total-Count", strconv.Itoa(count))
	// w.WriteHeader(code)
	// w.Write([]byte(response))

	SetNewHeader(ctx, "X-Total-Count", strconv.Itoa(count))
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(code, payload)
}

// RespondErrorMessage make error response with payload.
func RespondErrorMessage(ctx *gin.Context, code int, msg string) {
	// RespondJSON(ctx, code, map[string]string{"error": msg})
	ctx.AbortWithStatusJSON(code, msg)
}

// RespondError check error type and Write to ResponseWriter.
func RespondError(ctx *gin.Context, err error) {
	switch err.(type) {
	case *errors.ValidationError:
		RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
	case *errors.HTTPError:
		httpError := err.(*errors.HTTPError)
		RespondErrorMessage(ctx, httpError.HTTPStatus, httpError.ErrorKey)
	default:
		RespondErrorMessage(ctx, http.StatusInternalServerError, err.Error())
	}
}

// SetNewHeader will expose and set the given headerName and value
//
//	SetNewHeader(w,"total","10") will set header "total" : "10"
func SetNewHeader(ctx *gin.Context, headerName, value string) {
	// w.Header().Add("Access-Control-Expose-Headers", headerName)
	// w.Header().Set(headerName, value)

	// ctx.Header("Access-Control-Expose-Headers", headerName)
	ctx.Header(headerName, value)
}
