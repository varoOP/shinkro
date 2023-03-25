package notifications

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	Color_completed = 40704
	Color_watching  = 49087
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

func (d *Discord) SendNotification() error {
	p, err := json.Marshal(d.Webhook)
	if err != nil {
		return err
	}

	_, err = http.Post(d.Url, "application/json", bytes.NewBuffer(p))
	if err != nil {
		return err
	}

	return nil
}
