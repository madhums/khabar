package config

import (
	"os"

	"gopkg.in/simversity/gottp.v3/conf"
)

type config struct {
	Gottp  conf.GottpSettings
	Khabar struct {
		DbUrl                string
		DbName               string
		TranslationDirectory string
		Debug                bool
	}
}

func (self *config) MakeConfig(configPath string) {
	self.Gottp.Listen = ":8911"

	self.Khabar.DbUrl = getEnv("MONGODB_URL", "mongodb://localhost/notifications_testing")
	self.Khabar.DbName = getEnv("MONGODB_NAME", "notifications_testing")
	self.Khabar.TranslationDirectory = getEnv("TRANSLATION_DIRECTORY", "")

	if configPath != "" {
		conf.MakeConfig(configPath, self)
	}
}

func getEnv(key string, defaultVal string) string {
	if env := os.Getenv(key); env != "" {
		return env
	} else {
		return defaultVal
	}
}

func (self *config) GetGottpConfig() *conf.GottpSettings {
	return &self.Gottp
}

var Settings config
