package handlers

import (
	"log"
	"net/http"

	"gopkg.in/bulletind/khabar.v1/core"
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/defaults"
	"gopkg.in/bulletind/khabar.v1/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/simversity/gottp.v3"
)

type Defaults struct {
	gottp.BaseHandler
}

func (self *Defaults) Post(request *gottp.Request) {
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
	defaultPref.PrepareSave()

	defaultPref.Channels = []string{channelIdent}

	request.ConvertArguments(defaultPref)
	defaultPref.Topic = topicIdent

	if !utils.ValidateAndRaiseError(request, defaultPref) {
		return
	}

	if defaults.IsDefaultExists(defaultPref.Organization, defaultPref.Topic, channelIdent, defaultPref.Enabled) {
		request.Raise(gottp.HttpError{http.StatusConflict,
			"Already Exists."})
		return
	}

	if defaults.IsDefaultExists(defaultPref.Organization, defaultPref.Topic, channelIdent, !defaultPref.Enabled) {
		request.Raise(gottp.HttpError{http.StatusConflict,
			"Already Set to Opposite. Delete it and retry."})
		return
	}

	err, existingObj := defaults.Get(defaultPref.Organization, defaultPref.Topic, defaultPref.Enabled)

	if err != nil {
		if err == mgo.ErrNotFound {
			defaults.Insert(defaultPref)
			request.Write(utils.R{StatusCode: http.StatusCreated, Data: nil,
				Message: "Created"})
			return
		}
		request.Raise(gottp.HttpError{http.StatusInternalServerError,
			"Unable to complete db operation."})
		return
	}

	err = defaults.AddChannel(existingObj.Topic, channelIdent, existingObj.Organization, existingObj.Enabled)

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

	err := defaults.RemoveChannel(defaultPref.Topic, channelIdent, defaultPref.Organization, defaultPref.Enabled)

	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError,
			"Unable to complete db operation."})
		return
	}

	request.Write(utils.R{StatusCode: http.StatusNoContent, Data: nil,
		Message: "NoContent."})
	return

}
