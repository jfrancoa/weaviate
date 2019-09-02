//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2019 SeMI Holding B.V. (registered @ Dutch Chamber of Commerce no 75221632). All rights reserved.
//  LICENSE WEAVIATE OPEN SOURCE: https://www.semi.technology/playbook/playbook/contract-weaviate-OSS.html
//  LICENSE WEAVIATE ENTERPRISE: https://www.semi.technology/playbook/contract-weaviate-enterprise.html
//  CONCEPT: Bob van Luijt (@bobvanluijt)
//  CONTACT: hello@semi.technology
//

// Code generated by go-swagger; DO NOT EDIT.

package actions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/semi-technologies/weaviate/entities/models"
)

// ActionsReferencesUpdateHandlerFunc turns a function with the right signature into a actions references update handler
type ActionsReferencesUpdateHandlerFunc func(ActionsReferencesUpdateParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn ActionsReferencesUpdateHandlerFunc) Handle(params ActionsReferencesUpdateParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// ActionsReferencesUpdateHandler interface for that can handle valid actions references update params
type ActionsReferencesUpdateHandler interface {
	Handle(ActionsReferencesUpdateParams, *models.Principal) middleware.Responder
}

// NewActionsReferencesUpdate creates a new http.Handler for the actions references update operation
func NewActionsReferencesUpdate(ctx *middleware.Context, handler ActionsReferencesUpdateHandler) *ActionsReferencesUpdate {
	return &ActionsReferencesUpdate{Context: ctx, Handler: handler}
}

/*ActionsReferencesUpdate swagger:route PUT /actions/{id}/references/{propertyName} actions actionsReferencesUpdate

Replace all references to a class-property.

Replace all references to a class-property.

*/
type ActionsReferencesUpdate struct {
	Context *middleware.Context
	Handler ActionsReferencesUpdateHandler
}

func (o *ActionsReferencesUpdate) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewActionsReferencesUpdateParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
