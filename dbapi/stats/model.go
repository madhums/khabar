package stats

import (
	"github.com/changer/khabar/db"
)

const (
	StatsCollection = "last_seen_at"
)

type LastSeen struct {
	db.BaseModel `bson:",inline"`
	User         string `json:"user" bson:"user" required:"true"`
	Organization string `json:"org" bson:"org"`
	AppName      string `json:"app_name" bson:"app_name"`
	Timestamp    int64  `json:"timestamp" bson:"timestamp" required:"true"`
}
