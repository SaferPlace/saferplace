package database

import (
	"context"
	"errors"
	"time"

	"api.safer.place/incident/v1"
	"api.safer.place/viewer/v1"

	// Register accepted sql databases
	_ "github.com/mattn/go-sqlite3"
)

var (
	// ErrAlreadyExists is returned when we try to add the incident but it
	// already exists
	ErrAlreadyExists = errors.New("database: already exists")
	// ErrDoesNotExist is returned when we try to update a record but it
	// does not exist.
	ErrDoesNotExist = errors.New("database: doesn't exist")
)

// Database defines the interface that a database needs to implement to be
// used. It is primarly designed to be write heavy.
type Database interface {
	SaveIncident(context.Context, *incident.Incident) error
	SaveReview(context.Context, string, incident.Resolution, *incident.Comment) error
	ViewIncident(context.Context, string) (*incident.Incident, error)
	IncidentsWithoutReview(context.Context) ([]*incident.Incident, error)
	IncidentsInRadius(context.Context, *incident.Coordinates, float64) ([]*incident.Incident, error)
	IncidentsInRegion(context.Context, time.Time, *viewer.Region) ([]*incident.Incident, error)
	SaveSession(context.Context, string) error
	IsValidSession(context.Context, string) error
	AlertingIncidents(context.Context, time.Time, *viewer.Region) ([]*incident.Incident, error)
}
