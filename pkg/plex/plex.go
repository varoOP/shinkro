package plex

import (
	"encoding/json"
)

type PlexWebhook struct {
	Rating  float32 `json:"rating,omitempty"`
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
		LibrarySectionType    string  `json:"librarySectionType,omitempty"`
		RatingKey             string  `json:"ratingKey,omitempty"`
		Key                   string  `json:"key,omitempty"`
		SkipParent            bool    `json:"skipParent,omitempty"`
		ParentRatingKey       string  `json:"parentRatingKey,omitempty"`
		GrandparentRatingKey  string  `json:"grandparentRatingKey,omitempty"`
		GUID                  string  `json:"guid,omitempty"`
		ParentGUID            string  `json:"parentGuid,omitempty"`
		GrandparentGUID       string  `json:"grandparentGuid,omitempty"`
		Type                  string  `json:"type,omitempty"`
		Title                 string  `json:"title,omitempty"`
		GrandparentKey        string  `json:"grandparentKey,omitempty"`
		ParentKey             string  `json:"parentKey,omitempty"`
		LibrarySectionTitle   string  `json:"librarySectionTitle,omitempty"`
		LibrarySectionID      int     `json:"librarySectionID,omitempty"`
		LibrarySectionKey     string  `json:"librarySectionKey,omitempty"`
		GrandparentTitle      string  `json:"grandparentTitle,omitempty"`
		ParentTitle           string  `json:"parentTitle,omitempty"`
		ContentRating         string  `json:"contentRating,omitempty"`
		Summary               string  `json:"summary,omitempty"`
		Index                 int     `json:"index,omitempty"`
		ParentIndex           int     `json:"parentIndex,omitempty"`
		Rating                float64 `json:"rating,omitempty"`
		Year                  int     `json:"year,omitempty"`
		Thumb                 string  `json:"thumb,omitempty"`
		Art                   string  `json:"art,omitempty"`
		GrandparentThumb      string  `json:"grandparentThumb,omitempty"`
		GrandparentArt        string  `json:"grandparentArt,omitempty"`
		OriginallyAvailableAt string  `json:"originallyAvailableAt,omitempty"`
		AddedAt               int     `json:"addedAt,omitempty"`
		UpdatedAt             int     `json:"updatedAt,omitempty"`
		Director              []struct {
			ID     int    `json:"id,omitempty"`
			Filter string `json:"filter,omitempty"`
			Tag    string `json:"tag,omitempty"`
		} `json:"Director,omitempty"`
		Writer []struct {
			ID     int    `json:"id,omitempty"`
			Filter string `json:"filter,omitempty"`
			Tag    string `json:"tag,omitempty"`
		} `json:"Writer,omitempty"`
		Producer []struct {
			ID     int    `json:"id,omitempty"`
			Filter string `json:"filter,omitempty"`
			Tag    string `json:"tag,omitempty"`
		} `json:"Producer,omitempty"`
	} `json:"Metadata,omitempty"`
}

func NewPlexWebhook(payload string) (*PlexWebhook, error) {
	p := &PlexWebhook{}

	err := json.Unmarshal([]byte(payload), p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
