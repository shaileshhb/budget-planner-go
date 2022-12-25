package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/log"
	userModal "github.com/shaileshhb/budget-planner-go/budgetplanner/models/user"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/security"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/user/service"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/web"
)

type AuthenticationController interface {
	RegisterRoutes(router *gin.Engine)
	register(ctx *gin.Context)
	login(ctx *gin.Context)
}

type authenticationController struct {
	service service.AuthenticationService
	log     log.Logger
	auth    *security.Authentication
}

func NewAuthenticationController(ser service.AuthenticationService, log log.Logger, auth *security.Authentication) AuthenticationController {
	return &authenticationController{
		service: ser,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes will register routes for authentication controller.
func (c *authenticationController) RegisterRoutes(router *gin.Engine) {
	router.POST("/register", c.register)
	router.POST("/login", c.login)

	c.log.Info("User auth routes registered.")
}

// register will register new user in the system.
func (c *authenticationController) register(ctx *gin.Context) {
	// parser := web.NewParser(ctx)
	user := userModal.User{}

	err := web.UnmarshalJSON(ctx.Request, &user)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = user.Validate()
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.Register(&user)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(ctx, http.StatusAccepted, nil)
}

// login will verify user details and login into the system
func (c *authenticationController) login(ctx *gin.Context) {
	login := userModal.Login{}
	auth := userModal.Authentication{}

	err := web.UnmarshalJSON(ctx.Request, &login)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = login.Validate()
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.Login(&login, &auth)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(ctx, http.StatusOK, auth)
}
