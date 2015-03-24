package handlers

import (
	"log"
	"net/http"

	"github.com/changer/khabar/core"
	"github.com/changer/khabar/db"
	"github.com/changer/khabar/dbapi/topics"

	"github.com/changer/khabar/dbapi/pending"
	sentApi "github.com/changer/khabar/dbapi/sent"
	"github.com/changer/khabar/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/simversity/gottp.v2"
)

type Notifications struct {
	gottp.BaseHandler
}

func (self *Notifications) Get(request *gottp.Request) {
	var args struct {
		Organization string `json:"org"`
		AppName      string `json:"app_name"`
		User         string `json:"user" required:"true"`
	}

	request.ConvertArguments(&args)
	paginator := request.GetPaginator()

	all, err := sentApi.GetAll(db.Conn, paginator, args.User, args.AppName,
		args.Organization)

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

func (self *Notifications) Put(request *gottp.Request) {
	var args struct {
		Organization string `json:"org"`
		AppName      string `json:"app_name"`
		User         string `json:"user" required:"true"`
	}

	request.ConvertArguments(&args)

	err := sentApi.MarkRead(db.Conn, args.User, args.AppName,
		args.Organization)

	if err != nil {
		log.Println(err)
		request.Raise(gottp.HttpError{
			http.StatusInternalServerError,
			"Unable to insert.",
		})

		return
	}

	request.Write(utils.R{StatusCode: http.StatusNoContent, Data: nil, Message: "NoContent"})
	return
}

func (self *Notifications) Post(request *gottp.Request) {
	pending := new(pending.PendingItem)
	request.ConvertArguments(pending)

	if request.GetArgument("topic") == nil {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Please provide topic for notification.",
		})

		return
	}

	pending.Topic = request.GetArgument("topic").(string)
	pending.IsRead = false

	pending.PrepareSave()

	if !utils.ValidateAndRaiseError(request, pending) {
		return
	}

	if !pending.IsValid() {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Context is required while inserting.",
		})

		return
	}

	topic, err := topics.Find(db.Conn, pending.User, pending.AppName, pending.Organization, pending.Topic)
	if err != nil {
		if err != mgo.ErrNotFound {
			log.Println(err)
			request.Raise(gottp.HttpError{
				http.StatusInternalServerError,
				"Unable to fetch data, Please try again later.",
			})

		} else {
			request.Raise(gottp.HttpError{http.StatusNotFound, "Not Found."})
		}
		return
	}

	core.SendNotification(db.Conn, pending, topic)
	request.Write(utils.R{StatusCode: http.StatusCreated, Data: topic.Id, Message: "Created"})
	return
}
