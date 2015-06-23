package config

import (
	"gopkg.in/simversity/gottp.v3/conf"
)

type config struct {
	Gottp  conf.GottpSettings
	Khabar struct {
		DBName               string
		DBAddress            string
		TranslationDirectory string
		Debug                bool
		DBUsername           string
		DBPassword           string
	}
}

func (self *config) MakeConfig(configPath string) {
	self.Gottp.Listen = "127.0.0.1:8911"
	self.Khabar.DBAddress = "127.0.0.1:27017"
	self.Khabar.DBName = "notifications_testing"
	if configPath != "" {
		conf.MakeConfig(configPath, self)
	}
}

func (self *config) GetGottpConfig() *conf.GottpSettings {
	return &self.Gottp
}

var Settings config
