package main

import (
	"github.com/Simversity/gottp"
)

func registerHandler(name string, pattern string, handler gottp.Handler) {
	gottp.NewUrl(name, pattern, handler)
}
