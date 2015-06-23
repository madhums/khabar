package handlers

import (
	"log"
	"net/http"

	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/user_locale"
	"gopkg.in/bulletind/khabar.v1/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/simversity/gottp.v3"
)

type UserLocale struct {
	gottp.BaseHandler
}

func (self *UserLocale) Put(request *gottp.Request) {
	inputUserLocale := new(db.UserLocale)
	request.ConvertArguments(inputUserLocale)

	if !inputUserLocale.IsValid() {
		request.Raise(gottp.HttpError{http.StatusBadRequest,
			"user, region_id and language_id must be present."})
		return
	}

	updateParams := make(utils.M)
	updateParams["timezone"] = inputUserLocale.TimeZone
	updateParams["locale"] = inputUserLocale.Locale

	err := user_locale.Update(inputUserLocale.User, &updateParams)

	if err != nil {
		log.Println(err)
		request.Raise(gottp.HttpError{http.StatusInternalServerError,
			"Unable to update."})
		return
	}

	request.Write(utils.R{Data: nil, Message: "NoContent",
		StatusCode: http.StatusNoContent})
	return
}

func (self *UserLocale) Post(request *gottp.Request) {
	userLocale := new(db.UserLocale)
	request.ConvertArguments(userLocale)
	userLocale.PrepareSave()

	if !userLocale.IsValid() {
		request.Raise(gottp.HttpError{http.StatusBadRequest,
			"user, region_id and language_id must be present."})
		return
	}

	if !utils.ValidateAndRaiseError(request, userLocale) {
		return
	}

	dblocale, err := user_locale.Get(userLocale.User)

	if err != nil {
		if err != mgo.ErrNotFound {
			log.Println(err)
			request.Raise(gottp.HttpError{http.StatusInternalServerError,
				"Unable to fetch data, Please try again later."})
			return
		}
	}

	if dblocale != nil {
		request.Raise(gottp.HttpError{http.StatusConflict,
			"User locale information already exists"})
		return
	}

	user_locale.Insert(userLocale)

	request.Write(utils.R{Data: userLocale.Id, Message: "Created",
		StatusCode: http.StatusCreated})
	return
}
