package handlers

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/simversity/gottp.v2"

	"github.com/changer/khabar/db"
	sentApi "github.com/changer/khabar/dbapi/sent"
	"github.com/changer/khabar/utils"
)

type Notification struct {
	gottp.BaseHandler
}

func (self *Notification) Put(request *gottp.Request) {
	sent_item := new(sentApi.SentItem)
	_id := request.GetArgument("_id").(string)

	if !bson.IsObjectIdHex(_id) {
		request.Raise(gottp.HttpError{http.StatusBadRequest, "_id is not a valid Hex object."})
		return
	}

	sent_item.Id = bson.ObjectIdHex(_id)
	sentApi.Update(db.Conn, sent_item.Id, &utils.M{"is_read": true})

	request.Raise(gottp.HttpError{http.StatusNoContent, "True"})
}
