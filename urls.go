package main

import (
	"github.com/changer/khabar/handlers"
	"gopkg.in/simversity/gottp.v2"
)

func registerHandlers() {
	gottp.NewUrl("notifications", "^/notifications/?$",
		new(handlers.Notifications))

	gottp.NewUrl("stats", "^/notifications/stats/?$",
		new(handlers.Stats))

	gottp.NewUrl("notification", "^/notification/(?P<_id>\\w+)/?$",
		new(handlers.Notification))

	gottp.NewUrl("channel", "^/channel/(?P<ident>\\w+)/?$",
		new(handlers.Gully))

	gottp.NewUrl("topic_channel",
		"^/(?P<ident>\\w+)/channel/(?P<channel>\\w+)/?$",
		new(handlers.TopicChannel))

	gottp.NewUrl("topic", "^/topic/(?P<topic>\\w+)/?$",
		new(handlers.Topic))

	gottp.NewUrl("user_locale", "^/locale/(?P<user>\\w+)/?$",
		new(handlers.UserLocale))
}
