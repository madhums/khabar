package models

type ChannelConfiguration struct {
	UserID         string      `json:"user_id"`
	OrganizationID string      `json:"org_id"`
	ApplicationID  string      `json:"app_id"`
	ChannelData    interface{} `json:"channel_data"`
}
