// Copyright 2022 SaferPlace

package report

import (
	"errors"
	"fmt"

	"api.safer.place/incident/v1"
)

var (
	errMissingDescription = errors.New("missing description")
	errMissingCoordinates = errors.New("missing coordinates")
)

type ValidatorFunc func(i *incident.Incident) error

type Validator interface {
	Validate(i *incident.Incident) error
}

type MultiValidator struct {
	validators []ValidatorFunc
}

func NewMultiValidator(validators ...ValidatorFunc) *MultiValidator {
	return &MultiValidator{validators: validators}
}

// Validate the incident based on the list of validators
func (v *MultiValidator) Validate(i *incident.Incident) error {
	var errs []error

	for _, fn := range v.validators {
		errs = append(errs, fn(i))
	}

	return errors.Join(errs...)
}

func validateDescription(i *incident.Incident) error {
	if i.Description == "" {
		return errMissingDescription
	}

	return nil
}

// validateCoordinates checks are the coordinates valid, but it
func validateCoordinates(i *incident.Incident) error {
	// If the incident happened on a mode of transportation, we can ignore empty
	// coordinates, otherwise still ensure they are valid
	if i.Location == incident.Location_LOCATION_TRANSPORTATION && i.Coordinates == nil {
		return nil
	}

	if i.Coordinates == nil {
		return errMissingCoordinates
	}

	// TODO: Only accept incidents in Ireland
	if !(-90 <= i.Coordinates.Lat && i.Coordinates.Lat <= 90) {
		return fmt.Errorf("lattitude %w", CoordinateError{-90, 90})
	}
	if !(-180 <= i.Coordinates.Lon && i.Coordinates.Lon <= 180) {
		return fmt.Errorf("longitude %w", CoordinateError{-180, 180})
	}

	return nil
}
