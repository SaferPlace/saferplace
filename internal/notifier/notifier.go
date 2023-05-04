package notifier

import (
	"context"

	"api.safer.place/incident/v1"
)

// Notifier sends a notification about an incident
type Notifier interface {
	Notify(context.Context, *incident.Incident) error
}
