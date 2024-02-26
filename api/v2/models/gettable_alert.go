// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// GettableAlert gettable alert
//
// swagger:model gettableAlert
type GettableAlert struct {

	// annotations
	// Required: true
	Annotations LabelSet `json:"annotations"`

	// ends at
	// Required: true
	// Format: date-time
	EndsAt *strfmt.DateTime `json:"endsAt"`

	// fingerprint
	// Required: true
	Fingerprint *string `json:"fingerprint"`

	// receivers
	// Required: true
	Receivers []*Receiver `json:"receivers"`

	// starts at
	// Required: true
	// Format: date-time
	StartsAt *strfmt.DateTime `json:"startsAt"`

	// status
	// Required: true
	Status *AlertStatus `json:"status"`

	// updated at
	// Required: true
	// Format: date-time
	UpdatedAt *strfmt.DateTime `json:"updatedAt"`

	Alert
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *GettableAlert) UnmarshalJSON(raw []byte) error {
	// AO0
	var dataAO0 struct {
		Annotations LabelSet `json:"annotations"`

		EndsAt *strfmt.DateTime `json:"endsAt"`

		Fingerprint *string `json:"fingerprint"`

		Receivers []*Receiver `json:"receivers"`

		StartsAt *strfmt.DateTime `json:"startsAt"`

		Status *AlertStatus `json:"status"`

		UpdatedAt *strfmt.DateTime `json:"updatedAt"`
	}
	if err := swag.ReadJSON(raw, &dataAO0); err != nil {
		return err
	}

	m.Annotations = dataAO0.Annotations

	m.EndsAt = dataAO0.EndsAt

	m.Fingerprint = dataAO0.Fingerprint

	m.Receivers = dataAO0.Receivers

	m.StartsAt = dataAO0.StartsAt

	m.Status = dataAO0.Status

	m.UpdatedAt = dataAO0.UpdatedAt

	// AO1
	var aO1 Alert
	if err := swag.ReadJSON(raw, &aO1); err != nil {
		return err
	}
	m.Alert = aO1

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m GettableAlert) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	var dataAO0 struct {
		Annotations LabelSet `json:"annotations"`

		EndsAt *strfmt.DateTime `json:"endsAt"`

		Fingerprint *string `json:"fingerprint"`

		Receivers []*Receiver `json:"receivers"`

		StartsAt *strfmt.DateTime `json:"startsAt"`

		Status *AlertStatus `json:"status"`

		UpdatedAt *strfmt.DateTime `json:"updatedAt"`
	}

	dataAO0.Annotations = m.Annotations

	dataAO0.EndsAt = m.EndsAt

	dataAO0.Fingerprint = m.Fingerprint

	dataAO0.Receivers = m.Receivers

	dataAO0.StartsAt = m.StartsAt

	dataAO0.Status = m.Status

	dataAO0.UpdatedAt = m.UpdatedAt

	jsonDataAO0, errAO0 := swag.WriteJSON(dataAO0)
	if errAO0 != nil {
		return nil, errAO0
	}
	_parts = append(_parts, jsonDataAO0)

	aO1, err := swag.WriteJSON(m.Alert)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this gettable alert
func (m *GettableAlert) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAnnotations(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEndsAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFingerprint(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateReceivers(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStartsAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdatedAt(formats); err != nil {
		res = append(res, err)
	}

	// validation for a type composition with Alert
	if err := m.Alert.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GettableAlert) validateAnnotations(formats strfmt.Registry) error {

	if err := validate.Required("annotations", "body", m.Annotations); err != nil {
		return err
	}

	if m.Annotations != nil {
		if err := m.Annotations.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("annotations")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("annotations")
			}
			return err
		}
	}

	return nil
}

func (m *GettableAlert) validateEndsAt(formats strfmt.Registry) error {

	if err := validate.Required("endsAt", "body", m.EndsAt); err != nil {
		return err
	}

	if err := validate.FormatOf("endsAt", "body", "date-time", m.EndsAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *GettableAlert) validateFingerprint(formats strfmt.Registry) error {

	if err := validate.Required("fingerprint", "body", m.Fingerprint); err != nil {
		return err
	}

	return nil
}

func (m *GettableAlert) validateReceivers(formats strfmt.Registry) error {

	if err := validate.Required("receivers", "body", m.Receivers); err != nil {
		return err
	}

	for i := 0; i < len(m.Receivers); i++ {
		if swag.IsZero(m.Receivers[i]) { // not required
			continue
		}

		if m.Receivers[i] != nil {
			if err := m.Receivers[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("receivers" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("receivers" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *GettableAlert) validateStartsAt(formats strfmt.Registry) error {

	if err := validate.Required("startsAt", "body", m.StartsAt); err != nil {
		return err
	}

	if err := validate.FormatOf("startsAt", "body", "date-time", m.StartsAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *GettableAlert) validateStatus(formats strfmt.Registry) error {

	if err := validate.Required("status", "body", m.Status); err != nil {
		return err
	}

	if m.Status != nil {
		if err := m.Status.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("status")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("status")
			}
			return err
		}
	}

	return nil
}

func (m *GettableAlert) validateUpdatedAt(formats strfmt.Registry) error {

	if err := validate.Required("updatedAt", "body", m.UpdatedAt); err != nil {
		return err
	}

	if err := validate.FormatOf("updatedAt", "body", "date-time", m.UpdatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this gettable alert based on the context it is used
func (m *GettableAlert) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAnnotations(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateReceivers(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateStatus(ctx, formats); err != nil {
		res = append(res, err)
	}

	// validation for a type composition with Alert
	if err := m.Alert.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GettableAlert) contextValidateAnnotations(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Annotations.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("annotations")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("annotations")
		}
		return err
	}

	return nil
}

func (m *GettableAlert) contextValidateReceivers(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Receivers); i++ {

		if m.Receivers[i] != nil {
			if err := m.Receivers[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("receivers" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("receivers" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *GettableAlert) contextValidateStatus(ctx context.Context, formats strfmt.Registry) error {

	if m.Status != nil {
		if err := m.Status.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("status")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("status")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *GettableAlert) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GettableAlert) UnmarshalBinary(b []byte) error {
	var res GettableAlert
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
