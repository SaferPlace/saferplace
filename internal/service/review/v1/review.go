// Copyright 2022 SaferPlace

package review

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"go.uber.org/zap"
	"safer.place/realtime/internal/database"

	"api.safer.place/incident/v1"
	pb "api.safer.place/review/v1"
	connectpb "api.safer.place/review/v1/reviewconnect"
)

// Service is the review service
type Service struct {
	db  database.Database
	log *zap.Logger
}

// Register the review service
func Register(
	db database.Database,
	log *zap.Logger,
) func() (string, http.Handler) {
	return func() (string, http.Handler) {
		return connectpb.NewReviewServiceHandler(&Service{
			db:  db,
			log: log,
		})
	}
}

// ReviewIncident accepts the reviewers comments and resolution
func (s *Service) ReviewIncident(
	ctx context.Context,
	req *connect.Request[pb.ReviewIncidentRequest],
) (
	*connect.Response[pb.ReviewIncidentResponse],
	error,
) {
	s.log.Info("review received",
		zap.String("id", req.Msg.Id),
		zap.String("resolution", req.Msg.Resolution.String()),
	)

	comment := &incident.Comment{
		// TODO: Replace with authenticated user
		AuthorId:  "voy", // TODO: Deduct from the authentication credentials.
		Timestamp: time.Now().Unix(),
		Message:   req.Msg.Comment,
	}

	if err := s.db.SaveReview(
		ctx,
		req.Msg.Id,
		req.Msg.Resolution,
		comment,
	); err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, err)
	}

	return connect.NewResponse(&pb.ReviewIncidentResponse{}), nil
}

// ViewIncident shows the incident information
func (s *Service) ViewIncident(
	ctx context.Context,
	req *connect.Request[pb.ViewIncidentRequest],
) (
	*connect.Response[pb.ViewIncidentResponse],
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
	return connect.NewResponse(&pb.ViewIncidentResponse{
		Incident: inc,
	}), nil
}

// IncidentsWithoutReview shows all the incidents that are not reviewed
func (s *Service) IncidentsWithoutReview(
	ctx context.Context,
	req *connect.Request[pb.IncidentsWithoutReviewRequest],
) (
	*connect.Response[pb.IncidentsWithoutReviewResponse],
	error,
) {
	incidents, err := s.db.IncidentsWithoutReview(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	basicIncidents := make([]*pb.BasicIncidentDetails, 0, len(incidents))
	for _, inc := range incidents {
		basicIncidents = append(basicIncidents, &pb.BasicIncidentDetails{
			Id:          inc.Id,
			Description: inc.Description,
			Timestamp:   inc.Timestamp.Seconds,
		})
	}
	return connect.NewResponse(&pb.IncidentsWithoutReviewResponse{
		Incidents: basicIncidents,
	}), nil
}
