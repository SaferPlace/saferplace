// Copyright 2022 SaferPlace

package review

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/trace"

	"api.safer.place/incident/v1"
	pb "api.safer.place/review/v1"
	connectpb "api.safer.place/review/v1/reviewconnect"
	"safer.place/internal/database"
	"safer.place/internal/log"
	"safer.place/internal/service"
)

// Service is the review service
type Service struct {
	tracer trace.Tracer
	db     database.Review
	log    log.Logger
}

// Register the review service
func Register(
	opts ...Option,
) service.Service {
	return func(interceptors ...connect.Interceptor) (string, http.Handler) {
		s := &Service{}

		for _, opt := range opts {
			opt(s)
		}

		return connectpb.NewReviewServiceHandler(s, connect.WithInterceptors(interceptors...))
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
	s.log.Info(ctx, "review received",
		slog.String("id", req.Msg.Id),
		slog.String("resolution", req.Msg.Resolution.String()),
	)

	comment := &incident.Comment{
		// TODO: Actually perform authentication and authorization and not just blindly accept this.
		AuthorId:  req.Header().Get("email"),
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
	s.log.Debug(ctx, "listing incidents without review")
	incidents, err := s.db.IncidentsWithoutReview(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&pb.IncidentsWithoutReviewResponse{
		Incidents: incidents,
	}), nil
}
