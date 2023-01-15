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

// TransactionController service provides methods to update, delete, add, get method for TransactionController.
type TransactionController interface {
	RegisterRoutes(router *gin.RouterGroup)
	addTransaction(ctx *gin.Context)
	updateTransaction(ctx *gin.Context)
	deleteTransaction(ctx *gin.Context)
	getUserTransaction(ctx *gin.Context)
}

// transactionController.
type transactionController struct {
	service service.TransactionService
	log     log.Logger
	auth    *security.Authentication
}

// NewTransactionController create new TransactionController
func NewTransactionController(ser service.TransactionService, log log.Logger,
	auth *security.Authentication) TransactionController {
	return &transactionController{
		service: ser,
		log:     log,
		auth:    auth,
	}
}

// RegisterRoutes will register routes for authentication controller.
func (c *transactionController) RegisterRoutes(router *gin.RouterGroup) {

	guarded := router.Group("/users", c.auth.Middleware())

	guarded.POST("/:userID/transactions", c.addTransaction)
	guarded.PUT("/:userID/transactions/:transactionID", c.updateTransaction)
	guarded.DELETE("/:userID/transactions/:transactionID", c.deleteTransaction)
	guarded.GET("/:userID/transactions", c.getUserTransaction)
}

// addTransaction will add new transaction for user.
func (c *transactionController) addTransaction(ctx *gin.Context) {

	transaction := envelopModel.Transaction{}
	parser := web.NewParser(ctx)

	err := web.UnmarshalJSON(ctx.Request, &transaction)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	transaction.UserID, err = parser.GetUUID("userID")
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// transaction.CreatedBy = transaction.UserID

	err = transaction.Validate()
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.AddTransaction(&transaction)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(ctx, http.StatusCreated, nil)
}

// updateTransaction will update specified transaction of user.
func (c *transactionController) updateTransaction(ctx *gin.Context) {

	transaction := envelopModel.Transaction{}
	parser := web.NewParser(ctx)

	err := web.UnmarshalJSON(ctx.Request, &transaction)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	transaction.UserID, err = parser.GetUUID("userID")
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	transaction.ID, err = parser.GetUUID("transactionID")
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// transaction.UpdatedBy = transaction.UserID

	err = transaction.Validate()
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.UpdateTransaction(&transaction)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(ctx, http.StatusAccepted, nil)
}

// deleteTransaction will delete specified transaction of user.
func (c *transactionController) deleteTransaction(ctx *gin.Context) {

	transaction := envelopModel.Transaction{}
	parser := web.NewParser(ctx)
	var err error

	transaction.UserID, err = parser.GetUUID("userID")
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	transaction.ID, err = parser.GetUUID("transactionID")
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// transaction.DeletedBy = transaction.UserID

	err = c.service.DeleteTransaction(&transaction)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSON(ctx, http.StatusAccepted, nil)
}

// getUserTransaction will fetch transactions of user.
func (c *transactionController) getUserTransaction(ctx *gin.Context) {

	transactions := []envelopModel.TransactionDTO{}
	parser := web.NewParser(ctx)

	userID, err := parser.GetUUID("userID")
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var totalCount int64

	err = c.service.GetUserTransaction(&transactions, userID, &totalCount, parser)
	if err != nil {
		c.log.Error(err)
		web.RespondErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	web.RespondJSONWithXTotalCount(ctx, http.StatusOK, int(totalCount), transactions)
}
