package viewer

import (
	"fmt"
	"testing"

	"api.safer.place/viewer/v1"
)

func TestValidateRegion(t *testing.T) {
	testCases := map[*viewer.Region]error{
		// good cases
		{}:                         nil,
		{North: 5335, South: 5334}: nil,
		{North: 5335, South: 5334, West: -630, East: -629}: nil,
		// out of bounds
		{North: -19000}: &RegionError{"north", -19000, errOutOfBounds},
		{North: 19000}:  &RegionError{"north", 19000, errOutOfBounds},
		{South: -19000}: &RegionError{"south", -19000, errOutOfBounds},
		{South: 19000}:  &RegionError{"south", 19000, errOutOfBounds},
		// invalid bounds
		{North: 5334, South: 5335}: &RegionError{"north-south", -1, errInvalidBounds},
		{East: -1}:                 &RegionError{"east-west", -1, errInvalidBounds},
		// region too big
		{North: 5300, South: 5200}: &RegionError{"north-south", 100, errTooBig},
		{North: 5330, South: 5320}: &RegionError{"north-south", 10, errTooBig},
		{North: 5335, South: 5333}: &RegionError{"north-south", 2, errTooBig},
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
