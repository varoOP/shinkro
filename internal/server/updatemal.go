package server

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/varoOP/shinkuro/pkg/plex"
)

func UpdateMal(p *plex.PlexWebhook, client *http.Client, db *sql.DB) {

	s := NewShow(p.Metadata.GUID)
	malid := s.GetMalID(db)
	fmt.Printf("%+v", s)
	fmt.Println("malid:", malid)

	if s.Ep.Season == 1 {

		endpoint := fmt.Sprintf("https://api.myanimelist.net/v2/anime/%v/my_list_status", malid)

		payload := url.Values{
			"num_watched_episodes": {strconv.Itoa(s.Ep.No)},
			"status":               {"watching"},
		}

		req, err := http.NewRequest("PUT", endpoint, strings.NewReader(payload.Encode()))
		if err != nil {
			fmt.Printf("Unable to perform POST request:\n%v", err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(body))
	}
}
