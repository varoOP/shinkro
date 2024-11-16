package domain

type Config struct {
	Version         string
	ConfigPath      string
	Host            string `koanf:"Host"`
	Port            int    `koanf:"Port"`
	BaseUrl         string `koanf:"BaseUrl"`
	SessionSecret   string `koanf:"SessionSecret"`
	LogLevel        string `koanf:"LogLevel"`
	LogPath         string `koanf:"LogPath"`
	LogMaxSize      int    `koanf:"LogMaxSize"`
	LogMaxBackups   int    `koanf:"LogMaxBackups"`
	CheckForUpdates bool   `koanf:"CheckForUpdates"`
}

type ConfigUpdate struct {
	Host            *string `json:"host,omitempty"`
	Port            *int    `json:"port,omitempty"`
	LogLevel        *string `json:"log_level,omitempty"`
	LogPath         *string `json:"log_path,omitempty"`
	BaseURL         *string `json:"base_url,omitempty"`
	CheckForUpdates *bool   `json:"check_for_updates,omitempty"`
}
