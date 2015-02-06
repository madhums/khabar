package handlers

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/models"
	"github.com/parthdesai/sc-notifications/utils"
	"gopkg.in/simversity/gottp.v1"
	"net/http"
)

type UserLocalHandler struct {
	gottp.BaseHandler
}

func (self *UserLocalHandler) Put(request *gottp.Request) {
	userLocale := new(models.UserLocale)
	request.ConvertArguments(userLocale)

	if !userLocale.IsValid() {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "user_id, region_id and language_id must be present."})
		return
	}

	hasData := userLocale.GetFromDatabase(db.DbConnection)
	userLocalId := userLocale.Id
	request.ConvertArguments(userLocale)

	if !hasData {
		request.Raise(gottp.HttpError{http.StatusNotFound, "Unable to find user locale for user id."})
		return
	}

	userLocale.Id = userLocalId
	userLocale.Update(db.DbConnection)
}

func (self *UserLocalHandler) Post(request *gottp.Request) {
	userLocale := new(models.UserLocale)
	request.ConvertArguments(userLocale)
	userLocale.PrepareSave()

	if !userLocale.IsValid() {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "user_id, region_id and language_id must be present."})
		return
	}

	if !utils.ValidateAndRaiseError(request, userLocale) {
		return
	}

	hasData := userLocale.GetFromDatabase(db.DbConnection)

	if hasData {
		request.Raise(gottp.HttpError{http.StatusConflict, "User locale information already exists"})
		return
	}

	userLocale.InsertIntoDatabase(db.DbConnection)
}
