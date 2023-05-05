package discordnotifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"

	"api.safer.place/incident/v1"
)

// Notifier sends a notification to discord about an incident.
type Notifier struct {
	client   *http.Client
	endpoint string
}

// New creates a new discord notifier
func New(c *http.Client) (*Notifier, error) {
	var cfg Config
	if err := envconfig.Process("discord", &cfg); err != nil {
		return nil, fmt.Errorf("unable to process discord settings: %w", err)
	}

	return &Notifier{
		client:   c,
		endpoint: cfg.Endpoint,
	}, nil
}

// Notify sends the discord webhook notification
func (n *Notifier) Notify(ctx context.Context, i *incident.Incident) error {
	msg := fmt.Sprintf(messageFmt,
		i.Id, i.Coordinates.Lat, i.Coordinates.Lon, i.Description, i.Id)

	data := discordgo.WebhookParams{
		Content:    msg,
		Components: []discordgo.MessageComponent{
			// discordgo.Button{
			// 	Label: "Review Incident",
			// 	// Style: discordgo.LinkButton,
			// 	// TODO: Maybe this can be customized?
			// 	URL: "https://review.safer.place/incident/" + i.Id,
			// },
		},
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(data); err != nil {
		return fmt.Errorf("unable to encode webhook body: %w", err)
	}

	log.Println(body)

	req, err := http.NewRequest(http.MethodPost, n.endpoint, body)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to send discord notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response %q: %s", string(respBody), resp.Status)
	}

	return nil
}

// Config used to parse the notifier confiuration
type Config struct {
	Endpoint string `required:"true"`
}

var messageFmt = `
New Incident for review: %s
Lat: %.6f
Lon: %.6f
Description: %s

REVIEW: https://review.safer.place/incident/%s
`
