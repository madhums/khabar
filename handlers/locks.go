package handlers

import (
	"log"
	"net/http"

	"github.com/bulletind/khabar/dbapi/topics"

	"github.com/bulletind/khabar/core"
	"github.com/bulletind/khabar/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/simversity/gottp.v3"
)

type Locks struct {
	gottp.BaseHandler
}

func (self *Locks) Post(request *gottp.Request) {
	channelName := request.GetArgument("channel").(string)
	ident := request.GetArgument("ident").(string)
	org := request.GetArgument("org").(string)

	if !core.IsChannelAvailable(channelName) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	err := topics.InsertOrUpdateTopic(org, ident, channelName, "Locked", true, "")

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

func (self *Locks) Delete(request *gottp.Request) {
	channelName := request.GetArgument("channel").(string)
	ident := request.GetArgument("ident").(string)
	org := request.GetArgument("org").(string)

	if !core.IsChannelAvailable(channelName) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	err := topics.InsertOrUpdateTopic(org, ident, channelName, "Locked", false, "")

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
