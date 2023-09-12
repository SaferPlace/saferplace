package review

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"api.safer.place/incident/v1"

	"safer.place/internal/database"
	"safer.place/internal/log"
	"safer.place/internal/notifier"
	"safer.place/internal/queue"
)

// Review is a big wrapper around incoming reviews
type Review struct {
	incoming       queue.Consumer[*incident.Incident]
	reviewNotifier notifier.Notifier
	db             database.Database

	log log.Logger
}

// New review handler
func New(
	log log.Logger,
	incoming queue.Consumer[*incident.Incident],
	db database.Database,
	reviewNotifier notifier.Notifier,
) *Review {
	return &Review{
		log:            log,
		reviewNotifier: reviewNotifier,
		db:             db,
		incoming:       incoming,
	}
}

// Run the review process
func (r *Review) Run(ctx context.Context) error {
	r.log.Info(ctx, "listening for incoming reviews")
	for {
		if err := r.handleIncoming(ctx); err != nil {
			r.log.Error(ctx, "incident handling failed",
				log.Error(err),
			)
			// TODO: For debugging. We typically don't actually want to shut
			//       down.
			return err
		}
	}
}

func (r *Review) handleIncoming(ctx context.Context) (err error) {
	msg, err := r.incoming.Consume(ctx)
	if err != nil {
		return fmt.Errorf("unable to receive: %w", err)
	}
	defer func() {
		if err != nil {
			r.log.Debug(ctx, "nacking incident",
				log.Error(err),
			)
			msg.Nack()
			return
		}
		r.log.Debug(ctx, "acking incident")
		msg.Ack()
	}()

	inc := msg.Body()

	// Save to database, proceed on if already exists. This means something
	// went wrong and it got requeued.
	if err := r.db.SaveIncident(ctx, inc); err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			r.log.Info(ctx, "incident already exists",
				slog.String("id", inc.Id),
			)
			return nil
		}
		return fmt.Errorf("unable to save incident: %w", err)
	}

	// Notify about incoming review
	if err := r.reviewNotifier.Notify(ctx, inc); err != nil {
		return fmt.Errorf("unable to notify about incoming review: %w", err)
	}

	return nil
}
