package core

import (
	"bytes"

	"log"
	"net/http"

	"gopkg.in/simversity/gottp.v3/utils"

	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/saved_item"
)

const PARSE_URL = "https://api.parse.com/1/push"

func pushHandler(
	item *db.PendingItem,
	text string,
	settings map[string]interface{},
) {
	log.Println("Sending Push Notification...")

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
	body["message"] = text
	body["entity"] = item.Entity
	body["organization"] = item.Organization
	body["app_name"] = item.AppName
	body["topic"] = item.Topic
	body["created_on"] = item.CreatedOn

	data := map[string]interface{}{}
	data["data"] = body
	data["channels"] = []string{"USER_" + item.User}

	var jsonStr = utils.Encoder(&data)

	req, err := http.NewRequest("POST", PARSE_URL, bytes.NewBuffer(jsonStr))

	req.Header.Set("X-Parse-Application-Id", application_id)
	req.Header.Set("X-Parse-REST-API-Key", api_key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	saved_item.Insert(db.SavedPushCollection, &db.SavedItem{Data: data, Details: *item})

	log.Println("Push notification status:", resp.Status)
}
