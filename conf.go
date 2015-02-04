package main

import (
	"gopkg.in/simversity/gottp.v1/conf"
)

type config struct {
	Gottp            conf.GottpSettings
	Sc_Notifications struct {
		DBName    string
		DBAddress string
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

var settings config
