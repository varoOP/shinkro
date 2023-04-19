package domain

type Config struct {
	ConfigPath        string
	Host              string   `koanf:"host"`
	Port              int      `koanf:"port"`
	PlexUser          string   `koanf:"plexUser"`
	AnimeLibraries    []string `koanf:"animeLibraries"`
	BaseUrl           string   `koanf:"baseUrl"`
	CustomMapPath     string   `koanf:"customMapPath"`
	DiscordWebHookURL string   `koanf:"discordWebhookUrl"`
	LogLevel          string   `koanf:"logLevel"`
	LogMaxSize        int      `koanf:"logMaxSize"`
	LogMaxBackups     int      `koanf:"logMaxBackups"`
}
