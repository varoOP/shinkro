package domain

type Config struct {
	ConfigPath        string
	Host              string   `koanf:"host"`
	Port              int      `koanf:"port"`
	PlexUser          string   `koanf:"plexUser"`
	PlexUrl           string   `koanf:"plexUrl"`
	PlexToken         string   `koanf:"plexToken"`
	AnimeLibraries    []string `koanf:"animeLibraries"`
	ApiKey            string   `koanf:"apiKey"`
	BaseUrl           string   `koanf:"baseUrl"`
	CustomMapTVDB     bool
	CustomMapTVDBPath string
	CustomMapTMDB     bool
	CustomMapTMDBPath string
	TVDBMalMap        *AnimeTVDBMap
	TMDBMalMap        *AnimeMovies
	DiscordWebHookURL string `koanf:"discordWebhookUrl"`
	LogLevel          string `koanf:"logLevel"`
	LogMaxSize        int    `koanf:"logMaxSize"`
	LogMaxBackups     int    `koanf:"logMaxBackups"`
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
