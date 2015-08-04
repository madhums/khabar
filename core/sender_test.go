package core

import (
	"testing"

	"github.com/bulletind/khabar/db"
)

const (
	validName   = "incidentapp"
	invalidName = "sheldon-cooper"
)

func init() {
	db.Conn = db.GetConn("notifications_testing", "localhost")
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
}

func TestGetParseKeys(t *testing.T) {
	channelData := getParseKeys(validName)

	if _, ok := channelData["parse_application_id"]; !ok {
		t.Error("parse_application_id is not set for category", validName)
	}

	if _, ok := channelData["parse_rest_api_key"]; !ok {
		t.Error("parse_rest_api_key is not set for category", validName)
	}
}
