package module

import (
	"github.com/shaileshhb/budget-planner-go/budgetplanner"
	envelopcontroller "github.com/shaileshhb/budget-planner-go/budgetplanner/envelop/controller"
	envelopservice "github.com/shaileshhb/budget-planner-go/budgetplanner/envelop/service"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/repository"
)

// registerEnvelopRoutes will register all routes of envelops.
func registerEnvelopRoutes(app *budgetplanner.App, repo repository.Repository) {
	defer app.WG.Done()

	envelopService := envelopservice.NewEnvelopService(app.DB, repo, app.Auth)
	enevlopController := envelopcontroller.NewEnvelopController(envelopService, app.Log, app.Auth)

	app.RegisterControllerRoutes([]budgetplanner.Controller{enevlopController})
}
