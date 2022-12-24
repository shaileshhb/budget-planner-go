package web

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Parser helps in parsing the data from the URL params.
type Parser struct {
	Params map[string]string
	Form   url.Values
}

// NewParser will call request.ParseForm() and create a new instance of parser.
func NewParser(r *http.Request) *Parser {
	r.ParseForm()
	return &Parser{
		Params: mux.Vars(r),
		Form:   r.Form,
	}
}

// GetUUID will get uuid from the given paramName in URL params.
func (p *Parser) GetUUID(paramName string) (uuid.UUID, error) {
	return uuid.Parse(p.Params[paramName])
	// if err != nil {
	// 	return uuid.Nil, err
	// }
	// return id, nil
}

// GetParameter will get parameter from the given paramName in URL params.
func (p *Parser) GetParameter(paramName string) string {
	paramString := p.Params[paramName]
	return paramString
}

// GetTenantID will get "tenantID" param in URL params.
func (p *Parser) GetTenantID() (uuid.UUID, error) {
	return uuid.Parse(p.Params["tenantID"])

	// idString := p.Params["tenantID"]
	// id, err := util.ParseUUID(idString)
	// if err != nil {
	// 	return uuid.Nil, err
	// }
	// return id, nil
}

// ParseLimitAndOffset will parse limit and offset from query params.
func (p *Parser) ParseLimitAndOffset() (limit, offset int) {
	limitparam := p.Form.Get("limit")
	offsetparam := p.Form.Get("offset")
	var err error
	limit = 5
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
