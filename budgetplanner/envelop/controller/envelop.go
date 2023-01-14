package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/envelop/service"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/log"
	envelopModel "github.com/shaileshhb/budget-planner-go/budgetplanner/models/envelop"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/security"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/web"
)

// EnvelopController service provides methods to update, delete, add, get method for EnvelopController.
type EnvelopController interface {
	RegisterRoutes(router *gin.RouterGroup)
	addEnvelop(ctx *gin.Context)
	updateEnvelop(ctx *gin.Context)
	deleteEnvelop(ctx *gin.Context)
	getEnvelops(ctx *gin.Context)
}

// envelopController.
type envelopController struct {
	service service.EnvelopService
	log     log.Logger
	auth    *security.Authentication
}

// NewEnvelopController create new EnvelopController
func NewEnvelopController(ser service.EnvelopService, log log.Logger,
	auth *security.Authentication) EnvelopController {
	return &envelopController{
		service: ser,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes will register routes for authentication controller.
func (c *envelopController) RegisterRoutes(router *gin.RouterGroup) {

	guarded := router.Group("/users", c.auth.Middleware())

	guarded.POST("/:userID/envelops", c.addEnvelop)
	guarded.PUT("/:userID/envelops/:envelopID", c.updateEnvelop)
	guarded.DELETE("/:userID/envelops/:envelopID", c.deleteEnvelop)
	guarded.GET("/:userID/envelops", c.getEnvelops)
}

// addEnvelop will add new envelop for specified user.
func (c *envelopController) addEnvelop(ctx *gin.Context) {

	envelop := envelopModel.Envelop{}
	parser := web.NewParser(ctx)

	err := web.UnmarshalJSON(ctx.Request, &envelop)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	envelop.UserID, err = parser.GetUUID("userID")
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	envelop.CreatedBy = envelop.UserID

	err = envelop.Validate()
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.AddEnvelop(&envelop)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(ctx, http.StatusCreated, nil)
}

// updateEnvelop will update specified envelop.
func (c *envelopController) updateEnvelop(ctx *gin.Context) {

	envelop := envelopModel.Envelop{}
	parser := web.NewParser(ctx)

	err := web.UnmarshalJSON(ctx.Request, &envelop)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	envelop.UserID, err = parser.GetUUID("userID")
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	envelop.ID, err = parser.GetUUID("envelopID")
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	envelop.UpdatedBy = envelop.UserID

	err = envelop.Validate()
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.UpdateEnvelop(&envelop)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(ctx, http.StatusAccepted, nil)
}

// deleteEnvelop will delete specified envelop.
func (c *envelopController) deleteEnvelop(ctx *gin.Context) {

	envelop := envelopModel.Envelop{}
	parser := web.NewParser(ctx)
	var err error

	envelop.UserID, err = parser.GetUUID("userID")
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	envelop.ID, err = parser.GetUUID("envelopID")
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	envelop.DeletedBy = envelop.UserID

	err = c.service.DeleteEnvelop(&envelop)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(ctx, http.StatusAccepted, nil)
}

// getEnvelops will fetch all the envelops for specifed user.
func (c *envelopController) getEnvelops(ctx *gin.Context) {

	var envelops []envelopModel.EnvelopDTO
	parser := web.NewParser(ctx)

	userID, err := parser.GetUUID("userID")
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.GetEnvelops(&envelops, userID)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(ctx, http.StatusOK, envelops)
}
