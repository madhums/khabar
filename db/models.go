package db

const (
	DELETE_OPERATION = 1
	INSERT_OPERATION = 2
	UPDATE_OPERATION = 3

	SentCollection       = "sent_notifications"
	StatsCollection      = "last_seen_at"
	TopicCollection      = "topics"
	GullyCollection      = "gullys"
	UserLocaleCollection = "user_locales"
)

type SentItem struct {
	BaseModel      `bson:",inline"`
	CreatedBy      string                 `json:"created_by" bson:"created_by" required:"true"`
	Organization   string                 `json:"org" bson:"org" required:"true"`
	AppName        string                 `json:"app_name" bson:"app_name" required:"true"`
	Topic          string                 `json:"topic" bson:"topic" required:"true"`
	User           string                 `json:"user" bson:"user" required:"true"`
	DestinationUri string                 `json:"destination_uri" bson:"destination_uri" required:"true"`
	Text           string                 `json:"text" bson:"text" required:"true"`
	IsRead         bool                   `json:"is_read" bson:"is_read"`
	Context        map[string]interface{} `json:"context" bson:"context"`
	Entity         string                 `json:"entity" bson:"entity" required:"true"`
}

func (self *SentItem) IsValid() bool {
	if len(self.Text) == 0 {
		return false
	}
	return true
}

type SavedItem struct {
	BaseModel `bson:",inline"`
	Data      interface{} `bson:"data"`
}

type LastSeen struct {
	BaseModel    `bson:",inline"`
	User         string `json:"user" bson:"user" required:"true"`
	Organization string `json:"org" bson:"org"`
	AppName      string `json:"app_name" bson:"app_name"`
	Timestamp    int64  `json:"timestamp" bson:"timestamp" required:"true"`
}

type Gully struct {
	BaseModel    `bson:",inline"`
	User         string                 `json:"user" bson:"user"`
	Organization string                 `json:"org" bson:"org"`
	AppName      string                 `json:"app_name" bson:"app_name"`
	Data         map[string]interface{} `json:"data" bson:"data" required:"true"`
	Ident        string                 `json:"ident" bson:"ident" required:"true"`
}

func (self *Gully) IsValid(op_type int) bool {

	if len(self.Ident) == 0 {
		return false
	}

	if op_type == INSERT_OPERATION {
		if len(self.Data) == 0 {
			return false
		}

	}

	return true
}

type UserLocale struct {
	BaseModel `bson:",inline"`
	User      string `json:"user" bson:"user" required:"true"`
	Locale    string `json:"locale" bson:"locale" required:"true"`
	TimeZone  string `json:"timezone" bson:"timezone" required:"true"`
}

func (self *UserLocale) IsValid() bool {
	if len(self.Locale) == 0 || len(self.User) == 0 || len(self.TimeZone) == 0 {
		return false
	}
	return true
}
