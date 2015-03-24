package handlers

import (
	"log"
	"net/http"

	"github.com/changer/khabar/db"
	"github.com/changer/khabar/dbapi/gully"
	"github.com/changer/khabar/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/simversity/gottp.v2"
)

type Gully struct {
	gottp.BaseHandler
}

func (self *Gully) Post(request *gottp.Request) {

	inputGully := new(gully.Gully)
	request.ConvertArguments(inputGully)
	inputGully.PrepareSave()

	if !inputGully.IsValid(db.INSERT_OPERATION) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user, org and app_name must be present."})
		return
	}

	if !utils.ValidateAndRaiseError(request, inputGully) {
		return
	}

	gly, err := gully.Get(db.Conn, inputGully.User, inputGully.AppName, inputGully.Organization, inputGully.Ident)

	if err != nil {
		if err != mgo.ErrNotFound {
			log.Println(err)
			request.Raise(gottp.HttpError{http.StatusInternalServerError, "Unable to fetch data, Please try again later."})
		} else {
			request.Raise(gottp.HttpError{http.StatusNotFound, "Not Found."})
		}
		return
	}

	if gly != nil {
		request.Raise(gottp.HttpError{http.StatusConflict, "Channel already exists"})
		return
	}

	gully.Insert(db.Conn, inputGully)

	request.Write(utils.R{StatusCode: http.StatusCreated, Data: inputGully.Id, Message: "Created"})
	return
}

func (self *Gully) Delete(request *gottp.Request) {
	gly := new(gully.Gully)
	request.ConvertArguments(gly)
	if !gly.IsValid(dbapi.DELETE_OPERATION) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user, org and app_name must be present."})
		return
	}

	err := gully.Delete(db.Conn, &utils.M{"app_name": gly.AppName,
		"org": gly.Organization, "user": gly.User, "ident": gly.Ident})
	if err != nil {
		log.Println(err)
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Unable to delete."})
		return
	}

	request.Write(utils.R{StatusCode: http.StatusNoContent, Data: nil, Message: "NoContent"})
	return
}

type Gullys struct {
	gottp.BaseHandler
}

func (self *Gullys) Get(request *gottp.Request) {
	var args struct {
		Organization string `json:"org"`
		AppName      string `json:"app_name"`
		User         string `json:"user"`
	}

	request.ConvertArguments(&args)

	all, err := gully.GetAll(db.Conn, args.User, args.AppName, args.Organization)

	if err != nil {
		if err != mgo.ErrNotFound {
			log.Println(err)
			request.Raise(gottp.HttpError{http.StatusInternalServerError, "Unable to fetch data, Please try again later."})
		} else {
			request.Raise(gottp.HttpError{http.StatusNotFound, "Not Found."})
		}
		return
	}

	request.Write(all)
	return
}
