package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type PlexRepo interface {
	Store(ctx context.Context, plex *Plex) error
	FindAll(ctx context.Context) ([]*Plex, error)
	Get(ctx context.Context, req *GetPlexRequest) (*Plex, error)
	
	Delete(ctx context.Context, req *DeletePlexRequest) error
}

type Plex struct {
	ID        int64     `json:"id"`
	Rating    float32   `json:"rating"`
	TimeStamp time.Time `json:"-"`
	Event     string    `json:"event"`
	User      bool      `json:"user"`
	Source    string    `json:"source"`
	Owner     bool      `json:"owner"`
	Account   struct {
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

type Metadata struct {
	RatingGlobal          float32 `json:"rating"`
	RatingKey             string  `json:"ratingKey"`
	Key                   string  `json:"key"`
	ParentRatingKey       string  `json:"parentRatingKey"`
	GrandparentRatingKey  string  `json:"grandparentRatingKey"`
	GUID                  GUID    `json:"guid"`
	ParentGUID            string  `json:"parentGuid"`
	GrandparentGUID       string  `json:"grandparentGuid"`
	Type                  string  `json:"type"`
	Title                 string  `json:"title"`
	GrandparentKey        string  `json:"grandparentKey"`
	ParentKey             string  `json:"parentKey"`
	LibrarySectionTitle   string  `json:"librarySectionTitle"`
	LibrarySectionID      int     `json:"librarySectionID"`
	LibrarySectionKey     string  `json:"librarySectionKey"`
	GrandparentTitle      string  `json:"grandparentTitle"`
	ParentTitle           string  `json:"parentTitle"`
	OriginalTitle         string  `json:"originalTitle"`
	ContentRating         string  `json:"contentRating"`
	Summary               string  `json:"summary"`
	Index                 int     `json:"index"`
	ParentIndex           int     `json:"parentIndex"`
	AudienceRating        float64 `json:"audienceRating"`
	UserRating            float64 `json:"userRating"`
	LastRatedAt           int     `json:"lastRatedAt"`
	Year                  int     `json:"year"`
	Thumb                 string  `json:"thumb"`
	Art                   string  `json:"art"`
	GrandparentThumb      string  `json:"grandparentThumb"`
	GrandparentArt        string  `json:"grandparentArt"`
	Duration              int     `json:"duration"`
	OriginallyAvailableAt string  `json:"originallyAvailableAt"`
	AddedAt               int     `json:"addedAt"`
	UpdatedAt             int     `json:"updatedAt"`
	AudienceRatingImage   string  `json:"audienceRatingImage"`
	Media                 []struct {
		ID              int     `json:"id"`
		Duration        int     `json:"duration"`
		Bitrate         int     `json:"bitrate"`
		Width           int     `json:"width"`
		Height          int     `json:"height"`
		AspectRatio     float64 `json:"aspectRatio"`
		AudioChannels   int     `json:"audioChannels"`
		AudioCodec      string  `json:"audioCodec"`
		VideoCodec      string  `json:"videoCodec"`
		VideoResolution string  `json:"videoResolution"`
		Container       string  `json:"container"`
		VideoFrameRate  string  `json:"videoFrameRate"`
		AudioProfile    string  `json:"audioProfile"`
		VideoProfile    string  `json:"videoProfile"`
		Part            []struct {
			ID           int    `json:"id"`
			Key          string `json:"key"`
			Duration     int    `json:"duration"`
			File         string `json:"file"`
			Size         int    `json:"size"`
			AudioProfile string `json:"audioProfile"`
			Container    string `json:"container"`
			Indexes      string `json:"indexes"`
			VideoProfile string `json:"videoProfile"`
			Stream       []struct {
				ID                   int     `json:"id"`
				StreamType           int     `json:"streamType"`
				Default              bool    `json:"default"`
				Codec                string  `json:"codec"`
				Index                int     `json:"index"`
				Bitrate              int     `json:"bitrate,omitempty"`
				BitDepth             int     `json:"bitDepth,omitempty"`
				ChromaLocation       string  `json:"chromaLocation,omitempty"`
				ChromaSubsampling    string  `json:"chromaSubsampling,omitempty"`
				CodedHeight          int     `json:"codedHeight,omitempty"`
				CodedWidth           int     `json:"codedWidth,omitempty"`
				ColorPrimaries       string  `json:"colorPrimaries,omitempty"`
				ColorRange           string  `json:"colorRange,omitempty"`
				ColorSpace           string  `json:"colorSpace,omitempty"`
				ColorTrc             string  `json:"colorTrc,omitempty"`
				FrameRate            float64 `json:"frameRate,omitempty"`
				HasScalingMatrix     bool    `json:"hasScalingMatrix,omitempty"`
				Height               int     `json:"height,omitempty"`
				Level                int     `json:"level,omitempty"`
				Profile              string  `json:"profile,omitempty"`
				RefFrames            int     `json:"refFrames,omitempty"`
				ScanType             string  `json:"scanType,omitempty"`
				Width                int     `json:"width,omitempty"`
				DisplayTitle         string  `json:"displayTitle"`
				ExtendedDisplayTitle string  `json:"extendedDisplayTitle"`
				Selected             bool    `json:"selected,omitempty"`
				Channels             int     `json:"channels,omitempty"`
				Language             string  `json:"language,omitempty"`
				LanguageTag          string  `json:"languageTag,omitempty"`
				LanguageCode         string  `json:"languageCode,omitempty"`
				AudioChannelLayout   string  `json:"audioChannelLayout,omitempty"`
				SamplingRate         int     `json:"samplingRate,omitempty"`
				Title                string  `json:"title,omitempty"`
			} `json:"Stream"`
		} `json:"Part"`
	} `json:"Media"`
	Rating []struct {
		Image string  `json:"image"`
		Value float64 `json:"value"`
		Type  string  `json:"type"`
	} `json:"Rating"`
	Director []struct {
		ID     int    `json:"id"`
		Filter string `json:"filter"`
		Tag    string `json:"tag"`
		TagKey string `json:"tagKey"`
	} `json:"Director"`
	Writer []struct {
		ID     int    `json:"id"`
		Filter string `json:"filter"`
		Tag    string `json:"tag"`
		TagKey string `json:"tagKey"`
		Thumb  string `json:"thumb"`
	} `json:"Writer"`
	Role []struct {
		ID     int    `json:"id"`
		Filter string `json:"filter"`
		Tag    string `json:"tag"`
		TagKey string `json:"tagKey"`
		Role   string `json:"role"`
		Thumb  string `json:"thumb,omitempty"`
	} `json:"Role"`
}

type PlexPayloadSource string

const (
	PlexWebhook PlexPayloadSource = "Plex Webhook"
	Tautulli    PlexPayloadSource = "Tautulli"
)

func NewPlexWebhook(payload []byte) (*Plex, error) {
	p := &Plex{}
	p.Source = string(PlexWebhook)
	err := json.Unmarshal(payload, p)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal plex payload")
	}

	return p, nil
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

type GetPlexRequest struct {
	Id int
}

type DeletePlexRequest struct {
	OlderThan int
}
