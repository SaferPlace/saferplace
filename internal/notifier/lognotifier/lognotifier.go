package lognotifier

import (
	"context"
	"fmt"
	"log/slog"

	"api.safer.place/incident/v1"
	"safer.place/internal/log"
)

type Notifier struct {
	log log.Logger
}

func New(log log.Logger) *Notifier {
	return &Notifier{log: log}
}

func (n *Notifier) Notify(ctx context.Context, inc *incident.Incident) error {
	n.log.Info(ctx, "incident for review",
		slog.String("url", fmt.Sprintf("https://review.safer.place/incident/%s", inc.Id)),
	)
	return nil
}
