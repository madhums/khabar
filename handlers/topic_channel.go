package handlers

import (
	"net/http"

	"github.com/bulletind/khabar/core"
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/topics"
	"github.com/bulletind/khabar/utils"
	"gopkg.in/simversity/gottp.v3"
)

type TopicChannel struct {
	gottp.BaseHandler
}

func (self *TopicChannel) Delete(request *gottp.Request) {
	channelName := request.GetArgument("channel").(string)
	if !core.IsChannelAvailable(channelName) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	topic := new(db.Topic)
	request.ConvertArguments(topic)

	err := topics.RemoveChannel(topic.Ident, channelName, topic.User, topic.Organization)
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
	channelName := request.GetArgument("channel").(string)
	if !core.IsChannelAvailable(channelName) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	topic := new(db.Topic)
	request.ConvertArguments(topic)

	err := topics.AddChannel(topic.Ident, channelName, topic.User, topic.Organization)
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
