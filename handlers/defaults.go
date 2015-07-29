package handlers

import (
	"log"
	"net/http"

	"github.com/bulletind/khabar/core"
	"github.com/bulletind/khabar/dbapi/topics"
	"github.com/bulletind/khabar/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/simversity/gottp.v3"
)

type Defaults struct {
	gottp.BaseHandler
}

func (self *Defaults) Post(request *gottp.Request) {
	channel := request.GetArgument("channel").(string)
	ident := request.GetArgument("ident").(string)
	org := request.GetArgument("org").(string)

	if !core.IsChannelAvailable(channel) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	err := topics.InsertOrUpdateTopic(org, ident, channel, "Default")

	if err != nil {
		log.Println(err)
		request.Raise(gottp.HttpError{http.StatusInternalServerError,
			"Unable to complete db operation."})
		return
	}

	request.Write(utils.R{
		Data:       nil,
		Message:    "true",
		StatusCode: http.StatusNoContent,
	})

	return

}

func (self *Defaults) Delete(request *gottp.Request) {
	channel := request.GetArgument("channel").(string)
	ident := request.GetArgument("ident").(string)
	org := request.GetArgument("org").(string)

	if !core.IsChannelAvailable(channel) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	err := topics.InsertOrUpdateTopic(org, ident, channel, "Default")

	if err != nil {
		if err != mgo.ErrNotFound {
			request.Raise(gottp.HttpError{http.StatusInternalServerError,
				"Unable to complete db operation."})
			return
		}
	}

	request.Write(utils.R{StatusCode: http.StatusNoContent, Data: nil,
		Message: "NoContent."})
	return

}
