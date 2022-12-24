package module

import (
	"github.com/shaileshhb/budget-planner-go/budgetplanner"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/models/user"
)

// Configure will migrate all the tables.
func Configure(app *budgetplanner.App) {
	userModule := user.NewUserModuleConfig(app.DB)

	app.MigrateTables([]budgetplanner.ModuleConfig{userModule})
}
