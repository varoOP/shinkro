package config

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/varoOP/shinkro/internal/api"
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
		Host:              "127.0.0.1",
		Port:              7011,
		PlexUser:          "",
		PlexUrl:           "",
		PlexToken:         "",
		AnimeLibraries:    []string{},
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
		SessionSecret:     api.GenerateSecureToken(16),
	}
}

var configTemplate = `###Example config.toml for shinkro
###[shinkro]
###Discord webhook, and BaseUrl are optional.
###LogLevel can be set to any one of the following: "INFO", "ERROR", "DEBUG", "TRACE"
###LogxMaxSize is in MB.
###[plex]
###PlexUser and AnimeLibraries must be set to the correct values. 
###AnimeLibraries is a list of your plex library names that contain anime - the ones shinkro will use to update your MAL account.
###Example: AnimeLibraries = ["Anime", "Anime Movies"]
###Url and Token are optional - only required if you have anime libraries that use the plex agents.

[shinkro]
Host = "127.0.0.1"
Port = 7011
#BaseUrl = "/shinkro"
#DiscordWebhookUrl = ""
LogLevel = "INFO"
LogMaxSize = 50
LogMaxBackups = 3
SessionSecret = "{{ .sessionSecret }}"

[plex]
PlexUsername = ""
AnimeLibraries = []
#Url = "http://127.0.0.1:32400"
#Token = "<Value of X-Plex-Token>"
`

func (c *AppConfig) writeConfig(cfgPath string) error {
	// check if config exists, if not create it
	if _, err := os.Stat(cfgPath); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(cfgPath)
		if err != nil { // perm 0666
			// handle failed create
			log.Printf("error creating file: %q", err)
			return err
		}
		defer f.Close()

		// setup text template to inject variables into
		tmpl, err := template.New("config").Parse(configTemplate)
		if err != nil {
			return errors.Wrap(err, "could not create config template")
		}

		tmplVars := map[string]string{
			"sessionSecret": c.Config.SessionSecret,
		}

		var buffer bytes.Buffer
		if err = tmpl.Execute(&buffer, &tmplVars); err != nil {
			return errors.Wrap(err, "could not write torrent url template output")
		}

		if _, err = f.WriteString(buffer.String()); err != nil {
			log.Printf("error writing contents to file: %v %q", cfgPath, err)
			return err
		}

		return f.Sync()
	}

	return nil
}

func (c *AppConfig) checkConfig(dir string) {

	if c.Config.PlexUser == "" {
		c.log.Fatal().Msgf("plex.PlexUsername not set in %v/config.toml", dir)
	}

	if len(c.Config.AnimeLibraries) < 1 {
		c.log.Fatal().Msgf("plex.AnimeLibraries not set in %v/config.toml", dir)
	}

	for i, v := range c.Config.AnimeLibraries {
		c.Config.AnimeLibraries[i] = strings.TrimSpace(v)
	}
}

func (c *AppConfig) parseConfig(dir string) error {
	cfgPath := filepath.Join(dir, "config.toml")

	if _, err := os.Stat(cfgPath); err != nil {
		err = c.writeConfig(cfgPath)
		if err != nil {
			c.log.Fatal().Msg("unable to write shinkro configuration file")
		}

		c.log.Fatal().Err(errors.New("shinkro configuration file not found")).Msgf("No config.toml found, example config.toml created at %v. Edit and run shinkro again", cfgPath)
	}

	k := koanf.New(".")
	if err := k.Load(file.Provider(cfgPath), toml.Parser()); err != nil {
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
