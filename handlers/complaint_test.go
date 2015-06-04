package handlers

import (
	"testing"

	"gopkg.in/simversity/gottp.v3/utils"
)

func TestComplaintMessage(t *testing.T) {
	jsonStr := `
	{
		"Type" : "Notification",
		"MessageId" : "f6cdc099-1020-554c-9f24-8d3f8067f821",
		"TopicArn" : "arn:aws:sns:eu-west-1:646538960088:gallery-upload",
		"Subject" : "Amazon S3 Notification",
		"Message" : "{\"notificationType\":\"Complaint\",\"complaint\":{\"complainedRecipients\":[{\"emailAddress\":\"email_to_check\"}],\"timestamp\":\"2015-04-20T22:04:04.000Z\",\"feedbackId\":\"0000014cd8dc1928-2916b250-e7a9-11e4-b632-1f8199fde76d-000000\"},\"mail\":{\"timestamp\":\"2015-04-20T21:48:25.000Z\",\"source\":\"notify@siminars.com\",\"messageId\":\"0000014cd8cdc387-a26f1294-af94-4470-9449-14cf425460ec-000000\",\"destination\":[\"diegoromo100@hotmail.com\"]}}",
		"Timestamp" : "2015-04-14T03:48:23.584Z",
		"SignatureVersion" : "1",
		"Signature" : "liP1M+gnXDSo5A4mZJ/lO8Ah0rsC1ThfU0cmU5QmLezGB/VRq5G9V1QObO5phohsWLhMZiNTLVDWe9KCg9zKx+/X1S880Ytjd+Dyj4y1G29zATG3hzuRI1Ernp0dqHyIMwvbLrh6mqge65EPA/dzWVUjIehlGnLCeM9fSWrHqpPdyCT0egeC21eA98TxCvs5aWoND9pIcfUh0zSH6J7CT+QxEcjBKIb2dHhARdE75lrfDyM5QkVg6kEvQ/M9LEJExXC5KCXHVpKlsEQwL/qU4YKSTlIzkU2RpJrbMEtNjWFottBxr5WzkV98/CxkxjysEXtFW7xF7kyVsGOjmkKzHQ=="
	}
	`

	var errs *[]error

	email := "email_to_check"

	n := snsNotice{}
	utils.Decoder([]byte(jsonStr), &n)

	errs = utils.Validate(&n)

	if len(*errs) != 0 {
		t.Error("Could not parse JSON message")
		return
	}

	msg := complaintMessage{}
	utils.Decoder([]byte(n.Message), &msg)

	errs = utils.Validate(&msg)

	if len(*errs) != 0 {
		t.Error("Could not parse Records")
		return
	}

	complaint := msg.Complaint.Recipients[0]

	if complaint.Email != email {
		t.Error("Email should have been", email, "Found", complaint.Email)
	}
}
