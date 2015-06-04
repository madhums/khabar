package handlers

import (
	"testing"

	"gopkg.in/simversity/gottp.v3/utils"
)

func TestBounceMessage(t *testing.T) {
	jsonStr := `
	{
		"Type" : "Notification",
		"MessageId" : "f6cdc099-1020-554c-9f24-8d3f8067f821",
		"TopicArn" : "arn:aws:sns:eu-west-1:646538960088:gallery-upload",
		"Subject" : "Amazon S3 Notification",
		"Message" : "{\"notificationType\":\"Bounce\",\"bounce\":{\"bounceSubType\":\"General\",\"bounceType\":\"Transient\",\"bouncedRecipients\":[{\"emailAddress\":\"email_to_check\"}],\"timestamp\":\"2015-05-06T14:01:11.000Z\",\"feedbackId\":\"0000014d2987c258-5a97950e-f3f8-11e4-ad72-eb371f28d002-000000\"},\"mail\":{\"timestamp\":\"2015-05-06T13:53:59.000Z\",\"source\":\"notify@siminars.com\",\"messageId\":\"0000014d2981266d-970f4d92-5c32-460d-8f7e-77bd5fd5c235-000000\",\"destination\":[\"waawesome@gmail.com\"]}}",
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

	msg := bounceMessage{}
	utils.Decoder([]byte(n.Message), &msg)

	errs = utils.Validate(&msg)

	if len(*errs) != 0 {
		t.Error("Could not parse Records")
		return
	}

	bounce := msg.Bounce.Recipients[0]

	if bounce.Email != email {
		t.Error("Email should have been", email, "Found", bounce.Email)
	}
}
