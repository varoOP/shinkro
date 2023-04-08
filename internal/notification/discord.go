package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	ColorCompleted = 40704
	ColorWatching  = 49087
)

type Discord struct {
	Webhook DiscordWebhook
	Url     string
}

type DiscordWebhook struct {
	Embeds []Embeds `json:"embeds"`
}

type Embeds struct {
	Title  string   `json:"title"`
	URL    string   `json:"url"`
	Color  int      `json:"color"`
	Fields []Fields `json:"fields"`
	Image  Image    `json:"image"`
}

type Fields struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Image struct {
	URL string `json:"url"`
}

func NewDicord(url string) *Discord {
	return &Discord{
		Url: url,
	}
}

func (d *Discord) SendNotification(ctx context.Context, content map[string]string) error {
	d.buildWebhook(content)
	p, err := json.Marshal(d.Webhook)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, d.Url, "application/json", bytes.NewBuffer(p))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New("something went wrong with sending discord notification")
	}

	return nil
}

func (d *Discord) buildWebhook(content map[string]string) {
	var (
		title    = content["title"]
		url      = content["url"]
		status   = content["status"]
		imageUrl = content["image_url"]
	)

	color := ColorWatching
	if status == "completed" {
		color = ColorCompleted
	}

	d.Webhook = DiscordWebhook{
		Embeds: []Embeds{
			{
				Title:  title,
				URL:    url,
				Color:  color,
				Fields: buildFields(content),
				Image: Image{
					URL: imageUrl,
				},
			},
		},
	}
}

func buildFields(content map[string]string) []Fields {
	var (
		event          = content["event"]
		status         = content["status"]
		score          = content["score"]
		startDate      = content["start_date"]
		finishDate     = content["finish_date"]
		totalEps       = content["total_eps"]
		watchedEps     = content["watched_eps"]
		timesRewatched = content["times_rewatched"]
		f              []Fields
	)

	if event == "media.rate" {
		event = "Update Score"
	}

	if event == "media.scrobble" {
		event = "Update Status"
	}

	f = append(f, Fields{
		Name:   "Event",
		Value:  event,
		Inline: true,
	})

	f = append(f, Fields{
		Name:   "Status",
		Value:  cases.Title(language.Und).String(status),
		Inline: true,
	})

	if totalEps == "0" {
		totalEps = "?"
	}

	f = append(f, Fields{
		Name:   "Episodes Seen",
		Value:  fmt.Sprintf("%v / %v", watchedEps, totalEps),
		Inline: true,
	})

	if score == "0" {
		score = "Not Scored"
	} else {
		score = fmt.Sprintf("%v / %v", score, "10")
	}

	f = append(f, Fields{
		Name:   "Score",
		Value:  score,
		Inline: true,
	})

	if startDate != "" {
		f = append(f, Fields{
			Name:   "Start Date",
			Value:  startDate,
			Inline: true,
		})
	}

	if finishDate != "" {
		f = append(f, Fields{
			Name:   "Finish Date",
			Value:  finishDate,
			Inline: true,
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
