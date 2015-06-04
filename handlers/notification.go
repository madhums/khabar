package handlers

import (
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/simversity/gottp.v3"

	"gopkg.in/bulletind/khabar.v1/db"
	sentApi "gopkg.in/bulletind/khabar.v1/dbapi/sent"
	"gopkg.in/bulletind/khabar.v1/utils"
)

type Notification struct {
	gottp.BaseHandler
}

func (self *Notification) Put(request *gottp.Request) {
	sent_item := new(db.SentItem)
	_id := request.GetArgument("_id").(string)

	if !bson.IsObjectIdHex(_id) {
		request.Raise(gottp.HttpError{http.StatusBadRequest,
			"_id is not a valid Hex object."})
		return
	}

	sent_item.Id = bson.ObjectIdHex(_id)

	err := sentApi.Update(sent_item.Id, &utils.M{"is_read": true})

	if err != nil {
		log.Println(err)
		request.Raise(gottp.HttpError{http.StatusInternalServerError,
			"Unable to insert."})
		return
	}

	request.Write(utils.R{StatusCode: http.StatusNoContent, Data: nil,
		Message: "NoContent"})
	return
}
