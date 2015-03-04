package handlers

import (
	"log"
	"net/http"

	"github.com/changer/khabar/db"
	"github.com/changer/khabar/dbapi"
	"github.com/changer/khabar/dbapi/topics"
	"github.com/changer/khabar/utils"
	"gopkg.in/simversity/gottp.v2"
)

type TopicChannel struct {
	gottp.BaseHandler
}

func (self *TopicChannel) Post(request *gottp.Request) {
	topic := new(topics.Topic)

	channelIdent := request.GetArgument("channel").(string)
	topic.Ident = request.GetArgument("ident").(string)

	request.ConvertArguments(topic)

	topic.AddChannel(channelIdent)

	topic = topics.Get(db.Conn, topic.User, topic.AppName, topic.Organization, topic.Ident)

	hasData := true

	if topic == nil {
		hasData = false
		log.Println("Creating new document")

		topic.PrepareSave()
		if !topic.IsValid(dbapi.INSERT_OPERATION) {
			request.Raise(gottp.HttpError{http.StatusBadRequest, "Atleast one of the user, org and app_name must be present."})
			return
		}

	} else {
		topic.AddChannel(channelIdent)
	}

	if !utils.ValidateAndRaiseError(request, topic) {
		log.Println("Validation Failed")
		return
	}

	var err error
	if hasData {
		err = topics.Update(db.Conn, topic.User, topic.AppName, topic.Organization, topic.Ident, &db.M{
			"channels": topic.Channels,
		})
	} else {
		log.Println("Successfull call: Inserting document")
		topics.Insert(db.Conn, topic)
	}

	if err != nil {
		log.Println("Error while inserting document :" + err.Error())
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Internal server error."})
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
		err = topics.Delete(db.Conn, &db.M{"app_name": topic.AppName,
			"org": topic.Organization, "user": topic.User, "type": topic.Ident})
	} else {
		log.Println("Updating...")
		err = topics.Update(db.Conn, topic.User, topic.AppName, topic.Organization, topic.Ident, &db.M{
			"channels": topic.Channels,
		})
	}

	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Internal server error."})
	}

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
	err := topics.Delete(db.Conn, &db.M{"app_name": topic.AppName,
		"org": topic.Organization, "user": topic.User, "type": topic.Ident})
	if err != nil {
		request.Raise(gottp.HttpError{http.StatusInternalServerError, "Unable to delete."})
	}
}
