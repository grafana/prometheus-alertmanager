// Code generated by go-swagger; DO NOT EDIT.

package alert

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// PostAlertsOKCode is the HTTP code returned for type PostAlertsOK
const PostAlertsOKCode int = 200

/*
PostAlertsOK Create alerts response

swagger:response postAlertsOK
*/
type PostAlertsOK struct {
}

// NewPostAlertsOK creates PostAlertsOK with default headers values
func NewPostAlertsOK() *PostAlertsOK {

	return &PostAlertsOK{}
}

// WriteResponse to the client
func (o *PostAlertsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// PostAlertsBadRequestCode is the HTTP code returned for type PostAlertsBadRequest
const PostAlertsBadRequestCode int = 400

/*
PostAlertsBadRequest Bad request

swagger:response postAlertsBadRequest
*/
type PostAlertsBadRequest struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewPostAlertsBadRequest creates PostAlertsBadRequest with default headers values
func NewPostAlertsBadRequest() *PostAlertsBadRequest {

	return &PostAlertsBadRequest{}
}

// WithPayload adds the payload to the post alerts bad request response
func (o *PostAlertsBadRequest) WithPayload(payload string) *PostAlertsBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post alerts bad request response
func (o *PostAlertsBadRequest) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostAlertsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// PostAlertsInternalServerErrorCode is the HTTP code returned for type PostAlertsInternalServerError
const PostAlertsInternalServerErrorCode int = 500

/*
PostAlertsInternalServerError Internal server error

swagger:response postAlertsInternalServerError
*/
type PostAlertsInternalServerError struct {

	/*
	  In: Body
	*/
	Payload string `json:"body,omitempty"`
}

// NewPostAlertsInternalServerError creates PostAlertsInternalServerError with default headers values
func NewPostAlertsInternalServerError() *PostAlertsInternalServerError {

	return &PostAlertsInternalServerError{}
}

// WithPayload adds the payload to the post alerts internal server error response
func (o *PostAlertsInternalServerError) WithPayload(payload string) *PostAlertsInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post alerts internal server error response
func (o *PostAlertsInternalServerError) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostAlertsInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
