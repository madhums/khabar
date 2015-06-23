package handlers

import (
	"log"
	"net/http"

	"gopkg.in/bulletind/khabar.v1/dbapi/topics"
	"gopkg.in/bulletind/khabar.v1/utils"
	"gopkg.in/simversity/gottp.v3"
)

type Topic struct {
	gottp.BaseHandler
}

func (self *Topics) Delete(request *gottp.Request) {
	var args struct {
		Ident string `json:"ident" required:"true"`
	}

	request.ConvertArguments(&args)

	if !utils.ValidateAndRaiseError(request, args) {
		log.Println("Validation Failed")
		return
	}

	err := topics.DeleteTopic(args.Ident)
	if err != nil {
		log.Println(err)
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
