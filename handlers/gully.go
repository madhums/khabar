package handlers

import (
	"log"
	"net/http"

	"gopkg.in/bulletind/khabar.v1/core"
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/gully"
	"gopkg.in/bulletind/khabar.v1/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/simversity/gottp.v3"
)

type Gully struct {
	gottp.BaseHandler
}

func (self *Gully) Post(request *gottp.Request) {

	inputGully := new(db.Gully)
	request.ConvertArguments(inputGully)
	inputGully.PrepareSave()

	log.Println("Input :", inputGully)

	if !core.IsChannelAvailable(inputGully.Ident) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	if !inputGully.IsValid(db.INSERT_OPERATION) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Atleast one of the user, org and app_name must be present.",
		})

		return
	}

	if !utils.ValidateAndRaiseError(request, inputGully) {
		return
	}

	gly, err := gully.Get(
		inputGully.User,
		inputGully.AppName,
		inputGully.Organization,
		inputGully.Ident,
	)

	if err != nil {
		if err != mgo.ErrNotFound {
			log.Println(err)
			request.Raise(gottp.HttpError{
				http.StatusInternalServerError,
				"Unable to fetch data, Please try again later.",
			})
			return
		}
	}

	if gly != nil {
		request.Raise(gottp.HttpError{
			http.StatusConflict,
			"Channel already exists",
		})

		return
	}

	gully.Insert(inputGully)

	log.Println("Saving :", inputGully)

	request.Write(utils.R{
		StatusCode: http.StatusCreated,
		Data:       inputGully.Id,
		Message:    "Created",
	})

	return
}

func (self *Gully) Delete(request *gottp.Request) {
	gly := new(db.Gully)
	request.ConvertArguments(gly)
	if !gly.IsValid(db.DELETE_OPERATION) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Atleast one of the user, org and app_name must be present.",
		})

		return
	}

	err := gully.Delete(&utils.M{
		"app_name": gly.AppName,
		"org":      gly.Organization,
		"user":     gly.User,
		"ident":    gly.Ident})

	if err != nil {
		log.Println(err)
		request.Raise(gottp.HttpError{
			http.StatusInternalServerError,
			"Unable to delete.",
		})

		return
	}

	request.Write(utils.R{
		StatusCode: http.StatusNoContent,
		Data:       nil,
		Message:    "NoContent",
	})

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

	all, err := gully.GetAll(args.User,
		args.AppName, args.Organization)

	if err != nil {
		if err != mgo.ErrNotFound {
			log.Println(err)
			request.Raise(gottp.HttpError{
				http.StatusInternalServerError,
				"Unable to fetch data, Please try again later.",
			})

		} else {
			request.Raise(gottp.HttpError{
				http.StatusNotFound,
				"Not Found.",
			})
		}

		return
	}

	request.Write(all)
	return
}
