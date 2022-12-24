package module

import "github.com/shaileshhb/budget-planner-go/budgetplanner"

func Configure(app *budgetplanner.App) {

	app.MigrateTables([]budgetplanner.ModuleConfig{})
}
