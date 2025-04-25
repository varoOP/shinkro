package domain

import (
	"encoding/json"
	"strconv"
	"time"
)

type Tautulli struct {
	Account struct {
		Title string `json:"title"`
	} `json:"Account"`
	Metadata struct {
		GrandparentKey      string        `json:"grandparentKey"`
		GrandparentTitle    string        `json:"grandparentTitle"`
		GUID                GUID          `json:"guid"`
		Index               string        `json:"index"`
		LibrarySectionTitle string        `json:"librarySectionTitle"`
		ParentIndex         string        `json:"parentIndex"`
		Title               string        `json:"title"`
		Type                PlexMediaType `json:"type"`
	} `json:"Metadata"`
	Event PlexEvent `json:"event"`
}

func NewTautulli(b []byte) (*Tautulli, error) {
	t := &Tautulli{}
	err := json.Unmarshal(b, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func ToPlex(b []byte) (*Plex, error) {
	t, err := NewTautulli(b)
	if err != nil {
		return nil, err
	}

	parentIndex, err := strconv.Atoi(t.Metadata.ParentIndex)
	if err != nil {
		return nil, err
	}

	index, err := strconv.Atoi(t.Metadata.Index)
	if err != nil {
		return nil, err
	}

	return &Plex{
		Event:     t.Event,
		Source:    TautulliWebhook,
		TimeStamp: time.Now(),
		Account: struct {
			Id           int    `json:"id"`
			ThumbnailUrl string `json:"thumb"`
			Title        string `json:"title"`
		}{
			Title: t.Account.Title,
		},
		Metadata: Metadata{
			GrandparentKey:      t.Metadata.GrandparentKey,
			GrandparentTitle:    t.Metadata.GrandparentTitle,
			GUID:                t.Metadata.GUID,
			Index:               index,
			LibrarySectionTitle: t.Metadata.LibrarySectionTitle,
			ParentIndex:         parentIndex,
			Title:               t.Metadata.Title,
			Type:                t.Metadata.Type,
		},
	}, nil
}
