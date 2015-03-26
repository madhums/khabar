package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bulletind/khabar/config"
	"github.com/bulletind/khabar/db"
	"github.com/nicksnyder/go-i18n/i18n"
	"gopkg.in/simversity/gottp.v2"
)

func sysInit() {
	<-(gottp.SysInitChan) //Buffered Channel to receive the server upstart boolean

	db.Conn = db.GetConn(config.Settings.Khabar.DBName,
		config.Settings.Khabar.DBAddress, config.Settings.Khabar.DBUsername,
		config.Settings.Khabar.DBPassword)

	log.Println("Database Connected :" + config.Settings.Khabar.DBName + " " +
		"at address:" + config.Settings.Khabar.DBAddress)

	transDir := config.Settings.Khabar.TranslationDirectory

	if len(transDir) == 0 {
		transDir = os.Getenv("PWD") + "/translations"
	}

	log.Println("Directory for translation :" + transDir)

	filepath.Walk(transDir, func(path string, _ os.FileInfo, err error) error {
		fileExt := filepath.Ext(path)
		if fileExt == ".json" && err == nil {
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
	log.Println("Translation has been parsed.")
}

func main() {
	go sysInit()

	registerHandlers()

	gottp.MakeServer(&config.Settings)
}
