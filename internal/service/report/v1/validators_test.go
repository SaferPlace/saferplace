// Copyright 2023 SaferPlace

package report

import (
	"errors"
	"testing"

	"api.safer.place/incident/v1"
)

func TestValidators(t *testing.T) {
	testCases := map[string]struct {
		inc *incident.Incident
		fn  ValidatorFunc
		err error
	}{
		"missing description": {
			inc: &incident.Incident{
				Description: "",
			},
			fn:  validateDescription,
			err: errMissingDescription,
		},
		"invalid lattitude": {
			inc: &incident.Incident{
				Coordinates: &incident.Coordinates{
					Lat: -200,
				},
			},
			fn:  validateCoordinates,
			err: CoordinateError{-90, 90},
		},
		"invalid longitude": {
			inc: &incident.Incident{
				Coordinates: &incident.Coordinates{
					Lon: 200,
				},
			},
			fn:  validateCoordinates,
			err: CoordinateError{-180, 180},
		},
		"empty coordinates on transportation": {
			inc: &incident.Incident{
				Location: incident.Location_LOCATION_TRANSPORTATION,
			},
			fn:  validateCoordinates,
			err: nil,
		},
		"coordinates on transportation invalid longitude": {
			inc: &incident.Incident{
				Location: incident.Location_LOCATION_TRANSPORTATION,
				Coordinates: &incident.Coordinates{
					Lon: 200,
				},
			},
			fn:  validateCoordinates,
			err: CoordinateError{-180, 180},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if err := tc.fn(tc.inc); !errors.Is(err, tc.err) {
				t.Errorf("%T(%v) = %v, want %v", tc.fn, tc.inc, err, tc.err)
			}
		})
	}
}
