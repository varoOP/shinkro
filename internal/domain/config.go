package domain

type Config struct {
	Host           string   `koanf:"Host"`
	Port           int      `koanf:"Port"`
	BaseUrl        string   `koanf:"BaseUrl"`
	SessionSecret  string   `koanf:"SessionSecret"`
	LogLevel       string   `koanf:"LogLevel"`
	LogMaxSize     int      `koanf:"LogMaxSize"`
	LogMaxBackups  int      `koanf:"LogMaxBackups"`
}

// func (cfg *Config) LocalMapsExist() {
// 	cfg.CustomMapTMDB = false
// 	if fileExists(cfg.CustomMapTMDBPath) {
// 		cfg.CustomMapTMDB = true
// 	}

// 	cfg.CustomMapTVDB = false
// 	if fileExists(cfg.CustomMapTVDBPath) {
// 		cfg.CustomMapTVDB = true
// 	}
// }

// func (cfg *Config) isPlexClient() bool {
// 	return cfg.PlexToken != ""
// }

// func fileExists(path string) bool {
// 	_, err := os.Open(path)
// 	return err == nil
// }
