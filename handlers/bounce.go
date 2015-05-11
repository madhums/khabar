package handlers

import (
	"log"
	"net/http"

	"github.com/bulletind/khabar/core"
	"github.com/bulletind/khabar/dbapi/available_topics"
	"github.com/bulletind/khabar/dbapi/saved_item"
	"github.com/bulletind/khabar/dbapi/topics"
	"github.com/bulletind/khabar/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/simversity/gottp.v2"
	gottpUtils "gopkg.in/simversity/gottp.v2/utils"
)

type Bounce struct {
	gottp.BaseHandler
}

const BounceNotification = "Bounce"

func (self *Bounce) Post(request *gottp.Request) {

	type snsNotice struct {
		Type      string `json:"Type"`
		MessageId string `json:"MessageId" required:"true"`
		TopicArn  string `json:"TopicArn"`
		Subject   string `json:"Subject"`
		Message   string `json:"Message" required:"true"`
		Timestamp string `json:"Timestamp"`
		Signature string `json:"Signature" required:"true"`
	}

	type bounceMessage struct {
		Type   string `json:"notificationType" required:"true"`
		Bounce struct {
			Recipients []struct {
				Email string `json:"emailAddress" required:"true"`
			} `json:"bouncedRecipients" required:"true"`
		} `json:"bounce" required:"true"`
	}

	args := new(snsNotice)

	request.ConvertArguments(&args)

	if !utils.ValidateAndRaiseError(request, args) {
		log.Println("Invalid Request", request.GetArguments())
		return
	}

	msg := bounceMessage{}
	gottpUtils.Decoder([]byte(args.Message), &msg)

	errs := gottpUtils.Validate(&msg)
	if len(*errs) > 0 {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			ConcatenateErrors(errs),
		})

		return
	}

	if msg.Type != BounceNotification {
		log.Println("Invalid Bounce Request", request.GetArguments())

		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Invalid Bounce Request",
		})

		return
	}

	for _, entry := range msg.Bounce.Recipients {
		if !DisableBounceEmail(entry.Email, request) {
			break
		}
	}

	request.Write(utils.R{
		StatusCode: http.StatusOK,
		Data:       nil,
	})

}

func DisableBounceEmail(email string, request *gottp.Request) bool {
	notification, err := saved_item.Get("saved_"+core.EMAIL,
		&utils.M{"details.context.email": email})

	if err != nil {
		if err == mgo.ErrNotFound {
			log.Println("Bounced Email not found")
			return true
		} else {
			request.Raise(gottp.HttpError{
				http.StatusInternalServerError,
				"Unable to fetch data, Please try again later.",
			})
			return false
		}
	}

	sentItem := notification.Details

	topicList := available_topics.GetAppTopics(sentItem.AppName, sentItem.Organization)

	for _, topic := range *topicList {
		disabled, err := topics.Get(sentItem.User, sentItem.Organization, topic)
		if err != nil && err != mgo.ErrNotFound {
			request.Raise(gottp.HttpError{
				http.StatusInternalServerError,
				"Unable to fetch data, Please try again later.",
			})
			return false
		}

		if disabled != nil {
			disabled.Channels = append(disabled.Channels, core.EMAIL)
			utils.RemoveDuplicates(&disabled.Channels)
			err := topics.Update(sentItem.User, sentItem.Organization, topic,
				&utils.M{"channels": disabled.Channels})
			if err != nil {
				request.Raise(gottp.HttpError{
					http.StatusInternalServerError,
					"Unable to fetch data, Please try again later.",
				})
				return false
			}
		} else {
			disabled = &topics.Topic{User: sentItem.User, Organization: sentItem.Organization,
				Ident: topic, Channels: []string{core.EMAIL}}
			disabled.PrepareSave()
			topics.Insert(disabled)
		}
	}
	return true
}
