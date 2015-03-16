package handlers

import (
	"net/http"

	"github.com/changer/khabar/db"

	statsApi "github.com/changer/khabar/dbapi/stats"
	"gopkg.in/simversity/gottp.v2"
	gottp_utils "gopkg.in/simversity/gottp.v2/utils"
)

func ConcatenateErrors(errs *[]error) string {
	var errString string
	for i := 0; i < len(*errs); i++ {
		errString += (*errs)[i].Error()
		if (len(*errs) - i) > 1 {
			errString += ","
		}
	}
	return errString
}

type Stats struct {
	gottp.BaseHandler
}

func (self *Stats) Get(request *gottp.Request) {
	var args struct {
		Organization string `json:"org"`
		AppName      string `json:"app_name"`
		User         string `json:"user" required:"true"`
	}

	request.ConvertArguments(&args)

	err := gottp_utils.Validate(&args)
	if len(*err) > 0 {
		request.Raise(gottp.HttpError{http.StatusBadRequest, ConcatenateErrors(err)})
		return
	}

	stats := statsApi.Get(db.Conn, args.User, args.AppName, args.Organization)

	request.Write(stats)

}

func (self *Stats) Post(request *gottp.Request) {
	var args struct {
		Organization string `json:"org"`
		AppName      string `json:"app_name"`
		User         string `json:"user" required:"true"`
	}

	request.ConvertArguments(&args)

	err := gottp_utils.Validate(&args)
	if len(*err) > 0 {
		request.Raise(gottp.HttpError{http.StatusBadRequest, ConcatenateErrors(err)})
		return
	}

	statsApi.Save(db.Conn, args.User, args.AppName, args.Organization)

	request.Raise(gottp.HttpError{http.StatusCreated, "Created"})
}
