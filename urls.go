package main

import (
	"gopkg.in/simversity/gottp.v1"
)

func registerHandler(name string, pattern string, handler gottp.Handler) {
	gottp.NewUrl(name, pattern, handler)
}
