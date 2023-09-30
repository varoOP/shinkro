package domain

type Config struct {
	ConfigPath        string
	Host              string   `koanf:"Host"`
	Port              int      `koanf:"Port"`
	PlexUser          string   `koanf:"UserName"`
	PlexUrl           string   `koanf:"Url"`
	PlexToken         string   `koanf:"Token"`
	AnimeLibraries    []string `koanf:"AnimeLibraries"`
	ApiKey            string   `koanf:"ApiKey"`
	BaseUrl           string   `koanf:"BaseUrl"`
	CustomMapTVDB     bool
	CustomMapTVDBPath string
	CustomMapTMDB     bool
	CustomMapTMDBPath string
	TVDBMalMap        *AnimeTVDBMap
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
