package main

import (
	"encoding/json"
	"flag"
	"github.com/changer/khabar/config"
	"github.com/changer/khabar/db"
	"github.com/changer/khabar/dbapi/gully"
	"log"
)

func main() {
	defGully := new(gully.Gully)

	defGully.PrepareSave()

	emailConfig := flag.String("emailConfig", "", "Default Email channel configuration")

	configFile := flag.String("config", "", "Configuration file used by khabar")

	flag.Parse()

	config.Settings.MakeConfig(*configFile)

	jsonBytes := []byte(*emailConfig)
	err := json.Unmarshal(jsonBytes, defGully)

	if err != nil {
		log.Println(err)
		return
	}

	conn := db.GetConn(config.Settings.Khabar.DBName, config.Settings.Khabar.DBAddress,
		config.Settings.Khabar.DBUsername, config.Settings.Khabar.DBPassword)

	gully.Insert(conn, defGully)
	return
}
