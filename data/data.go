// Package data contains all the files located in the data directory
package data

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

var initFuncs = []func() error{
	parseGardaLocations,
	parseGardaCrimes,
}

func init() {
	for _, f := range initFuncs {
		if err := f(); err != nil {
			panic(err)
		}
	}
}

// Coordinates define the X and Y location, X being the lattitude and Y being
// the longitude.
type Coordinates struct {
	X, Y float64
}

// NamedCoordinates are a set of coordinates with list of names associated with
// it.
type NamedCoordinates struct {
	Names []string
	Coordinates
}

//go:embed garda_locations.csv
var gardaLocationsFile string
var gardaLocations map[string]Coordinates

func parseGardaLocations() error {
	records, err := csv.NewReader(
		strings.NewReader(gardaLocationsFile)).
		ReadAll()
	if err != nil {
		return fmt.Errorf("unable to read garda station location file: %w", err)
	}

	gardaLocations = make(map[string]Coordinates)

	// Remove the first line as its just the field headers
	records = records[1:]
	for _, line := range records {
		x, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return fmt.Errorf("unable to parse X coordinate: %w", err)
		}
		y, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return fmt.Errorf("unable to parse Y coordinate: %w", err)
		}

		gardaLocations[line[0]] = Coordinates{X: x, Y: y}
	}

	return nil
}

// GardaLocations contains a list of names of the Garda stations and where they
// are
func GardaLocations() map[string]Coordinates {
	return gardaLocations
}

//go:embed garda_crimes.csv
var gardaCrimesFile string
var gardaCrimes map[string]CrimesPerType

// CrimeType determines
type CrimeType string

var (
	// MurderAndAssaultCrime 2004 - 2016
	MurderAndAssaultCrime CrimeType = "Attempts or threats to murder, assaults, " +
		"harassments and related offences"
	// DangerousActsCrime 2003 - 2016
	DangerousActsCrime CrimeType = "Dangerous or negligent acts"
	// KidnappingCrime 2003-2016
	KidnappingCrime CrimeType = "Kidnapping and related offences"
	// RobberyCrime 2003-2016
	RobberyCrime CrimeType = "Robbery, extortion and hijacking offences"
	// BurglaryCrime 2003 - 2016
	BurglaryCrime CrimeType = "Burglary and related offences"
	// TheftCrime 2003 - 2016
	TheftCrime CrimeType = "Theft and related offences"
	//FraudCrime 2003-2016
	FraudCrime CrimeType = "Fraud, deception and related offences"
	// DrugCrime 2003-2016
	DrugCrime CrimeType = "Controlled drug offences"
	// WeaponsCrime 2003-2016
	WeaponsCrime CrimeType = "Weapons and Explosives Offences 2003"
	// PropertyDamageCrime 2003-2016
	PropertyDamageCrime CrimeType = "Damage to property and to the environment"
	// PublicOrderCrime 2003 - 2016
	PublicOrderCrime CrimeType = "Public order and other social code offences"
	// OrganisedCrime 2003 - 2016
	OrganisedCrime CrimeType = "Offences against government, justice procedures " +
		"and organisation of crime 2003"
)

// CrimesPerType maps the crime type and the number of crimes commited in the
// given year.
type CrimesPerType map[CrimeType]CrimesInYear

// CrimesInYear contains a mapping from a year to the number of crimes committed
// in that year.
type CrimesInYear map[int]int

func parseGardaCrimes() error {
	records, err := csv.NewReader(
		strings.NewReader(gardaCrimesFile)).
		ReadAll()
	if err != nil {
		return fmt.Errorf("unable to read garda station location file: %w", err)
	}

	gardaCrimes = make(map[string]CrimesPerType)

	// remove the csv header
	records = records[1:]
	for _, line := range records {
		station := line[1]

		vs := matoi(line[5:])
		gardaCrimes[station] = CrimesPerType{
			MurderAndAssaultCrime: mapToYear(vs[0:13], 2004),
			DangerousActsCrime:    mapToYear(vs[13:27], 2003),
			KidnappingCrime:       mapToYear(vs[27:41], 2003),
			RobberyCrime:          mapToYear(vs[41:55], 2003),
			BurglaryCrime:         mapToYear(vs[55:69], 2003),
			TheftCrime:            mapToYear(vs[69:83], 2003),
			FraudCrime:            mapToYear(vs[83:97], 2003),
			DrugCrime:             mapToYear(vs[97:111], 2003),
			WeaponsCrime:          mapToYear(vs[111:125], 2003),
			PropertyDamageCrime:   mapToYear(vs[125:139], 2003),
			PublicOrderCrime:      mapToYear(vs[139:153], 2003),
			OrganisedCrime:        mapToYear(vs[153:167], 2003),
		}
	}

	return nil
}

// CrimesPerStation defines the number of crimes that were reported by the
// given station.
func CrimesPerStation() map[string]CrimesPerType {
	return gardaCrimes
}

// matoi is multi atoi
func matoi(in []string) []int {
	out := make([]int, len(in))
	for i, in := range in {
		v, err := strconv.Atoi(in)
		if err != nil {
			panic(err) // we can just panic here
		}
		out[i] = v
	}
	return out
}

func mapToYear(in []int, startYear int) CrimesInYear {
	res := make(CrimesInYear, len(in))
	for i, v := range in {
		res[startYear+i] = v
	}
	return res
}
