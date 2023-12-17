package tautulli

import (
	"encoding/json"
	"strconv"

	"github.com/varoOP/shinkro/internal/domain"
)

type Tautulli struct {
	Account struct {
		Title string `json:"title"`
	} `json:"Account"`
	Metadata struct {
		GrandparentKey      string      `json:"grandparentKey"`
		GrandparentTitle    string      `json:"grandparentTitle"`
		GUID                domain.GUID `json:"guid"`
		Index               string      `json:"index"`
		LibrarySectionTitle string      `json:"librarySectionTitle"`
		ParentIndex         string      `json:"parentIndex"`
		Title               string      `json:"title"`
		Type                string      `json:"type"`
	} `json:"Metadata"`
	Event string `json:"event"`
}

func NewTautulli(b []byte) (*Tautulli, error) {
	t := &Tautulli{}
	err := json.Unmarshal(b, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func ToPlex(b []byte) (*domain.Plex, error) {
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

	return &domain.Plex{
		Event: t.Event,
		Account: struct {
			Id           int    `json:"id"`
			ThumbnailUrl string `json:"thumb"`
			Title        string `json:"title"`
		}{
			Title: t.Account.Title,
		},
		Metadata: domain.Metadata{
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
