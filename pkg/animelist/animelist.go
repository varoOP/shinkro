package animelist

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strconv"
)

type AnimeList struct {
	XMLName xml.Name `xml:"anime-list"`
	Text    string   `xml:",chardata"`
	Anime   []struct {
		Text              string `xml:",chardata"`
		Anidbid           string `xml:"anidbid,attr"`
		Tvdbid            string `xml:"tvdbid,attr"`
		Defaulttvdbseason string `xml:"defaulttvdbseason,attr"`
		Name              string `xml:"name"`
		SupplementalInfo  struct {
			Text   string `xml:",chardata"`
			Studio string `xml:"studio"`
		} `xml:"supplemental-info"`
	} `xml:"anime"`
}

func NewAnimeList() *AnimeList {

	al := &AnimeList{}

	resp, err := http.Get("https://raw.githubusercontent.com/Anime-Lists/anime-lists/master/anime-list.xml")
	if err != nil {
		log.Fatalf("Error getting anime list xml: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body of anime list xml: %v", err)
	}

	err = xml.Unmarshal(body, al)
	if err != nil {
		log.Fatalf("Error unmarshalling anime list xml: %v", err)
	}

	return al
}

func (al *AnimeList) AnidDbTvDbmap() map[string]string {

	m := make(map[string]string)

	for _, v := range al.Anime {
		m[v.Anidbid] = v.Tvdbid
	}

	return m
}

func (al *AnimeList) GetTvdbID(anidbid int, m map[string]string) int {

	tvdbids := m[strconv.Itoa(anidbid)]

	if tvdbid, err := strconv.Atoi(tvdbids); err == nil {

		return tvdbid

	}

	return 0
}
