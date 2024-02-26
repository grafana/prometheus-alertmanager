// Code generated by go-swagger; DO NOT EDIT.

package silence

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/prometheus/alertmanager/api/v2/models"
)

// GetSilenceOKCode is the HTTP code returned for type GetSilenceOK
const GetSilenceOKCode int = 200

/*
GetSilenceOK Get silence response

swagger:response getSilenceOK
*/
type GetSilenceOK struct {

	/*
	  In: Body
	*/
	Payload *models.GettableSilence `json:"body,omitempty"`
}

// NewGetSilenceOK creates GetSilenceOK with default headers values
func NewGetSilenceOK() *GetSilenceOK {

	return &GetSilenceOK{}
}

// WithPayload adds the payload to the get silence o k response
func (o *GetSilenceOK) WithPayload(payload *models.GettableSilence) *GetSilenceOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get silence o k response
func (o *GetSilenceOK) SetPayload(payload *models.GettableSilence) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSilenceOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetSilenceNotFoundCode is the HTTP code returned for type GetSilenceNotFound
const GetSilenceNotFoundCode int = 404

/*
GetSilenceNotFound A silence with the specified ID was not found

swagger:response getSilenceNotFound
*/
type GetSilenceNotFound struct {
}

// NewGetSilenceNotFound creates GetSilenceNotFound with default headers values
func NewGetSilenceNotFound() *GetSilenceNotFound {

	return &GetSilenceNotFound{}
}

// WriteResponse to the client
func (o *GetSilenceNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// GetSilenceInternalServerErrorCode is the HTTP code returned for type GetSilenceInternalServerError
const GetSilenceInternalServerErrorCode int = 500

/*
GetSilenceInternalServerError Internal server error

swagger:response getSilenceInternalServerError
*/
type GetSilenceInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewGetSilenceInternalServerError creates GetSilenceInternalServerError with default headers values
func NewGetSilenceInternalServerError() *GetSilenceInternalServerError {

	return &GetSilenceInternalServerError{}
}

// WithPayload adds the payload to the get silence internal server error response
func (o *GetSilenceInternalServerError) WithPayload(payload string) *GetSilenceInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get silence internal server error response
func (o *GetSilenceInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSilenceInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
