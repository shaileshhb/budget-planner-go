package web

import (
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shaileshhb/budget-planner-go/budgetplanner/errors"
)

// Parser helps in parsing the data from the URL params.
type Parser struct {
	Params gin.Params
	Form   url.Values
}

// NewParser will call request.ParseForm() and create a new instance of parser.
func NewParser(ctx *gin.Context) *Parser {
	ctx.Request.ParseForm()
	return &Parser{
		Params: ctx.Params,
		Form:   ctx.Request.Form,
	}
}

// GetUUID will get uuid from the given paramName in URL params.
func (p *Parser) GetUUID(paramName string) (uuid.UUID, error) {
	param, ok := p.Params.Get(paramName)
	if !ok {
		return uuid.Nil, errors.NewValidationError(param + " not found")
	}
	return uuid.Parse(param)
}

// GetParameter will get parameter from the given paramName in URL params.
func (p *Parser) GetParameter(paramName string) string {
	paramString, ok := p.Params.Get(paramName)
	if !ok {
		return ""
	}
	return paramString
}

// GetTenantID will get "tenantID" param in URL params.
func (p *Parser) GetTenantID() (uuid.UUID, error) {
	param, ok := p.Params.Get("tenantID")
	if !ok {
		return uuid.Nil, errors.NewValidationError("tenantID not found")
	}
	return uuid.Parse(param)
}

// ParseLimitAndOffset will parse limit and offset from query params.
func (p *Parser) ParseLimitAndOffset() (limit, offset int) {
	limitparam := p.Form.Get("limit")
	offsetparam := p.Form.Get("offset")
	var err error
	limit = 30
	if len(limitparam) > 0 {
		limit, err = strconv.Atoi(limitparam)
		if err != nil {
			return
		}
	}
	if len(offsetparam) > 0 {
		offset, err = strconv.Atoi(offsetparam)
		if err != nil {
			return
		}
	}
	return
}
