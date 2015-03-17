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
	"gopkg.in/simversity/gottp.v2"
	gottp_utils "gopkg.in/simversity/gottp.v2/utils"
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

	all := sentApi.GetAll(db.Conn, paginator, args.User, args.AppName,
		args.Organization)

	request.Write(all)
}

func (self *Notifications) Put(request *gottp.Request) {
	var args struct {
		Organization string `json:"org"`
		AppName      string `json:"app_name"`
		User         string `json:"user" required:"true"`
	}

	request.ConvertArguments(&args)

	sentApi.MarkRead(db.Conn, args.User, args.AppName,
		args.Organization)

	request.Raise(gottp.HttpError{http.StatusNoContent, "True"})
}

func (self *Notifications) Post(request *gottp.Request) {
	ntfInst := new(pending.PendingItem)
	request.ConvertArguments(ntfInst)

	if request.GetArgument("topic") == nil {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Please provide topic for notification."})
		return
	}

	ntfInst.Topic = request.GetArgument("topic").(string)
	ntfInst.IsRead = false

	ntfInst.PrepareSave()

	if !utils.ValidateAndRaiseError(request, ntfInst) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Entity is invalid."})
		return
	}

	if !ntfInst.IsValid() {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "Context is required while inserting."})
		return
	}

	topic := topics.Find(db.Conn, ntfInst.User, ntfInst.AppName, ntfInst.Organization, ntfInst.Topic)

	if topic == nil {
		log.Println("Unable to find suitable notification setting.")
		request.Raise(gottp.HttpError{http.StatusNotFound, "Unable to find suitable notification setting."})
		return
	}

	core.SendNotification(db.Conn, ntfInst, topic)
	request.Raise(gottp.HttpError{http.StatusCreated, string(gottp_utils.Encoder(ntfInst))})
}
