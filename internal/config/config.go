package config

import (
	"fmt"
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
	Addr      string
	User      string
	BaseUrl   string
	CustomMap string
}

func NewConfig(dir string) *Config {
	if dir == "" {
		log.Println("path to configuration not set")
		log.Println("Run: shinkuro help, for the help message.")
		os.Exit(1)
	}

	c := &Config{}
	c.joinPaths(dir)
	c.checkConfig()

	err := c.parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func (c *Config) createConfig() error {
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
	f, err := os.Create(c.Config)
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

func (c *Config) checkConfig() {
	if _, err := os.Stat(c.Config); err != nil {
		err = c.createConfig()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Example config.toml created at", c.Config)
		log.Printf("Edit %v before running shinkuro or shinkuro malauth again!\n", c.Config)
		os.Exit(0)
	}
}

func (c *Config) joinPaths(dir string) {
	dsn := filepath.Join(dir, "shinkuro.db")
	c.Dsn = fmt.Sprintf("file:%v?cache=shared&mode=rwc&_journal_mode=WAL", dsn)
	c.Config = filepath.Join(dir, "config.toml")
}

func (c *Config) parseConfig() error {
	k := koanf.New(".")

	if err := k.Load(file.Provider(c.Config), toml.Parser()); err != nil {
		return err
	}

	c.CustomMap = k.String("custom_map")

	c.BaseUrl = "/"
	if b := k.String("base_url"); b != "" {
		c.BaseUrl = b
	}

	c.Addr = fmt.Sprintf("%v:%v", k.String("host"), k.Int("port"))

	c.User = k.String("plex_user")

	return nil
}
