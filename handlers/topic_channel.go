package handlers

import (
	"log"
	"net/http"

	"gopkg.in/bulletind/khabar.v1/core"
	"gopkg.in/bulletind/khabar.v1/db"
	"gopkg.in/bulletind/khabar.v1/dbapi/topics"
	"gopkg.in/bulletind/khabar.v1/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/simversity/gottp.v2"
)

type TopicChannel struct {
	gottp.BaseHandler
}

func (self *TopicChannel) Delete(request *gottp.Request) {
	intopic := new(topics.Topic)

	channelIdent := request.GetArgument("channel").(string)

	if !core.IsChannelAvailable(channelIdent) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	intopic.Ident = request.GetArgument("ident").(string)

	request.ConvertArguments(intopic)

	topic, err := topics.Get(
		intopic.User,
		intopic.Organization,
		intopic.Ident,
	)

	if err != nil && err != mgo.ErrNotFound {
		log.Println(err)
		request.Raise(gottp.HttpError{
			http.StatusInternalServerError,
			"Unable to fetch data, Please try again later.",
		})

		return

	}

	var hasData bool

	if topic == nil {
		log.Println("Creating new document")
		intopic.AddChannel(channelIdent)

		intopic.PrepareSave()
		if !intopic.IsValid(db.INSERT_OPERATION) {
			request.Raise(gottp.HttpError{
				http.StatusBadRequest,
				"Atleast one of the user or org must be present.",
			})

			return
		}

		if !utils.ValidateAndRaiseError(request, intopic) {
			log.Println("Validation Failed")
			return
		}

		topic = intopic

	} else {
		hasData = true

		for _, ident := range topic.Channels {
			if ident == channelIdent {
				request.Raise(gottp.HttpError{
					http.StatusConflict,
					"You have already unsubscribed this channel",
				})

				return
				break
			}
		}

		topic.AddChannel(channelIdent)
	}

	if hasData {
		err = topics.Update(
			topic.User,
			topic.Organization,
			topic.Ident,
			&utils.M{"channels": topic.Channels},
		)

		if err != nil {
			log.Println("Error while inserting document :" + err.Error())
			request.Raise(gottp.HttpError{
				http.StatusInternalServerError,
				"Internal server error.",
			})

			return
		} else {
			request.Write(utils.R{
				Data:       nil,
				Message:    "true",
				StatusCode: http.StatusNoContent,
			})

			return
		}
	} else {
		log.Println("Successfull call: Inserting document")
		topics.Insert(topic)
		request.Write(utils.R{
			Data:       nil,
			Message:    "true",
			StatusCode: http.StatusNoContent,
		})

		return
	}
}

func (self *TopicChannel) Post(request *gottp.Request) {
	topic := new(topics.Topic)

	channelIdent := request.GetArgument("channel").(string)
	topic.Ident = request.GetArgument("ident").(string)

	if !core.IsChannelAvailable(channelIdent) {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Channel is not supported",
		})

		return
	}

	request.ConvertArguments(topic)

	topic, err := topics.Get(
		topic.User,
		topic.Organization,
		topic.Ident,
	)

	if err != nil {
		if err != mgo.ErrNotFound {
			log.Println(err)
			request.Raise(gottp.HttpError{
				http.StatusInternalServerError,
				"Unable to fetch data, Please try again later.",
			})

		} else {
			request.Write(utils.R{
				Data:       nil,
				Message:    "true",
				StatusCode: http.StatusNoContent,
			})
		}

		return
	}

	topic.RemoveChannel(channelIdent)
	log.Println(topic.Channels)

	if len(topic.Channels) == 0 {
		log.Println("Deleting from database, since channels are now empty.")
		err = topics.Delete(

			&utils.M{
				"org":   topic.Organization,
				"user":  topic.User,
				"ident": topic.Ident,
			},
		)

	} else {
		log.Println("Updating...")

		err = topics.Update(
			topic.User, topic.Organization,
			topic.Ident, &utils.M{"channels": topic.Channels},
		)
	}

	if err != nil {
		request.Raise(gottp.HttpError{
			http.StatusInternalServerError,
			"Unable to delete.",
		})

		return
	}

	request.Write(utils.R{
		Data:       nil,
		Message:    "true",
		StatusCode: http.StatusNoContent,
	})

	return
}
