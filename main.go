// Package main is the CLI.
// You can use the CLI via Terminal.
package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bulletind/khabar/config"
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/utils"
	"github.com/nicksnyder/go-i18n/i18n"
	"gopkg.in/simversity/gottp.v3"
)

func sysInit() {
	<-(gottp.SysInitChan) //Buffered Channel to receive the server upstart boolean

	config.InitTracer()
	log.Println("Initialized GoTracer")

	db.Conn = db.GetConn(config.Settings.Khabar.DbUrl, config.Settings.Khabar.DbName)

	log.Println("Database Connected: " + config.Settings.Khabar.DbName)

	db.Conn.InitIndexes()

	transDir := config.Settings.Khabar.TranslationDirectory

	if len(transDir) == 0 {
		cwd := os.Getenv("PWD")
		transDir = cwd + "/translations"
		config.Settings.Khabar.TranslationDirectory = cwd
	}

	log.Println("Directory for translation :" + transDir)

	filepath.Walk(transDir, func(path string, _ os.FileInfo, err error) error {
		fileExt := filepath.Ext(path)
		if fileExt == ".json" && strings.HasSuffix(path, "base_email.json") == false && err == nil {
			log.Println("Loading translation file:" + path)
			i18n.MustLoadTranslationFile(path)
		} else {
			log.Print("Skipping translation file:" + path + " " +
				"File Extension:" + fileExt + " ")
			if err != nil {
				log.Print("Error:" + err.Error())
			}
		}
		return nil
	})

	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.Println("Translations have been parsed.")
}

func initCleaner() {
	ticker := time.NewTicker(time.Hour * 1)
	go func() {
		for _ = range ticker.C {
			utils.CleanupDownloads()
		}
	}()
}

func main() {

	go sysInit()
	go initCleaner()

	registerHandlers()

	gottp.MakeServer(&config.Settings)
}
