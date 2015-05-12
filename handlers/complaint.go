package handlers

import (
	"log"
	"net/http"

	"github.com/bulletind/khabar/utils"
	"gopkg.in/simversity/gottp.v2"
	gottpUtils "gopkg.in/simversity/gottp.v2/utils"
)

type Complaint struct {
	gottp.BaseHandler
}

const ComplaintNotification = "Complaint"

func (self *Complaint) Post(request *gottp.Request) {

	type snsNotice struct {
		Type      string `json:"Type"`
		MessageId string `json:"MessageId" required:"true"`
		TopicArn  string `json:"TopicArn"`
		Subject   string `json:"Subject"`
		Message   string `json:"Message" required:"true"`
		Timestamp string `json:"Timestamp"`
		Signature string `json:"Signature" required:"true"`
	}

	type complaintMessage struct {
		Type      string `json:"notificationType" required:"true"`
		Complaint struct {
			Recipients []struct {
				Email string `json:"emailAddress" required:"true"`
			} `json:"complainedRecipients" required:"true"`
		} `json:"complaint" required:"true"`
	}

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
