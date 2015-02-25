package handlers

import (
	"github.com/changer/sc-notifications/db"
	"github.com/changer/sc-notifications/dbapi"
	"github.com/changer/sc-notifications/dbapi/gully"
	"github.com/changer/sc-notifications/utils"
	"gopkg.in/simversity/gottp.v2"
	"net/http"
)

type Gully struct {
	gottp.BaseHandler
}

func (self *Gully) Post(request *gottp.Request) {

	inputGully := new(gully.Gully)
	request.ConvertArguments(inputGully)
	inputGully.PrepareSave()

	if !inputGully.IsValid(dbapi.INSERT_OPERATION) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user, org and app_name must be present."})
		return
	}

	if !utils.ValidateAndRaiseError(request, inputGully) {
		return
	}

	gly := gully.Get(db.Conn, inputGully.User, inputGully.AppName, inputGully.Organization, inputGully.Ident)

	if gly != nil {
		request.Raise(gottp.HttpError{http.StatusConflict, "Channel already exists"})
		return
	}

	gully.Insert(db.Conn, inputGully)
	request.Write(inputGully)
}

func (self *Gully) Delete(request *gottp.Request) {
	gly := new(gully.Gully)
	request.ConvertArguments(gly)
	if !gly.IsValid(dbapi.DELETE_OPERATION) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user, org and app_name must be present."})
		return
	}
	err := gully.Delete(db.Conn, &db.M{"app_name": gly.AppName,
		"org": gly.Organization, "user": gly.User, "ident": gly.Ident})
	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Unable to delete."})
	}
}
