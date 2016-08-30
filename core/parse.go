package core

import (
	"bytes"
	"os"

	"log"
	"net/http"

	gottpUtils "gopkg.in/simversity/gottp.v3/utils"

	"github.com/bulletind/khabar/config"
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/saved_item"
	khabarUtils "github.com/bulletind/khabar/utils"
)

const (
	PARSE_URL  = "https://api.parse.com/1/push"
	PUSH_SOUND = "default"
)

var parseKeys = []Parse{
	Parse{"APP_ID", "parse_application_id"},
	Parse{"API_KEY", "parse_rest_api_key"},
}

type Parse struct {
	Name string
	Key  string
}

func parseHandler(item *db.PendingItem, text string, locale string, appName string) {
	log.Println("Sending Push Notification...")

	settings := getParseKeys(appName)

	application_id, ok := settings["parse_application_id"].(string)
	if !ok {
		log.Println("parse_application_id is a required parameter.")
		return
	}

	api_key, ok := settings["parse_rest_api_key"].(string)
	if !ok {
		log.Println("parse_rest_api_key is a required parameter.")
		return
	}

	subject, ok := item.Context["subject"].(string)
	if !ok || subject == "" {
		subject = item.Topic
	}

	body := map[string]interface{}{}
	body["alert"] = subject
	body["title"] = subject
	body["message"] = text
	body["entity"] = item.Entity
	body["organization"] = item.Organization
	body["app_name"] = item.AppName
	body["topic"] = item.Topic
	body["created_on"] = item.CreatedOn
	body["sound"] = PUSH_SOUND
	body["badge"] = "Increment"

	data := map[string]interface{}{}
	data["data"] = body
	data["channels"] = []string{"USER_" + item.User}

	var jsonStr = gottpUtils.Encoder(&data)

	req, err := http.NewRequest("POST", PARSE_URL, bytes.NewBuffer(jsonStr))

	req.Header.Set("X-Parse-Application-Id", application_id)
	req.Header.Set("X-Parse-REST-API-Key", api_key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		config.Tracer.Notify()
	}
	defer resp.Body.Close()

	saved_item.Insert(db.SavedPushCollection, &db.SavedItem{Data: data, Details: *item})

	log.Println("Push notification status:", resp.Status)
}

// getParseKeys returns map of parse api key and app id
// It gets the values from the enviroment variables
func getParseKeys(appName string) khabarUtils.M {
	doc := khabarUtils.M{}

	// Set the Parse api key and id
	for _, parse := range parseKeys {
		envKey := "PARSE_" + appName + "_" + parse.Name
		doc[parse.Key] = os.Getenv(envKey)
		if len(os.Getenv(envKey)) == 0 {
			log.Println(envKey, "is empty. Make sure you set this env variable")
		}
	}
	return doc
}
