package animedb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nstratos/go-myanimelist/mal"
)

type SeasonMap struct {
	Anime []AnimeS `json:"anime"`
}

type AnimeS struct {
	MalID   int      `json:"mal_id"`
	Title   string   `json:"title"`
	Seasons []string `json:"seasons"`
}

func (s *SeasonMap) GetSeasonMap() {
	file, err := os.Open("./season-mal-map.json")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *SeasonMap) AddtoSeasonMap(ctx context.Context, malid, season int, title string, client *mal.Client) {

	var seasons []string
	m := malid
	for i := 2; i <= season; i++ {
		f, _, err := client.Anime.Details(ctx, malid, mal.Fields{"related_anime"})
		if err != nil {
			log.Println(err)
		}
		for _, v := range f.RelatedAnime {
			if v.RelationType == "sequel" {
				s := fmt.Sprintf("Season_%v:%v", i, v.Node.ID)
				seasons = append(seasons, s)
				malid = v.Node.ID
			}
		}
	}

	s.Anime = append(s.Anime, AnimeS{
		m, title, seasons})
}

func (s *SeasonMap) ModifySeasonMap(ctx context.Context, malid, season int, title string, client *mal.Client) {

}
