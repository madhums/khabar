package handlers

import (
	"log"
	"net/http"

	"github.com/changer/khabar/db"
	"github.com/changer/khabar/dbapi"
	"github.com/changer/khabar/dbapi/topics"
	"github.com/changer/khabar/utils"
	"gopkg.in/simversity/gottp.v2"
	gottp_utils "gopkg.in/simversity/gottp.v2/utils"
)

type TopicChannel struct {
	gottp.BaseHandler
}

func (self *TopicChannel) Post(request *gottp.Request) {
	intopic := new(topics.Topic)

	channelIdent := request.GetArgument("channel").(string)
	intopic.Ident = request.GetArgument("ident").(string)

	request.ConvertArguments(intopic)

	topic := topics.Get(db.Conn, intopic.User, intopic.AppName, intopic.Organization, intopic.Ident)

	hasData := true

	if topic == nil {
		hasData = false
		log.Println("Creating new document")
		intopic.AddChannel(channelIdent)

		intopic.PrepareSave()
		if !intopic.IsValid(dbapi.INSERT_OPERATION) {
			request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user, org and app_name must be present."})
			return
		}

		if !utils.ValidateAndRaiseError(request, intopic) {
			log.Println("Validation Failed")
			return
		}

		topic = intopic

	} else {
		topic.AddChannel(channelIdent)
	}

	var err error
	if hasData {
		err = topics.Update(db.Conn, topic.User, topic.AppName, topic.Organization, topic.Ident, &utils.M{
			"channels": topic.Channels,
		})
		if err != nil {
			log.Println("Error while inserting document :" + err.Error())
			request.Raise(gottp.HttpError{http.StatusInternalServerError, "Internal server error."})
			return
		} else {
			request.Write(gottp.HttpError{http.StatusNoContent, "NoContent"})
		}
	} else {
		log.Println("Successfull call: Inserting document")
		topics.Insert(db.Conn, topic)
		request.Write(gottp.HttpError{http.StatusCreated, string(gottp_utils.Encoder(topic))})
	}
}

func (self *TopicChannel) Delete(request *gottp.Request) {
	topic := new(topics.Topic)

	channelIdent := request.GetArgument("channel").(string)
	topic.Ident = request.GetArgument("ident").(string)

	request.ConvertArguments(topic)

	topic = topics.Get(db.Conn, topic.User, topic.AppName, topic.Organization, topic.Ident)

	if topic == nil {
		request.Raise(gottp.HttpError{http.StatusNotFound, "topics setting does not exists."})
		return
	}

	topic.RemoveChannel(channelIdent)
	log.Println(topic.Channels)

	var err error

	if len(topic.Channels) == 0 {
		log.Println("Deleting from database, since channels are now empty.")
		err = topics.Delete(db.Conn, &utils.M{"app_name": topic.AppName,
			"org": topic.Organization, "user": topic.User, "ident": topic.Ident})
	} else {
		log.Println("Updating...")
		err = topics.Update(db.Conn, topic.User, topic.AppName, topic.Organization, topic.Ident, &utils.M{
			"channels": topic.Channels,
		})
	}

	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Internal server error."})
		return
	}

	request.Write(gottp.HttpError{http.StatusNoContent, "NoContent"})

}

type Topic struct {
	gottp.BaseHandler
}

func (self *Topic) Delete(request *gottp.Request) {
	topic := new(topics.Topic)
	request.ConvertArguments(topic)
	if !topic.IsValid(dbapi.DELETE_OPERATION) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user, org and app_name must be present."})
		return
	}
	err := topics.Delete(db.Conn, &utils.M{"app_name": topic.AppName,
		"org": topic.Organization, "user": topic.User, "ident": topic.Ident})
	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Unable to delete."})
		return
	}

	request.Write(gottp.HttpError{http.StatusNoContent, "NoContent"})
}

type Topics struct {
	gottp.BaseHandler
}

func (self *Topics) Get(request *gottp.Request) {
	var args struct {
		Organization string `json:"org"`
		AppName      string `json:"app_name"`
		User         string `json:"user"`
	}

	request.ConvertArguments(&args)

	all := topics.GetAll(db.Conn, args.User, args.AppName, args.Organization)

	request.Write(all)
}
