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

// AuthenticationController service provides methods to update, delete, add, get method for AuthenticationController.
type AuthenticationController interface {
	RegisterRoutes(router *gin.RouterGroup)
	register(ctx *gin.Context)
	login(ctx *gin.Context)
	getUser(ctx *gin.Context)
}

// authenticationController.
type authenticationController struct {
	service service.AuthenticationService
	log     log.Logger
	auth    *security.Authentication
}

// NewAuthenticationController create new AuthenticationController
func NewAuthenticationController(ser service.AuthenticationService, log log.Logger,
	auth *security.Authentication) AuthenticationController {
	return &authenticationController{
		service: ser,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes will register routes for authentication controller.
func (c *authenticationController) RegisterRoutes(router *gin.RouterGroup) {

	// unguarded := router.Group("")
	router.POST("/register", c.register)
	router.POST("/login", c.login)

	guarded := router.Group("/users", c.auth.Middleware())
	guarded.PUT("/:userID", c.updateUser)
	guarded.GET("/:userID", c.getUser)
}

// register will register new user in the system.
func (c *authenticationController) register(ctx *gin.Context) {
	// parser := web.NewParser(ctx)
	user := userModal.User{}
	auth := userModal.Authentication{}

	err := web.UnmarshalJSON(ctx.Request, &user)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = user.ValidateRegistration()
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.Register(&user, &auth)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(ctx, http.StatusOK, auth)
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

// updateUser will update user details.
func (c *authenticationController) updateUser(ctx *gin.Context) {
	user := userModal.User{}
	parser := web.NewParser(ctx)

	err := web.UnmarshalJSON(ctx.Request, &user)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user.ID, err = parser.GetUUID("userID")
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user.UpdatedBy, err = c.auth.ExtractUserID(ctx)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = user.ValidateUser()
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.UpdateUser(&user)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(ctx, http.StatusAccepted, nil)
}

// getUser will fetch specified user details.
func (c *authenticationController) getUser(ctx *gin.Context) {
	user := userModal.UserDTO{}
	parser := web.NewParser(ctx)
	var err error

	user.ID, err = parser.GetUUID("userID")
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.GetUser(&user)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(ctx, http.StatusOK, user)
}
