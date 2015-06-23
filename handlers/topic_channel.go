package handlers

import (
	"net/http"

	"gopkg.in/bulletind/khabar.v1/core"
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/topics"
	"gopkg.in/bulletind/khabar.v1/utils"
	"gopkg.in/simversity/gottp.v3"
)

type TopicChannel struct {
	gottp.BaseHandler
}

func (self *TopicChannel) Delete(request *gottp.Request) {
	channelIdent := request.GetArgument("channel").(string)
	if !core.IsChannelAvailable(channelIdent) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	intopic := new(db.Topic)
	request.ConvertArguments(intopic)

	err := topics.RemoveChannel(intopic.Ident, channelIdent, intopic.User, intopic.Organization)
	if err != nil {
		request.Raise(gottp.HttpError{http.StatusBadRequest, err.Error()})
		return
	}

	request.Write(utils.R{
		Data:       nil,
		Message:    "true",
		StatusCode: http.StatusNoContent,
	})

	return
}

func (self *TopicChannel) Post(request *gottp.Request) {
	channelIdent := request.GetArgument("channel").(string)
	if !core.IsChannelAvailable(channelIdent) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	intopic := new(db.Topic)
	request.ConvertArguments(intopic)

	err := topics.AddChannel(intopic.Ident, channelIdent, intopic.User, intopic.Organization)
	if err != nil {
		request.Raise(gottp.HttpError{http.StatusBadRequest, err.Error()})
		return
	}

	request.Write(utils.R{
		Data:       nil,
		Message:    "true",
		StatusCode: http.StatusNoContent,
	})

	return
}
