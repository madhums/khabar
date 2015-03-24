package db

const (
	EMAIL = "email"
	WEB   = "web"
	SMS   = "sms"
)

var Channels = map[string]bool{
	EMAIL: true,
	WEB:   true,
	SMS:   false,
}
