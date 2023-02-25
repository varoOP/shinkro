package config

import "fmt"

type ConfigPath struct {
	Dsn    string
	Config string
	Token  string
	Log    string
}

func NewConfigPath(dir string) *ConfigPath {
	p := &ConfigPath{}
	p.Dsn = fmt.Sprintf("file:%v/shinkuro.db?cache=shared&mode=rwc&_journal_mode=WAL", dir)
	p.Config = makePath(dir, "config.toml")
	p.Token = makePath(dir, "token.json")
	p.Log = makePath(dir, "shinkuro.log")
	return p
}

func makePath(dir, base string) string {
	return fmt.Sprintf("%v/%v", dir, base)
}
