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
	"github.com/spf13/pflag"
)

type Config struct {
	Dsn    string
	Config string
	Token  string
	Log    string
	Addr   string
	Map    string
	K      *koanf.Koanf
	Logger *os.File
}

func NewConfig() *Config {

	c := &Config{}

	var (
		dir string
		err error
	)

	c.K = koanf.New(".")

	pflag.StringVar(&dir, "config", ".", "Absolute path to shinkuro's configuration directory")

	pflag.Parse()

	dsn := filepath.Join(dir, "shinkuro.db")
	c.Dsn = fmt.Sprintf("file:%v?cache=shared&mode=rwc&_journal_mode=WAL", dsn)
	c.Config = filepath.Join(dir, "config.toml")
	c.Token = filepath.Join(dir, "token.json")
	c.Log = filepath.Join(dir, "shinkuro.log")
	c.Map = filepath.Join(dir, "tvdb-mal.yml")

	if err := c.K.Load(file.Provider(c.Config), toml.Parser()); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	c.Addr = fmt.Sprintf("%v:%v", c.K.String("host"), c.K.Int("port"))

	c.Logger, err = os.OpenFile(c.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatalln(err)
	}

	mw := io.MultiWriter(os.Stdout, c.Logger)

	log.SetOutput(mw)

	return c
}
