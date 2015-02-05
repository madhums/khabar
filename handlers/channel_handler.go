package handlers

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/models"
	"github.com/parthdesai/sc-notifications/utils"
	"gopkg.in/simversity/gottp.v1"
	"net/http"
)

type ChannelHandler struct {
	gottp.BaseHandler
}

func (self *ChannelHandler) Post(request *gottp.Request) {

	channel := new(models.Channel)
	request.ConvertArguments(channel)
	channel.PrepareSave(db.DbConnection)

	if !channel.IsValid() {
		request.Raise(gottp.HttpError{http.StatusPreconditionFailed, "Atleast one of the user_id, org_id and app_id must be present."})
		return
	}

	if !utils.ValidateAndRaiseError(request, channel) {
		return
	}

	hasData := channel.GetFromDatabase(db.DbConnection)

	if hasData {
		request.Raise(gottp.HttpError{http.StatusPreconditionFailed, "Channel already exists"})
		return
	}

	db.DbConnection.Insert("channels", channel)
	request.Write(channel)
}

func (self *ChannelHandler) Delete(request *gottp.Request) {
	channel := new(models.Channel)
	request.ConvertArguments(channel)
	if !channel.IsValid() {
		request.Raise(gottp.HttpError{http.StatusPreconditionFailed, "Atleast one of the user_id, org_id and app_id must be present."})
		return
	}
	err := channel.DeleteFromDatabase(db.DbConnection)
	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Unable to delete."})
	}
}
