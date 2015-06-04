package handlers

import (
	"log"
	"net/http"

	"gopkg.in/bulletind/khabar.v1/utils"
	"gopkg.in/simversity/gottp.v3"
	gottpUtils "gopkg.in/simversity/gottp.v3/utils"
)

type SnsComplaint struct {
	gottp.BaseHandler
}

const ComplaintNotification = "Complaint"

type complaintMessage struct {
	Type      string `json:"notificationType" required:"true"`
	Complaint struct {
		Recipients []struct {
			Email string `json:"emailAddress" required:"true"`
		} `json:"complainedRecipients" required:"true"`
	} `json:"complaint" required:"true"`
}

func (self *SnsComplaint) Post(request *gottp.Request) {

	args := new(snsNotice)

	request.ConvertArguments(&args)

	if !utils.ValidateAndRaiseError(request, args) {
		log.Println("Invalid Request", request.GetArguments())
		return
	}

	msg := complaintMessage{}
	gottpUtils.Decoder([]byte(args.Message), &msg)

	errs := gottpUtils.Validate(&msg)
	if len(*errs) > 0 {
		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			ConcatenateErrors(errs),
		})

		return
	}

	if msg.Type != ComplaintNotification {
		log.Println("Invalid Complaint Request", request.GetArguments())

		request.Raise(gottp.HttpError{
			http.StatusBadRequest,
			"Invalid Complaint Request",
		})

		return
	}

	for _, entry := range msg.Complaint.Recipients {
		DisableBounceEmail(entry.Email, request)
	}

	request.Write(utils.R{
		StatusCode: http.StatusOK,
		Data:       nil,
	})

}
