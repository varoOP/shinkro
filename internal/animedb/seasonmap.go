package animedb

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type SeasonMap struct {
	Anime []AnimeMT `yaml:"anime"`
}

type AnimeMT struct {
	Title   string `yaml:"title"`
	Seasons []struct {
		Season int `yaml:"season"`
		MalID  int `yaml:"mal-id"`
		Start  int `yaml:"start,omitempty"`
	} `yaml:"seasons"`
}

func (s *SeasonMap) GetSeasonMap(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(body, s)
	if err != nil {
		log.Fatalln(err)
	}
}
