package config

import (
	"crypto/rand"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/domain"
)

type AppConfig struct {
	log    zerolog.Logger
	Config *domain.Config
}

func NewConfig(dir string) *AppConfig {
	if dir == "" {
		log.Println("path to configuration not set")
		log.Fatal("Run: shinkro help, for the help message.")
	}

	c := &AppConfig{
		log: zerolog.New(
			zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.DateTime,
			}).With().Timestamp().Logger(),
	}

	c.defaultConfig(dir)
	err := c.parseConfig(dir)
	if err != nil {
		c.log.Fatal().Err(err).Msg("unable to parse config.toml")
	}

	c.checkConfig(dir)
	return c
}

func (c *AppConfig) defaultConfig(dir string) {
	c.Config = &domain.Config{
		ConfigPath:        filepath.Join(dir, "config.toml"),
		Host:              "127.0.0.1",
		Port:              7011,
		PlexUser:          "",
		PlexUrl:           "",
		PlexToken:         "",
		AnimeLibraries:    []string{},
		ApiKey:            genApikey(),
		BaseUrl:           "/",
		CustomMapTVDB:     false,
		CustomMapTVDBPath: filepath.Join(dir, "tvdb-mal.yaml"),
		CustomMapTMDB:     false,
		CustomMapTMDBPath: filepath.Join(dir, "tmdb-mal.yaml"),
		TMDBMalMap:        nil,
		TVDBMalMap:        nil,
		DiscordWebHookURL: "",
		LogLevel:          "INFO",
		LogMaxSize:        50,
		LogMaxBackups:     3,
	}
}

func (c *AppConfig) createConfig(dir string) error {
	var config = `###Example config.toml for shinkro
###[shinkro]
### Username and Password is required for basic authentication.
###Discord webhook, and BaseUrl are optional.
###LogLevel can be set to any one of the following: "INFO", "ERROR", "DEBUG", "TRACE"
###LogxMaxSize is in MB.
###[plex]
###PlexUser and AnimeLibraries must be set to the correct values. 
###AnimeLibraries is a list of your plex library names that contain anime - the ones shinkro will use to update your MAL account.
###Example: AnimeLibraries = ["Anime", "Anime Movies"]
###Url and Token are optional - only required if you have anime libraries that use the plex agents.

[shinkro]
Username = ""
Password = ""
Host = "127.0.0.1"
Port = 7011
ApiKey = "` + c.Config.ApiKey + `"
#BaseUrl = "/shinkro"
#DiscordWebhookUrl = ""
LogLevel = "INFO"
LogMaxSize = 50
LogMaxBackups = 3

[plex]
PlexUser = ""
AnimeLibraries = []
#Url = "http://127.0.0.1:32400"
#Token = "<Value of X-Plex-Token>"
`

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(c.Config.ConfigPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(config)
	if err != nil {
		return err
	}

	return nil
}

func (c *AppConfig) checkConfig(dir string) {
	if c.Config.ApiKey == "" {
		c.log.Fatal().Msgf("shinkro.ApiKey not set in %v/config.toml", dir)
	}

	if c.Config.Username == "" {
		c.log.Fatal().Msgf("shinkro.Username not set in %v/config.toml", dir)
	}

	if c.Config.Password == "" {
		c.log.Fatal().Msgf("shinkro.Password not set in %v/config.toml", dir)
	}

	if c.Config.PlexUser == "" {
		c.log.Fatal().Msgf("plex.PlexUser not set in %v/config.toml", dir)
	}

	if len(c.Config.AnimeLibraries) < 1 {
		c.log.Fatal().Msgf("plex.AnimeLibraries not set in %v/config.toml", dir)
	}

	for i, v := range c.Config.AnimeLibraries {
		c.Config.AnimeLibraries[i] = strings.TrimSpace(v)
	}
}

func (c *AppConfig) parseConfig(dir string) error {
	if _, err := os.Stat(c.Config.ConfigPath); err != nil {
		err = c.createConfig(dir)
		if err != nil {
			c.log.Fatal().Msg("unable to write shinkro configuration file")
		}

		c.log.Fatal().Err(errors.New("shinkro configuration file not found")).Msgf("No config.toml found, example config.toml created at %v. Edit and run shinkro again", c.Config.ConfigPath)
	}

	k := koanf.New(".")
	if err := k.Load(file.Provider(c.Config.ConfigPath), toml.Parser()); err != nil {
		return err
	}

	err := k.Unmarshal("shinkro", c.Config)
	if err != nil {
		return err
	}

	err = k.Unmarshal("plex", c.Config)
	if err != nil {
		return err
	}

	c.Config.LocalMapsExist()
	return nil
}

func genApikey() string {
	allowed := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, 32)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(allowed))))
		if err != nil {
			log.Fatal(err)
		}

		b[i] = allowed[n.Int64()]
	}

	return string(b)
}
