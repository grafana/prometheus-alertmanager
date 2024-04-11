// Code generated by go-swagger; DO NOT EDIT.

package alert

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetAlertsHandlerFunc turns a function with the right signature into a get alerts handler
type GetAlertsHandlerFunc func(GetAlertsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetAlertsHandlerFunc) Handle(params GetAlertsParams) middleware.Responder {
	return fn(params)
}

// GetAlertsHandler interface for that can handle valid get alerts params
type GetAlertsHandler interface {
	Handle(GetAlertsParams) middleware.Responder
}

// NewGetAlerts creates a new http.Handler for the get alerts operation
func NewGetAlerts(ctx *middleware.Context, handler GetAlertsHandler) *GetAlerts {
	return &GetAlerts{Context: ctx, Handler: handler}
}

/*
	GetAlerts swagger:route GET /alerts alert getAlerts

Get a list of alerts
*/
type GetAlerts struct {
	Context *middleware.Context
	Handler GetAlertsHandler
}

func (o *GetAlerts) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetAlertsParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
