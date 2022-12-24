package module

import (
	"github.com/shaileshhb/budget-planner-go/budgetplanner"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/repository"
	usercontroller "github.com/shaileshhb/budget-planner-go/budgetplanner/user/controller"
	userservice "github.com/shaileshhb/budget-planner-go/budgetplanner/user/service"
)

func registerUserRoutes(app *budgetplanner.App, repo repository.Repository) {
	defer app.WG.Done()

	authService := userservice.NewAuthenticationService(app.DB, repo, app.Auth)
	authController := usercontroller.NewAuthenticationController(authService, app.Log, app.Auth)

	app.RegisterControllerRoutes([]budgetplanner.Controller{authController})
}
