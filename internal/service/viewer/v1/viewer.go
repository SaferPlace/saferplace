// Copyright 2023 SaferPlace

package viewer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"

	"api.safer.place/viewer/v1"
	"api.safer.place/viewer/v1/viewerconnect"
	"safer.place/internal/database"
	"safer.place/internal/log"
	"safer.place/internal/service"
)

const (
	// RegionIncrements specifies how granular the requests for region can be. Too broad
	// and we can't correcty cache, too granular and they are breaching our location privacy.
	RegionIncrements = 1 // ~1.11km at the equator
)

// Service is the viewer service
type Service struct {
	db  database.Database
	log log.Logger
}

// Register the viewer service
func Register(
	db database.Database,
	log log.Logger,
) service.Service {
	return func(interceptors ...connect.Interceptor) (string, http.Handler) {
		return viewerconnect.NewViewerServiceHandler(
			&Service{
				db:  db,
				log: log,
			},
			connect.WithInterceptors(interceptors...),
		)
	}
}

// ViewInRadius gets incidents in the specified radius
// Deprecated: Use [ViewInRegion] instead as respects privacy better.
func (s *Service) ViewInRadius(
	ctx context.Context,
	req *connect.Request[viewer.ViewInRadiusRequest],
) (
	*connect.Response[viewer.ViewInRadiusResponse],
	error,
) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("deprecated"))
}

// ViewInRegion shows all incidents in the specified region.
func (s *Service) ViewInRegion(
	ctx context.Context,
	req *connect.Request[viewer.ViewInRegionRequest],
) (
	*connect.Response[viewer.ViewInRegionResponse],
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

	s.log.Info(ctx, "viewing incidents in region",
		slog.Any("region", req.Msg.Region),
		slog.Time("since", since),
	)

	inc, err := s.db.IncidentsInRegion(ctx, since, req.Msg.Region)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, err)
	}

	return connect.NewResponse(&viewer.ViewInRegionResponse{
		Incidents: inc,
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
	s.log.Info(ctx, "view incident",
		slog.String("id", req.Msg.Id),
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

	s.log.Info(ctx, "viewing alerting incidents",
		slog.Any("region", req.Msg.Region),
		slog.Time("since", since),
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
	errOutOfBounds   = errors.New("not a valid earth coordinate")
	errInvalidBounds = errors.New("invalid bounds")
	errTooBig        = errors.New("region is too big")
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
	if -9000 > region.North || region.North > 9000 {
		return &RegionError{"north", region.North, errOutOfBounds}
	}
	if -9000 > region.South || region.South > 9000 {
		return &RegionError{"south", region.South, errOutOfBounds}
	}
	if region.North < region.South {
		return &RegionError{"north-south", region.North - region.South, errInvalidBounds}
	}

	// east and west
	if 18000 > region.East && region.East > 18000 {
		return &RegionError{"east", region.East, errOutOfBounds}
	}
	if -9000 > region.West && region.West > 9000 {
		return &RegionError{"west", region.West, errOutOfBounds}
	}
	if region.East < region.West {
		return &RegionError{"east-west", region.East - region.West, errInvalidBounds}
	}

	// validate that they are at most RegionDegreesIncrement apart
	if diff := region.North - region.South; diff > RegionIncrements {
		return &RegionError{"north-south", diff, errTooBig}
	}
	if diff := region.East - region.West; diff > RegionIncrements {
		return &RegionError{"east-west", diff, errTooBig}
	}

	return nil
}
