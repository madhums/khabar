package core

import (
	"testing"

	"github.com/bulletind/khabar/db"
)

const (
	validName   = "incidentapp"
	invalidName = "sheldon-cooper"
	dbName      = "notifications_test"
)

func init() {
	db.Conn = db.GetConn(dbName, "localhost")
	setup()
}

func setup() {
	var t *testing.T
	available := db.AvailableTopic{
		Ident:    "high_prio_log_incoming",
		AppName:  validName,
		Channels: []string{"email", "web", "push"},
	}
	available.PrepareSave()
	err := db.Conn.Session.DB(dbName).C(db.AvailableTopicCollection).Insert(available)
	if err != nil {
		t.Error("Unable to setup db for testing")
	}
}

func TestValidCategory(t *testing.T) {
	if validCategory(invalidName) {
		t.Error(invalidName, "is not a valid category/app_name")
		return
	}

	if !validCategory(validName) {
		t.Error(validName, "is not a valid category/app_name")
		return
	}

	defer cleanup()
}

func TestGetParseKeys(t *testing.T) {
	channelData := getParseKeys(validName)

	if _, ok := channelData["parse_application_id"]; !ok {
		t.Error("parse_application_id is not set for category", validName)
	}

	if _, ok := channelData["parse_rest_api_key"]; !ok {
		t.Error("parse_rest_api_key is not set for category", validName)
	}

	defer cleanup()
}

func cleanup() {
	db.Conn.Session.DB(dbName).DropDatabase()
}
