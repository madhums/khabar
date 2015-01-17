package main

import (
	"github.com/Simversity/gottp"
	"log"
)

func sysInit() {
	<-(gottp.SysInitChan) //Buffered Channel to receive the server upstart boolean
	log.Println("System is ready to Serve")
}

func main() {
	go sysInit()
	registerHandler("hello", "/hello/\\w{3,5}/?$", new(SampleHandler))
	gottp.MakeServer(&settings)
}
