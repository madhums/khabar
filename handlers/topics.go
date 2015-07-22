package handlers

import (
	"log"
	"net/http"

	"gopkg.in/bulletind/khabar.v1/core"
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/available_topics"
	"gopkg.in/bulletind/khabar.v1/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/simversity/gottp.v3"
)

type Topics struct {
	gottp.BaseHandler
}

func (self *Topics) Get(request *gottp.Request) {
	var args struct {
		AppName      string `json:"app_name" required:"true"`
		Organization string `json:"org"`
		User         string `json:"user"`
	}
	var iter map[string]available_topics.ChotaTopic
	var err error

	request.ConvertArguments(&args)

	if !utils.ValidateAndRaiseError(request, args) {
		log.Println("Validation Failed")
		return
	}

	appTopics, err := available_topics.GetAppTopics(args.AppName, args.Organization)

	if args.Organization == "" {
		request.Write(appTopics)
		return
	}

	channels := []string{}
	for ident, _ := range core.ChannelMap {
		channels = append(channels, ident)
	}

	if args.User == "" {
		iter, err = available_topics.GetOrgTopics(args.Organization, appTopics, &channels)
	} else {
		iter, err = available_topics.GetUserTopics(args.User, args.Organization, appTopics, &channels)
	}

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

	request.Write(iter)
	return
}

func (self *Topics) Post(request *gottp.Request) {
	newTopic := new(db.AvailableTopic)
	newTopic.Channels = []string{"email", "web", "push"}

	request.ConvertArguments(newTopic)

	newTopic.PrepareSave()

	if !utils.ValidateAndRaiseError(request, newTopic) {
		log.Println("Validation Failed")
		return
	}

	if _, err := available_topics.Get(newTopic.Ident); err == nil {
		request.Raise(gottp.HttpError{
			http.StatusConflict,
			"Topic already exists"})
		return
	} else {
		if err != mgo.ErrNotFound {
			log.Println(err)
			request.Raise(gottp.HttpError{
				http.StatusInternalServerError,
				"Unable to fetch data, Please try again later.",
			})
			return
		}
	}

	available_topics.Insert(newTopic)

	request.Write(utils.R{
		StatusCode: http.StatusCreated,
		Data:       newTopic.Id,
		Message:    "Created",
	})
	return
}
