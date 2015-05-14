package config

import "gopkg.in/simversity/gotracer.v1"

var Tracer gotracer.Tracer

func InitTracer() {
	Tracer = gotracer.Tracer{
		Dummy:         Settings.Gottp.EmailDummy,
		EmailHost:     Settings.Gottp.EmailHost,
		EmailPort:     Settings.Gottp.EmailPort,
		EmailPassword: Settings.Gottp.EmailPassword,
		EmailUsername: Settings.Gottp.EmailUsername,
		EmailSender:   Settings.Gottp.EmailSender,
		EmailFrom:     Settings.Gottp.EmailFrom,
		ErrorTo:       Settings.Gottp.ErrorTo,
	}
}
