package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/log"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/security"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/user/service"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/web"
)

type AuthenticationController interface {
	RegisterRoutes(router *gin.Engine)
	register(ctx *gin.Context)
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
	router.GET("/register", c.register)

	c.log.Info("User auth routes registered.")
}

func (c *authenticationController) register(ctx *gin.Context) {
	fmt.Println("================ register called ================")
	web.RespondJSON(ctx, http.StatusOK, "hello world!!")
}
