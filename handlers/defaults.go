package handlers

import (
	"log"
	"net/http"

	"gopkg.in/bulletind/khabar.v1/core"
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/defaults"
	"gopkg.in/bulletind/khabar.v1/dbapi/topics"
	"gopkg.in/bulletind/khabar.v1/utils"
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

	err := topics.InsertOrUpdateTopic(org, ident, channel)

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
	channelIdent := request.GetArgument("channel").(string)
	topicIdent := request.GetArgument("ident").(string)

	if !core.IsChannelAvailable(channelIdent) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	defaultPref := new(db.Defaults)
	request.ConvertArguments(defaultPref)
	defaultPref.Topic = topicIdent

	if !defaults.IsDefaultExists(defaultPref.Organization, defaultPref.Topic, channelIdent, defaultPref.Enabled) {
		request.Raise(gottp.HttpError{http.StatusNotFound,
			"Does not Exists."})
		return
	}

	err, dbDefault := defaults.Get(defaultPref.Organization, defaultPref.Topic, defaultPref.Enabled)

	if err != nil {
		if err != mgo.ErrNotFound {
			request.Raise(gottp.HttpError{http.StatusInternalServerError,
				"Unable to complete db operation."})
			return
		}
	} else {
		if len(dbDefault.Channels) == 1 {
			defaults.Delete(&utils.M{"_id": dbDefault.Id})
		} else {
			defaults.RemoveChannel(dbDefault.Topic, channelIdent,
				dbDefault.Organization, dbDefault.Enabled)
		}
	}

	request.Write(utils.R{StatusCode: http.StatusNoContent, Data: nil,
		Message: "NoContent."})
	return

}
