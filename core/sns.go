package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/bulletind/khabar/db"
	"github.com/bulletind/khabar/dbapi/devices"
	"github.com/bulletind/khabar/dbapi/saved_item"
	push "github.com/changer/pushnotification"
	"github.com/streamrail/concurrent-map"
)

const (
	PUSH_SOUND_SNS = "default"
)

var snsSettings cmap.ConcurrentMap

func init() {
	snsSettings = cmap.New()
}

func snsHandler(item *db.PendingItem, text string, locale string, appName string) {
	log.Println("Sending Push Notification using SNS...")
	service := getService(appName)

	subject, ok := item.Context["subject"].(string)
	if !ok || subject == "" {
		subject = item.Topic
	}

	custom := getCustomData(item)
	custom["data"] = getCustomData(item)

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

		data := &push.Data{
			Alert:   aws.String(text),
			Subject: aws.String(subject),
			Sound:   aws.String(PUSH_SOUND_SNS),
			Badge:   aws.Int(1),
			Data:    custom,
		}

		if len(getSound(userDevice.Type)) > 0 {
			data.Sound = aws.String(getSound(userDevice.Type))
		}

		// custom apps have separate certificate for iOS
		if userDevice.AppVariant != "" && userDevice.Type == "ios" {
			pushDevice.Type = handleAppVariantType(service, appName, userDevice.AppVariant)
		}

		err = service.Send(pushDevice, data)

		if err != nil {
			// try again once
			err = service.Send(pushDevice, data)
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Printf("Sent to token '%v' (%v) - app variant '%v' - '%v'", userDevice.Token, dbDevice.Type, userDevice.AppVariant, pushDevice.Type)
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

func getService(appName string) *push.Service {
	// try to return existing settings 1st
	if settings, ok := snsSettings.Get(appName); ok {
		return settings.(*push.Service)
	}

	awsKey := getKey("SNS_KEY", true)
	awsSecret := getKey("SNS_SECRET", true)
	awsRegion := getKey("SNS_REGION", true)

	base := getKey("SNS_BASE", true)   // 'arn:aws:sns:eu-west-1:123454678:app'
	env := getKey("ENVIRONMENT", true) // 'testing'

	//arn:aws:sns:eu-west-1:123454678:app/APNS/testing_inspectionApp
	awsAPNS := fmt.Sprintf("%v/%v/%v_%v", base, "APNS", env, appName)
	awsGCM := fmt.Sprintf("%v/%v/%v_%v", base, "GCM", env, appName)

	pushService := &push.Service{
		Key:    awsKey,
		Secret: awsSecret,
		Region: awsRegion,
		APNS:   awsAPNS,
		//APNSSandbox: awsAPNSSandbox,
		GCM:       awsGCM,
		Platforms: map[string]string{},
	}

	snsSettings.Set(appName, pushService)
	return pushService
}

func getSound(deviceType string) string {
	// try to return existing settings 1st
	if settings, ok := snsSettings.Get(deviceType); ok {
		return settings.(string)
	}

	sound := getKey("SNS_SOUND_"+strings.ToUpper(deviceType), false)
	snsSettings.Set(deviceType, sound)
	return sound
}

/// return special devicetype to target the right platform
func handleAppVariantType(service *push.Service, appName string, appVariant string) string {
	// custom apps have separate certificate
	base := getKey("SNS_BASE", true) // 'arn:aws:sns:eu-west-1:123454678:app'
	env := getKey("ENVIRONMENT", true)
	specialType := "ios_" + appVariant
	if _, ok := service.Platforms[specialType]; !ok {
		service.Platforms[specialType] = fmt.Sprintf("%v/%v/%v_%v_%v", base, "APNS", env, appName, appVariant)
	}
	return specialType
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
