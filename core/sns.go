package core

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/devices"
	"github.com/bulletind/khabar/dbapi/saved_item"
	push "github.com/changer/pushnotification"
)

const (
	PUSH_SOUND_SNS = "default"
)

func snsHandler(item *db.PendingItem, text string, locale string, appName string) {
	log.Println("Sending Push Notification using SNS...")
	service := getService(appName)

	subject, ok := item.Context["subject"].(string)
	if !ok || subject == "" {
		subject = item.Topic
	}

	custom := getCustomData(item)
	custom["data"] = getCustomData(item)

	data := &push.Data{
		Alert:   aws.String(text),
		Subject: aws.String(subject),
		Sound:   aws.String(PUSH_SOUND_SNS),
		Badge:   aws.Int(1),
		Data:    custom,
	}

	for _, userDevice := range item.DeviceTokens {
		// only send if appname is provided or app name is same
		if len(userDevice.AppName) == 0 || userDevice.AppName != appName {
			continue
		}

		exists := false
		err, dbDevice := devices.Get(userDevice.Token)

		if err != nil {
			dbDevice.Token = userDevice.Token
			dbDevice.Type = userDevice.Type
		} else {
			exists = true
		}

		pushDevice := &push.Device{
			Token:       dbDevice.Token,
			Type:        dbDevice.Type,
			EndpointArn: dbDevice.EndpointArn,
		}

		err = service.Send(pushDevice, data)

		if err != nil {
			log.Println(err)
		} else {
			dbDevice.EndpointArn = pushDevice.EndpointArn
			if !exists {
				devices.Insert(dbDevice)
			} else if pushDevice.IsCreated() {
				devices.Update(dbDevice)
			}
		}
	}

	store(item, text, subject)
}

func getCustomData(item *db.PendingItem) map[string]interface{} {
	return map[string]interface{}{
		"entity":       item.Entity,
		"organization": item.Organization,
		"app_name":     item.AppName,
		"topic":        item.Topic,
		"created_on":   item.CreatedOn,
	}
}

func getService(appName string) push.Service {
	awsKey := getKey("SNS_KEY", true)
	awsSecret := getKey("SNS_SECRET", true)
	awsRegion := getKey("SNS_REGION", true)
	awsAPNS := getKey("SNS_APNS_"+appName, false)
	//awsAPNSSandbox := getKey("SNS_APNSSANDBOX_"+appName, false)
	awsGCM := getKey("SNS_GCM_"+appName, false)

	return push.Service{
		Key:    awsKey,
		Secret: awsSecret,
		Region: awsRegion,
		APNS:   awsAPNS,
		//APNSSandbox: awsAPNSSandbox,
		GCM: awsGCM,
	}
}

func getKey(key string, required bool) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		if required {
			log.Fatal(key, " is empty. Make sure you set this env variable")
		} else {
			log.Print(key, " is empty. Plz check if this was intended")
		}
	}
	return value
}

func store(item *db.PendingItem, text string, subject string) {
	body := map[string]interface{}{}
	body["alert"] = subject
	body["title"] = subject
	body["message"] = text
	body["entity"] = item.Entity
	body["organization"] = item.Organization
	body["app_name"] = item.AppName
	body["topic"] = item.Topic
	body["created_on"] = item.CreatedOn
	body["sound"] = PUSH_SOUND_SNS
	body["badge"] = 1
	body["devices"] = item.DeviceTokens

	data := map[string]interface{}{}
	data["data"] = body
	data["channels"] = []string{"USER_" + item.User}

	saved_item.Insert(db.SavedPushCollection, &db.SavedItem{Data: data, Details: *item})
}
