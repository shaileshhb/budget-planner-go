package module

import (
	"github.com/shaileshhb/budget-planner-go/budgetplanner"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/repository"
)

// CreateRouterInstance will create and register the routes of all routers.
func CreateRouterInstance(app *budgetplanner.App, repository repository.Repository) {
	// app.WG.Add(1)
	// go registerCommunityRoutes(app, repository)

	// log := app.Log

	app.InitializeRouter()

	app.WG.Add(2)

	go registerUserRoutes(app, repository)
	go registerEnvelopRoutes(app, repository)

	app.WG.Wait()
}
