package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	Dsn       string
	Config    string
	Log       string
	Addr      string
	User      string
	BaseUrl   string
	CustomMap bool
	K         *koanf.Koanf
	Logger    *os.File
}

func NewConfig(dir string) *Config {
	if dir == "" {
		log.Println("path to configuration not set")
		log.Println("Run: shinkuro help")
		os.Exit(1)
	}

	var err error
	c := &Config{}
	c.K = koanf.New(".")

	dsn := filepath.Join(dir, "shinkuro.db")
	c.Dsn = fmt.Sprintf("file:%v?cache=shared&mode=rwc&_journal_mode=WAL", dsn)
	c.Config = filepath.Join(dir, "config.toml")
	c.Log = filepath.Join(dir, "shinkuro.log")

	if _, err = os.Stat(c.Config); err != nil {
		createConfig(c)
		log.Println("Example config.toml created at", c.Config)
		log.Println("Edit before running shinkuro or shinkuro malauth again!")
		os.Exit(0)
	}

	if err := c.K.Load(file.Provider(c.Config), toml.Parser()); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	c.CustomMap = false
	if mapPath := c.K.String("custom_map"); mapPath != "" {
		c.CustomMap = true
	}

	c.BaseUrl = "/"
	if b := c.K.String("base_url"); b != "" {
		c.BaseUrl = b
	}

	c.Addr = fmt.Sprintf("%v:%v", c.K.String("host"), c.K.Int("port"))

	c.User = c.K.String("plex_user")

	c.Logger, err = os.OpenFile(c.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	mw := io.MultiWriter(os.Stdout, c.Logger)

	log.SetOutput(mw)
	return c
}

func createConfig(cfg *Config) {
	const config = `###Sample shinkuro config

host = "127.0.0.1"

port = 7011

plex_user = "Your_Plex_account_Title_EDIT_REQUIRED" 

############################################
#Default base_url = "/" (Optional setting) #
############################################
#base_url = "/shinkuro"

##########################################################
#Default is set to the community map. (Optional setting) #
##########################################################
#custom_map = "Absolute path to custom-mapping.yaml"
`
	f, err := os.Create(cfg.Config)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(config)
	if err != nil {
		log.Fatal(err)
	}
}
