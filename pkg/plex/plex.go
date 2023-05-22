package plex

import (
	"encoding/json"
	"fmt"
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
	Metadata struct {
		GUID                GUID   `json:"guid"`
		Type                string `json:"type"`
		Title               string `json:"title"`
		GrandparentTitle    string `json:"grandparentTitle"`
		LibrarySectionTitle string `json:"librarySectionTitle"`
	} `json:"Metadata"`
}

type GUID struct {
	GUIDS []struct {
		ID string `json:"id"`
	}

	GUID string
}

func (g *GUID) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a string first
	if err := json.Unmarshal(data, &g.GUID); err == nil {
		return nil
	}

	// If it's not a string, try to unmarshal as an anonymous slice of struct
	if err := json.Unmarshal(data, &g.GUIDS); err == nil {
		return nil
	}

	return fmt.Errorf("guid: cannot unmarshal %q", data)
}

func NewPlexWebhook(payload []byte) (*PlexWebhook, error) {
	p := &PlexWebhook{}
	err := json.Unmarshal(payload, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
