package plex

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type PlexWebhook struct {
	Rating  float32 `json:"rating"`
	Event   string  `json:"event"`
	User    bool    `json:"user"`
	Owner   bool    `json:"owner"`
	Account struct {
		Id           int    `json:"id"`
		ThumbnailUrl string `json:"thumb"`
		Title        string `json:"title"`
	} `json:"Account"`
	Server struct {
		Title string `json:"title"`
		UUID  string `json:"uuid"`
	} `json:"Server"`
	Player struct {
		Local         bool   `json:"local"`
		PublicAddress string `json:"publicAddress"`
		Title         string `json:"title"`
		UUID          string `json:"uuid"`
	} `json:"Player"`
	Metadata Metadata `json:"Metadata"`
}

func NewPlexWebhook(payload []byte) (*PlexWebhook, error) {
	p := &PlexWebhook{}
	err := json.Unmarshal(payload, p)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal plex payload")
	}

	return p, nil
}
