package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/varoOP/shinkuro/internal/domain"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	ColorCompleted = 40704
	ColorWatching  = 49087
	ColorError     = 12517376
)

type Discord struct {
	Webhook DiscordWebhook
	Url     string
	log     *zerolog.Logger
}

type DiscordWebhook struct {
	Embeds []Embeds `json:"embeds"`
}

type Embeds struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Color       int      `json:"color"`
	Description string   `json:"description"`
	Fields      []Fields `json:"fields"`
	Image       Image    `json:"image"`
}

type Fields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Image struct {
	URL string `json:"url"`
}

func NewDicord(url string, n *domain.NotificationPayload, log *zerolog.Logger) *Discord {
	d := &Discord{
		Url: url,
		log: log,
	}

	d.buildWebhook(n)
	return d
}

func (d *Discord) SendNotification(ctx context.Context) error {
	p, err := json.Marshal(d.Webhook)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.Url, bytes.NewBuffer(p))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		d.log.Trace().RawJSON("discordResponse", body).Msg("discord response dump")
		return errors.New("something went wrong with sending discord notification")
	}

	d.log.Trace().Msg("sent discord notification")
	return nil
}

func (d *Discord) buildWebhook(n *domain.NotificationPayload) {
	var (
		title    = n.Title
		url      = n.Url
		status   string
		imageUrl = n.ImageUrl
		color    int
		fields   []Fields
	)

	switch status = n.Status; status {
	case "watching":
		color = ColorWatching
	case "completed":
		color = ColorCompleted
	case "":
		color = ColorError
	}

	if n.Message == "" {
		fields = buildFields(n)
	}

	d.Webhook = DiscordWebhook{
		Embeds: []Embeds{
			{
				Title:       title,
				URL:         url,
				Color:       color,
				Description: n.Message,
				Fields:      fields,
				Image: Image{
					URL: imageUrl,
				},
			},
		},
	}
}

func buildFields(n *domain.NotificationPayload) []Fields {
	var (
		event          = n.Event
		status         = n.Status
		score          = strconv.Itoa(n.Score)
		startDate      = n.StartDate
		finishDate     = n.FinishDate
		totalEps       = strconv.Itoa(n.TotalEps)
		watchedEps     = strconv.Itoa(n.WatchedEps)
		timesRewatched = strconv.Itoa(n.TimesRewatched)
		f              []Fields
	)

	switch event {
	case "media.rate":
		event = "Update Score"
	case "media.scrobble":
		event = "Update Status"
	}

	f = append(f, Fields{
		Name:   "Event",
		Value:  event,
		Inline: false,
	})

	f = append(f, Fields{
		Name:   "Status",
		Value:  cases.Title(language.Und).String(status),
		Inline: false,
	})

	if totalEps == "0" {
		totalEps = "?"
	}

	f = append(f, Fields{
		Name:   "Episodes Seen",
		Value:  fmt.Sprintf("%v / %v", watchedEps, totalEps),
		Inline: false,
	})

	if score == "0" {
		score = "Not Scored"
	} else {
		score = fmt.Sprintf("%v / %v", score, "10")
	}

	f = append(f, Fields{
		Name:   "Score",
		Value:  score,
		Inline: false,
	})

	if startDate != "" {
		f = append(f, Fields{
			Name:   "Start Date",
			Value:  startDate,
			Inline: false,
		})
	}

	if finishDate != "" {
		f = append(f, Fields{
			Name:   "Finish Date",
			Value:  finishDate,
			Inline: false,
		})
	}

	if timesRewatched != "0" {
		f = append(f, Fields{
			Name:   "Times Rewatched",
			Value:  timesRewatched,
			Inline: false,
		})
	}

	return f
}
