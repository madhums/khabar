package handlers

import (
	"github.com/parthdesai/sc-notifications/db"
	"github.com/parthdesai/sc-notifications/dbapi"
	"github.com/parthdesai/sc-notifications/dbapi/gully"
	"github.com/parthdesai/sc-notifications/utils"
	"gopkg.in/simversity/gottp.v1"
	"net/http"
)

type GullyHandler struct {
	gottp.BaseHandler
}

func (self *GullyHandler) Post(request *gottp.Request) {

	inputGully := new(gully.Gully)
	request.ConvertArguments(inputGully)
	inputGully.PrepareSave()

	if !inputGully.IsValid(dbapi.INSERT_OPERATION) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user_id, org_id and app_id must be present."})
		return
	}

	if !utils.ValidateAndRaiseError(request, inputGully) {
		return
	}

	gly := gully.GetFromDatabase(db.DbConnection, inputGully.UserID, inputGully.ApplicationID, inputGully.OrganizationID, inputGully.Ident)

	if gly != nil {
		request.Raise(gottp.HttpError{http.StatusConflict, "Channel already exists"})
		return
	}

	gully.InsertIntoDatabase(db.DbConnection, inputGully)
	request.Write(inputGully)
}

func (self *GullyHandler) Delete(request *gottp.Request) {
	gly := new(gully.Gully)
	request.ConvertArguments(gly)
	if !gly.IsValid(dbapi.DELETE_OPERATION) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user_id, org_id and app_id must be present."})
		return
	}
	err := gully.DeleteFromDatabase(db.DbConnection, gly)
	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Unable to delete."})
	}
}
