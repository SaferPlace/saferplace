// Copyright 2022 SaferPlace

package report

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	ipb "api.safer.place/incident/v1"
	pb "api.safer.place/report/v1"
	connectpb "api.safer.place/report/v1/reportconnect"
	"safer.place/internal/log"
	"safer.place/internal/queue"
	"safer.place/internal/service"
)

// Service is the report service
type Service struct {
	queue queue.Producer[*ipb.Incident]
	log   log.Logger

	validator Validator
}

// Register creates a new service and and returns the
func Register(q queue.Producer[*ipb.Incident], log log.Logger) service.Service {
	return func(interceptors ...connect.Interceptor) (string, http.Handler) {
		return connectpb.NewReportServiceHandler(
			&Service{
				queue: q,
				log:   log,
				validator: NewMultiValidator(
					validateDescription,
					validateCoordinates,
				),
			},
			connect.WithInterceptors(interceptors...),
		)
	}
}

// SendReport receives the report and pushes it to the provided queue
func (s *Service) SendReport(
	ctx context.Context,
	req *connect.Request[pb.SendReportRequest],
) (
	*connect.Response[pb.SendReportResponse],
	error,
) {
	incident := req.Msg.Incident
	// Override the ID no matter what its set to.
	incident.Id = strings.ReplaceAll(uuid.New().String(), "-", "_")
	incident.Timestamp = timestamppb.Now()

	if err := s.validator.Validate(incident); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	s.log.Info(ctx, "received report",
		slog.String("id", incident.Id),
	)

	if err := s.queue.Produce(ctx, incident); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.SendReportResponse{
		Id: incident.Id,
	}), nil
}

// CoordinateError is returned when the provided coordinate does not match the
// max and min
type CoordinateError struct {
	min, max float64
}

func (e CoordinateError) Error() string {
	return fmt.Sprintf("coordinate must be between %.4f and %.4f", e.min, e.max)
}
