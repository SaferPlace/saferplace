package database

import (
	"context"
	"errors"
	"time"

	"api.safer.place/incident/v1"
	"api.safer.place/viewer/v1"
)

var (
	ErrDatabaseExists       = errors.New("database already exists")
	ErrDatabaseDoesNotExist = errors.New("databased does not exist")
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
	Review
	Incidents
	Sessions
}

type Review interface {
	SaveIncident(context.Context, *incident.Incident) error
	SaveReview(context.Context, string, incident.Resolution, *incident.Comment) error
	IncidentsWithoutReview(context.Context) ([]*incident.Incident, error)
	ViewIncident(context.Context, string) (*incident.Incident, error)
}

type Incidents interface {
	ViewIncident(context.Context, string) (*incident.Incident, error)
	IncidentsInRegion(context.Context, time.Time, *viewer.Region) ([]*incident.Incident, error)
	AlertingIncidents(context.Context, time.Time, *viewer.Region) ([]*incident.Incident, error)
}

type Sessions interface {
	SaveSession(context.Context, string) error
	IsValidSession(context.Context, string) error
}
