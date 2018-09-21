package db

const BLANK = ""

const (
	DELETE_OPERATION = 1
	INSERT_OPERATION = 2
	UPDATE_OPERATION = 3

	SentCollection       = "sent_notifications"
	StatsCollection      = "last_seen_at"
	TopicCollection      = "topics"
	UserLocaleCollection = "user_locales"
	DeviceCollection     = "devices"

	SavedEmailCollection     = "saved_email"
	SavedPushCollection      = "saved_push"
	SavedWebCollection       = "saved_web"
	AvailableTopicCollection = "topics_available"

	ProcessedCollection = "processed"
)

type Processed struct {
	BaseModel    `bson:",inline"`
	User         string `bson:"user"`
	Organization string `bson:"org"`
}

type AvailableTopic struct {
	BaseModel `bson:",inline"`
	Ident     string   `json:"ident" bson:"ident" required:"true" form:"ident" binding:"required"`
	AppName   string   `json:"app_name" bson:"app_name" required:"true" form:"app_name" binding:"required"`
	Channels  []string `json:"channels" bson:"channels" required:"true" form:"channels" binding:"required"`
}

type SentItem struct {
	BaseModel      `bson:",inline"`
	CreatedBy      string                 `json:"created_by" bson:"created_by" required:"true"`
	Organization   string                 `json:"org" bson:"org" required:"true"`
	AppName        string                 `json:"app_name" bson:"app_name" required:"true"`
	Topic          string                 `json:"topic" bson:"topic" required:"true"`
	User           string                 `json:"user" bson:"user" required:"true"`
	DestinationUri string                 `json:"destination_uri" bson:"destination_uri"`
	Text           string                 `json:"text" bson:"text" required:"true"`
	IsRead         bool                   `json:"is_read" bson:"is_read"`
	Context        map[string]interface{} `json:"context" bson:"context"`
	Entity         string                 `json:"entity" bson:"entity" required:"true"`
}

type SavedItem struct {
	BaseModel `bson:",inline"`
	Data      interface{} `bson:"data"`
	Details   PendingItem `bson:"details"`
}

type PendingItem struct {
	BaseModel      `bson:",inline"`
	CreatedBy      string                 `json:"created_by" bson:"created_by" required:"true"`
	Organization   string                 `json:"org" bson:"org" required:"true"`
	AppName        string                 `json:"app_name" bson:"app_name" required:"true"`
	Topic          string                 `json:"topic" bson:"topic" required:"true"`
	IsPending      bool                   `json:"is_pending" bson:"is_pending" required:"true"`
	User           string                 `json:"user" bson:"user" required:"true"`
	DestinationUri string                 `json:"destination_uri" bson:"destination_uri"`
	Context        map[string]interface{} `json:"context" bson:"context" required:"true"`
	IsRead         bool                   `json:"is_read" bson:"is_read"`
	Entity         string                 `json:"entity" bson:"entity" required:"true"`
	DeviceTokens   []DeviceToken          `json:"device_tokens" bson:"device_tokens"`
	Download       Attachment             `json:"download" bson:"download"`
	Attachments    []Attachment           `json:"attachments" bson:"attachments"`
}

type DeviceToken struct {
	Type         string `json:"type" bson:"type"`
	Token        string `json:"token" bson:"token"`
	AppName      string `json:"app_name" bson:"app_name"`
	AppVariant   string `bson:"app_variant" json:"app_variant"`
	AppVersion   string `bson:"app_version" json:"app_version"`
	AppUniversal bool   `bson:"app_universal" json:"app_universal"`
}

func (self *PendingItem) IsValid() bool {
	if len(self.Context) == 0 {
		return false
	}
	return true
}

type Attachment struct {
	Url          string `json:"url" bson:"url"`
	ThumbnailUrl string `json:"thumbnail_url" bson:"thumbnail_url"`
	Name         string `json:"name" bson:"name"`
	Type         string `json:"type" bson:"type"`
	IsPrivate    bool   `json:"is_private" bson:"is_private"`
}

type LastSeen struct {
	BaseModel    `bson:",inline"`
	User         string `json:"user" bson:"user" required:"true"`
	Organization string `json:"org" bson:"org"`
	AppName      string `json:"app_name" bson:"app_name"`
	Timestamp    int64  `json:"timestamp" bson:"timestamp" required:"true"`
}

type Channel struct {
	Name    string `json:"name" bson:"name"`
	Enabled bool   `json:"enabled" bson:"enabled"`
	Default bool   `json:"default" bson:"default"`
	Locked  bool   `json:"locked" bson:"locked"`
}

type Topic struct {
	BaseModel    `bson:",inline"`
	User         string    `json:"user" bson:"user"`
	Organization string    `json:"org" bson:"org"`
	Channels     []Channel `json:"channels" bson:"channels" required:"true"`
	Ident        string    `json:"ident" bson:"ident" required:"true" form:"ident" binding:"required"`
}

type Device struct {
	BaseModel   `bson:",inline"`
	Token       string `json:"token" bson:"token" required:"true"`
	Type        string `json:"type" bson:"type" required:"true"`
	EndpointArn string `json:"endpoint_arn" bson:"endpoint_arn" required:"true"`
}
