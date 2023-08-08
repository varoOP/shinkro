package config

import (
	"crypto/rand"
	"log"
	"math/big"
	"os"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/varoOP/shinkro/internal/domain"
)

type AppConfig struct {
	Config *domain.Config
}

func NewConfig(dir string) *AppConfig {
	if dir == "" {
		log.Println("path to configuration not set")
		log.Fatal("Run: shinkro help, for the help message.")
	}

	c := &AppConfig{}
	c.defaultConfig(dir)
	c.checkConfig(dir)

	err := c.parseConfig()
	if err != nil {
		log.Println(err)
		log.Fatal("unable to parse config.toml")
	}

	return c
}

func (c *AppConfig) defaultConfig(dir string) {
	c.Config = &domain.Config{
		ConfigPath:        filepath.Join(dir, "config.toml"),
		Host:              "127.0.0.1",
		Port:              7011,
		PlexUser:          "",
		AnimeLibraries:    []string{""},
		ApiKey:            genApikey(),
		BaseUrl:           "/",
		CustomMapPath:     "",
		DiscordWebHookURL: "",
		LogLevel:          "INFO",
		LogMaxSize:        50,
		LogMaxBackups:     3,
	}
}

func (c *AppConfig) createConfig(dir string) error {
	var config = `#Sample shinkro config

host = "127.0.0.1"

port = 7011

plexUser = "Your_Plex_account_Title_EDIT_REQUIRED"

animeLibraries = ["Your", "Anime", "Library", "Names", "Edit", "This"]

apiKey = "` + c.Config.ApiKey + `"

#baseUrl = "/shinkro"

#customMapPath = ""

#discordWebhookUrl = ""

#logLevel = ""

#logMaxSize = 50

#logMaxBackups = 3
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
	if _, err := os.Stat(c.Config.ConfigPath); err != nil {
		err = c.createConfig(dir)
		if err != nil {
			log.Fatal("unable to write shinkro configuration file")
		}

		log.Println("shinkro configuration file not found")
		log.Fatalf("example config.toml created at %v", c.Config.ConfigPath)
	}
}

func (c *AppConfig) parseConfig() error {
	k := koanf.New(".")

	if err := k.Load(file.Provider(c.Config.ConfigPath), toml.Parser()); err != nil {
		return err
	}

	err := k.Unmarshal("", c.Config)
	if err != nil {
		return err
	}

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
