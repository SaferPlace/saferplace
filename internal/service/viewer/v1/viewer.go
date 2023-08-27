// Copyright 2023 SaferPlace

package viewer

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"api.safer.place/viewer/v1"
	"api.safer.place/viewer/v1/viewerconnect"
	"connectrpc.com/connect"
	"go.uber.org/zap"
	"safer.place/internal/database"
)

const (
	// RegionDegreesMin specifies how granular the requests for region can be. Too broad
	// and we can't correcty cache, too granular and they are breaching our location privacy.
	RegionDegreesIncrement = 0.01 // ~1.11km at the equator
)

// Service is the viewer service
type Service struct {
	db  database.Database
	log *zap.Logger
}

// Register the viewer service
func Register(
	db database.Database,
	log *zap.Logger,
) func() (string, http.Handler) {
	return func() (string, http.Handler) {
		return viewerconnect.NewViewerServiceHandler(&Service{
			db:  db,
			log: log,
		})
	}
}

func (s *Service) ViewInRadius(
	ctx context.Context,
	req *connect.Request[viewer.ViewInRadiusRequest],
) (
	*connect.Response[viewer.ViewInRadiusResponse],
	error,
) {
	s.log.Info("getting incidents in radius",
		zap.Float64("radius", req.Msg.Radius),
		// Lattitude and Longitude logs omitted on purpose to avoid
		// PII (location data) in logs.
	)

	if req.Msg.Center == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			errors.New("missing center"),
		)
	}

	incidents, err := s.db.IncidentsInRadius(ctx, req.Msg.Center, req.Msg.Radius)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, err)
	}

	return connect.NewResponse(&viewer.ViewInRadiusResponse{
		Incidents: incidents,
	}), nil
}

// ViewIncident shows the incident information
func (s *Service) ViewIncident(
	ctx context.Context,
	req *connect.Request[viewer.ViewIncidentRequest],
) (
	*connect.Response[viewer.ViewIncidentResponse],
	error,
) {
	s.log.Info("view incident",
		zap.String("id", req.Msg.Id),
	)

	inc, err := s.db.ViewIncident(ctx, req.Msg.Id)
	if err != nil {
		if errors.Is(err, database.ErrDoesNotExist) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&viewer.ViewIncidentResponse{
		Incident: inc,
	}), nil
}

// ViewAlerting shows all incidents alerting in the provided area. This is so that we can ensure
// privacy without collecting too much information.
func (s *Service) ViewAlerting(
	ctx context.Context,
	req *connect.Request[viewer.ViewAlertingRequest],
) (
	*connect.Response[viewer.ViewAlertingResponse],
	error,
) {

	if err := validateRegion(req.Msg.Region); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("invalid region: %w", err),
		)
	}

	// Default to using one week in the past.
	// TODO: Allow to be configurable using configuration.
	// TODO: Handle alerting incidents long in the past.
	since := req.Msg.Since.AsTime()
	if since.Unix() == 0 {
		since = time.Now().Add(-7 * 24 * time.Hour)
	}

	s.log.Info("viewing alerting incidents",
		zap.Any("region", req.Msg.Region),
		zap.String("since", since.String()),
	)

	inc, err := s.db.AlertingIncidents(ctx, since, req.Msg.Region)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, err)
	}

	return connect.NewResponse(&viewer.ViewAlertingResponse{
		Incidents: inc,
	}), nil
}

var (
	errOutOfBounds     = errors.New("not a valid earth coordinate")
	errNotInIncrements = errors.New("not in specified increments")
	errInvalidBounds   = errors.New("invalid bounds")
	errTooBig          = errors.New("region is too big")
)

// RegionError describes errors which are caused by invalid regions.
type RegionError struct {
	direction string
	value     float64
	cause     error
}

func (e RegionError) Error() string {
	return fmt.Sprintf("%s (%.4f): %v", e.direction, e.value, e.cause)
}

func (e RegionError) Unwrap() error {
	return e.cause
}

// validateRegion ensures that the region is specified in the correct format:
//   - on the planet Earth
//   - in increments of `RegionDegreesIncrement` (or rounded if slightly inaccurate, up to a 1/10 of
//     the increment)
func validateRegion(region *viewer.Region) error {
	// north and south
	if -90 > region.North || region.North > 90 {
		return &RegionError{"north", region.North, errOutOfBounds}
	}
	if -90 > region.South || region.South > 90 {
		return &RegionError{"south", region.South, errOutOfBounds}
	}
	if region.North < region.South {
		return &RegionError{"north-south", region.North - region.South, errInvalidBounds}
	}

	// east and west
	if 180 > region.East && region.East > 180 {
		return &RegionError{"east", region.East, errOutOfBounds}
	}
	if -90 > region.West && region.West > 90 {
		return &RegionError{"west", region.West, errOutOfBounds}
	}
	if region.East < region.West {
		return &RegionError{"east-west", region.East - region.West, errInvalidBounds}
	}

	// validate each coordinate is within the cardinality, round if within acceptable offset.
	if !isIncrement(region.North) {
		return &RegionError{"north", region.North, errNotInIncrements}
	}
	if !isIncrement(region.South) {
		return &RegionError{"south", region.South, errNotInIncrements}
	}
	if !isIncrement(region.East) {
		return &RegionError{"east", region.East, errNotInIncrements}
	}
	if !isIncrement(region.West) {
		return &RegionError{"west", region.West, errNotInIncrements}
	}

	// validate that they are at most RegionDegreesIncrement apart
	if diff := region.North - region.South; diff > RegionDegreesIncrement {
		return &RegionError{"north-south", diff, errTooBig}
	}
	if diff := region.East - region.West; diff > RegionDegreesIncrement {
		return &RegionError{"east-west", diff, errTooBig}
	}

	return nil
}

// isIncrement validates that the provided lattitude/longitude is within the provided increments.
// TODO: Fix this. This method is stupid inefficent but it's late.
func isIncrement(x float64) bool {
	s := strings.TrimRight(fmt.Sprint(x), "0")
	return len(strings.TrimPrefix(s, fmt.Sprintf("%.0f", x))) == 0 ||
		len(strings.TrimPrefix(s, fmt.Sprintf("%.1f", x))) == 0 ||
		len(strings.TrimPrefix(s, fmt.Sprintf("%.2f", x))) == 0
}
