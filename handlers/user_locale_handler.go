package handlers

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/dbapi/user_locale"
	"github.com/parthdesai/sc-notifications/utils"
	"gopkg.in/simversity/gottp.v1"
	"net/http"
)

type UserLocalHandler struct {
	gottp.BaseHandler
}

func (self *UserLocalHandler) Put(request *gottp.Request) {
	userLocale := new(user_locale.UserLocale)
	request.ConvertArguments(userLocale)

	if !userLocale.IsValid() {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "user, region_id and language_id must be present."})
		return
	}

	userLocale = user_locale.GetFromDatabase(db.DbConnection, userLocale.User)

	if userLocale == nil {
		request.Raise(gottp.HttpError{http.StatusNotFound, "Unable to find user locale for user id."})
		return
	}

	userLocaleId := userLocale.Id

	request.ConvertArguments(userLocale)

	userLocale.Id = userLocaleId
	user_locale.Update(db.DbConnection, userLocale)
}

func (self *UserLocalHandler) Post(request *gottp.Request) {
	userLocale := new(user_locale.UserLocale)
	request.ConvertArguments(userLocale)
	userLocale.PrepareSave()

	if !userLocale.IsValid() {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "user, region_id and language_id must be present."})
		return
	}

	if !utils.ValidateAndRaiseError(request, userLocale) {
		return
	}

	if user_locale.GetFromDatabase(db.DbConnection, userLocale.User) != nil {
		request.Raise(gottp.HttpError{http.StatusConflict, "User locale information already exists"})
		return
	}

	user_locale.InsertIntoDatabase(db.DbConnection, userLocale)
}
