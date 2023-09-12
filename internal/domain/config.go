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
	CustomMapPath     string   `koanf:"customMapPath"`
	DiscordWebHookURL string   `koanf:"discordWebhookUrl"`
	LogLevel          string   `koanf:"logLevel"`
	LogMaxSize        int      `koanf:"logMaxSize"`
	LogMaxBackups     int      `koanf:"logMaxBackups"`
}
