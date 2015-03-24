package handlers

import (
	"net/http"

	"log"

	statsApi "github.com/changer/khabar/dbapi/stats"
	"github.com/changer/khabar/utils"
	"gopkg.in/mgo.v2"
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
	args := statsApi.RequestArgs{}

	request.ConvertArguments(&args)

	err := gottp_utils.Validate(&args)
	if len(*err) > 0 {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			ConcatenateErrors(err),
		})

		return
	}

	stats, getErr := statsApi.Get(&args)

	if getErr != nil {
		if getErr != mgo.ErrNotFound {
			log.Println(getErr)
			request.Raise(gottp.HttpError{
				http.StatusInternalServerError,
				"Unable to fetch data, Please try again later.",
			})

		} else {
			request.Raise(gottp.HttpError{http.StatusNotFound, "Not Found."})
		}
		return
	}

	request.Write(stats)
	return

}

func (self *Stats) Post(request *gottp.Request) {
	args := statsApi.RequestArgs{}

	request.ConvertArguments(&args)

	err := gottp_utils.Validate(&args)
	if len(*err) > 0 {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			ConcatenateErrors(err),
		})

		return
	}

	insErr := statsApi.Save(&args)

	if insErr != nil {
		log.Println(insErr)
		request.Raise(gottp.HttpError{
			http.StatusInternalServerError,
			"Unable to insert.",
		})

		return
	}

	request.Write(utils.R{Data: nil, Message: "Created", StatusCode: http.StatusCreated})
	return
}
