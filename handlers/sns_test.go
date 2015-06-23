package handlers

import (
	"testing"

	"gopkg.in/simversity/gottp.v3/utils"
)

func TestSNSNoticeFail(t *testing.T) {
	jsonStr := `{
		"Type" : "Notification",
		"MessageId" : "f6cdc099-1020-554c-9f24-8d3f8067f821",
		"TopicArn" : "arn:aws:sns:eu-west-1:646538960088:gallery-upload",
		"Subject" : "Amazon S3 Notification"
	}`

	n := snsNotice{}
	utils.Decoder([]byte(jsonStr), &n)

	errs := utils.Validate(&n)
	if len(*errs) == 0 {
		t.Error("Message should be absent.")
		return
	}
}
