package core

import (
	"testing"

	"github.com/bulletind/khabar/db"
)

const (
	validName   = "incidentapp"
	invalidName = "sheldon-cooper"
	dbName      = "notifications_test"
	dbUrl       = "mongodb://localhost/notifications_test"
)

func init() {
	db.Conn = db.GetConn(dbUrl, dbName)
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

	for _, p := range parseKeys {
		if _, ok := channelData[p.Key]; !ok {
			t.Error(p.Key, "is not set for category", validName)
		}
	}
}

// func TestGetEmailKeys(t *testing.T) {
// 	loadConfig()
//
// 	checkKey(t, settings.SMTP.HostName)
// 	checkKey(t, settings.SMTP.UserName)
// 	checkKey(t, settings.SMTP.Password)
// 	checkKey(t, settings.SMTP.Port)
// 	checkKey(t, settings.SMTP.From)
// }
//
// func checkKey(t *testing.T, key string) {
// 	if key == "" {
// 		t.Error("smtp_" + key + "is not set")
// 	}
// }

func cleanup() {
	db.Conn.Session.DB(dbName).DropDatabase()
}
