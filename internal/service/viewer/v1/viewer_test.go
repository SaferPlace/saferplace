package viewer

import (
	"fmt"
	"testing"

	"api.safer.place/viewer/v1"
)

func TestValidateRegion(t *testing.T) {
	testCases := map[*viewer.Region]error{
		// good cases
		{}:                           nil,
		{North: 53.35, South: 53.34}: nil,
		{North: 53.35, South: 53.34, West: -6.30, East: -6.29}: nil,
		// out of bounds
		{North: -190}: &RegionError{"north", -190, errOutOfBounds},
		{North: 190}:  &RegionError{"north", 190, errOutOfBounds},
		{South: -190}: &RegionError{"south", -190, errOutOfBounds},
		{South: 190}:  &RegionError{"south", 190, errOutOfBounds},
		// invalid bounds
		{North: 53.34, South: 53.35}: &RegionError{"north-south", -0.01, errInvalidBounds},
		{East: -0.01}:                &RegionError{"east-west", -0.01, errInvalidBounds},
		// not in increments
		{North: 1.1, South: 1.1}:     nil,
		{North: 1.11, South: 1.11}:   nil,
		{North: 1.111, South: 1.111}: &RegionError{"north", 1.111, errNotInIncrements},
		// region too big
		{North: 53, South: 52}:       &RegionError{"north-south", 1, errTooBig},
		{North: 53.3, South: 53.2}:   &RegionError{"north-south", 0.1, errTooBig},
		{North: 53.35, South: 53.33}: &RegionError{"north-south", 0.02, errTooBig},
	}

	for in, want := range testCases {
		t.Run(fmt.Sprintf("%+v", in), func(t *testing.T) {
			// TODO: Fix the plain string error checking since errors.Is doesn't seem to be working

			if got := validateRegion(in); fmt.Sprint(got) != fmt.Sprint(want) {
				t.Errorf("validateRegion(%v) = %v; want %v", in, got, want)
			}
		})
	}
}
