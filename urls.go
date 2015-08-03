package main

import (
	"github.com/bulletind/khabar/handlers"
	"gopkg.in/simversity/gottp.v3"
)

func registerHandlers() {
	// List all notifications
	gottp.NewUrl("notifications", "^/notifications/?$", new(handlers.Notifications))

	// Mark a notification as read
	gottp.NewUrl("notification", "^/notifications/(?P<_id>\\w+)/?$", new(handlers.Notification))

	// Get notification stats
	gottp.NewUrl("stats", "^/notifications/stats/?$", new(handlers.Stats))

	// Set/Unset user preference
	gottp.NewUrl("topic_channel", "^/topics/(?P<ident>\\w+)/channels/(?P<channel>\\w+)/?$", new(handlers.TopicChannel))

	// List user/org preferences
	gottp.NewUrl("topics", "^/topics/?$", new(handlers.Topics))

	// Set/Unset org defaults
	gottp.NewUrl("defaultTopics", "^/topics/defaults/(?P<ident>\\w+)/channels/(?P<channel>\\w+)/?$", new(handlers.Defaults))

	// Set/Unset org locked
	gottp.NewUrl("lockedTopics", "^/topics/locked/(?P<ident>\\w+)/channels/(?P<channel>\\w+)/?$", new(handlers.Locks))

	// Store and update user locale
	gottp.NewUrl("user_locale", "^/locales/(?P<user>\\w+)/?$", new(handlers.UserLocale))

	gottp.NewUrl("snsBounce", "^/sns/bounce/?$", new(handlers.SnsBounce))
	gottp.NewUrl("snsComplain", "^/sns/complaint/?$", new(handlers.SnsComplaint))
	gottp.NewUrl("mandrillBounce", "^/mandrill/bounce/?$", new(handlers.MandrillBounce))

	gottp.NewUrl("channels", "^/channels/?$", new(handlers.Gullys))
	gottp.NewUrl("channel", "^/channels/(?P<ident>\\w+)/?$", new(handlers.Gully))

	gottp.NewUrl("topic", "^/topics/(?P<ident>\\w+)/?$", new(handlers.Topic))

}
