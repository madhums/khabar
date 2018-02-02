package config

import (
	"os"

	"gopkg.in/simversity/gottp.v3/conf"
)

type config struct {
	Gottp conf.GottpSettings
	S3    struct {
		AccessKey string
		SecretKey string
		Bucket    string
		Region    string
	}
	Moire struct {
		DbUrl                  string
		DbName                 string
		TranslationDirectory   string
		Debug                  bool
		FFmpeg                 string
		SignRequests           bool
		ImageTimeout           int
		StaticPath             string
		PublicKey              string
		PrivateKey             string
		SentryDSN              string
		UploadUrlExpiry        int64
		GetUrlExpiry           int64
		RedirectUrlCacheExpiry int64
	}
}

func (self *config) MakeConfig(configPath string) {
	self.Gottp.Listen = "127.0.0.1:8811"

	self.Moire.DbUrl = getEnv("MONGODB_URL", "mongodb://localhost/gallery")
	self.Moire.DbName = getEnv("MONGODB_NAME", "gallery")
	self.S3.AccessKey = getEnv("S3_ACCESS_KEY", "")
	self.S3.SecretKey = getEnv("S3_SECRET_KEY", "")
	self.S3.Bucket = getEnv("S3_BUCKET", "moire-gallery")

	self.Moire.Debug = false
	self.Moire.FFmpeg = "ffmpeg"
	self.Moire.SignRequests = false
	self.Moire.ImageTimeout = 15
	self.Moire.StaticPath = "https://cdn.safetychanger.com/statics"

	self.Moire.PublicKey = DefaultPublicKey
	self.Moire.PrivateKey = DefaultPrivateKey
	self.Moire.UploadUrlExpiry = 7200      // 5 days (60 * 24 * 5)
	self.Moire.GetUrlExpiry = 60           // 1 hour
	self.Moire.RedirectUrlCacheExpiry = 45 // 45 minutes, must be lower than GetUrlExpiry

	self.S3.Region = "eu-west-1"

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
