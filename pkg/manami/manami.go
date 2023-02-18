package manami

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Manami struct {
	License struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	} `json:"license,omitempty"`
	Repository string       `json:"repository,omitempty"`
	LastUpdate string       `json:"lastUpdate,omitempty"`
	Data       []ManamiData `json:"data,omitempty"`
}

type ManamiData struct {
	Sources     []string `json:"sources,omitempty"`
	Title       string   `json:"title,omitempty"`
	Type        string   `json:"type,omitempty"`
	Episodes    int      `json:"episodes,omitempty"`
	Status      string   `json:"status,omitempty"`
	AnimeSeason struct {
		Season string `json:"season,omitempty"`
		Year   int    `json:"year,omitempty"`
	} `json:"animeSeason,omitempty"`
	Picture   string   `json:"picture,omitempty"`
	Thumbnail string   `json:"thumbnail,omitempty"`
	Synonyms  []string `json:"synonyms,omitempty"`
	Relations []string `json:"relations,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

func NewManami() *Manami {

	m := &Manami{}

	resp, err := http.Get("https://github.com/manami-project/anime-offline-database/raw/master/anime-offline-database.json")
	if err != nil {
		log.Fatalf("Error getting manami offline database: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body of manami offline database: %v", err)
	}

	err = json.Unmarshal(body, m)
	if err != nil {
		log.Fatalf("Error unmarshalling manami database: %v", err)
	}

	return m
}

func (a *ManamiData) SortToIDs(re string) int {

	for _, val := range a.Sources {

		r := regexp.MustCompile(re)

		match := r.FindStringSubmatch(val)

		if len(match) > 1 {

			ID, err := strconv.Atoi(match[1])
			if err != nil {
				log.Fatalf("Error converting id to int: %v", err)
			}

			return ID
		}
	}
	return 0
}
