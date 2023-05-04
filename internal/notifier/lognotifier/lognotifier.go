package lognotifier

import (
	"context"
	"fmt"

	"api.safer.place/incident/v1"
	"go.uber.org/zap"
)

type Notifier struct {
	log *zap.Logger
}

func New(log *zap.Logger) *Notifier {
	return &Notifier{log: log}
}

func (n *Notifier) Notify(ctx context.Context, inc *incident.Incident) error {
	n.log.Info("incident for review",
		zap.String("url", fmt.Sprintf("https://review.safer.place/incident/%s", inc.Id)),
	)
	return nil
}
