package report

import (
	"errors"
	"testing"

	incident "api.safer.place/incident/v1"
)

func TestValidateReport(t *testing.T) {
	testCases := map[*incident.Incident]error{
		{
			Description: "",
		}: errMissingDescription,
		{
			Description: "invalid lattitude",
			Lat:         -200,
		}: CoordinateError{-90, 90},
		{
			Description: "invalid longitude",
			Lon:         200,
		}: CoordinateError{-180, 180},
	}

	for in, want := range testCases {
		t.Run(in.String(), func(t *testing.T) {
			if err := validateIncident(in); !errors.Is(err, want) {
				t.Errorf("validateIncident() = %v, want %v", err, want)
			}
		})
	}
}
