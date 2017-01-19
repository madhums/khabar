package handlers

import (
	"log"
	"net/http"

	"github.com/bulletind/khabar/core"

	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/pending"
	sentApi "github.com/bulletind/khabar/dbapi/sent"
	"github.com/bulletind/khabar/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/simversity/gottp.v3"
)

type Notifications struct {
	gottp.BaseHandler
}

// List all notifications
func (self *Notifications) Get(request *gottp.Request) {
	var args struct {
		Organization string `json:"org"`
		AppName      string `json:"app_name"`
		User         string `json:"user" required:"true"`
	}

	request.ConvertArguments(&args)
	paginator := request.GetPaginator()

	all, err := sentApi.GetAll(paginator, args.User, args.AppName, args.Organization)

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

// Mark all notifications as read
func (self *Notifications) Put(request *gottp.Request) {
	var args struct {
		Organization string `json:"org"`
		AppName      string `json:"app_name"`
		User         string `json:"user" required:"true"`
	}

	request.ConvertArguments(&args)

	err := sentApi.MarkRead(args.User, args.AppName, args.Organization)

	if err != nil {
		log.Println(err)
		request.Raise(gottp.HttpError{
			http.StatusInternalServerError,
			"Unable to insert.",
		})

		return
	}

	request.Write(utils.R{StatusCode: http.StatusNoContent,
		Data: nil, Message: "NoContent"})
	return
}

// Send a notification to the user depending on preferences
func (self *Notifications) Post(request *gottp.Request) {
	pending_item := new(db.PendingItem)
	request.ConvertArguments(pending_item)

	if request.GetArgument("topic") == nil {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Please provide topic for notification.",
		})

		return
	}

	if pending.Throttled(pending_item) {
		MSG := "Repeated Notifications are Blocked. Skipping."

		log.Println(MSG)
		request.Raise(gottp.HttpError{http.StatusBadRequest, MSG})
		return
	}

	pending_item.Topic = request.GetArgument("topic").(string)
	pending_item.IsRead = false

	pending_item.PrepareSave()

	if !utils.ValidateAndRaiseError(request, pending_item) {
		return
	}

	if !pending_item.IsValid() {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Context is required for sending notification.",
		})

		return
	}

	core.SendNotification(pending_item)
	request.Write(utils.R{StatusCode: http.StatusNoContent,
		Data: nil, Message: "true"})
	return
}
