package config

import (
	"fmt"
	"path/filepath"
)

type ConfigPath struct {
	Dsn    string
	Config string
	Token  string
	Log    string
}

func NewConfigPath(dir string) *ConfigPath {
	p := &ConfigPath{}
	dsn := filepath.Join(dir, "shinkuro.db")
	p.Dsn = fmt.Sprintf("file:%v?cache=shared&mode=rwc&_journal_mode=WAL", dsn)
	p.Config = filepath.Join(dir, "config.toml")
	p.Token = filepath.Join(dir, "token.json")
	p.Log = filepath.Join(dir, "shinkuro.log")
	return p
}
