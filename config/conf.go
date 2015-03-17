package config

import (
	"gopkg.in/simversity/gottp.v2/conf"
)

type config struct {
	Gottp  conf.GottpSettings
	Khabar struct {
		DBName               string
		DBAddress            string
		TranslationDirectory string
		Debug                bool
	}
}

func (self *config) MakeConfig(configPath string) {
	self.Gottp.Listen = "127.0.0.1:8911"
	if configPath != "" {
		conf.MakeConfig(configPath, self)
	}
}

func (self *config) GetGottpConfig() *conf.GottpSettings {
	return &self.Gottp
}

var Settings config
