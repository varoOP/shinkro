package domain

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type PlexRepo interface {
	Store(ctx context.Context, plex *Plex) error
	FindAll(ctx context.Context) ([]*Plex, error)
	FindAllWithFilters(ctx context.Context, params PlexPayloadQueryParams) (*FindPlexPayloadsResponse, error)
	Get(ctx context.Context, req *GetPlexRequest) (*Plex, error)
	Delete(ctx context.Context, req *DeletePlexRequest) error
	CountScrobbleEvents(ctx context.Context) (int, error)
	CountRateEvents(ctx context.Context) (int, error)
	GetRecent(ctx context.Context, limit int) ([]*Plex, error)
	UpdateStatus(ctx context.Context, plexID int64, success *bool, errorType PlexErrorType, errorMsg string) error
}

type Plex struct {
	ID        int64             `json:"id"`
	Rating    float32           `json:"rating"`
	TimeStamp time.Time         `json:"timestamp"`
	Event     PlexEvent         `json:"event"`
	User      bool              `json:"user"`
	Source    PlexPayloadSource `json:"source"`
	Owner     bool              `json:"owner"`
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
	// Status fields (consolidated from plex_status table)
	Success   *bool         `json:"success,omitempty"`
	ErrorType PlexErrorType `json:"errorType,omitempty"`
	ErrorMsg  string        `json:"errorMsg,omitempty"`
}

type PlexErrorType string

const (
	PlexErrorAgentNotSupported PlexErrorType = "AGENT_NOT_SUPPORTED"
	PlexErrorExtractionFailed  PlexErrorType = "EXTRACTION_FAILED"
	PlexErrorUnknown           PlexErrorType = "UNKNOWN_ERROR"
)

type Metadata struct {
	RatingGlobal          float32       `json:"rating"`
	RatingKey             string        `json:"ratingKey"`
	Key                   string        `json:"key"`
	ParentRatingKey       string        `json:"parentRatingKey"`
	GrandparentRatingKey  string        `json:"grandparentRatingKey"`
	GUID                  GUID          `json:"guid"`
	ParentGUID            string        `json:"parentGuid"`
	GrandparentGUID       string        `json:"grandparentGuid"`
	Type                  PlexMediaType `json:"type"`
	Title                 string        `json:"title"`
	GrandparentKey        string        `json:"grandparentKey"`
	ParentKey             string        `json:"parentKey"`
	LibrarySectionTitle   string        `json:"librarySectionTitle"`
	LibrarySectionID      int           `json:"librarySectionID"`
	LibrarySectionKey     string        `json:"librarySectionKey"`
	GrandparentTitle      string        `json:"grandparentTitle"`
	ParentTitle           string        `json:"parentTitle"`
	OriginalTitle         string        `json:"originalTitle"`
	ContentRating         string        `json:"contentRating"`
	Summary               string        `json:"summary"`
	Index                 int           `json:"index"`
	ParentIndex           int           `json:"parentIndex"`
	AudienceRating        float64       `json:"audienceRating"`
	UserRating            float64       `json:"userRating"`
	LastRatedAt           int           `json:"lastRatedAt"`
	Year                  int           `json:"year"`
	Thumb                 string        `json:"thumb"`
	Art                   string        `json:"art"`
	GrandparentThumb      string        `json:"grandparentThumb"`
	GrandparentArt        string        `json:"grandparentArt"`
	Duration              int           `json:"duration"`
	OriginallyAvailableAt string        `json:"originallyAvailableAt"`
	AddedAt               int           `json:"addedAt"`
	UpdatedAt             int           `json:"updatedAt"`
	AudienceRatingImage   string        `json:"audienceRatingImage"`
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

type PlexResponse struct {
	MediaContainer struct {
		Size                int        `json:"size"`
		AllowSync           bool       `json:"allowSync"`
		Identifier          string     `json:"identifier"`
		LibrarySectionID    int        `json:"librarySectionID"`
		LibrarySectionTitle string     `json:"librarySectionTitle"`
		LibrarySectionUUID  string     `json:"librarySectionUUID"`
		MediaTagPrefix      string     `json:"mediaTagPrefix"`
		MediaTagVersion     int        `json:"mediaTagVersion"`
		Metadata            []Metadata `json:"metadata"`
	} `json:"MediaContainer"`
}

type GUID struct {
	GUIDS []struct {
		ID string `json:"id"`
	}

	GUID string
}

type PlexPayloadSource string

const (
	PlexWebhook     PlexPayloadSource = "Plex Webhook"
	TautulliWebhook PlexPayloadSource = "Tautulli"
)

type PlexEvent string

const (
	PlexScrobbleEvent PlexEvent = "media.scrobble"
	PlexRateEvent     PlexEvent = "media.rate"
)

type PlexMediaType string

const (
	PlexEpisode PlexMediaType = "episode"
	PlexMovie   PlexMediaType = "movie"
)

type GetPlexRequest struct {
	Id int
}

type DeletePlexRequest struct {
	Id        int64
	OlderThan int
}

type PlexSupportedAgents string

const (
	MALAgent  PlexSupportedAgents = "mal"
	HAMA      PlexSupportedAgents = "hama"
	PlexAgent PlexSupportedAgents = "plex"
)

type PlexSupportedDBs string

const (
	TVDB  PlexSupportedDBs = "tvdb"
	TMDB  PlexSupportedDBs = "tmdb"
	AniDB PlexSupportedDBs = "anidb"
	MAL   PlexSupportedDBs = "myanimelist"
)

func NewPlexWebhook(payload []byte) (*Plex, error) {
	p := &Plex{}
	p.Source = PlexWebhook
	p.TimeStamp = time.Now()
	err := json.Unmarshal(payload, p)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal plex payload")
	}

	return p, nil
}

