package animelist

import (
	"encoding/xml"
	"io"
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
		Tmdbid            string `xml:"tmdbid,attr"`
		Defaulttvdbseason string `xml:"defaulttvdbseason,attr"`
		Name              string `xml:"name"`
		SupplementalInfo  struct {
			Text   string `xml:",chardata"`
			Studio string `xml:"studio"`
		} `xml:"supplemental-info"`
	} `xml:"anime"`
}

func NewAnimeList() (*AnimeList, error) {
	al := &AnimeList{}
	resp, err := http.Get("https://raw.githubusercontent.com/Anime-Lists/anime-lists/master/anime-list.xml")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(body, al)
	if err != nil {
		return nil, err
	}

	return al, nil
}

func (al *AnimeList) AnidDbTvmDbmap() map[string][]string {
	m := make(map[string][]string)
	for _, v := range al.Anime {
		m[v.Anidbid] = []string{v.Tvdbid, v.Tmdbid}
	}

	return m
}

func (al *AnimeList) GetTvmdbID(anidbid int, m map[string][]string) (int, int) {
	ids := m[strconv.Itoa(anidbid)]
	var tvdbids, tmdbids string
	if len(ids) >= 1 {
		tvdbids = ids[0]
		tmdbids = ids[1]
	}

	tvdbid, err := strconv.Atoi(tvdbids)
	if err != nil {
		tvdbid = -1
	}

	tmdbid, err := strconv.Atoi(tmdbids)
	if err != nil {
		tmdbid = -1
	}

	return tvdbid, tmdbid
}
