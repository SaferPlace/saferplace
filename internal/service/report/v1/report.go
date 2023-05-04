// Copyright 2022 SaferPlace

package report

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"safer.place/realtime/internal/queue"

	ipb "api.safer.place/incident/v1"
	pb "api.safer.place/report/v1"
	connectpb "api.safer.place/report/v1/reportconnect"
)

// Service is the report service
type Service struct {
	queue queue.Producer[*ipb.Incident]
	log   *zap.Logger
}

// Register creates a new service and and returns the
func Register(q queue.Producer[*ipb.Incident], log *zap.Logger) func() (string, http.Handler) {
	return func() (string, http.Handler) {
		return connectpb.NewReportServiceHandler(&Service{
			queue: q,
			log:   log,
		})
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
	incident.Id = uuid.New().String()
	incident.Timestamp = time.Now().Unix()

	if err := validateIncident(incident); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	s.log.Info("received report",
		zap.String("id", incident.Id),
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

var (
	errMissingDescription = errors.New("missing description")
)

func validateIncident(i *ipb.Incident) error {
	// Do not accept incidents in the pass, but still allow for slow network
	// connections and times being set incorrectly
	if i.Description == "" {
		return errMissingDescription
	}
	// TODO: Only accept incidents in Ireland
	if !(-90 <= i.Lat && i.Lat <= 90) {
		return fmt.Errorf("lattitude %w", CoordinateError{-90, 90})
	}
	if !(-180 <= i.Lon && i.Lon <= 180) {
		return fmt.Errorf("longitude %w", CoordinateError{-180, 180})
	}
	return nil
}