func (p *Plex) GetPlexMediaType() PlexMediaType {
	switch p.Metadata.Type {
	case PlexEpisode:
		return PlexEpisode
	case PlexMovie:
		return PlexMovie
	}

	return ""
}

func (p *Plex) GetPlexEvent() PlexEvent {
	switch p.Event {
	case PlexScrobbleEvent:
		return PlexScrobbleEvent
	case PlexRateEvent:
		return PlexRateEvent
	}

	return ""
}

func (p *Plex) IsEventAllowed() bool {
	return p.Event == PlexRateEvent || p.Event == PlexScrobbleEvent
}

func (p *Plex) IsRatingAllowed() bool {
	return p.Rating >= 0
}

func (p *Plex) IsPlexUserAllowed(ps *PlexSettings) bool {
	return p.Account.Title == ps.PlexUser
}

func (p *Plex) IsAnimeLibrary(ps *PlexSettings) bool {
	for _, library := range ps.AnimeLibraries {
		if library == p.Metadata.LibrarySectionTitle {
			return true
		}
	}
	return false
}

func (p *Plex) IsMediaTypeAllowed() bool {
	return p.Metadata.Type == PlexEpisode || p.Metadata.Type == PlexMovie
}

func (p *Plex) SetAnimeFields(source PlexSupportedDBs, id int) AnimeUpdate {
	// Extract title from Plex metadata
	var title string
	if p.Metadata.Type == PlexMovie {
		title = p.Metadata.Title
	} else {
		// For episodes, the show title is in GrandparentTitle
		title = p.Metadata.GrandparentTitle
	}

	// Initialize ListDetails with title from Plex payload
	listDetails := ListDetails{
		Title: title,
	}

	if p.Metadata.Type == PlexMovie {
		return AnimeUpdate{
			PlexId:      p.ID,
			Plex:        p,
			SourceId:    id,
			SourceDB:    source,
			Timestamp:   time.Now(),
			SeasonNum:   1,
			EpisodeNum:  1,
			ListDetails: listDetails,
		}
	}
	return AnimeUpdate{
		PlexId:      p.ID,
		Plex:        p,
		SourceId:    id,
		SourceDB:    source,
		Timestamp:   time.Now(),
		SeasonNum:   p.Metadata.ParentIndex,
		EpisodeNum:  p.Metadata.Index,
		ListDetails: listDetails,
	}
}

func (p *Plex) IsMetadataAgentAllowed() (bool, PlexSupportedAgents) {
	agents := map[string]PlexSupportedAgents{
		"agents.hama": HAMA,
		"myanimelist": MALAgent,
		"plex://":     PlexAgent,
	}

	for key, value := range agents {
		if strings.Contains(p.Metadata.GUID.GUID, key) || strings.Contains(p.Metadata.GrandparentGUID, key) {
			return true, value
		}
	}

	return false, ""
}

func (g *GUID) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &g.GUID); err == nil {
		return nil
	}

	// If it's not a string, try to unmarshal as an anonymous slice of struct
	if err := json.Unmarshal(data, &g.GUIDS); err == nil {
		return nil
	}

	return errors.Errorf("guid: cannot unmarshal %q", data)
}

func (g *GUID) HamaMALAgent(agent PlexSupportedAgents) (PlexSupportedDBs, int, error) {
	var agentRegExMap = map[PlexSupportedAgents]string{
		HAMA:     `//(.* ?)-(\d+ ?)`,
		MALAgent: `.(m.*)://(\d+ ?)`,
	}

	guid := g.GUID
	r := regexp.MustCompile(agentRegExMap[agent])
	if !r.MatchString(guid) {
		return "", -1, errors.Errorf("unable to parse GUID: %v", guid)
	}

	mm := r.FindStringSubmatch(guid)
	source := mm[1]
	id, err := strconv.Atoi(mm[2])
	if err != nil {
		return "", -1, errors.Wrap(err, "conversion of id failed")
	}

	return PlexSupportedDBs(source), id, nil
}

func (g *GUID) PlexAgent(mediaType PlexMediaType) (PlexSupportedDBs, int, error) {
	for _, gid := range g.GUIDS {
		dbid := strings.Split(gid.ID, "://")
		if (mediaType == PlexEpisode && dbid[0] == "tvdb") || (mediaType == PlexMovie && dbid[0] == "tmdb") {
			id, err := strconv.Atoi(dbid[1])
			if err != nil {
				return "", -1, errors.Wrap(err, "id conversion failed")
			}

			return PlexSupportedDBs(dbid[0]), id, nil
		}
	}

	return "", -1, errors.New("no supported online database found")
}

type PlexHistoryRequest struct {
	Limit int `json:"limit"`
}

type PlexHistoryItem struct {
	Plex        *Plex        `json:"plex"`
	AnimeUpdate *AnimeUpdate `json:"animeUpdate,omitempty"`
}

type PlexPayloadQueryParams struct {
	Limit   uint64
	Offset  uint64
	Search  string
	Filters struct {
		Event  PlexEvent
		Source PlexPayloadSource
		Status *bool // nil = no filter, true = success, false = failed
	}
}

type PlexPayloadListItem struct {
	Plex *Plex `json:"plex"`
}

type FindPlexPayloadsResponse struct {
	Data       []PlexPayloadListItem `json:"data"`
	TotalCount int                   `json:"count"`
}
