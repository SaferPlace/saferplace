package surreal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/surrealdb/surrealdb.go"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"

	"api.safer.place/incident/v1"
	"api.safer.place/viewer/v1"
	"safer.place/internal/config/secret"
	"safer.place/internal/database"
	"safer.place/internal/log"
)

// SurrealDB endpoint
type Config struct {
	Endpoint  string        `yaml:"endpoint" default:"ws://localhost:8000/rpc"`
	Namespace string        `yaml:"namespace" default:"saferplace"`
	Database  string        `yaml:"database" default:"saferplace"`
	Username  string        `yaml:"username"`
	Password  secret.Secret `yaml:"password"`
}

type Database struct {
	tracer trace.Tracer
	logger log.Logger
	db     *surrealdb.DB
}

func New(cfg *Config, opts ...Option) (*Database, error) {
	var err error

	db := &Database{}

	for _, opt := range opts {
		opt(db)
	}

	db.db, err = surrealdb.New(cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to SurrealDB Endpoint: %w", err)
	}

	if _, err := db.db.Signin(map[string]any{
		"user": cfg.Username,
		"pass": cfg.Password,
	}); err != nil {
		return nil, fmt.Errorf("unable to sign in: %w", err)
	}

	if _, err = db.db.Use(cfg.Namespace, cfg.Database); err != nil {
		return nil, fmt.Errorf("unable to use namespace/database: %w", err)
	}

	return db, nil
}

func (db *Database) SaveIncident(ctx context.Context, inc *incident.Incident) error {
	ctx, span := db.tracer.Start(ctx, "SaveIncident")
	defer span.End()

	if exists, err := db.hasIncident(ctx, inc.Id); err != nil {
		return fmt.Errorf("unable to check does the incident exist: %w", err)
	} else if exists {
		return database.ErrAlreadyExists
	}

	if _, err := db.db.Create("incident", inc); err != nil {
		return fmt.Errorf("unable to create incident: %w", err)
	}

	db.logger.Info(ctx, "created new incident", slog.String("id", inc.Id))

	return nil
}

func (db *Database) SaveReview(
	ctx context.Context,
	id string,
	res incident.Resolution,
	comment *incident.Comment,
) error {
	ctx, span := db.tracer.Start(ctx, "SaveReview")
	defer span.End()

	if exists, err := db.hasIncident(ctx, id); err != nil {
		return fmt.Errorf("unable to check does the incident exist: %w", err)
	} else if !exists {
		return database.ErrDoesNotExist
	}

	incData, err := db.db.Select("incident:" + id)
	if err != nil {
		return fmt.Errorf("unable to get incident data: %w", err)
	}
	inc := new(incident.Incident)
	if err := surrealdb.Unmarshal(incData, inc); err != nil {
		return fmt.Errorf("unable to unmarshal incident: %w", err)
	}

	inc.Resolution = res
	inc.ReviewerComments = append(inc.ReviewerComments, comment)

	if _, err := db.db.Update(inc.Id, inc); err != nil {
		return fmt.Errorf("unable to update incident with comment: %w", err)
	}

	return nil
}

func (db *Database) ViewIncident(ctx context.Context, id string) (*incident.Incident, error) {
	_, span := db.tracer.Start(ctx, "ViewIncident")
	defer span.End()

	data, err := db.db.Select("incident:" + id)
	if err != nil {
		if errors.Is(err, surrealdb.ErrNoRow) {
			return nil, database.ErrDoesNotExist
		}
		return nil, fmt.Errorf("unable to get incident data: %w", err)
	}
	inc := new(incident.Incident)
	if err := surrealdb.Unmarshal(data, &inc); err != nil {
		return nil, fmt.Errorf("unable to unmarshal incident: %w", err)
	}

	inc.Id = strings.TrimPrefix(inc.Id, "incident:")
	sort.Sort(ByCommentTimestamp(inc.ReviewerComments))

	return inc, nil
}

var incidentsWithoutReviewQuery = `
SELECT * FROM incident WHERE resolution = nil
`

func (db *Database) IncidentsWithoutReview(ctx context.Context) ([]*incident.Incident, error) {
	_, span := db.tracer.Start(ctx, "IncidentsWithoutReview")
	defer span.End()

	results, err := db.db.Query(incidentsWithoutReviewQuery, map[string]any{})
	if err != nil {
		return nil, fmt.Errorf("unable to query for incidents without resolution: %w", err)
	}

	incs, err := surrealdb.SmartUnmarshal[[]*incident.Incident](results, nil)
	for i := range incs {
		incs[i].Id = strings.TrimPrefix(incs[i].Id, "incident:")
	}
	return incs, err
}

var incidentsInRegionQuery = `
SELECT *
FROM incident
WHERE
	array::matches($resolutions, resolution)
AND
	coordinates.lat < $north
AND
	coordinates.lat > $south
AND
	coordinates.lon > $west
AND
	coordinates.lon < $east
`

func (db *Database) IncidentsInRegion(
	ctx context.Context, since time.Time, region *viewer.Region,
) ([]*incident.Incident, error) {
	_, span := db.tracer.Start(ctx, "IncidentsInRegion")
	defer span.End()

	fmt.Println(incidentsInRegionQuery)

	results, err := db.db.Query(incidentsInRegionQuery, map[string]any{
		"resolutions": []incident.Resolution{
			incident.Resolution_RESOLUTION_ACCEPTED,
			incident.Resolution_RESOLUTION_ALERTED,
		},
		"since": timestamppb.New(since), // for type compat
		"north": region.North / 100,
		"south": region.South / 100,
		"west":  region.West / 100,
		"east":  region.East / 100,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to query for incidents without resolution: %w", err)
	}

	incs, err := surrealdb.SmartUnmarshal[[]*incident.Incident](results, nil)
	for i := range incs {
		incs[i].Id = strings.TrimPrefix(incs[i].Id, "incident:")
	}
	return incs, err
}

// TODO: Change the query so it uses different resolutions.
func (db *Database) AlertingIncidents(
	ctx context.Context, since time.Time, region *viewer.Region,
) ([]*incident.Incident, error) {
	_, span := db.tracer.Start(ctx, "AlertingIncidents")
	defer span.End()

	results, err := db.db.Query(incidentsInRegionQuery, map[string]any{
		"resolutions": []incident.Resolution{incident.Resolution_RESOLUTION_ALERTED},
		"since":       timestamppb.New(since), // for type compat
		"north":       region.North / 100,
		"south":       region.South / 100,
		"west":        region.West / 100,
		"east":        region.East / 100,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to query for incidents without resolution: %w", err)
	}

	incs, err := surrealdb.SmartUnmarshal[[]*incident.Incident](results, nil)
	for i := range incs {
		incs[i].Id = strings.TrimPrefix(incs[i].Id, "incident:")
	}
	return incs, err
}

func (db *Database) SaveSession(_ context.Context, _ string) error {
	return errors.New("unsupported")
}

func (db *Database) IsValidSession(_ context.Context, _ string) error {
	return errors.New("unsupported")
}

func (db *Database) hasIncident(ctx context.Context, id string) (bool, error) {
	_, span := db.tracer.Start(ctx, "hasIncident")
	defer span.End()

	data, err := db.db.Select("incident:" + id)
	if errors.Is(err, surrealdb.ErrNoRow) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("unable to check for existance of incident: %w", err)
	}
	inc := new(incident.Incident)
	if err := surrealdb.Unmarshal(data, &inc); err != nil {
		return false, nil
	}

	return true, nil
}
