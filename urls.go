package main

import (
	"github.com/bulletind/khabar/handlers"
	"gopkg.in/simversity/gottp.v3"
)

func registerHandlers() {

	// Notifications
	//
	// GET - List all notifications
	// PUT - Mark all notifications as read
	gottp.NewUrl("notifications", "^/notifications/?$", new(handlers.Notifications))

	// Get notification stats
	// Add/Update timestamp in last_seen_at collection
	gottp.NewUrl("stats", "^/notifications/stats/?$", new(handlers.Stats))

	// Mark a notification as read
	gottp.NewUrl("notification", "^/notifications/(?P<_id>\\w+)/?$", new(handlers.Notification))

	// Preferences
	//
	// List user/org preferences
	gottp.NewUrl("topics", "^/topics/?$", new(handlers.Topics))

	// Set/Unset user preference
	gottp.NewUrl("topic_channel", "^/topics/(?P<ident>\\w+)/channels/(?P<channel>\\w+)/?$", new(handlers.TopicChannel))

	// Set/Unset org defaults or locked
	gottp.NewUrl("default_locked", "^/topics/(?P<type>\\w+)/(?P<ident>\\w+)/channels/(?P<channel>\\w+)/?$", new(handlers.DefaultLocked))

	// Store and update user locale
	gottp.NewUrl("user_locale", "^/locales/(?P<user>\\w+)/?$", new(handlers.UserLocale))

	gottp.NewUrl("snsBounce", "^/sns/bounce/?$", new(handlers.SnsBounce))
	gottp.NewUrl("snsComplain", "^/sns/complaint/?$", new(handlers.SnsComplaint))
	gottp.NewUrl("mandrillBounce", "^/mandrill/bounce/?$", new(handlers.MandrillBounce))
}
