package plex

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type Plex struct {
	ID        int64     `json:"id"`
	Rating    float32   `json:"rating"`
	TimeStamp time.Time `json:"timestamp"`
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

type ServerResponse struct {
	Servers []Server
}

type Connection struct {
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	URI      string `json:"uri"`
	Local    bool   `json:"local"`
	Relay    bool   `json:"relay"`
	IPv6     bool   `json:"IPv6"`
}

type Server struct {
	Name                   string       `json:"name"`
	Product                string       `json:"product"`
	ProductVersion         string       `json:"productVersion"`
	Platform               *string      `json:"platform"`
	PlatformVersion        *string      `json:"platformVersion"`
	Device                 *string      `json:"device"`
	ClientIdentifier       string       `json:"clientIdentifier"`
	CreatedAt              time.Time    `json:"createdAt"`
	LastSeenAt             time.Time    `json:"lastSeenAt"`
	Provides               string       `json:"provides"`
	OwnerID                *int         `json:"ownerId"`
	SourceTitle            *string      `json:"sourceTitle"`
	PublicAddress          string       `json:"publicAddress"`
	AccessToken            string       `json:"accessToken"`
	Owned                  bool         `json:"owned"`
	Home                   bool         `json:"home"`
	Synced                 bool         `json:"synced"`
	Relay                  bool         `json:"relay"`
	Presence               bool         `json:"presence"`
	HTTPSRequired          bool         `json:"httpsRequired"`
	PublicAddressMatches   bool         `json:"publicAddressMatches"`
	DNSRebindingProtection bool         `json:"dnsRebindingProtection"`
	NATLoopbackSupported   bool         `json:"natLoopbackSupported"`
	Connections            []Connection `json:"connections"`
}

type LibraryResponse struct {
	MediaContainer struct {
		Size      int         `json:"size"`
		AllowSync bool        `json:"allowSync"`
		Title1    string      `json:"title1"`
		Directory []Directory `json:"Directory"`
	} `json:"MediaContainer"`
}

type Directory struct {
	AllowSync        bool       `json:"allowSync"`
	Art              string     `json:"art"`
	Composite        string     `json:"composite"`
	Filters          bool       `json:"filters"`
	Refreshing       bool       `json:"refreshing"`
	Thumb            string     `json:"thumb"`
	Key              string     `json:"key"`
	Type             string     `json:"type"`
	Title            string     `json:"title"`
	Agent            string     `json:"agent"`
	Scanner          string     `json:"scanner"`
	Language         string     `json:"language"`
	UUID             string     `json:"uuid"`
	UpdatedAt        int64      `json:"updatedAt"`
	CreatedAt        int64      `json:"createdAt"`
	ScannedAt        int64      `json:"scannedAt"`
	Content          bool       `json:"content"`
	Directory        bool       `json:"directory"`
	ContentChangedAt int64      `json:"contentChangedAt"`
	Hidden           int        `json:"hidden"`
	Location         []Location `json:"Location"`
}

type Location struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
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
