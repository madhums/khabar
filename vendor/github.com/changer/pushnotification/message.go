package pushnotification

import "encoding/json"

type message struct {
	APNS        string `json:"APNS"`
	APNSSandbox string `json:"APNS_SANDBOX"`
	Default     string `json:"default"`
	GCM         string `json:"GCM"`
}

type iosPush struct {
	APS iosAps `json:"aps"`
}

type iosAps struct {
	Alert iosAlert `json:"alert,omitempty"`
	Sound *string  `json:"sound,omitempty"`
	Badge *int     `json:"badge,omitempty"`
}

type iosAlert struct {
	Body  *string `json:"body,omitempty"`
	Title *string `json:"title,omitempty"`
}

type gcmPush struct {
	Message *string     `json:"message,omitempty"`
	Title   *string     `json:"title,omitempty"`
	Custom  interface{} `json:"custom"`
	Badge   *int        `json:"badge,omitempty"`
}

type gcmPushWrapper struct {
	Data gcmPush `json:"data"`
}

func newMessageJSON(data *Data) (m string, err error) {
	b, err := json.Marshal(iosPush{
		APS: iosAps{
			Alert: iosAlert{
				Body:  data.Alert,
				Title: data.Subject,
			},
			Sound: data.Sound,
			Badge: data.Badge,
		},
	})
	if err != nil {
		return
	}
	payload := string(b)

	b, err = json.Marshal(gcmPushWrapper{
		Data: gcmPush{
			Message: data.Alert,
			Custom:  data.Data,
			Badge:   data.Badge,
			Title:   data.Subject,
		},
	})
	if err != nil {
		return
	}
	gcm := string(b)

	pushData, err := json.Marshal(message{
		Default:     *data.Alert,
		APNS:        payload,
		APNSSandbox: payload,
		GCM:         gcm,
	})
	if err != nil {
		return
	}
	m = string(pushData)
	return
}
