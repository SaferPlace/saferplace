// Copyright 2023 SaferPlace

package viewer

import (
	"context"
	"errors"
	"net/http"

	"api.safer.place/viewer/v1"
	"api.safer.place/viewer/v1/viewerconnect"
	"connectrpc.com/connect"
	"go.uber.org/zap"
	"safer.place/realtime/internal/database"
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
