package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"os"

	"golang.org/x/oauth2"
)

func jsonUnmarshal(path string, a interface{}) error {

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	body, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, a)
	if err != nil {
		return err
	}

	return nil
}

func randomString(l int) string {
	random := make([]byte, l)
	_, err := rand.Read(random)
	if err != nil {
		log.Fatalln(err)
	}

	return base64.URLEncoding.EncodeToString(random)[:l]
}

func saveToken(token *oauth2.Token) {

	tokenJ, err := os.Create("./token.json")
	if err != nil {
		log.Fatalln(err)
	}

	defer tokenJ.Close()

	e := json.NewEncoder(tokenJ)
	e.SetIndent("", "  ")
	e.Encode(token)
	log.Println("Token saved.")

}
