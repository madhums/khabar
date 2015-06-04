package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/nicksnyder/go-i18n/i18n"
	"gopkg.in/bulletind/khabar.v1/config"
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/simversity/gottp.v3"
)

func sysInit() {
	<-(gottp.SysInitChan) //Buffered Channel to receive the server upstart boolean

	config.InitTracer()
	log.Println("Initialized GoTracer")

	db.Conn = db.GetConn(
		config.Settings.Khabar.DBName,
		config.Settings.Khabar.DBAddress,
		config.Settings.Khabar.DBUsername,
		config.Settings.Khabar.DBPassword,
	)

	log.Println("Database Connected :" + config.Settings.Khabar.DBName + " " +
		"at address:" + config.Settings.Khabar.DBAddress)

	transDir := config.Settings.Khabar.TranslationDirectory

	if len(transDir) == 0 {
		cwd := os.Getenv("PWD")
		transDir = cwd + "/translations"
		config.Settings.Khabar.TranslationDirectory = cwd
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
	log.Println("Translations have been parsed.")
}

func main() {
	cores := runtime.NumCPU()
	log.Println("Setting no. of Cores as", cores)
	runtime.GOMAXPROCS(cores)

	go sysInit()

	registerHandlers()

	gottp.MakeServer(&config.Settings)
}
