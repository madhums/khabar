package config

import (
	"gopkg.in/simversity/gottp.v2/conf"
)

type config struct {
	Gottp            conf.GottpSettings
	Khabar struct {
		DBName               string
		DBAddress            string
		TranslationDirectory string
	}
}

func (self *config) MakeConfig(configPath string) {
	if configPath != "" {
		conf.MakeConfig(configPath, self)
	}
}

func (self *config) GetGottpConfig() *conf.GottpSettings {
	return &self.Gottp
}

var Settings config
