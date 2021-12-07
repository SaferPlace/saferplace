// Package data contains all the files located in the data directory
package data

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

func init() {
	if err := parseRoughPrefixCoordinates(); err != nil {
		panic(err)
	}
}

func parseRoughPrefixCoordinates() error {
	records, err := csv.NewReader(
		strings.NewReader(roughPrefixCoordinatesFile)).
		ReadAll()
	if err != nil {
		return fmt.Errorf("unable to read rough prefix coordinates file: %w",
			err)
	}

	roughPrefixCoordinates = make(map[string]PrefixCoordinates, len(records))

	for _, line := range records {
		x, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return fmt.Errorf("unable to parse X coordinate: %w", err)
		}
		y, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			return fmt.Errorf("unable to parse Y coordinate: %w", err)
		}

		roughPrefixCoordinates[line[0]] = PrefixCoordinates{
			Names: strings.Split(line[1], "|"),
			X:     x,
			Y:     y,
		}
	}
	return nil
}

//go:embed rough_prefix_coordinates.csv
var roughPrefixCoordinatesFile string
var roughPrefixCoordinates map[string]PrefixCoordinates

type PrefixCoordinates struct {
	Names []string
	X, Y  float64
}

func RoughPrefixCoordinates() map[string]PrefixCoordinates {
	return roughPrefixCoordinates
}
