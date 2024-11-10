package domain

import "os"

type Config struct {
	Host              string   `koanf:"Host"`
	Port              int      `koanf:"Port"`
	PlexUser          string   `koanf:"PlexUsername"`
	PlexUrl           string   `koanf:"Url"`
	PlexToken         string   `koanf:"Token"`
	AnimeLibraries    []string `koanf:"AnimeLibraries"`
	BaseUrl           string   `koanf:"BaseUrl"`
	CustomMapTVDB     bool
	CustomMapTVDBPath string
	CustomMapTMDB     bool
	SessionSecret     string `koanf:"SessionSecret"`
	CustomMapTMDBPath string
	TVDBMalMap        *AnimeTVShows
	TMDBMalMap        *AnimeMovies
	DiscordWebHookURL string `koanf:"DiscordWebhookUrl"`
	LogLevel          string `koanf:"LogLevel"`
	LogMaxSize        int    `koanf:"LogMaxSize"`
	LogMaxBackups     int    `koanf:"LogMaxBackups"`
}

func (cfg *Config) LocalMapsExist() {
	cfg.CustomMapTMDB = false
	if fileExists(cfg.CustomMapTMDBPath) {
		cfg.CustomMapTMDB = true
	}

	cfg.CustomMapTVDB = false
	if fileExists(cfg.CustomMapTVDBPath) {
		cfg.CustomMapTVDB = true
	}
}

func (cfg *Config) isPlexClient() bool {
	return cfg.PlexToken != ""
}

func fileExists(path string) bool {
	_, err := os.Open(path)
	return err == nil
}
