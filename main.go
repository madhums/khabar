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
	registerHandler("hello", "/hello/(?P<user_name>\\w+$)", new(SampleHandler))
	gottp.MakeServer(&settings)
}
